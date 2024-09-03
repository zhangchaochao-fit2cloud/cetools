package cmd

import "github.com/spf13/cobra"

var ha = &cobra.Command{
	Use:   "ha",
	Short: "CloudExplorer 高可用部署",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ha)
}
