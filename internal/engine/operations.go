/*
All Rights Reversed (ɔ)
*/

package engine

import (
	"fmt"
	"path/filepath"

	"github.com/ThisaruGuruge/bestow/internal/file"
	"github.com/ThisaruGuruge/bestow/internal/output"
)

type ResolveStrategy string

const (
	ResolveSkip        ResolveStrategy = "skip"
	ResolveForce       ResolveStrategy = "force"
	ResolveAdopt       ResolveStrategy = "adopt"
	ResolveBackup      ResolveStrategy = "backup"
	ResolveInteractive ResolveStrategy = "interactive"
)

type Operation struct {
	Source      string
	Destination string
	BackupPath  string
	Action      FileAction
	Strategy    ResolveStrategy
}

// TODO: Need to verify if two operations have the same destination.
// Which should be an error; We should catch it here before proceesing to
// execute the operations
func (e *Engine) populateOperations(ctx *CommandContext) ([]Operation, error) {
	result := []Operation{}
	packageList, err := e.populatePackageList(ctx.Args)
	if err != nil {
		return nil, err
	}
	for _, pkg := range packageList {
		e.Logger.Debug("populating operations for package", "pacakge", pkg)
		pacakgeOperations, err := e.getPackageOperation(pkg, ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, pacakgeOperations...)
	}
	return result, nil
}

func (e *Engine) getPackageOperation(pkg string, ctx *CommandContext) ([]Operation, error) {
	sourceFileList, err := e.FileSystem.ListAllFilesInDir(e.Source, pkg)
	if err != nil {
		return nil, &EngineError{
			Message: "failed to read the package contents",
			Cause:   err,
		}
	}
	operations := []Operation{}
	for _, fileName := range sourceFileList {
		doIgnore, err := e.Ignore.shouldIgnore(fileName, pkg)
		if err != nil {
			return nil, err
		}
		if doIgnore {
			e.Logger.Debug("ignoring the file", "fileName", fileName, "package", pkg)
			continue
		}
		srcFile := filepath.Join(e.Source, fileName)
		relativePath := filepath.Join(e.FileSystem.GetPathSegments(fileName)[1:]...)
		destFile := filepath.Join(e.Destination, relativePath)

		operations = append(operations, Operation{
			Source:      srcFile,
			Destination: destFile,
			Action:      FileActionLink,
			Strategy:    ctx.ConflictStrategy,
		})
	}
	return operations, nil
}

// TODO: Handle interactive/non-interactive modes.
// Pass config here to check interactivity and conflict resolution strategy
// Should return error in any invalid scenario.
func (e *Engine) resolveStowOperations(operations *[]Operation) ([]Operation, error) {
	for i := range *operations {
		if err := e.resolveStowOperation(&(*operations)[i]); err != nil {
			return *operations, err
		}
	}
	return *operations, nil
}

