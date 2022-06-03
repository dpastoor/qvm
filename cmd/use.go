package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type useCmd struct {
	cmd  *cobra.Command
	opts useOpts
}

type useOpts struct {
}

func newUse(useOpts useOpts, version string) error {
	iv, err := config.GetInstalledVersions()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	versions := maps.Keys(iv)
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	if len(iv) == 0 {
		return errors.New("no installed versions found, please install a version first")
	}
	if version == "" {

		err := survey.AskOne(&survey.Select{
			Message: "Which version do you want to install?",
			Options: versions,
		}, &version)
		if err != nil {
			return err
		}
	}
	if version == "latest" {
		version = versions[0]
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
