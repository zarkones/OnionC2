# Available Commands
The following commands are available in OnionC2.

Executes a shell command on the remote system without waiting for it to complete:
- /run|<COMMAND>

Downloads a file from C2 server to the remote system via Tor. Requires an operator’s download request to ensure security, preventing enumeration, as once a file is served it would require a new download request which is specific to an agent:
- /download-file|<FILE_NAME_IN_C2s_DOWNLOAD_DIRECTORY>

Uploads a file from the remote system to C2 via Tor. Secured by requiring a unique, long UUID for each file, making the upload endpoint difficult to abuse:
- /upload-file|<FILE_PATH>

Searches for files or directories on the remote system via search terms:
- /find-files|<STARTING_DIR_PATH>|<COMMA_SEPARATED_SEARCH_TERMS>

Returns detailed information about the remote system, including RAM memory, CPU temperature, network interfaces, and more:
- /system-details

Returns clipboard data or an error if the clipboard is empty or cannot be accessed:
- /read-clipboard

Any other input would be executed as a shell command and the agent would send back the result.