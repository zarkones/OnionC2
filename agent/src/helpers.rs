use crate::config;

use tokio::io::{
    AsyncWriteExt,
    AsyncReadExt,
    Result,
    BufReader,
};
use std::path::PathBuf;
use tokio::fs::File;
use goldberg::goldberg_stmts;
use tokio::time::{sleep, Duration};

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
        Ok(_) => {
            return Some(contents);
        },
        Err(_) => {
            return None;
        },
    };
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