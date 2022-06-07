package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type initCmd struct {
	cmd  *cobra.Command
	opts initOpts
}

type initOpts struct {
}

func newInit(initOpts initOpts) error {
	activeDir := config.GetPathToActiveBinDir()
	activeIndex := strings.Index(os.Getenv("PATH"), activeDir)
	if activeIndex == -1 {
		fmt.Println("please add the active bin directory to your path")
		fmt.Println(`you can dynamically query the active bin directory by running: "qvm path active" 
		if you would like to script this, you can add the following to your ~/.bashrc or ~/.zshrc:
		export PATH="$(qvm path add)")
		`)

	} else {
		fmt.Println("already detect qvm active directory on your path - you're good to go!")
	}
	return nil
}

func setInitOpts(initOpts *initOpts) {

}

func (opts *initOpts) Validate() error {
	return nil
}

func newInitCmd() *initCmd {
	root := &initCmd{opts: initOpts{}}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize qvm",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setInitOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("init-opts")
			if err := newInit(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
