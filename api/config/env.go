package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	DbPath                 = flag.String("db-path", "api.db", "path to database")
	ApiSockPath            = flag.String("sock-path", "api.sock", "path to a unix socket of the c2 api")
	ApiHost                = flag.String("api-host", "127.0.0.1:8080", "<host>:<port> of the user facing api")
	OnionServicePath       = flag.String("onion-service-path", "./onionservice", "path to onion service's directory")
	UploadsDirectoryPath   = flag.String("uploads-dir-path", "./agent-upload-dir", "path to directory where agents would upload files an operator requests of them")
	DownloadsDirectoryPath = flag.String("downloads-dir-path", "./agent-download-dir", "path to directory where agents would download files from")
	AllowedOrigins         = flag.String("allowed-origins", "http://127.0.0.1:3000", "comma separated list of allowed origins for CORS requests, use '*' to allow all origins")

	// Insert user action.
	UserInsertAdmin = flag.Bool("create-admin", false, "creates an administrative operator account with all permissions, requires --username and --password arguments")
	UserName        = flag.String("username", "", "specify username, used by --create-admin")
)

func init() {
	// Parse command line arguments declared in "config" package.
	flag.Parse()
}

// Validate checks if command line arguments are valid enoguh to run the program.
// In a case they aren't then the program would print an error an exit.
func Validate() {
	if len(*ApiSockPath) == 0 {
		fmt.Println("api sock path is not specified")
		os.Exit(1)
	}
	if !strings.HasSuffix(*ApiSockPath, ".sock") {
		fmt.Println("api sock path must have '.sock' suffix")
		os.Exit(1)
	}
	if *AllowedOrigins == "*" {
		fmt.Println("warning: allowed-origins is set to '*'")
	}
}
