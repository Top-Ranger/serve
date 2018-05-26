// Copyright (c) 2017,2018 Marcus Soll
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/Top-Ranger/serve/modules"

	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	fileTemplate = "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Files shared</title></head><body></body><h1>Files shared</h1><ul>{{range .}}<li><a href=\"./{{.}}\">{{.}}</a></li>{{end}}</ul></html>"
)

// Default options
var (
	port         int  = 8080
	getIp        bool = true
	showFilelist bool = true
)

func main() {
	flag.IntVar(&port, "port", port, "Port on which the server is listening")
	flag.BoolVar(&getIp, "getip", getIp, "Enables / disables public ip lookup")
	flag.BoolVar(&showFilelist, "filelist", showFilelist, "Enables / disables file list at root")

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalln("No files to serve - exiting")
	}

	// Parse files
	args := flag.Args()
	fileSystem := modules.NewSelectedFileSystem(flag.NArg())
	for _, arg := range args {
		err := fileSystem.AddFile(arg)
		if err != nil {
			log.Fatal("Fatal error - can not start server: ", err)
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
	log.Print("Server reachable at http://localhost", portString, "/")
	if getIp {
		go func() {
			ip, err := modules.GetPublicIp()
			if err != nil {
				log.Println("Can not look up public IP:", err)
			} else {
				log.Print("Server publicly reachable at http://", ip, portString, "/")
			}
		}()
	} else {
		log.Println("Public IP lookup disabled")
	}

	log.Print(http.ListenAndServe(portString, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if showFilelist && r.URL.EscapedPath() == "/" {
			log.Println(r.RemoteAddr, ":", r.URL, "(", r.Method, ")", "; Action: Showing file list")
			templateStruct := template.New("")
			templateStruct, err := templateStruct.Parse(fileTemplate)
			if err != nil {
				log.Panicln(err)
			}

			err = templateStruct.Execute(w, fileSystem.GetFiles())
			if err != nil {
				log.Panicln(err)
			}
		} else {
			log.Println(r.RemoteAddr, ":", r.URL, "(", r.Method, ")")
			fileServer.ServeHTTP(w, r)
		}
	})))
}
