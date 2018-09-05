/*
Copyright 2018 Paul Sitoh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package gopath is a library to navigate GOPATH
// and create Go packages
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

// ProjectPaths return paths to all projects
func ProjectPaths() ([]string, error) {

	gopath, err := Path()
	if err != nil {
		return []string{}, err
	}

	gopathSrc := filepath.Join(gopath, "src")

	fileList := []string{}
	err = filepath.Walk(gopathSrc, func(path string, fileInfo os.FileInfo, err error) error {

		if gopathSrc != path {
			if fileInfo.IsDir() {
				fileList = append(fileList, path)
			}
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return fileList, nil
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
// e.g. $GOPATH/src/github.com/user/repo
func CreateProject(packages ...string) (string, error) {
	gopath, err := Path()
	if err != nil {
		return "", err
	}

	fullPath := []string{gopath, "src"}
	for _, item := range packages {
		if !isValidName(item) {
			return "", fmt.Errorf("Invalid package name %s", item)
		}
		fullPath = append(fullPath, item)
	}

	path := filepath.Join(fullPath...)
	if err := os.MkdirAll(path, 0777); err != nil {
		return "", errors.New("Unable to create project")
	}

	return path, nil
}

// Search GOPATH for path that matches term
func Search(term string) ([]string, error) {

	gopath, err := Path()
	if err != nil {
		return []string{}, err
	}

	fileList := []string{}
	err = filepath.Walk(gopath, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	result := []string{}
	for _, file := range fileList {
		if strings.Contains(file, term) {
			result = append(result, file)
		}
	}

	return result, nil
}
