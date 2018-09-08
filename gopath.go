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
func Path() string {
	goPath := os.Getenv("GOPATH")

	if _, err := os.Stat(goPath); os.IsNotExist(err) {
		return ""
	}

	return goPath
}

// CreateProject creates repository according to https://golang.org/doc/code.html
// convention.
// e.g. $GOPATH/src/github/user/repo/package
func CreateProject(hostsite string, author string, project string, packages ...string) (string, error) {

	goPath := os.Getenv("GOPATH")
	if len(goPath) == 0 {
		return "", fmt.Errorf("GOPATH is not set")
	}

	if hostsite == "" {
		return "", fmt.Errorf("Host site not specified")
	}

	if author == "" {
		return "", fmt.Errorf("Author not specified")
	}

	if project == "" {
		return "", fmt.Errorf("Project not specified")
	}

	fullPath := []string{goPath, "src", hostsite, author, project}
	for _, item := range packages {
		if !isValidName(item) {
			fullPath = []string{}
			return "", fmt.Errorf("Invalid package name %s", item)
		}
		fullPath = append(fullPath, item)
	}

	path := filepath.Join(fullPath...)
	if err := os.MkdirAll(path, 0777); err != nil {
		return "", fmt.Errorf("Unable to create project")
	}

	return path, nil
}

// ProjectPaths return paths to all projects in src folder
func ProjectPaths() []string {

	goPath := Path()
	if len(goPath) == 0 {
		return []string{}
	}

	gopathSrc := filepath.Join(goPath, "src")

	fileList := []string{}
	err := filepath.Walk(gopathSrc, func(path string, fileInfo os.FileInfo, err error) error {
		if gopathSrc != path {
			if fileInfo.IsDir() {
				fileList = append(fileList, path)
			}
		}
		return nil
	})

	if err != nil {
		return []string{}
	}

	return fileList
}

// Search GOPATH for path that matches term
func Search(term string) []string {

	goPath := Path()
	if len(goPath) == 0 {
		return []string{}
	}

	result := []string{}
	err := filepath.Walk(goPath, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, term) {
			fmt.Println(path)
			result = append(result, path)
		}
		return nil
	})

	if err != nil {
		return []string{}
	}

	return result
}
