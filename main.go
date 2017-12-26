package main

import (
	"github.com/Top-Ranger/serve/modules"

	"flag"
	"fmt"
	"log"
	"net/http"
)

// Default options
var (
	port int = 8080
)

func main() {
	flag.IntVar(&port, "port", port, "Port on which the server is listening")

	flag.Parse()

	// Parse files
	args := flag.Args()
	fileSystem := modules.NewSelectedFileSystem(flag.NArg())
	for _, arg := range args {
		err := fileSystem.AddFile(arg)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Check parameter
	if port < 1024 || port > 65535 {
            log.Fatal("Please specify a valid port between 1024 and 65535")
        }

	// Prepare server
	portString := fmt.Sprint(":", port)
	fileServer := http.FileServer(fileSystem)

	// Output of server info
	log.Println("Server listening at", portString)

	log.Print(http.ListenAndServe(portString, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, ":", r.URL, "(", r.Method, ")")
		fileServer.ServeHTTP(w, r)
	})))
}
