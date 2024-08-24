package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "detriot",
		Short:        "detriot server",
		Long:         "detriot server",
		SilenceUsage: true,
	}
)

// init 初始化命令行工具
func init() {
	rootCmd.AddCommand(masterCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(proxyCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute 执行命令行解析
func Execute() {
	_ = rootCmd.Execute()
}
