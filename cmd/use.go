package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type useCmd struct {
	cmd  *cobra.Command
	opts useOpts
}

type useOpts struct {
}

func newUse(useOpts useOpts, version string) error {
	iv, err := config.GetInstalledVersions()
	if err != nil {
		return err
	}
	quartopath, ok := iv[version]
	if !ok {
		return fmt.Errorf("version %s not found", version)
	}
	err = os.MkdirAll(config.GetPathToActiveBinDir(), 0755)
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(config.GetPathToActiveBinDir(), "quarto"))
	if err != nil {
		return err
	}
	return os.Symlink(
		quartopath,
		filepath.Join(config.GetPathToActiveBinDir(), "quarto"),
	)
}

func setUseOpts(useOpts *useOpts) {

}

func (opts *useOpts) Validate() error {
	return nil
}

func newUseCmd() *useCmd {
	root := &useCmd{opts: useOpts{}}

	cmd := &cobra.Command{
		Use:   "use",
		Short: "use a version of quarto",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setUseOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("use-opts")
			if err := newUse(root.opts, args[0]); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
