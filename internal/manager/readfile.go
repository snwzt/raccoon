package manager

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snwzt/raccoon/internal/models"
	"github.com/snwzt/raccoon/internal/views"
	"github.com/snwzt/raccoon/pkg/download"
	"github.com/snwzt/raccoon/pkg/fileutil"
	"github.com/spf13/cobra"
)

type readFileCmd struct {
	cmd  *cobra.Command
	opts readFileOpts
}

type readFileOpts struct {
	directory   string
	connections int64
}

func newReadFileCmd() *readFileCmd {
	dir, _ := os.Getwd()

	root := &readFileCmd{}
	cmd := &cobra.Command{
		Use:           "readfile <file path>",
		Aliases:       []string{"f"},
		Short:         "Reads all URL(s) in file and downloads from all the URL(s)",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if root.opts.connections > 16 {
				root.opts.connections = 16
			}

			content, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			lines := bytes.Split(content, []byte{'\n'})

			errChan := make(chan error, len(lines)+1)
			defer close(errChan)

			status := make([]*models.Status, len(lines))

			ctx, cancel := context.WithCancel(context.Background())

			var wg sync.WaitGroup

			for idx, line := range lines {
				status[idx] = &models.Status{}

				wg.Add(1)
				go fileInstance(&wg, ctx, status[idx], string(line), errChan, root.opts)
			}

			views := tea.NewProgram(views.InitialModel(cancel, status))

			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := views.Run()
				if err != nil {
					errChan <- err
					return
				}
			}()

			wg.Wait()

			select {
			case err := <-errChan:
				return err
			default:
			}

			return nil
		},
	}

	root.cmd = cmd
	root.cmd.Flags().StringVarP(&root.opts.directory, "directory", "d",
		filepath.ToSlash(fmt.Sprintf("%v/", dir)), "Downloads file to the specified directory")
	root.cmd.Flags().Int64VarP(&root.opts.connections, "connections", "c", 4,
		"Number of connections to create for download acceleration")
	return root
}

func fileInstance(wg *sync.WaitGroup, ctx context.Context, status *models.Status, url string, errChan chan error, opts readFileOpts) {
	defer wg.Done()

	file, filepath, err := fileutil.CreateFile(url, opts.directory)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	status.Path = filepath
	status.Parts = make([]int64, opts.connections)

	downloader, err := download.NewDownloadInstance(ctx, url, opts.connections, file, status)
	if err != nil {
		errChan <- err
		fileutil.DeleteFile(filepath)
		return
	}

	if err := downloader.Download(); err != nil {
		errChan <- err
		fileutil.DeleteFile(filepath)
		return
	}
}
