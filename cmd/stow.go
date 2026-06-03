/*
All Rights Reversed (ɔ)
*/

package cmd

import (
	"fmt"

	"github.com/ThisaruGuruge/bestow/internal/engine"
	"github.com/ThisaruGuruge/bestow/internal/output"
	"github.com/spf13/cobra"
)

var stowCmd = &cobra.Command{
	Use:     "stow [packages...]",
	Short:   stowShort,
	Long:    stowLong,
	Example: stowExamples,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadConfig(cmd); err != nil {
			return err
		}
		appLogger.Debug("running stow command", "args", args)
		var force, adopt, backup bool
		var err error
		force, err = cmd.Flags().GetBool(FlagForce)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", FlagForce, err)
		}
		adopt, err = cmd.Flags().GetBool(FlagAdopt)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", FlagAdopt, err)
		}
		backup, err = cmd.Flags().GetBool(FlagBackup)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", FlagBackup, err)
		}

		flagValues := []boolFlagValue{
			{name: FlagForce, value: force, strategy: engine.ResolveForce},
			{name: FlagAdopt, value: adopt, strategy: engine.ResolveAdopt},
			{name: FlagBackup, value: backup, strategy: engine.ResolveBackup},
		}

		conflictResolution, err := conflictResolve(flagValues)
		if err != nil {
			return err
		}
		dryrun, err := cmd.Flags().GetBool(FlagDryRun)
		if err != nil {
			return fmt.Errorf("parse flag %s: %w", FlagDryRun, err)
		}
		ctx := engine.CommandContext{
			Action:           engine.ActionStow,
			Args:             args,
			ConflictStrategy: conflictResolution,
		}
		eng, err := engine.NewEngine(cfg, dryrun, appLogger)
		if err != nil {
			return err
		}
		summary, err := eng.Execute(&ctx)
		if err != nil {
			return err
		}
		output.PrintSummary(summary)
		return nil
	},
}

func init() {
	addOperationFlags(stowCmd.Flags())
	addConflictResolutionFlags(stowCmd.Flags())

	rootCmd.AddCommand(stowCmd)
}
