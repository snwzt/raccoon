package manager

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	cmd    *cobra.Command
	debug  bool
	logger *zerolog.Logger
	exit   func(int)
}

func newRootCmd(customlogger *zerolog.Logger, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}

	cmd := &cobra.Command{
		Use:           "raccoon",
		Short:         "Stupid simple Download Accelerator",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if root.debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&root.debug, "debug", false, "Enable debug mode")
	cmd.AddCommand(
		newURLCmd().cmd,
		newReadFileCmd().cmd,
	)

	root.cmd = cmd
	root.logger = customlogger

	return root
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(commander(cmd.cmd, args))

	err := cmd.cmd.Execute()
	if err != nil {
		cmd.logger.Info().Msg("download failed")
		cmd.logger.Debug().Err(err).Msg("details")
		cmd.exit(1) // exits with code 1, i.e. general error
	} else {
		cmd.logger.Info().Msg("download finished")
	}
}

func commander(cmd *cobra.Command, args []string) []string {
	set := map[string]bool{
		"-h":        true,
		"--help":    true,
		"--version": true,
		"help":      true,
	}

	xmd, _, _ := cmd.Find(args)

	if xmd != nil {
		if len(args) > 1 && args[1] == "help" {
			args[1] = "--help"
		}
		return args
	}

	if len(args) > 0 &&
		(args[0] == "completion" ||
			args[0] == cobra.ShellCompRequestCmd ||
			args[0] == cobra.ShellCompNoDescRequestCmd) {
		return args
	}

	if len(args) == 0 || (len(args) == 1 && set[args[0]]) {
		return args
	}

	return []string{"help"}
}

func Execute(customlogger *zerolog.Logger, exit func(int), args []string) {
	newRootCmd(customlogger, exit).Execute(args)
}
