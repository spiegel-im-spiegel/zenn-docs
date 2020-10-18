package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
	"github.com/spiegel-im-spiegel/gocli/rwi"
)

func newRootCmd(ui *rwi.RWI, args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hash",
		Short: "Hash functions",
		Long:  "Hash functions (detail)",
	}
	rootCmd.SilenceUsage = true
	rootCmd.SetArgs(args)            //arguments of command-line
	rootCmd.SetIn(ui.Reader())       //Stdin
	rootCmd.SetOut(ui.ErrorWriter()) //Stdout -> Stderr
	rootCmd.SetErr(ui.ErrorWriter()) //Stderr
	rootCmd.SetOutput(ui.ErrorWriter())
	rootCmd.AddCommand(
		newEncodeCmd(ui),
	)
	return rootCmd
}

func Execute(ui *rwi.RWI, args []string) exitcode.ExitCode {
	if err := newRootCmd(ui, args).Execute(); err != nil {
		return exitcode.Abnormal
	}
	return exitcode.Normal
}
