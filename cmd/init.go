package cmd

import (
	"github.com/spf13/cobra"
	"gog/git"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize git directory",
	Run: func(cmd *cobra.Command, args []string) {
		var repoPath string
		if len(args) == 0 {
			repoPath = "."
		} else {
			repoPath = args[0]
		}

		if _, err := git.CreateRepository(repoPath); err != nil {
			handleError(err.Error())
		}
	},
}

func init() {}
