package cmd

import (
	"fmt"
	"os"
	"time"

	"sync-ethereum/internal/app/scheduler"
	"sync-ethereum/pkg/util"

	"github.com/spf13/cobra"
)

var (
	_SchedulerCmd = &cobra.Command{
		Use:           "scheduler",
		Short:         "Start scheduler",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := scheduler.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Start, app.Stop, time.Duration(_Timeout)*time.Second)
		},
	}
)