func (e *Engine) resolveStowOperation(operation *Operation) error {
	destExists, _ := e.FileSystem.Exists(operation.Destination)
	if destExists {
		// TODO: Doesn't make any sense for the static resolver. But we need this when we have interactive mode
		existing, err := e.FileSystem.GetExistingFileType(operation.Source, operation.Destination)
		if err != nil {
			return &EngineError{
				Message: "failed to check exising file type",
				Cause:   err,
			}
		}
		resolver := StaticResolver{strategy: operation.Strategy}
		strategy, err := resolver.Resolve(e.Source, e.Destination, existing)
		if err := e.resolveFileAction(operation, strategy, existing); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) resolveUnstowOperations(operations *[]Operation) ([]Operation, error) {
	for i := range *operations {
		if err := e.resolveUnstowOperation(&(*operations)[i]); err != nil {
			return *operations, err
		}
	}
	return *operations, nil
}

func (e *Engine) resolveUnstowOperation(operation *Operation) error {
	destExists, err := e.FileSystem.Exists(operation.Destination)
	if err != nil {
		return &EngineError{
			Message: "failed to check the destination file",
			Cause:   err,
		}
	}
	if !destExists {
		e.Logger.Debug("file have not been stowed", "source", operation.Source)
		operation.Action = FileActionSkip
		return nil
	}
	existingType, err := e.FileSystem.GetExistingFileType(operation.Source, operation.Destination)
	if err != nil {
		return &EngineError{
			Message: "failed to check the destination file",
			Cause:   err,
		}
	}
	if existingType != file.ExistingManagedSymlink {
		operation.Action = FileActionSkip
		return nil
	}
	operation.Action = FileActionUnlink
	return nil
}

func (e *Engine) stow(operations []Operation) error {
	e.Logger.Debug("stowing files")
	// I'm sorry, but for sentimental reasons, I will not accept any AI-generated PRs here.
	operations, err := e.resolveStowOperations(&operations)
	if err != nil {
		return &EngineError{
			Message: "failed to resolve stow operation",
			Cause:   err,
		}
	}
	operations = filterSkipFiles(operations)
	for _, operation := range operations {
		if err := e.stowOperation(&operation); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) stowOperation(operation *Operation) error {
	switch operation.Action {
	case FileActionSkip:
		e.Logger.Debug("skipping file", "source", operation.Source, "destination", operation.Destination, "action", operation.Action)
		return nil
	case FileActionLink:
		return e.createLink(operation.Source, operation.Destination)
	case FileActionReplaceLink:
		return e.updateLink(operation.Source, operation.Destination)
	case FileActionBackupLink:
		return e.backupLink(operation.Source, operation.Destination)
	case FileActionAdoptLink:
		return e.adoptLink(operation.Source, operation.Destination)
	}
	return nil
}

// TODO: When unstowing, empty directories should be removed.
// This should be configurable using the config.yaml.
// Default behavior should be to remove the emplty directory.
func (e *Engine) unstow(operations []Operation) error {
	e.Logger.Debug("unstowing files")
	operations, err := e.resolveUnstowOperations(&operations)
	if err != nil {
		return err
	}
	operations = filterSkipFiles(operations)
	for _, operation := range operations {
		if err := e.unstowOperation(&operation); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) unstowOperation(operation *Operation) error {
	e.Logger.Debug("unstowing file", "source", operation.Source, "destination", operation.Destination)
	exists, err := e.FileSystem.Exists(operation.Destination)
	if err != nil {
		return &EngineError{
			Message: "failed to check the destination file",
			Cause:   err,
		}
	}
	if !exists {
		e.Logger.Debug("file not found for unstow", "path", operation.Destination)
		return nil
	}
	fileType, err := e.FileSystem.GetExistingFileType(operation.Source, operation.Destination)
	if err != nil {
		return &EngineError{
			Message: "failed to check the fily type",
			Cause:   err,
		}
	}
	if fileType != file.ExistingManagedSymlink {
		e.Logger.Warn("existing file is not managed by bestow", "existing_file", operation.Destination)
		return nil
	}
	err = e.FileSystem.Remove(operation.Destination)
	// TODO: Remove empty files; check config before
	// TODO: Fix the bug where subdirectories in the packages won't removed-identifying as not empty
	parent := filepath.Dir(operation.Destination)
	isEmpty, err := e.FileSystem.IsEmpty(parent)
	if err != nil {
		return err
	}
	if isEmpty {
		e.FileSystem.Remove(parent)
	}
	if err != nil {
		return &EngineError{
			Message: "failed to unstow the file",
			Cause:   err,
		}
	}
	e.Logger.Debug("successfully unstowed the file", "source", operation.Source, "destination", operation.Destination)
	return nil
}

func (e *Engine) createLink(src, dest string) error {
	if err := e.FileSystem.Link(src, dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	output.Success(fmt.Sprintf("created: %s", dest))
	return nil
}

func (e *Engine) updateLink(src, dest string) error {
	if err := e.FileSystem.Remove(dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	if err := e.createLink(src, dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	return nil
}

func (e *Engine) backupLink(src, dest string) error {
	if err := e.FileSystem.Backup(dest); err != nil {
		return &EngineError{
			Message: "failed to backup existing file",
			Cause:   err,
		}
	}
	if err := e.createLink(src, dest); err != nil {
		return err
	}
	return nil
}

func (e *Engine) adoptLink(src, dest string) error {
	if err := e.FileSystem.Copy(dest, src); err != nil {
		return &EngineError{
			Message: "failed to adopt file from destination",
			Cause:   err,
		}
	}
	if err := e.FileSystem.Link(src, dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	return nil
}

// TODO: Implement this or find a better alternative
func (e *Engine) removeEmptyDirs(path string) error {

	return nil
}
