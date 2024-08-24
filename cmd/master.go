package cmd

import (
	"github.com/spf13/cobra"

	"github.com/neverlee/detriot/cmd/master"
	"github.com/neverlee/detriot/lrpc/log"
)

var (
	masterConfigFilePath = ""

	masterCmd = &cobra.Command{
		Use:   "master",
		Short: "master running",
		Long:  "master running",
		Run: func(c *cobra.Command, args []string) {
			err := master.Run(masterConfigFilePath)
			if err != nil {
				log.Error("Run master error:", err)
			}
		},
	}
)

func init() {
	masterCmd.PersistentFlags().StringVarP(&masterConfigFilePath, "config", "c", "conf/master.yaml", "config file path")
}
