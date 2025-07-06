# 0.1.0
- Authentication migrated away from API keys to public-key scheme.
- Fine grained access control managemenet for operator accounts.
- New web based user interface.

# 0.0.14
- Upgrading Tor client version.

# 0.0.13
- Bug fix: hiding terminal window on shell command execution on Windows operating system.

# 0.0.12
- Introduction of optional active hours. Allowing an agent to communicate only within specific time frames.

# 0.0.11
- Introduction of a command "/read-clipboard" which returns clipboard data or an error if the clipboard is empty or cannot be accessed.

# 0.0.10
- Additional CPU information to "/system-details" command.

# 0.0.9
- Introduction of a command "/run|<COMMAND>" which executes a shell command, however, does not await it to finish.

# 0.0.8
- Introduction of a command of "/download-file|<FILE_NAME_IN_C2s_DOWNLOAD_DIRECTORY>" enabling file downloads via Tor. Requires operator's download request, meaning there is no public directory someone can enumerate.

# 0.0.7
- Introduction of a command "/upload-file|<FILE_PATH>" which uploads a file via Tor. It's difficult to abuse the file upload endpoint due to its limited attack surface as it requires a long UUID generated for each file.

# 0.0.6
- Introduction of a command "/find-files|<STARTING_DIR_PATH>|<COMMA_SEPARATED_SEARCH_TERMS>" which based on criteria returns absolute path of files/directories of interest. 

# 0.0.5
- Introduction of a command "/system-details" which returns information about RAM memory, CPU temperature, network interfaces, etc...

# 0.0.4
- Agent's binary size optimized. From 21.1MB down to 4.4MB for Windows build artifact.

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
