package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dpastoor/qvm/internal/config"
	"github.com/dpastoor/qvm/internal/gh"
	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type lsCmd struct {
	cmd  *cobra.Command
	opts lsOpts
}

type lsOpts struct {
	remote  bool
	num     int
	release string
	output  string
}

func newLs(lsOpts lsOpts) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if lsOpts.remote {
		client := gh.NewClient(os.Getenv("GITHUB_PAT"))
		releases, err := gh.GetReleases(client, lsOpts.num, lsOpts.release)
		if err != nil {
			return err
		}
		t.AppendHeader(table.Row{"version", "release date", "description", "type"})
		for _, r := range releases {
			createdAt := r.GetCreatedAt()
			var releaseType string
			if r.GetPrerelease() {
				releaseType = "pre-release"
			} else {
				releaseType = "release"
			}
			t.AppendRow(table.Row{r.GetTagName(), createdAt.Format("2006-01-02"), r.GetName(), releaseType})
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
		entries = entries[:lsOpts.num]
		t.AppendHeader(table.Row{"version", "active", "install date"})

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
				active := ""
				if activeVersion == e.Name() {
					active = "*"
				} else {
				}
				t.AppendRow(table.Row{name, active, dinfo.ModTime().Format("2006-01-02")})
			}
		}
	}
	switch lsOpts.output {
	case "default":
		t.Render()
	case "markdown":
		t.RenderMarkdown()
	case "csv":
		t.RenderCSV()
	case "html":
		t.RenderHTML()
	}
	return nil
}

func setLsOpts(lsOpts *lsOpts) {
	lsOpts.remote = viper.GetBool("remote")
	lsOpts.num = viper.GetInt("number")
	lsOpts.release = viper.GetString("release-type")
	lsOpts.output = viper.GetString("output")
}

func (opts *lsOpts) Validate() error {
	acceptedMethods := []string{"csv", "markdown", "html", "default"}
	for _, m := range acceptedMethods {
		if m == opts.output {
			return nil
		}
	}
	return fmt.Errorf("invalid output type, must be one of %s",
		strings.Join(acceptedMethods, ","))
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
	cmd.Flags().String("release-type", "", "list remote versions")
	viper.BindPFlag("release-type", cmd.Flags().Lookup("release-type"))
	cmd.Flags().IntP("number", "n", 10, "number of versions to list")
	viper.BindPFlag("number", cmd.Flags().Lookup("number"))
	cmd.Flags().String("output", "default", "list remote versions")
	viper.BindPFlag("output", cmd.Flags().Lookup("output"))
	root.cmd = cmd
	return root
}
