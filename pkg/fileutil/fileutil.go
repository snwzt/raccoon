package fileutil

import (
	"fmt"
	"os"
	"regexp"
)

const pattern = `[^\/]+$`

func CreateFile(url string, filedir string) (*os.File, string, error) {
	re := regexp.MustCompile(pattern)
	filename := re.FindString(url)

	_, err := os.Stat(filedir + filename)

	if !os.IsNotExist(err) {
		return nil, "", fmt.Errorf("%s already exists", filename)
	}

	out, err := os.Create(filedir + filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create %s", filename)
	}

	return out, out.Name(), nil
}

func DeleteFile(file string) {
	os.Remove(file)
}
