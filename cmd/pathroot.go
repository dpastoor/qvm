package cmd

import (
	"fmt"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pathRootCmd struct {
	cmd  *cobra.Command
	opts pathRootOpts
}

type pathRootOpts struct {
}

func newPathRoot(pathRootOpts pathRootOpts) error {
	fmt.Println(config.GetRootConfigPath())
	return nil
}

func setPathRootOpts(pathRootOpts *pathRootOpts) {

}

func (opts *pathRootOpts) Validate() error {
	return nil
}

func newPathRootCmd() *pathRootCmd {
	root := &pathRootCmd{opts: pathRootOpts{}}

	cmd := &cobra.Command{
		Use:   "root",
		Short: "get the root qvm path",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setPathRootOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("pathRoot-opts")
			if err := newPathRoot(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
