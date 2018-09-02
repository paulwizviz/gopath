package gopath

import (
	"os"
	"path/filepath"
	"testing"
)

func locationOfTestFixture(t *testing.T) string {
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error unable to locate working director: %v", err)
	}

	return filepath.Join(pwd, "go")
}

func createTestFixture(t *testing.T) {

	path := locationOfTestFixture(t)
	if err := os.MkdirAll(path, 0777); err != nil {
		t.Errorf("Error unable to create test-fixtures: %v", err)
	}
}

func removeTestFixture(t *testing.T) {

	path := locationOfTestFixture(t)
	if err := os.RemoveAll(path); err != nil {
		t.Errorf("Error unable to delete test-fixtures: %v", err)
	}
}

func TestPath(t *testing.T) {

	// Validate that if GOPATH is not set, the API return error and empty string
	os.Setenv("GOPATH", "")
	path, err := Path()
	if err == nil {
		t.Fatalf("Expected: <gopath not set> got: %v", err)
	}
	if path != "" {
		t.Fatalf("Expected: <empty> got: %s", path)
	}

	// Validate that if GOPATH is set, the API returns path value
	os.Setenv("GOPATH", locationOfTestFixture(t))
	path, err = Path()
	if err != nil {
		t.Fatalf("Error expected: <nil> got: %v", err)
	}
	if path != locationOfTestFixture(t) {
		t.Fatalf("Expected: %s Got: %s", locationOfTestFixture(t), path)
	}
}

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

func TestExist(t *testing.T) {

	os.Setenv("GOPATH", "")
	if result := Exists(); result != false {
		t.Fatalf("Expected: false Got: %t", result)
	}

	createTestFixture(t)
	os.Setenv("GOPATH", locationOfTestFixture(t))
	if result := Exists(); result != true {
		t.Fatalf("Expected: true Got: %t", result)
	}
	removeTestFixture(t)
}

func TestCreateProject(t *testing.T) {

	// Verify that no project is created when
	// GOPATH is not set
	os.Setenv("GOPATH", "")
	path, err := CreateProject("test", "test", "test")
	if path != "" && err != nil {
		t.Fatalf(`Expected: path created Got: %v %v`, path, err)
	}

	// Invalid user name
	os.Setenv("GOPATH", locationOfTestFixture(t))
	path, err = CreateProject("github.com", "user/", "project")
	if err == nil {
		t.Fatalf(`Expected: not nil Got: %v`, err)
	}
	if path != "" {
		t.Fatalf("Expected: empty path Got %s", path)
	}

	// Valid arguments
	os.Setenv("GOPATH", locationOfTestFixture(t))
	path, err = CreateProject("github.com", "user", "project")
	if err != nil {
		t.Fatalf(`Expected: not nil Got: %v`, err)
	}
	if path == "" {
		t.Fatalf("Expected: no path Got %s", path)
	}

	removeTestFixture(t)
}
