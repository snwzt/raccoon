package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/snwzt/raccoon/internal/models"
	"github.com/snwzt/raccoon/pkg/helper"
)

type Download struct {
	ListenContext context.Context
	URL           string
	Parts         [][2]int64
	Client        *http.Client
	File          *os.File
	status        *models.Status
	errChan       chan error
	doneChan      chan bool
}

func NewDownloadInstance(ctx context.Context, url string, numParts int64, file *os.File, status *models.Status) (*Download, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	segmentSize := resp.ContentLength / numParts
	segments := make([][2]int64, numParts)

	for idx := range segments {
		start := int64(idx) * segmentSize
		end := start + segmentSize - 1

		if idx == int(numParts-1) {
			end = resp.ContentLength - 1
		}

		segments[idx] = [2]int64{start, end}
	}

	status.FinalSize = resp.ContentLength

	errChan := make(chan error, numParts)
	doneChan := make(chan bool, numParts)

	return &Download{
		ListenContext: ctx,
		URL:           url,
		Parts:         segments,
		Client:        &http.Client{},
		File:          file,
		status:        status,
		errChan:       errChan,
		doneChan:      doneChan,
	}, nil
}

func (dm *Download) Download() error {
	defer close(dm.doneChan)
	defer close(dm.errChan)

	var wg sync.WaitGroup

	for idx, segment := range dm.Parts {
		wg.Add(1)

		go dm.worker(&wg, segment, &dm.status.Parts[idx])
	}

	wg.Wait()

	dm.status.Done = true

	select {
	case err := <-dm.errChan:
		dm.status.Err = err
		return err
	default:
	}

	return nil
}

func (dm *Download) worker(wg *sync.WaitGroup, segment [2]int64, parts *int64) {
	defer wg.Done()

	req, err := http.NewRequest("GET", dm.URL, nil)
	if err != nil {
		dm.errChan <- err
		dm.doneChan <- true
		return
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", segment[0], segment[1]))

	resp, err := dm.Client.Do(req)
	if err != nil {
		dm.errChan <- err
		dm.doneChan <- true
		return
	}

	defer resp.Body.Close()

	buf := make([]byte, 1024*1024) // 1 MB buffer

	for {
		select {
		case <-dm.doneChan:
			return
		case <-dm.ListenContext.Done():
			dm.errChan <- fmt.Errorf("download has been interrupted")
			dm.doneChan <- true
			return
		default:
			n, err := helper.ReadWithTimeout(resp, buf, 10*time.Second)
			if err != nil && err != io.EOF {
				dm.errChan <- err
				dm.doneChan <- true
				return
			}
			if n == 0 { // end of body
				return
			}
			_, err = dm.File.WriteAt(buf[:n], segment[0])
			if err != nil {
				dm.errChan <- err
				dm.doneChan <- true
				return
			}

			*parts += int64(n)
			segment[0] += int64(n)
		}
	}
}
