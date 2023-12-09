package helper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/snwzt/raccoon/internal/models"
)

func ReadWithTimeout(resp *http.Response, buf []byte, timeout time.Duration) (int, error) {
	resultChan := make(chan models.ReadData)

	go func() {
		n, err := resp.Body.Read(buf)
		resultChan <- models.ReadData{
			N:   n,
			Err: err,
		}
		close(resultChan)
	}()

	select {
	case result := <-resultChan:
		return result.N, result.Err
	case <-time.After(timeout):
		return 0, fmt.Errorf("timeout reached")
	}
}
