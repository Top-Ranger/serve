package modules

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// The SelectedFileSystem represents a pseudo-filesystem which only allows the access to prior registered files.
type SelectedFileSystem struct {
	files []string
}

// Creates a SelectedFileSystem instance.
func NewSelectedFileSystem(capacity int) SelectedFileSystem {
	return SelectedFileSystem{make([]string, 0, capacity)}
}

// Adds a file to the SelectedFileSystem.
// Will do some basic checks and return a non nil error if the tests are failed.
func (sfs *SelectedFileSystem) AddFile(name string) error {
	if sfs.files == nil {
		sfs.files = make([]string, 0, 0)
	}
	name = cleanPathName(name)
	if strings.Contains(name, "../") {
		err := errors.New(fmt.Sprint("Error adding file ", name, ": Final path must not contain '../'"))
		log.Println(err)
		return err
	}

	// Check if file exists
	stat, err := os.Stat(name)
	if err != nil {
		log.Println("Error adding file", name, ":", err)
		return err
	}

	// Check if it is a regular file
	if !stat.Mode().IsRegular() {
		err := errors.New(fmt.Sprint("Error adding file ", name, ": File is not a regular file"))
		log.Println(err)
		return err
	}

	log.Println("Adding", url.PathEscape(name))

	sfs.files = append(sfs.files, name)
	return nil
}

// Returns a file if the name is known, else returns an error not nil.
func (sfs SelectedFileSystem) Open(name string) (http.File, error) {
	name = cleanPathName(name)
	for _, path := range sfs.files {
		if name == path {
			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			return file, nil
		}
	}
	return nil, os.ErrNotExist
}

// Cleans the path.
func cleanPathName(path string) string {
	path = filepath.Clean(path)
	for strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	for strings.HasPrefix(path, "./") {
		path = path[2:]
	}
	return path
}
