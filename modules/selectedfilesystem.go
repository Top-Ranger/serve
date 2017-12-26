package modules

import (
	"log"
	"net/http"
	"os"
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
	log.Println("Adding", name)

	_, err := os.Stat(name)
	if err != nil {
		log.Println("Error adding file", name, ":", err)
		return err
	}

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
// Currently only cleans leading slashes.
func cleanPathName(path string) string {
	for strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	return path
}
