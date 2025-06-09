use crate::config;

use systemstat::{Platform, System};
use tokio::io::{
    AsyncWriteExt,
    AsyncReadExt,
    Result,
    BufReader,
};
use std::{collections::BTreeMap, env, path::PathBuf};
use tokio::fs::File;
use goldberg::goldberg_stmts;
use tokio::time::{sleep, Duration};

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

    let information = SystemInformation{
        memory: memory,
        is_ac_power: is_on_ac_power,
        uptime_seconds: uptime,
        networks: network_names,
        cpu_temperature: cpu_temperature,
    };

    return information;
}