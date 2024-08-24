package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	nodeConfigFilePath = ""

	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "node running",
		Long:  "node running",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	}
)

// init 解析命令参数
func init() {
	nodeCmd.PersistentFlags().StringVarP(&nodeConfigFilePath, "config", "c", "conf/node.yaml", "config file path")
}
