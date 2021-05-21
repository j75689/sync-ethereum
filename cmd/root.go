package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var (
	cfgFile string
	rootCmd = &cobra.Command{
		SilenceUsage: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(httpCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/default.config.yaml", "config file")
}
