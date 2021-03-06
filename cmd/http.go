package cmd

import (
	"fmt"
	"os"
	"time"

	"sync-ethereum/internal/app/http"
	"sync-ethereum/pkg/util"

	"github.com/spf13/cobra"
)

var (
	_HttpCmd = &cobra.Command{
		Use:           "http",
		Short:         "Start http server",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := http.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Start, app.Stop, time.Duration(_Timeout)*time.Second)
		},
	}
)
