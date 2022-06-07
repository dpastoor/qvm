package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pathAddCmd struct {
	cmd  *cobra.Command
	opts pathAddOpts
}

type pathAddOpts struct {
}

func newPathAdd(pathAddOpts pathAddOpts) error {
	path := os.Getenv("PATH")
	pathToAdd := config.GetPathToActiveBinDir()
	activeIndex := strings.Index(path, pathToAdd)
	if activeIndex == -1 {
		path = fmt.Sprintf("%s:%s", pathToAdd, path)
	}
	fmt.Println(path)
	return nil
}

func setPathAddOpts(pathAddOpts *pathAddOpts) {

}

func (opts *pathAddOpts) Validate() error {
	return nil
}

func newPathAddCmd() *pathAddCmd {
	root := &pathAddCmd{opts: pathAddOpts{}}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "add qvm directories to the path, if they don't exist",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setPathAddOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("pathAdd-opts")
			if err := newPathAdd(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
