/*
All Rights Reversed (ɔ)
*/

package cmd

import (
	"fmt"

	"github.com/ThisaruGuruge/bestow/internal/config"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setupLogging(cmd *cobra.Command) error {
	verbose, err := cmd.Flags().GetBool(FlagVerbose)
	if err != nil {
		return fmt.Errorf("parse flag %s: %w", FlagVerbose, err)
	}
	if verbose {
		logHandler.SetLevel(log.DebugLevel)
	}
	return nil
}

func loadConfig(cmd *cobra.Command) error {
	if initConfigError != nil {
		return fmt.Errorf("failed to read configs: %w", initConfigError)
	}
	bindOperationalFlags(cmd, viper.GetViper())
	if !cfgFileFound {
		appLogger.Warn("config file not found; using default values", "hint", "run 'bestow init' to initialize the configs")
	}
	var err error
	cfg, err = config.GetConfig(viper.GetViper(), appLogger)
	if err != nil {
		return err
	}
	return nil
}
