/*
All Rights Reversed (ɔ)
*/

package cmd

import (
	"fmt"

	"github.com/ThisaruGuruge/bestow/internal/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagSource      = "source"
	flagDestination = "destination"
)

const (
	flagVerbose    string = "verbose"
	flagQuiet      string = "quiet"
	flagDryRun     string = "dry-run"
	flagConfigFile string = "config-file"
	flagProfile    string = "profile"
	flagForce      string = "force"
	flagAdopt      string = "adopt"
	flagBackup     string = "backup"
)

type boolFlagValue struct {
	name     string
	value    bool
	strategy engine.ResolveStrategy
}

func addOperationFlags(fs *pflag.FlagSet) {
	fs.StringP(flagSource, "s", "", "root directory of the source files (Eg.: `dotfiles` repo)")
	fs.StringP(flagDestination, "d", "", "destination directory of the symlinks (Eg.: `$HOME` directory)")
}

func bindOperationalFlags(cmd *cobra.Command, v *viper.Viper) {
	if f := cmd.Flags().Lookup(flagProfile); f != nil {
		v.BindPFlag(flagProfile, f)
	}
	profile := v.GetString(flagProfile)
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

func addConflictResolutionFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP(flagForce, "f", false, "remove the existing file and create the symlink")
	cmd.Flags().BoolP(flagAdopt, "a", false, "move the existing file to the source and create the symlink")
	cmd.Flags().BoolP(flagBackup, "b", false, "rename the existing file to <filename>.bak and create the symlink")

	cmd.MarkFlagsMutuallyExclusive(flagForce, flagAdopt, flagBackup)
}
