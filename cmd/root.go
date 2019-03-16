package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const shortDesc = "gog is a git subset written in go"
const longDesc = `gog is a git subset written in go.
reference: https://wyag.thb.lt/`

func handleError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "gog",
	Short: shortDesc,
	Long:  longDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello gog")
	},
}

func init() {
	cobra.OnInitialize()

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(catFileCmd)
	rootCmd.AddCommand(hashObjectCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleError(err.Error())
	}
}
