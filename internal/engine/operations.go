package engine

import (
	"path/filepath"
	"slices"

	"github.com/ThisaruGuruge/bestow/internal/constant"
	"github.com/ThisaruGuruge/bestow/internal/file"
	"github.com/ThisaruGuruge/bestow/internal/log"
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
}

func (e *Engine) populateOperations() ([]Operation, error) {
	result := []Operation{}
	if slices.Contains(*e.PackageList, RootDir) {
		rootOperations, err := e.getRootOperation(result, e.Source, e.Destination)
		if err != nil {
			return nil, err
		}
		result = append(result, rootOperations...)
	}

	for _, pkg := range *e.PackageList {
		// Skip processing root directory as a package
		if pkg == RootDir {
			continue
		}
		pacakgeOperations, err := e.getPackageOperation(pkg)
		if err != nil {
			return nil, err
		}
		result = append(result, pacakgeOperations...)
	}
	return result, nil
}

func (e *Engine) getRootOperation(operations []Operation, src, dest string) ([]Operation, error) {
	rootFileList, err := file.ListFiles(e.Source)
	if err != nil {
		return nil, &EngineError{
			Message: "failed to read files from the source root",
			Package: ".",
			Cause:   err,
		}
	}
	for _, fileName := range rootFileList {
		doIgnore, err := e.Ignore.shouldIgnore(fileName, constant.RootPackageName)
		if err != nil {
			return nil, err
		}
		if doIgnore {
			log.Debug("ignoring the file", "fileName", fileName)
			continue
		}
		srcFile := filepath.Join(src, fileName)
		destFile := filepath.Join(dest, fileName)
		operations = append(operations, Operation{
			Source:      srcFile,
			Destination: destFile,
			Action:      FileActionLink,
		})
	}
	return operations, nil
}

func (e *Engine) getPackageOperation(pkg string) ([]Operation, error) {
	sourceFileList, err := file.ListAllFilesInDir(e.Source, pkg)
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
			log.Debug("ignoring the file", "fileName", fileName, "package", pkg)
			continue
		}
		srcFile := filepath.Join(e.Source, fileName)
		relativePath := filepath.Join(file.GetPathSegments(fileName)[1:]...)
		destFile := filepath.Join(e.Destination, relativePath)

		operations = append(operations, Operation{
			Source:      srcFile,
			Destination: destFile,
			Action:      FileActionLink,
		})
	}
	return operations, nil
}

// TODO: Handle interactive/non-interactive modes.
// Pass config here to check interactivity and conflict resolution strategy
// Should return error in any invalid scenario.
func (e *Engine) resolveOperations(operations *[]Operation) ([]Operation, error) {
	for i := range *operations {
		if err := e.resolveOperation(&(*operations)[i]); err != nil {
			return *operations, err
		}
	}
	return *operations, nil
}

func (e *Engine) resolveOperation(operation *Operation) error {
	destExists, _ := file.Exists(operation.Destination)
	if destExists {
		// TODO: Doesn't make any sense for the static resolver. But we need this when we have interactive mode
		existing, err := file.GetExistingFileType(operation.Source, operation.Destination)
		if err != nil {
			return &EngineError{
				Message: "failed to check exising file type",
				Command: e.Action,
				Cause:   err,
			}
		}
		resolver := StaticResolver{strategy: e.Strategy}
		strategy, _ := resolver.Resolve(e.Source, e.Destination, existing)
		if err := e.resolveFileAction(operation, strategy, existing); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) stow(operations []Operation) error {
	for _, operation := range operations {
		err := e.stowOperation(&operation)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) stowOperation(operation *Operation) error {
	switch operation.Action {
	case FileActionSkip:
		return nil
	case FileActionLink:
		return createLink(operation.Source, operation.Destination)
	case FileActionRemoveLink:
		return updateLink(operation.Source, operation.Destination)
	case FileActionBackupLink:
		return backupLink(operation.Source, operation.Destination)
	case FileActionAdoptLink:
		return adoptLink(operation.Source, operation.Destination)
	}
	return nil
}

func (e *Engine) unstow(operations []Operation) error {

	return nil
}

func createLink(src, dest string) error {
	if err := file.Link(src, dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	return nil
}

func updateLink(src, dest string) error {
	if err := file.Remove(dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	if err := createLink(src, dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	return nil
}

func backupLink(src, dest string) error {
	if err := file.Backup(dest); err != nil {
		return &EngineError{
			Message: "failed to backup existing file",
			Cause:   err,
		}
	}
	if err := createLink(src, dest); err != nil {
		return err
	}
	return nil
}

func adoptLink(src, dest string) error {
	if err := file.Copy(dest, src); err != nil {
		return &EngineError{
			Message: "failed to adopt file from destination",
			Cause:   err,
		}
	}
	if err := file.Link(src, dest); err != nil {
		return &EngineError{
			Message: "failed to stow the file",
			Cause:   err,
		}
	}
	return nil
}
