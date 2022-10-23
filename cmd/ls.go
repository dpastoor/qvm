package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/dpastoor/qvm/internal/config"
	"github.com/dpastoor/qvm/internal/gh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type lsCmd struct {
	cmd  *cobra.Command
	opts lsOpts
}

type lsOpts struct {
	remote bool
	num    int
}

func newLs(lsOpts lsOpts) error {
	if lsOpts.remote {
		client := gh.NewClient(os.Getenv("GITHUB_PAT"))
		releases, err := gh.GetReleases(client, lsOpts.num)
		if err != nil {
			return err
		}
		fmt.Println("version  | release date | description | type")
		for _, r := range releases {
			createdAt := r.GetCreatedAt()
			var releaseType string
			if r.GetPrerelease() {
				releaseType = "pre-release"
			} else {
				releaseType = "release"
			}
			fmt.Printf("%s |  %s  |   %s  | %s \n", r.GetTagName(), createdAt.Format("2006-01-02"), r.GetName(), releaseType)
		}
	} else {
		entries, err := os.ReadDir(config.GetPathToVersionsDir())
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No installed quarto versions found")
			return nil
		}
		if err != nil {
			return err
		}
		if len(entries) < lsOpts.num {
			lsOpts.num = len(entries)
		}
		entries = entries[:lsOpts.num-1]
		// TODO: replace with actual table
		fmt.Println("version           | install time")
		fmt.Println("--------------------------------")

		// modification time
		// sort.Slice(entries, func(i, j int) bool {
		// 	x, _ := entries[i].Info()
		// 	y, _ := entries[j].Info()
		// 	return x.ModTime().After(y.ModTime())
		// })
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name() > entries[j].Name()
		})
		// no need to worry about errors since just need to know version
		// for matching below and won't match if doesn't exist
		activeVersion, _ := config.GetActiveVersion()
		for _, e := range entries {
			if e.IsDir() {
				dinfo, _ := e.Info()
				name := e.Name()
				if activeVersion == e.Name() {
					name += " (active)"
				} else {
					name += "         "
				}
				fmt.Printf("%s | %s\n", name, dinfo.ModTime().Format("2006-01-02"))
			}
		}
	}

	return nil
}

func setLsOpts(lsOpts *lsOpts) {
	lsOpts.remote = viper.GetBool("remote")
	lsOpts.num = viper.GetInt("number")
}

func (opts *lsOpts) Validate() error {
	return nil
}

func newLsCmd() *lsCmd {
	root := &lsCmd{opts: lsOpts{}}

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "list versions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setLsOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("ls-opts")
			if err := newLs(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().Bool("remote", false, "list remote versions")
	viper.BindPFlag("remote", cmd.Flags().Lookup("remote"))
	cmd.Flags().IntP("number", "n", 10, "number of versions to list")
	viper.BindPFlag("number", cmd.Flags().Lookup("number"))
	root.cmd = cmd
	return root
}
