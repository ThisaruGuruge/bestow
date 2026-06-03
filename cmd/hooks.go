/*
All Rights Reversed (ɔ)
*/

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ThisaruGuruge/bestow/internal/config"
	"github.com/ThisaruGuruge/bestow/internal/output"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ErrIncompatibleFlags = errors.New("mutually exclusive flags")

func setupLogging(cmd *cobra.Command) error {
	verbose, err := cmd.Flags().GetBool(FlagVerbose)
	if err != nil {
		return fmt.Errorf("parse flag %s: %w", FlagVerbose, err)
	}
	quiet, err := cmd.Flags().GetBool(FlagQuiet)
	if err != nil {
		return fmt.Errorf("parse flag %s: %w", FlagQuiet, err)
	}
	if verbose && quiet {
		return fmt.Errorf("parse flag %s %s: %w", FlagVerbose, FlagQuiet, ErrIncompatibleFlags)
	}
	if verbose {
		logHandler.SetLevel(log.DebugLevel)
	}
	if quiet {
		appOutput = output.NewOutput(output.Quiet)
	} else {
		appOutput = output.NewOutput(output.Normal)
	}
	return nil
}

func loadConfig(cmd *cobra.Command) (*config.Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			appLogger.Warn("config file not found; using $HOME as the default destination")
		} else {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}
	bindOperationalFlags(cmd, viper.GetViper())
	cfg, err := config.GetConfig(viper.GetViper(), appLogger)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
