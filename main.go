package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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
	fmt.Println(dryRun)

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

	// create the directory that the files will be downloaded to
	// if this is not a dry run
	if !dryRun {
		// TODO check if dir is empty string
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}
	}

	// list file items
	// TODO skip over directories
	for _, item := range fileList {
		if dryRun {
			fmt.Println(item.Name, item.Size)
			continue
		}

		// create an empty file in the sub directory for this file
		// TODO check if dir is empty string
		out, err := os.Create(fmt.Sprintf("%s/%s", dir, item.Name))
		if err != nil {
			return fmt.Errorf("Error creating file: %s: %v", item.Name, err)
		}
		defer out.Close()

		// get the download url
		res, err := http.Get(item.DownloadURL)
		if err != nil {
			return fmt.Errorf("Error doing get request on file: %s: %v", item.Name, err)
		}
		defer res.Body.Close()

		// if the response code is not ok, bail out
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Unexpected response code for file: %s: %v", item.Name, err)
		}

		// copy the contents to the empty file
		_, err = io.Copy(out, res.Body)
		if err != nil {
			return fmt.Errorf("Error saving file contents for file: %s: %v", item.Name, err)
		}
	}

	// all gucci
	return nil
}
