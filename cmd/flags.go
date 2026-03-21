package cmd

import (
	"fmt"

	"github.com/ThisaruGuruge/bestow/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagSource      = "source"
	flagDestination = "destination"
)

func addOperationFlags(fs *pflag.FlagSet) {
	fs.StringP(flagSource, "s", "", "root directory of the source files (Eg.: `dotfiles` repo)")
	fs.StringP(flagDestination, "d", "", "destination directory of the symlinks (Eg.: `$HOME` directory)")
}

func bindOperationalFlags(cmd *cobra.Command, v *viper.Viper) {
	if f := cmd.Flags().Lookup(FlagProfile); f != nil {
		v.BindPFlag(FlagProfile, f)
	}
	profile := v.GetString(FlagProfile)
	if profile == "" {
		profile = "default"
	}
	prefix := fmt.Sprintf("profiles.%s", profile)
	if f := cmd.Flags().Lookup(flagSource); f != nil {
		v.BindPFlag(prefix+".source", f)
	}
	if f := cmd.Flags().Lookup(flagDestination); f != nil {
		v.BindPFlag(prefix+".destination", f)
	}
}

func checkVerbose(cmd *cobra.Command) error {
	verbose, err := cmd.Flags().GetBool(FlagVerbose)
	if err != nil {
		return err
	}
	if verbose {
		log.SetLevel(log.LevelDebug)
	}
	return nil
}
