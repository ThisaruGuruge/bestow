/*
All Rights Reversed (ɔ)
*/

package engine

import (
	"github.com/ThisaruGuruge/bestow/internal/file"
)

type FileAction string

const (
	FileActionLink        FileAction = "Link"
	FileActionReplaceLink FileAction = "ReplaceLink"
	FileActionBackupLink  FileAction = "BackupLink"
	FileActionAdoptLink   FileAction = "AdoptLink"
	FileActionUnlink      FileAction = "Unlink"
	FileActionSkip        FileAction = "Skip"
)

func (e *Engine) resolveFileAction(operation *Operation, strategy ResolveStrategy, existing file.ExistingType) error {
	e.Logger.Debug("Resolving file actions", "source", operation.Source, "destination", operation.Destination, "strategy", strategy, "existing_type", existing)
	switch existing {
	case file.ExistingManagedSymlink:
		e.Logger.Debug("symlink already exists, skipping", "destination", operation.Destination, "strategy", strategy, "existing_type", existing)
		operation.Action = FileActionSkip
	case file.ExistingDir:
		return e.resolveExistingDir(operation, strategy)
	case file.ExistingRegularFile:
		return e.resolveRegularFile(operation, strategy)
	case file.ExistingForeignSymlink:
		return e.resolveForeignSymlink(operation, strategy)
	}
	return nil
}

func (e *Engine) resolveExistingDir(operation *Operation, strategy ResolveStrategy) error {
	if strategy == ResolveSkip {
		e.Logger.Warn("destination is a directory; skipping the file", "destination", operation.Destination, "destination_type", "DIRECTORY", "strategy", strategy)
		operation.Action = FileActionSkip
		return nil
	}
	return &EngineError{
		Message: "cannot perform operation; destination is a directory",
	}
}

func (e *Engine) resolveRegularFile(operation *Operation, strategy ResolveStrategy) error {
	switch strategy {
	case ResolveSkip:
		e.Logger.Warn("destination exists; skipping the file", "destination", operation.Destination, "destination_type", "FILE", "strategy", strategy)
		operation.Action = FileActionSkip
	case ResolveForce:
		operation.Action = FileActionReplaceLink
	case ResolveAdopt:
		operation.Action = FileActionAdoptLink
	case ResolveBackup:
		operation.Action = FileActionBackupLink
	}
	return nil
}

func (e *Engine) resolveForeignSymlink(operation *Operation, strategy ResolveStrategy) error {
	switch strategy {
	case ResolveSkip:
		operation.Action = FileActionSkip
		e.Logger.Warn("destination exists, skipping the file", "destination", operation.Destination, "destination_type", "FOREIGN SYMLINK", "strategy", strategy)
	case ResolveForce:
		operation.Action = FileActionReplaceLink
	case ResolveAdopt:
		operation.Action = FileActionAdoptLink
	case ResolveBackup:
		operation.Action = FileActionBackupLink
	}
	return nil
}

func (e *Engine) resolveManagedSymlink(operation *Operation, strategy ResolveStrategy) error {
	return nil
}
