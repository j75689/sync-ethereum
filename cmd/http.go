package cmd

import (
	"fmt"
	"os"

	"sync-ethereum/internal/app/http"

	"github.com/spf13/cobra"
)

var (
	httpCmd = &cobra.Command{
		Use:           "http",
		Short:         "Start http server",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := http.Initialize(cfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.Start()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)
