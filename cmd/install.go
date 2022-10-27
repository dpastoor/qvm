package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/AlecAivazis/survey/v2"
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

func newInstall(installOpts installOpts, release string) (error, string) {
	iv, err := config.GetInstalledVersions()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err, ""
	}
	if release == "" {
		client := gh.NewClient(os.Getenv("GITHUB_PAT"))
		releases, err := gh.GetReleases(client, 100)
		if err != nil {
			return err, ""
		}
		versions := []string{}
		for _, r := range releases {
			versions = append(versions, r.GetTagName())
		}
		err = survey.AskOne(&survey.Select{
			Message: "Which version do you want to install?",
			Options: versions,
			Description: func(value string, index int) string {
				_, ok := iv[value]
				if ok {
					return "**installed**"
				}
				return ""
			},
		}, &release, survey.WithPageSize(10))
		if err != nil {
			return err, ""
		}
	}

	// github's latest release is literally their latest release,
	// not the latest tagged version
	if release == "release" {
		client := gh.NewClient(os.Getenv("GITHUB_PAT"))
		latestRelease, err := gh.GetLatestRelease(client)
		if err != nil {
			return err, ""
		}
		release = latestRelease.GetTagName()
	}

	if release == "latest" {
		client := gh.NewClient(os.Getenv("GITHUB_PAT"))
		releases, err := gh.GetReleases(client, 1)
		if err != nil {
			return err, ""
		}
		release = releases[0].GetTagName()
	}

	_, ok := iv[release]
	if ok {
		log.Infof("quarto version %s is already installed\n", release)
		return nil, release
	}
	log.Info("attempting to install quarto version: ", release)
	res, err := pipeline.DownloadReleaseVersion(release, runtime.GOOS, installOpts.progress)
	if err != nil {
		return err, ""
	}
	log.Infof("new quarto version %s installed\n", release)
	log.Debugf("new quarto version installed to %s\n", res)
	return nil, release
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
			// this will allow the autoprompt to kick in
			if len(args) == 0 {
				args = []string{""}
			}
			errChan := make(chan error, len(args))
			for _, arg := range args {
				wg.Add(1)
				go func(errc <-chan error, release string) {
					defer wg.Done()
					err, _ := newInstall(root.opts, release)
					errChan <- err
				}(errChan, arg)
			}
			log.Trace("install waiting")
			wg.Wait()
			log.Trace("install done waiting")
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
