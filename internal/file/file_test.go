package file

import (
	"fmt"
	"os"
	"testing"
)

var testDir *string

func TestMain(m *testing.M) {
	fmt.Println("Running File Tests")
	createTempDirStructure()

	code := m.Run()

	fmt.Println("Finished File Tests")
	os.Exit(code)
}

func TestListAllFilesInDir(t *testing.T) {
	fileList := []string{}

}

func createTempDirStructure() {
	tempDir, err := os.MkdirTemp(".", "bestow_test_*")
	if err != nil {
		fmt.Println("failed to create the temp directory for testing")
		return
	}
	defer os.Remove(tempDir)
	testDir = &tempDir
}
