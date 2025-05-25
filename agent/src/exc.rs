use std::process::Command;
use tokio::io;

#[inline]
pub fn run(cmd: &String) -> io::Result<String> {
    let mut command = if cfg!(target_os = "windows") {
        let mut c = Command::new("cmd");
        c.args(&["/C", cmd]);
        c
    } else {
        let mut c = Command::new("sh");
        c.args(&["-c", cmd]);
        c
    };

    let output = command.output()?;

    let stdout = String::from_utf8(output.stdout)
        .map_err(|e| io::Error::new(io::ErrorKind::InvalidData, e))?;

    let stderr = String::from_utf8(output.stderr)
        .map_err(|e| io::Error::new(io::ErrorKind::InvalidData, e))?;

    Ok(format!("{}{}", stdout, stderr))
}