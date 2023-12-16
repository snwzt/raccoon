package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/snwzt/raccoon/internal/models"
	"github.com/snwzt/raccoon/pkg/download"
	"github.com/snwzt/raccoon/pkg/fileutil"
)

func TestDownload(t *testing.T) {
	url := "http://speedtest.ftp.otenet.gr/files/test100Mb.db"
	dir := "./"
	filename := "test100Mb.db"
	filePath := dir + filename
	connections := int64(4)

	file, filename, _ := fileutil.CreateFile(url, dir)
	defer file.Close()

	status := []*models.Status{
		{
			Name:  filename,
			Parts: make([]int64, connections),
		},
	}

	downloader, err := download.NewDownloadInstance(context.Background(), url, connections, file, status[0])
	if err != nil {
		t.Errorf("download instance - %s", err.Error())
	}

	err = downloader.Download()
	if err != nil {
		t.Errorf("download error - %s", err.Error())
	}

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Error("did not download the file")
	}

	os.Remove(filePath)
}

func TestDownloadCancellation(t *testing.T) {
	url := "http://speedtest.ftp.otenet.gr/files/test1Gb.db"
	dir := "./"
	filename := "test1Gb.db"
	filePath := dir + filename
	connections := int64(4)

	file, filename, _ := fileutil.CreateFile(url, dir)
	defer file.Close()

	status := []*models.Status{
		{
			Name:  filename,
			Parts: make([]int64, connections),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	downloader, err := download.NewDownloadInstance(ctx, url, connections, file, status[0])
	if err != nil {
		t.Errorf("download instance - %s", err.Error())
	}

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	err = downloader.Download()
	if err != nil && err.Error() != "download has been interrupted" {
		t.Errorf("download error - expected cancellation error, got %s", err.Error())
	}

	os.Remove(filePath)
}
