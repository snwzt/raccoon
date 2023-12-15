package manager

import (
	"context"
	"fmt"
	"math"
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

type urlCmd struct {
	cmd  *cobra.Command
	opts urlOpts
}

type urlOpts struct {
	directory   string
	connections int64
}

func newURLCmd() *urlCmd {
	dir, _ := os.Getwd()

	root := &urlCmd{}
	cmd := &cobra.Command{
		Use:           "url <url>",
		Aliases:       []string{"u"},
		Short:         "Downloads from URL",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if root.opts.connections > 16 {
				root.opts.connections = 16
			}

			file, filename, err := fileutil.CreateFile(args[0], root.opts.directory)
			if err != nil {
				return err
			}
			defer file.Close()

			errChan := make(chan error, 2)
			defer close(errChan)

			ctx, cancel := context.WithCancel(context.Background())

			status := []*models.Status{
				{
					Name:      filename,
					Parts:     make([]int64, root.opts.connections),
					FinalSize: math.MaxInt, // to fix the division by zero issue in begining
				},
			}

			var wg sync.WaitGroup

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

			downloader, err := download.NewDownloadInstance(ctx, args[0], root.opts.connections,
				file, status[0])
			if err != nil {
				fileutil.DeleteFile(file.Name())
				return err
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := downloader.Download(); err != nil {
					errChan <- err
					fileutil.DeleteFile(file.Name())
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
	root.cmd.Flags().StringVarP(&root.opts.directory, "directory", "d", filepath.ToSlash(fmt.Sprintf("%v/", dir)),
		"Downloads file to the specified directory")
	root.cmd.Flags().Int64VarP(&root.opts.connections, "connections", "c", 4,
		"Number of connections to create for download acceleration")
	return root
}
