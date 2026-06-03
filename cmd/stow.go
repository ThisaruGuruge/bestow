/*
All Rights Reversed (ɔ)
*/

package cmd

import (
	"fmt"

	"github.com/ThisaruGuruge/bestow/internal/engine"
	"github.com/spf13/cobra"
)

var stowCmd = &cobra.Command{
	Use:     "stow [packages...]",
	Short:   stowShort,
	Long:    stowLong,
	Example: stowExamples,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(cmd)
		if err != nil {
			return err
		}
		appLogger.Debug("running stow command", "args", args)
		var force, adopt, backup bool
		force, err = cmd.Flags().GetBool(flagForce)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", flagForce, err)
		}
		adopt, err = cmd.Flags().GetBool(flagAdopt)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", flagAdopt, err)
		}
		backup, err = cmd.Flags().GetBool(flagBackup)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", flagBackup, err)
		}
		var strategy engine.ResolveStrategy
		if force {
			strategy = engine.ResolveForce
		}
		if adopt {
			strategy = engine.ResolveAdopt
		}
		if backup {
			strategy = engine.ResolveBackup
		}

		dryrun, err := cmd.Flags().GetBool(flagDryRun)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", flagDryRun, err)
		}
		ctx := engine.CommandContext{
			Action:           engine.ActionStow,
			Args:             args,
			ConflictStrategy: strategy,
		}
		eng, err := engine.NewEngine(cfg, dryrun, appLogger)
		if err != nil {
			return err
		}
		summary, err := eng.Execute(&ctx)
		if err != nil {
			return err
		}
		appOutput.PrintSummary(summary)
		return nil
	},
}

func init() {
	addOperationFlags(stowCmd.Flags())
	addConflictResolutionFlags(stowCmd)

	rootCmd.AddCommand(stowCmd)
}
