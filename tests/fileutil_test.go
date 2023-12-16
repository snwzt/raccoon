package tests

import (
	"os"
	"testing"

	"github.com/snwzt/raccoon/pkg/fileutil"
)

func TestCreateFile(t *testing.T) {
	filedir := "./"
	filename := "test.txt"
	filePath := filedir + filename

	file, _, err := fileutil.CreateFile("https://example.com/test.txt", filedir)
	if err != nil {
		t.Errorf("CreateFile Error: %s", err.Error())
	}
	defer file.Close()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("CreateFile did not create the file")
	}

	_, _, err = fileutil.CreateFile("https://example.com/test.txt", filedir)
	if err == nil {
		t.Error("No error for creating an already existing file")
	}

	os.Remove(filePath)
}

func TestDeleteFile(t *testing.T) {
	filedir := "./"
	filename := "test.txt"
	filePath := filedir + filename

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file for test: %s", err.Error())
	}

	file.Close()

	fileutil.DeleteFile(file.Name())

	_, err = os.Stat(filePath)
	if !os.IsNotExist(err) {
		t.Error("DeleteFile did not delete the file")
	}
}
