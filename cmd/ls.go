package cmd

import (
	"fmt"
	"os"

	"github.com/dpastoor/qvm/internal/gh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type lsCmd struct {
	cmd  *cobra.Command
	opts lsOpts
}

type lsOpts struct {
}

func newLs(lsOpts lsOpts) error {
	client := gh.NewClient(os.Getenv("GITHUB_PAT"))
	releases, err := gh.GetReleases(client, false)
	for _, r := range releases {
		createdAt := r.GetCreatedAt()
		fmt.Printf("%s - %s\n", r.GetName(), createdAt.Format("2006-01-02"))
	}
	return err
}

func setLsOpts(lsOpts *lsOpts) {

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
	root.cmd = cmd
	return root
}
