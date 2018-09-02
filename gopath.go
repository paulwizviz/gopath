package gopath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IsValidName ensure that repo and path names follow the following pattern:
// (a) first character must only be an alphabet, and
// (b) subsequent characters can be either a ".", "_", "-", alphabets, numbers,
func isValidName(argName string) bool {

	if argName == "" {
		return false
	}

	regex := regexp.MustCompile("([a-zA-Z][\x2Ea-zA-Z0-9_-]*)")
	matches := regex.FindAllString(argName, 1)
	if len(matches) == 0 {
		return false
	} else if matches[0] != argName {
		return false
	}
	return true
}

// Path returns the location of $GOPATH
func Path() (string, error) {

	gopath := os.Getenv("GOPATH")

	if len(gopath) == 0 {
		return "", errors.New("gopath not set")
	}

	return gopath, nil
}

// Exists determine if there is an actual GOPATH
func Exists() bool {

	path, err := Path()
	if err != nil {
		return false
	}

	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}
	return true
}

// CreateProject creates repository according to https://golang.org/doc/code.html
// convention.
// e.g. $GOPATH/src/<srcRepo>/<userFolder>/<projectFolder>
// where:
//   - srcRepo could be github.com, etc
//   - userFolder is the name associated with the source repo
//   - projectFolder is the name of folder containing the Go project
func CreateProject(srcRepo string, userFolder string, projectFolder string) (string, error) {
	gopath, err := Path()
	if err != nil {
		return "", err
	}

	if !isValidName(srcRepo) {
		return "", errors.New("Invalid source repo")
	}

	if !isValidName(userFolder) {
		return "", errors.New("Invalid user folder ")
	}

	if !isValidName(projectFolder) {
		return "", errors.New("Invalid project folder")
	}

	fullPaths := []string{gopath, "src", srcRepo, projectFolder}
	path := filepath.Join(fullPaths...)
	if err := os.MkdirAll(path, 0777); err != nil {
		return "", errors.New("Unable to create project")
	}

	return path, nil
}

// SearchSrc GOPATH for path that matches term
func SearchSrc(term string) ([]string, error) {

	gopath, err := Path()
	if err != nil {
		return []string{}, err
	}

	fileList := []string{}
	err = filepath.Walk(gopath, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	fmt.Println(gopath)

	result := []string{}
	for _, file := range fileList {
		if strings.Contains(file, term) {
			result = append(result, file)
		}
	}

	return result, nil
}
