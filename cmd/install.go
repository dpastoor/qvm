package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/dpastoor/qvm/internal/config"
	"github.com/dpastoor/qvm/internal/gh"
	"github.com/dpastoor/qvm/internal/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type installCmd struct {
	cmd  *cobra.Command
	opts installOpts
}

type installOpts struct {
	progress bool
}

func newInstall(installOpts installOpts, release string) error {
	if release == "latest" {
		client := gh.NewClient(os.Getenv("GITHUB_PAT"))
		latestRelease, err := gh.GetLatestRelease(client)
		if err != nil {
			return err
		}
		release = latestRelease.GetTagName()
	}
	iv, err := config.GetInstalledVersions()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	_, ok := iv[release]
	if ok {
		log.Infof("quarto version %s is already installed\n", release)
		return nil
	}
	log.Info("attempting to install quarto version: ", release)
	res, err := pipeline.DownloadReleaseVersion(release, runtime.GOOS, installOpts.progress)
	if err != nil {
		return err
	}
	log.Infof("new quarto version %s installed\n", release)
	log.Debugf("new quarto version installed to %s\n", res)
	return nil
}

func setInstallOpts(installOpts *installOpts) {
	installOpts.progress = !viper.GetBool("no-progress")
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
			wg := sync.WaitGroup{}
			errChan := make(chan error, len(args))
			for _, arg := range args {
				wg.Add(1)
				go func(errc <-chan error, release string) {
					err := newInstall(root.opts, release)
					errChan <- err
					wg.Done()
				}(errChan, arg)
			}
			wg.Wait()
			// make sure to close so the range will terminate
			close(errChan)
			anyErrors := false
			for err := range errChan {
				if err != nil {
					anyErrors = true
					log.Error(err)
				}
			}
			if anyErrors {
				log.Fatal("install failed for one or more releases")
			}
			return nil
		},
	}
	cmd.Flags().BoolP("no-progress", "", false, "do not print download progress")
	viper.BindPFlag("no-progress", cmd.Flags().Lookup("no-progress"))
	root.cmd = cmd
	return root
}
