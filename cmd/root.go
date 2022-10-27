package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type settings struct {
	// logrus log level
	loglevel string
}

type rootCmd struct {
	cmd *cobra.Command
	cfg *settings
}

func Execute(version string, args []string) {
	newRootCmd(version).Execute(args)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)
	if err := cmd.cmd.Execute(); err != nil {
		// if get to this point and don't fatally log in the subcommand,
		// the Usage help will be printed before the error,
		// which may or may not be the desired behavior
		log.Fatalf("failed with error: %s", err)
	}
}

func setGlobalSettings(cfg *settings) {
	cfg.loglevel = viper.GetString("loglevel")
	setLogLevel(cfg.loglevel)
}
func newRootCmd(version string) *rootCmd {
	root := &rootCmd{cfg: &settings{}}
	cmd := &cobra.Command{
		Use:   "qvm",
		Short: "quarto version manager",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// need to set the config values here as the viper values
			// will not be processed until Execute, so can't
			// set them in the initializer.
			// If persistentPreRun is used elsewhere, should
			// remember to setGlobalSettings in the initializer
			setGlobalSettings(root.cfg)
		},
	}
	cmd.Version = version
	// without this, the default version is like `cmd version <version>` so this
	// will just print the version for simpler parsing
	cmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
	cmd.PersistentFlags().String("loglevel", "info", "log level")
	viper.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))
	cmd.AddCommand(newDebugCmd(root.cfg))
	cmd.AddCommand(newManCmd().cmd)
	cmd.AddCommand(newLsCmd().cmd)
	cmd.AddCommand(newPathCmd().cmd)
	cmd.AddCommand(newInstallCmd().cmd)
	cmd.AddCommand(newUseCmd().cmd)
	cmd.AddCommand(newInitCmd().cmd)
	cmd.AddCommand(newRmCmd().cmd)
	root.cmd = cmd
	return root
}
