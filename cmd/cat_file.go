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
		if len(args) != 2 {
			handleError("type and sha-1 of the object should be provided")
		}
		kind := args[0]
		sha := args[1]

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
