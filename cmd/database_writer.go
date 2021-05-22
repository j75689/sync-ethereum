package cmd

import (
	"fmt"
	"os"
	"time"

	"sync-ethereum/internal/app/database_writer"
	"sync-ethereum/pkg/util"

	"github.com/spf13/cobra"
)

var (
	_WriterCmd = &cobra.Command{
		Use:           "writer",
		Short:         "Start database writer",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := database_writer.Initialize(_CfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Start, app.Stop, time.Duration(_Timeout)*time.Second)
		},
	}
)
