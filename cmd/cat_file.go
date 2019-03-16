package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gog/git"
	"os"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "show git object contents",
	Run: func(cmd *cobra.Command, args []string) {
		var kind, sha string
		if len(args) != 0 {
			_, _ = fmt.Fprintln(os.Stderr, "type and sha-1 of the object should be provided")
		}

		pwd, err := os.Getwd()
		if err != nil {
			handleError(err.Error())
		}
		repo, err := git.FindRepository(pwd, true)
		if err != nil {
			handleError(err.Error())
		}

		obj, err := repo.ReadObject(sha, kind)
		if err != nil {
			handleError(err.Error())
		}

		fmt.Println(obj)
	},
}

func init() {}
