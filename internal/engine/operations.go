package engine

import (
	"github.com/ThisaruGuruge/bestow/internal/file"
)

type ConflictResolution string

const (
	ConflictSkip        ConflictResolution = "skip"
	ConflictForce       ConflictResolution = "force"
	ConflictAdopt       ConflictResolution = "adopt"
	ConflictBackup      ConflictResolution = "backup"
	ConflictInteractive ConflictResolution = "interactive"
)

type Operation struct {
	Source      string
	Destination string
	Package     string
	Steps       []Step
}

type Step struct {
	SourceFilePath      string
	DestinationFilePath string
	Conflict            ConflictResolution
}

func populateOperations(actionCtx ActionContext, executionCtx ExecutionContext) ([]Operation, error) {
	result := []Operation{}
	var action Action = actionCtx.Action
	var source, destination string
	if action == ActionUnstow {
		source = executionCtx.Destination
		destination = executionCtx.Source
	} else {
		source = executionCtx.Source
		destination = executionCtx.Destination
	}

	rootOperation := Operation{
		Source:      source,
		Destination: destination,
		Package:     "",
	}
	if err := populateStepsForRoot(&rootOperation); err != nil {
		return nil, err
	}
	result = append(result, rootOperation)

	for _, pkg := range executionCtx.PackageList {
		operation := Operation{
			Source:      source,
			Destination: destination,
			Package:     pkg,
		}
		populateStepsForPackage(&operation)
		result = append(result, operation)

	}
	return result, nil
}

func populateStepsForRoot(operation *Operation) error {
	rootFileList, err := file.ListFiles(operation.Source)
	if err != nil {
		return &EngineError{
			Message: "failed to read files from the source root",
			Package: ".",
			Cause:   err,
		}
	}
	steps := []Step{}
	for _, fileName := range rootFileList {
		steps = append(steps, Step{
			SourceFilePath: fileName,
		})
	}
	operation.Steps = steps
	return nil
}

func populateStepsForPackage(operation *Operation) error {
	sourceFileList := []string{}
	// destinationFileList := []string{}

	err := file.ListAllFilesInDir(operation.Source, operation.Package, &sourceFileList)
	if err != nil {
		return &EngineError{
			Message: "failed to read the package contents",
			Cause:   err,
		}
	}
	steps := []Step{}
	for _, fileName := range sourceFileList {
		steps = append(steps, Step{
			SourceFilePath: fileName,
		})
	}
	operation.Steps = steps
	return nil
}

// TODO: Handle interactive/non-interactive modes.
// Pass config here to check interactivity and conflict resolution strategy
// Should return error in any invalid scenario.
func validateOperations(operations *[]Operation) error {
	for _, operation := range *operations {
		validateOperation(&operation)
	}
	return nil
}

func validateOperation(operation *Operation) error {
	// operation.Source
	// operation.Destination
	// operation.Package
	// operation.Steps
	//
	// operation.Steps[0]
	return nil
}
