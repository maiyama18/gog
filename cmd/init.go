package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gog/git"
	"os"
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
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
	},
}

func init() {}
