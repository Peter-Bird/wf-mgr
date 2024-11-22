package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"wf-mgr/handlers"
)

/*
	Usage:
	------
	go run main.go -port=8084
*/

func init() {
	// Set log prefix based on application name from args
	appName := filepath.Base(os.Args[0]) // Removes leading "./" if present
	log.SetPrefix("[" + appName + "] ")
	//log.SetFlags(0) // Optional: removes default date and time from log output
}

func main() {

	port := flag.String("port", "8084", "Port to run the server on")
	flag.Parse()

	address := fmt.Sprintf(":%s", *port)
	http.HandleFunc("/exec/", handlers.ExecHandler)

	log.Printf("Server running on port: %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
