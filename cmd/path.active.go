package cmd

import (
	"fmt"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pathActiveCmd struct {
	cmd  *cobra.Command
	opts pathActiveOpts
}

type pathActiveOpts struct {
}

func newPathActive(pathActiveOpts pathActiveOpts) error {
	fmt.Println(config.GetPathToActiveBinDir())
	return nil
}

func setPathActiveOpts(pathActiveOpts *pathActiveOpts) {

}

func (opts *pathActiveOpts) Validate() error {
	return nil
}

func newPathActiveCmd() *pathActiveCmd {
	root := &pathActiveCmd{opts: pathActiveOpts{}}

	cmd := &cobra.Command{
		Use:   "active",
		Short: "get the path to the active quarto bin dir",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setPathActiveOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("pathActive-opts")
			if err := newPathActive(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
