package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "获取版本信息",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.1")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
