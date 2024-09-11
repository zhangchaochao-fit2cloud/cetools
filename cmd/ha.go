package cmd

import (
	"github.com/spf13/cobra"
)

var (
	GlobalConfigFiles []string
)

var ha = &cobra.Command{
	Use:   "ha",
	Short: "CloudExplorer 高可用部署",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {

	//RootCmd.PersistentFlags().BoolVar(&GlobalErrorOnWarning, "error-on-warning", false, "Treat any warning as an error")
	RootCmd.AddCommand(ha)
}
