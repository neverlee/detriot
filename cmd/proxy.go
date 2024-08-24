package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	proxyConfigFilePath = ""

	proxyCmd = &cobra.Command{
		Use:   "proxy",
		Short: "proxy running",
		Long:  "proxy running",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	}
)

// init 解析命令参数
func init() {
	proxyCmd.PersistentFlags().StringVarP(&proxyConfigFilePath, "config", "c", "conf/proxy.yaml", "config file path")
}
