package cmd

import "github.com/spf13/cobra"

var initCmd = &cobra.Command{
	Use: "init",
	Short: "initialize git directory",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {

}
