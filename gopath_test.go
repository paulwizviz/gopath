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

package gopath

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Utilities
func locationOfTestFixture(t *testing.T) string {
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error unable to locate working director: %v", err)
	}

	return filepath.Join(pwd, "go")
}

func createTestFixture(t *testing.T) {
	path := locationOfTestFixture(t)
	t.Logf("Creating test fixture: %s", path)
	if err := os.MkdirAll(path, 0777); err != nil {
		t.Errorf("Error unable to create test-fixtures: %v", err)
	}
}

func removeTestFixture(t *testing.T) {

	path := locationOfTestFixture(t)
	t.Logf("Removing test fixture: %s", path)
	if err := os.RemoveAll(path); err != nil {
		t.Errorf("Error unable to delete test-fixtures: %v", err)
	}
}

func projectExists(t *testing.T, path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Tests
func TestIsValidName(t *testing.T) {

	testData := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"1A", false},
		{"github.com", true},
		{"test/", false},
		{"test-1", true},
		{"t1st_1", true},
		{".1", false},
	}

	for _, data := range testData {
		if isValidName(data.input) != data.expected {
			t.Errorf("Input: %s Expect: false Got: %t", data.input, data.expected)
		}
	}

}

func TestPath(t *testing.T) {

	// Validate that if GOPATH is not set, the API return error and empty string
	os.Setenv("GOPATH", "")
	goPath := Path()
	if len(goPath) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(goPath))
	}

	// Validate non-existence GOPATH
	expected := locationOfTestFixture(t)
	os.Setenv("GOPATH", expected)
	goPath = Path()
	if strings.Compare(expected, goPath) == 0 {
		t.Fatalf("Expected: %s Got: %s", expected, goPath)
	}

	// Validate an actual GOPATH
	expected = locationOfTestFixture(t)
	createTestFixture(t)
	os.Setenv("GOPATH", expected)
	goPath = Path()
	if strings.Compare(expected, goPath) != 0 {
		t.Fatalf("Expected: %s Got: %s", expected, goPath)
	}
	removeTestFixture(t)
}

func TestCreateProject(t *testing.T) {

	packages := []string{"package", "subpackage"}

	// GOPATH unset
	t.Log("GOPATH is not set")
	os.Setenv("GOPATH", "")
	project, err := CreateProject("github.com", "test", "test", packages...)
	if err == nil {
		t.Fatalf("Expected: Valid GOPATH Got: %v", err)
	}
	if len(project) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(project))
	}

	// GOPATH set
	goPath := locationOfTestFixture(t)
	os.Setenv("GOPATH", goPath)
	t.Logf("GOPATH set to %s", goPath)

	// Invalid package name
	invalidPackage := "package/"
	project, err = CreateProject("github.com", "test", "test", invalidPackage, packages[1])
	if err == nil {
		t.Fatalf("Expected: no error Got: %v", err)
	}
	if len(project) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(project))
	}

	// Valid package
	expected := filepath.Join(goPath, "src", "github.com", "test", "test", packages[0], packages[1])
	t.Logf("Valid package names %s", expected)
	project, err = CreateProject("github.com", "test", "test", packages...)
	if err != nil {
		t.Fatalf("Expected: no error Got: %v", err)
	}
	if strings.Compare(expected, project) != 0 {
		t.Fatalf("Expected: %s Got: %s", expected, project)
	}
	if !projectExists(t, project) {
		t.Fatal("Expected: true Got: false")
	}
	removeTestFixture(t)

}

func TestProjectPaths(t *testing.T) {

	os.Setenv("GOPATH", "")
	projects := ProjectPaths()
	if len(projects) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(projects))
	}

	goPath := locationOfTestFixture(t)
	os.Setenv("GOPATH", goPath)
	if len(projects) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(projects))
	}

	CreateProject("github.com", "test", "test")
	CreateProject("bitbucket.org", "test", "test")
	paths := ProjectPaths()
	if !projectExists(t, paths[0]) {
		t.Fatalf("Path[0] does not exists")
	}
	if !projectExists(t, paths[1]) {
		t.Fatalf("Path[1] does not exists")
	}
	removeTestFixture(t)

}

func TestSearch(t *testing.T) {

	os.Setenv("GOPATH", "")
	result := Search("test")
	if len(result) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(result))
	}

	goPath := locationOfTestFixture(t)
	os.Setenv("GOPATH", goPath)
	CreateProject("github.com", "test1", "test")
	CreateProject("bitbucket.org", "test", "test")
	CreateProject("something", "test1", "test1")
	result = Search("test1")

	if len(result) != 4 {
		t.Fatalf("Expected: 4 Got: %d", len(result))
	}

}
