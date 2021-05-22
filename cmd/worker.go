package cmd

import (
	"fmt"
	"os"
	"time"

	"sync-ethereum/internal/app/worker"
	"sync-ethereum/pkg/util"

	"github.com/spf13/cobra"
)

var (
	_WorkerCmd = &cobra.Command{
		Use:           "worker",
		Short:         "Start worker",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := worker.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Start, app.Stop, time.Duration(_Timeout)*time.Second)
		},
	}
)
