package cmd

import (
	"fmt"
	"os"

	"sync-ethereum/internal/app/migration"

	"github.com/spf13/cobra"
)

var (
	_MigrationCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migration tool",
	}

	_MigrationUpCmd = &cobra.Command{
		Use:           "up",
		Short:         "Migrate the DB to the most recent version available",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, args []string) {
			app, err := migration.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.MigrateUp()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.Stop()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}

	_MigrationUpToCmd = &cobra.Command{
		Use:           "up-to",
		Short:         "Migrate the DB to a specific VERSION",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, args []string) {
			app, err := migration.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			version := ""
			if len(args) > 0 {
				version = args[0]
			}
			err = app.MigrateUpTo(version)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.Stop()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}

	_MigrationDownCmd = &cobra.Command{
		Use:           "down",
		Short:         "Roll back all migrations",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, args []string) {
			app, err := migration.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.MigrateDown()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.Stop()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}

	_MigrationDownToCmd = &cobra.Command{
		Use:           "down-to",
		Short:         "Roll back to a specific VERSION",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, args []string) {
			app, err := migration.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			version := ""
			if len(args) > 0 {
				version = args[0]
			}
			err = app.MigrateDownTo(version)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.Stop()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	_MigrationCmd.AddCommand(_MigrationUpCmd, _MigrationUpToCmd, _MigrationDownCmd, _MigrationDownToCmd)
}
