package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gog/git"
)

var kind string
var write bool

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "compute object id and optionally creates a blob from a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			handleError("filename should be provided")
		}
		filePath := args[0]

		sha, err := git.ObjectHash(filePath, kind, !write)
		if err != nil {
			handleError(err.Error())
		}

		fmt.Println(sha)
	},
}

func init() {
	hashObjectCmd.Flags().StringVarP(&kind, "type", "t", "blob", "type of object (blob|commit|tag|tree)")
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "if set, create a blob")
}
