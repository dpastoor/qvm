package cmd

import (
	"fmt"
	"runtime"

	"github.com/dpastoor/qvm/internal/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type installCmd struct {
	cmd  *cobra.Command
	opts installOpts
}

type installOpts struct {
}

func newInstall(installOpts installOpts, release string) error {
	log.Info("attempting to install quarto version: ", release)
	res, err := pipeline.DownloadReleaseVersion(release, runtime.GOOS)
	if err != nil {
		return err
	}
	log.Infof("new quarto version %s installed\n", release)
	log.Debugf("new quarto version installed to %s\n", res)
	return nil
}

func setInstallOpts(installOpts *installOpts) {

}

func (opts *installOpts) Validate() error {
	return nil
}

func newInstallCmd() *installCmd {
	root := &installCmd{opts: installOpts{}}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "install a given quarto version",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setInstallOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("install-opts")
			if err := newInstall(root.opts, args[0]); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
