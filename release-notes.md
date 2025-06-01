# 0.0.3
- Optional agent mutex disallowing multiple instances running at the same time. (relies on unsafe Windows API)
- Shortcut based persistence mechanism.

# 0.0.2
- Classic registry based persistence mechanism for Windows operating system.

# v0.0.1
- Agent is able to execute of shell commands and send back the output.
- Agent is able to communicate with the C2 server via the Tor network.
- Obfuscation of agent's config.
- HTTP (unix sockets) based API interface for C2 communication.
- HTTP based API interface for integration with user interfaces with API key based authentication.
- SQLite based storage for the APIs.