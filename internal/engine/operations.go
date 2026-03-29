package engine

import (
	"path/filepath"

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
	Strategy    ResolveStrategy
}

func (e *Engine) populateOperations() error {
	result := []Operation{}
	var action Action = e.Action
	var source, destination string
	if action == ActionUnstow {
		source = e.Destination
		destination = e.Source
	} else {
		source = e.Source
		destination = e.Destination
	}

	rootOperation := Operation{
		Source:      source,
		Destination: destination,
	}
	if err := e.getRootOperation(); err != nil {
		return nil
	}
	result = append(result, rootOperation)

	for _, pkg := range *e.PackageList {
		e.populateStepsForPackage(pkg)
	}
	return nil
}

func (e *Engine) getRootOperation() error {
	rootFileList, err := file.ListFiles(e.Source)
	if err != nil {
		return &EngineError{
			Message: "failed to read files from the source root",
			Package: ".",
			Cause:   err,
		}
	}
	for _, fileName := range rootFileList {
		doIgnore, err := e.Ignore.shouldIgnore(fileName, constant.RootPackageName)
		if err != nil {
			return err
		}
		if doIgnore {
			log.Debug("ignoring the file", "fileName", fileName)
			continue
		}
		srcFile := filepath.Join(e.Source, fileName)
		destFile := filepath.Join(e.Destination, fileName)
		*e.Operations = append(*e.Operations, Operation{
			Source:      srcFile,
			Destination: destFile,
		})
	}
	return nil
}

func (e *Engine) populateStepsForPackage(pkg string) error {
	sourceFileList := []string{}
	err := file.ListAllFilesInDir(e.Source, pkg, &sourceFileList)
	if err != nil {
		return &EngineError{
			Message: "failed to read the package contents",
			Cause:   err,
		}
	}
	for _, fileName := range sourceFileList {
		doIgnore, err := e.Ignore.shouldIgnore(fileName, pkg)
		if err != nil {
			return err
		}
		if doIgnore {
			log.Debug("ignoring the file", "fileName", fileName, "package", pkg)
			continue
		}
		srcFile := filepath.Join(e.Source, fileName)
		// TODO: Remove package name from the file name
		relativePath := filepath.Join(file.GetPathSegments(fileName)[1:]...)
		destFile := filepath.Join(e.Destination, relativePath)
		*e.Operations = append(*e.Operations, Operation{
			Source:      srcFile,
			Destination: destFile,
		})
	}
	return nil
}

// TODO: Handle interactive/non-interactive modes.
// Pass config here to check interactivity and conflict resolution strategy
// Should return error in any invalid scenario.
func (e *Engine) validateOperations() error {
	for _, operation := range *e.Operations {
		e.validateOperation(&operation)
	}
	return nil
}

func (e *Engine) validateOperation(operation *Operation) error {
	destExists, _ := file.Exists(e.Destination)
	if destExists {
		existingType := getExistingType(operation)
		// TODO: Doesn't make any sense for the static resolver. But we need this when we have interactive mode
		resolver := StaticResolver{strategy: e.Strategy}
		strategy, _ := resolver.Resolve(e.Source, e.Destination, existingType)
		operation.Strategy = strategy
	}
	return nil
}

func getExistingType(operation *Operation) ExistingType {
	isDir, _ := file.IsDir(operation.Destination)
	if isDir {
		return ExistingDir
	}
	if file.IsSameFile(operation.Source, operation.Destination) {
		return ExistingManagedSymlink
	}
	return ExistingForeignSymlink
}
