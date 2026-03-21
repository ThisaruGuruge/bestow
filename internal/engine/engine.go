package engine

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ThisaruGuruge/bestow/internal/config"
	"github.com/ThisaruGuruge/bestow/internal/file"
	"github.com/ThisaruGuruge/bestow/internal/log"
	"github.com/ThisaruGuruge/bestow/internal/output"
)

type Action string

const (
	ActionStow   Action = "stow"
	ActionUnstow Action = "unstow"
)

type EngineError struct {
	Message string
	Command Action
	Package string
	Cause   error
}

func (e *EngineError) Error() string {
	msg := e.Message
	if e.Command != "" {
		msg += fmt.Sprint(": [%s]", e.Command)
	}
	if e.Package != "" {
		msg += fmt.Sprint(": [%s]", e.Package)
	}
	if e.Cause != nil {
		msg += fmt.Sprintf(": %v", e.Cause)
	}
	return msg
}

func (e *EngineError) Unwrap() error { return e.Cause }

type ActionContext struct {
	Action Action
	Args   []string
	DryRun bool
}

type ExecutionContext struct {
	Source      string
	Destination string
	PackageList []string
}

func Execute(ctx *ActionContext, cfg *config.Config) error {
	log.Debug("executing", "config", cfg, "context", ctx)
	packageList, err := getPackageList(ctx, cfg)
	if err != nil {
		return &EngineError{
			Message: "failed to read the files",
			Command: ctx.Action,
			Cause:   err,
		}
	}
	log.Debug("packages received", "packages", packageList)
	executionCtx := &ExecutionContext{
		Source:      cfg.Source,
		Destination: cfg.Destination,
		PackageList: *packageList,
	}
	log.Info("found candidates", "context", executionCtx)
	operations, err := populateOperations(*ctx, *executionCtx)
	if err != nil {
		return &EngineError{
			Message: "failed to opulate operations",
			Command: ctx.Action,
			Cause:   err,
		}
	}
	output.Success(fmt.Sprint("operation: ", ctx.Action))
	for _, operation := range operations {
		output.Success(fmt.Sprint("Pacakge: ", operation.Package))
		output.Success(fmt.Sprint("Steps:"))
		for i, step := range operation.Steps {
			output.Success(fmt.Sprint("Step No. ", i))
			output.Success(fmt.Sprint("    Source File Path: ", step.SourceFilePath))
		}
	}
	if err := validateOperations(&operations); err != nil {
		return &EngineError{
			Message: "invalid operation",
			Command: ctx.Action,
			Cause:   err,
		}
	}
	return nil
}

func getPackageList(ctx *ActionContext, cfg *config.Config) (*[]string, error) {
	source := cfg.Source
	ignoreList, err := newIgnoreList(source)
	if err != nil {
		//TODO: Custom error?
		return nil, err
	}
	log.Debug("retrieving package list", "source", source)
	var pkgCandidates []string
	if len(ctx.Args) == 0 {
		log.Warn("no packages provided, processing all the packages", "action", ctx.Action)
		pkgCandidates, err = file.ListAllDirectories(source)
		if err != nil {
			return nil, &EngineError{
				Message: "failed to read packages from source",
				Command: ctx.Action,
				Cause:   err,
			}
		}
	} else {
		for _, pkgCandidate := range ctx.Args {
			isDir, err := file.IsDir(filepath.Join(source, pkgCandidate))
			if err != nil {
				return nil, &EngineError{
					Message: "failed to read package",
					Command: ctx.Action,
					Package: pkgCandidate,
					Cause:   err,
				}
			}
			if !isDir {
				return nil, &EngineError{
					Message: "failed to read package",
					Package: pkgCandidate,
					Command: ctx.Action,
				}
			}
			pkgCandidates = append(pkgCandidates, pkgCandidate)
		}

		log.Debug("retrieved candidates for packages", "candidates", pkgCandidates)
	}
	packages, err := filterPackages(pkgCandidates, ignoreList.items)
	if err != nil {
		return nil, err
	}

	return &packages, nil
}

func filterPackages(candidates, filters []string) ([]string, error) {
	log.Debug("filtering packages", "candidates", candidates, "filter", filters)
	result := []string{}
	for _, candidate := range candidates {
		found := false
		for _, filter := range filters {
			//TODO: this works for packages. For file filtering the logic is different. Think when it comes to that
			if strings.HasSuffix(filter, "/") {
				filter = strings.TrimSuffix(filter, "/")
			}
			matched, err := filepath.Match(filter, candidate)
			if err != nil {
				return nil, &EngineError{
					Message: "failed to check ignore pattern",
					Cause:   err,
				}
			}
			if matched {
				found = true
				break
			}
		}
		if found {
			continue
		}
		result = append(result, candidate)
	}
	return result, nil
}
