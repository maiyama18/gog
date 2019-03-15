package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const shortDesc = "gog is a git subset written in go"
const longDesc = `gog is a git subset written in go.
reference: https://wyag.thb.lt/`

var rootCmd = &cobra.Command{
	Use: "gog",
	Short: shortDesc,
	Long: longDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello gog")
	},
}

func init() {
	cobra.OnInitialize()

	rootCmd.AddCommand(initCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}