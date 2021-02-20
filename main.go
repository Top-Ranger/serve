// SPDX-License-Identifier: Apache-2.0
// Copyright 2017,2018,2019 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Top-Ranger/serve/modules"
)

const (
	fileTemplate = "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Files shared</title></head><body></body><h1>Files shared</h1><ul>{{range .}}<li><a href=\"./{{.}}\">{{.}}</a></li>{{end}}</ul></html>"
)

// Default options
var (
	port         = 8080
	getIP        = true
	showFilelist = true
)

//go:embed favicon.ico
var favicon []byte

func main() {
	flag.IntVar(&port, "port", port, "Port on which the server is listening")
	flag.BoolVar(&getIP, "getip", getIP, "Enables / disables public ip lookup")
	flag.BoolVar(&showFilelist, "filelist", showFilelist, "Enables / disables file list at root")

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalln("No files to serve - exiting")
	}

	faviconServed := false

	// Parse files
	args := flag.Args()
	fileSystem := modules.NewSelectedFileSystem(flag.NArg())
	for _, arg := range args {
		log.Println(filepath.Rel("./", arg))
		if f, _ := filepath.Rel("./", arg); f == "favicon.ico" {
			log.Println("Found 'favicon.ico' file - disabled serving favicon")
			faviconServed = true
		}
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
	if getIP {
		go func() {
			ip, err := modules.GetPublicIP()
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
		} else if !faviconServed && r.URL.EscapedPath() == "/favicon.ico" {
			w.WriteHeader(http.StatusOK)
			w.Write(favicon)
		} else {
			log.Println(r.RemoteAddr, ":", r.URL, "(", r.Method, ")")
			fileServer.ServeHTTP(w, r)
		}
	})))
}
