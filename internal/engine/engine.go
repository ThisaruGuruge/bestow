package engine

import (
	"fmt"
	"path/filepath"

	"github.com/ThisaruGuruge/bestow/internal/config"
	"github.com/ThisaruGuruge/bestow/internal/constant"
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

type CommandContext struct {
	Action   Action
	Args     []string
	DryRun   bool
	Conflict ResolveStrategy
}

type ExecutionContext struct {
	Source      string
	Destination string
	PackageList []string
	Ignore      IgnoreList
}

type Engine struct {
	Source      string
	Destination string
	Ignore      IgnoreList
	Args        *[]string
	Action      Action
	PackageList *[]string
	Operations  *[]Operation
	Strategy    ResolveStrategy
}

func NewEngine(ctx *CommandContext, cfg *config.Config) (*Engine, error) {
	ignoreList, err := newIgnoreList(cfg.Source)
	if err != nil {
		return nil, &EngineError{
			Message: "failed to initialize the engine",
			Command: ctx.Action,
		}
	}
	return &Engine{
		Source:      cfg.Source,
		Destination: cfg.Destination,
		Ignore:      *ignoreList,
		Args:        &ctx.Args,
		Action:      ctx.Action,
		Operations:  &[]Operation{},
		Strategy:    ctx.Conflict,
	}, nil
}

func (e *Engine) Execute() error {
	if err := e.populatePackageList(); err != nil {
		return &EngineError{
			Message: "failed to read the files",
			Command: e.Action,
			Cause:   err,
		}
	}
	if err := e.populateOperations(); err != nil {
		return &EngineError{
			Message: "failed to opulate operations",
			Command: e.Action,
			Cause:   err,
		}
	}
	// I'm sorry, but for sentimental reasons, I will not accept any AI-generated PRs here.
	if err := e.validateOperations(); err != nil {
		return &EngineError{
			Message: "invalid operation",
			Command: e.Action,
			Cause:   err,
		}
	}
	for _, operation := range *e.Operations {
		output.Success(fmt.Sprintf("[Copy]: %s -> %s", operation.Source, operation.Destination))
	}
	return nil
}

func (e *Engine) populatePackageList() error {
	log.Debug("populating package list", "source", e.Source)
	var pkgCandidates []string
	if len(*e.Args) == 0 {
		log.Warn("no packages provided, processing all the packages", "action", e.Action)
		var err error
		pkgCandidates, err = file.ListAllDirectories(e.Source)
		if err != nil {
			return &EngineError{
				Message: "failed to read packages from source",
				Command: e.Action,
				Cause:   err,
			}
		}
	} else {
		for _, pkgCandidate := range *e.Args {
			isDir, err := file.IsDir(filepath.Join(e.Source, pkgCandidate))
			if err != nil {
				return &EngineError{
					Message: "failed to read package",
					Command: e.Action,
					Package: pkgCandidate,
					Cause:   err,
				}
			}
			if !isDir {
				return &EngineError{
					Message: "failed to read package",
					Package: pkgCandidate,
					Command: e.Action,
				}
			}
			pkgCandidates = append(pkgCandidates, pkgCandidate)
		}

		log.Debug("retrieved candidates for packages", "candidates", pkgCandidates)
	}
	packages, err := filterPackages(pkgCandidates, e.Ignore)
	if err != nil {
		return err
	}
	e.PackageList = &packages
	log.Debug("package list populated", "package_list", e.PackageList)
	return nil
}

func filterPackages(candidates []string, ignoreList IgnoreList) ([]string, error) {
	log.Debug("filtering packages", "candidates", candidates, "filter", ignoreList.items)
	result := []string{}
	for _, candidate := range candidates {
		shouldIgnore, err := ignoreList.shouldIgnore(candidate, constant.RootPackageName)
		if err != nil {
			return nil, err
		}
		if shouldIgnore {
			log.Debug("Ignoring package candidate", "candidate", candidate)
			continue
		}
		result = append(result, candidate)
	}
	return result, nil
}
