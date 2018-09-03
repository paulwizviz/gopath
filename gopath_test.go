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
	if err == nil && len(path) == 0 {
		t.Fatalf(`Expected: "" not <nil> path got: %v %v`, path, err)
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
		t.Fatalf("Expected: %t Got: %t", true, result)
	}
	removeTestFixture(t)
}

func TestCreateProject(t *testing.T) {

	// Verify that no project is created when
	// GOPATH is not set
	os.Setenv("GOPATH", "")
	_, err := CreateProject("test", "test", "test")
	if err != nil {
		t.Fatalf(`Expected: path created Got: %v`, err)
	}

	// Invalid user name
	os.Setenv("GOPATH", locationOfTestFixture(t))
	path, err := CreateProject("github.com", "user/", "project")
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

	expectedPath := filepath.Join(locationOfTestFixture(t), "src", "github.com", "user", "project")
	if path != expectedPath {
		t.Fatalf("Expected: %s Got: %s", expectedPath, path)
	}

	removeTestFixture(t)
}

func TestProjectPaths(t *testing.T) {

	os.Setenv("GOPATH", locationOfTestFixture(t))
	result, err := ProjectPaths()
	if err != nil {
		t.Fatalf("Expected: %v Got: %v", nil, err)
	}

	if len(result) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(result))
	}

	CreateProject("github.com")
	result, err = ProjectPaths()
	if err != nil {
		t.Fatalf("Expected: %v Got: %v", nil, err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected: 1 Got: %d", len(result))
	}

	CreateProject("github.com")
	result, err = ProjectPaths()
	if err != nil {
		t.Fatalf("Expected: %v Got: %v", nil, err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected: 1 Got: %d", len(result))
	}

	CreateProject("github.com", "user")
	result, err = ProjectPaths()
	if err != nil {
		t.Fatalf("Expected: %v Got: %v", nil, err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected: 2 Got: %d", len(result))
	}
	removeTestFixture(t)

}

func TestSearch(t *testing.T) {

	os.Setenv("GOPATH", "")
	_, err := Search("github.com")
	if err == nil {
		t.Fatalf("Expected: not %v Got: %v", nil, err)
	}

	os.Setenv("GOPATH", locationOfTestFixture(t))
	CreateProject("github.com", "test", "test")
	CreateProject("something")
	CreateProject("bitbucket.org", "user", "test")

	result, _ := Search("hello")
	if len(result) != 0 {
		t.Fatalf("Expected: 0 Got: %d", len(result))
	}

	result, _ = Search("test")
	if len(result) != 3 {
		t.Fatalf("Expected: 3 Got: %d", len(result))
	}

	removeTestFixture(t)
}
