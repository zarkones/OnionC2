![Promo Image 1](https://raw.githubusercontent.com/zarkones/OnionC2/production/promo/promo1.png)

# INTRODUCTION
OnionC2 is a command and control (C2) framework with communications over Tor network. It's packed with privacy & security features, and operational capabilities. It is simple to setup, and has a friendly user interface. It is cross-platform and supports collaboration between operators.

# AGENT'S FEATURES
- Tor integration (allows for end to end encryption, hiding the C2's IP address)
- Execution of shell commands.
- Obfuscating C2 configuration in the agent's binary.
- Registry based persistence on Windows.
- Shortcut takeover based persistence on Windows.
- Active hours, allowing an agent to communicate only within specific time frames.
- Command "/system-details" makes an agent return information about CPU, RAM, networks, etc...
- Command "/find-files|<STARTING_DIR_PATH>|<COMMA_SEPARATED_SEARCH_TERMS>" which based on criteria returns absolute path of files/directories of interest.
- Command "/upload-file|<FILE_PATH>" which uploads a file via Tor.
- Command "/download-file|<FILE_NAME_IN_C2s_DOWNLOAD_DIRECTORY>" which downloads a file via Tor.
- Command "/run|<SHELL_COMMAND>" which executes shell command without awaiting it. 
- Command "/read-clipboard" which returns clipboard data.

# SETUP
This guide assumes you have Go, Rust, and Node ready.

## Back-end Setup

### Administrator Account
OnionC2 supports multiple operator accounts. In this setup guide, you'll learn how to create an over-powered administrator account. You should use this administrator account only when required. During day-to-day operations it is recommended to use an account with less permissions.

Navigate to "api" directory and run the following command: go run . --create-admin --username <YOUR_USERNAME>

This command would print out your account's recovery word phrase and its private key. Save it somewhere secure, as without the private key you won't be able to authenticate with the C2 API, and without the recovery word phrase you won't be able to recover your private key in case you lose it.

### API Setup
Back-end service is composed of two APIs. Agents-facing API is listening on a Unix socket, while the Operator-facing API is listening by default on port 8080. To see additional configuration possibilities run the API with "-h" argument.

To run the API navigate into the "api" directory an execute: go run .

This would automatically create SQLite database and perform all the needed database migrations. Also, it would create a file named "torrc", this file describes out Onion service and allows Tor to route traffic to our agents-facing API.

To run the Onion service run the following command inside of the "api" directory: tor -f torrc

## Front-End Setup
First you need to build the user interface, prior to serving it, in order to do so, execute the following command inside of the "ui" directory: npm run build

This would generate static HTML/JS/CSS files in directory ".output/public"

You can serve files from that directory using a web server of your choice, or use the one provided by the OnionC2 by running the following command inside of the "ui" directory: go run serve.go

## Agent Setup
All configuration related to the behavior of agents is located in a file "agent/src/config.rs". Basic configuration requires you to set at least the Onion domain in the function named "get_address". Domain of your Onion service is located in "api/onionservice/hostname".

To build an agent run the following script inside of the "agent" directory: sh build.sh

To configure persistence or any other option refer to comments inside of the "config.rs" file. It's well documented with code comments. If something isn't clear reach out to our Discord server.

# SOCIAL
[Patreon](https://www.patreon.com/zarkones) |
[Discord](https://discord.gg/qjJwSh2TF9) |
[X.com](https://x.com/zarkones) |
[YouTube](https://www.youtube.com/channel/UCn-7I-L-ZpiELb8-6z7z_Ug) |
[Itch.io](https://zarkones.itch.io) |
[GitHub](https://github.com/zarkones)