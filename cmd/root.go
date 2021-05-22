package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var (
	_Timeout uint
	_CfgFile string
	_RootCmd = &cobra.Command{
		SilenceUsage: true,
	}
)

func Execute() {
	if err := _RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	_RootCmd.AddCommand(_HttpCmd, _SchedulerCmd, _CrawlerCmd, _WriterCmd, _MigrationCmd)
	_RootCmd.PersistentFlags().StringVar(&_CfgFile, "config", "config/default.config.yaml", "config file")
	_RootCmd.PersistentFlags().UintVar(&_Timeout, "timeout", 300, "graceful shutdown timeout (second)")
}
