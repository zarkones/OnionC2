package main

import (
	"flag"
	"net"
	"net/http"
)

func main() {
	port := flag.String("port", "3000", "Port on which to serve the application")
	host := flag.String("host", "", "Host on which to serve the application")
	pathToStaticFiles := flag.String("path", ".output/public", "Path to the built application and static files")

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*pathToStaticFiles)))
	http.ListenAndServe(net.JoinHostPort(*host, *port), nil)
}
