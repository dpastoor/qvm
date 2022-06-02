package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dpastoor/qvm/internal/config"
	"github.com/dpastoor/qvm/internal/gh"
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
	if version == "latest" {
		client := gh.NewClient(os.Getenv("GITHUB_PAT"))
		latestRelease, err := gh.GetLatestRelease(client)
		if err != nil {
			return err
		}
		version = latestRelease.GetTagName()
	}
	iv, err := config.GetInstalledVersions()
	if err != nil {
		return err
	}
	quartopath, ok := iv[version]
	if !ok {
		return fmt.Errorf("version %s not found", version)
	}
	err = os.MkdirAll(config.GetPathToActiveBinDir(), 0700)
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(config.GetPathToActiveBinDir(), "quarto"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.Symlink(
		quartopath,
		filepath.Join(config.GetPathToActiveBinDir(), "quarto"),
	)
	if err != nil {
		return err
	}
	log.Infof("now using quarto version: %s\n", version)
	return nil
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
