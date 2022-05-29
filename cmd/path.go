package cmd

import (
	"github.com/spf13/cobra"
)

type pathCmd struct {
	cmd *cobra.Command
}

func newPathCmd() *pathCmd {
	root := &pathCmd{}

	cmd := &cobra.Command{
		Use:   "path",
		Short: "get the path to the active quarto executable",
	}
	cmd.AddCommand(newPathRootCmd().cmd)
	cmd.AddCommand(newPathActiveCmd().cmd)
	cmd.AddCommand(newPathVersionsCmd().cmd)
	root.cmd = cmd
	return root
}
