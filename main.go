package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
)

var (
	dryRun bool
	repo   string
	user   string
	dir    string
)

// file represents a Github file
// see https://developer.github.com/v3/repos/contents/#response-if-content-is-a-file
type file struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"download_url"`
}

// set up the command line flags
func init() {
	flag.BoolVar(&dryRun, "dry-run", true, "Whether or not to execute the download")
	flag.StringVar(&repo, "repo", "", "The Github repository name")
	flag.StringVar(&user, "user", "", "The Github user")
	flag.StringVar(&dir, "dir", "", "The sub directory to fetch")
}

func main() {
	flag.Parse()

	err := getFiles()
	if err != nil {
		fmt.Println(err)
	}
}

// getFiles makes a request to the Github API to list and download the
// contents of a given sub directory in a Github repository.
// if dry run is true, the list of file names and sizes will be printed
// if dry run is false, these files will be downloaded to the current directory
func getFiles() error {
	if repo == "" || user == "" {
		return errors.New("Repository and username cannot be empty")
	}

	// list the files in the subdirectory
	path := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", user, repo, dir)
	resp, err := http.Get(path)
	if err != nil {
		return fmt.Errorf("Unable to list directory contents: %v", err)
	}
	defer resp.Body.Close()

	// decode the response body into a fileList struct
	var fileList []file
	err = json.NewDecoder(resp.Body).Decode(&fileList)
	if err != nil {
		return fmt.Errorf("Error decoding response body: %v", err)
	}

	// list file items
	for _, item := range fileList {
		if dryRun {
			fmt.Println(item.Name, item.Size)
			continue
		}
	}

	// all gucci
	return nil
}
