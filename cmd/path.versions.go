package cmd

import (
	"fmt"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pathVersionsCmd struct {
	cmd  *cobra.Command
	opts pathVersionsOpts
}

type pathVersionsOpts struct {
}

func newPathVersions(pathVersionsOpts pathVersionsOpts) error {
	fmt.Println(config.GetPathToVersionsDir())
	return nil
}

func setPathVersionsOpts(pathVersionsOpts *pathVersionsOpts) {

}

func (opts *pathVersionsOpts) Validate() error {
	return nil
}

func newPathVersionsCmd() *pathVersionsCmd {
	root := &pathVersionsCmd{opts: pathVersionsOpts{}}

	cmd := &cobra.Command{
		Use:   "versions",
		Short: "get the path to where the various quarto versions are stored",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setPathVersionsOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("pathVersions-opts")
			if err := newPathVersions(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
