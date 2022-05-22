package cmd

import (
	"fmt"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pathCmd struct {
	cmd  *cobra.Command
	opts pathOpts
}

type pathOpts struct {
}

func newPath(pathOpts pathOpts) error {
	fmt.Println(config.GetActiveQuartoPath())
	return nil
}

func setPathOpts(pathOpts *pathOpts) {

}

func (opts *pathOpts) Validate() error {
	return nil
}

func newPathCmd() *pathCmd {
	root := &pathCmd{opts: pathOpts{}}

	cmd := &cobra.Command{
		Use:   "path",
		Short: "get the path to the active quarto executable",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setPathOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("path-opts")
			if err := newPath(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
