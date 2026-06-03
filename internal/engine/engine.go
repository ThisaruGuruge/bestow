/*
All Rights Reversed (ɔ)
*/

package engine

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/ThisaruGuruge/bestow/internal/config"
	"github.com/ThisaruGuruge/bestow/internal/file"
	"github.com/ThisaruGuruge/bestow/internal/output"
)

type Action int

const (
	ActionStow Action = iota
	ActionUnstow
)

type CommandContext struct {
	Action           Action
	Args             []string
	ConflictStrategy ResolveStrategy
}

type Engine struct {
	source      string
	destination string
	ignore      *IgnoreList
	logger      *slog.Logger
	fileSystem  file.System
	actionLabel string
}

func NewEngine(cfg *config.Config, dryrun bool, l *slog.Logger) (*Engine, error) {
	var handler file.System
	if dryrun {
		handler = file.NewNoWriteHandler(l)
	} else {
		handler = file.NewHandler(l)
	}
	ignoreList, err := newIgnoreList(cfg.Source, handler, l)
	if err != nil {
		return nil, err
	}
	label := ""
	if dryrun {
		label = "[dryrun]"
	}
	return &Engine{
		source:      cfg.Source,
		destination: cfg.Destination,
		ignore:      ignoreList,
		logger:      l.With("component", "engine"),
		fileSystem:  handler,
		actionLabel: label,
	}, nil
}

func (e *Engine) Execute(ctx *CommandContext) error {
	actions, err := e.populateOperations(ctx)
	if err != nil {
		return err
	}
	summary, err := e.executeFileActions(actions)
	if err != nil {
		return err
	}
	output.PrintSummary(summary)
	return nil
}

// TODO: When skipping files;
// - in .bestowignore: debug log
// - skip because already stowed (due to state of the operation): include a summary
// - skip because conflict resolution strategy is set to skip: print as same as any other operation
func (e *Engine) executeFileActions(actions []FileAction) (*output.Summary, error) {
	summary := &output.Summary{}
	for _, action := range actions {
		if err := action.Execute(e.fileSystem, e.actionLabel); err != nil {
			return nil, err
		}
		actionType := action.Type()
		switch actionType {
		case UpToDate:
			summary.UpToDate += 1
		case Skip:
			summary.Skipped += 1
		case Link:
			summary.Stowed += 1
		case Replace:
			summary.Replaced += 1
		case Backup:
			summary.Backed += 1
		case Adopt:
			summary.Adopted += 1
		case Remove:
			summary.Unstowed += 1
		default:
			panic(fmt.Sprintf("undefined action %d", actionType))
		}

	}
	return summary, nil
}

func (e *Engine) populatePackageList(args []string) ([]string, error) {
	e.logger.Debug("populating package list", "source", e.source)
	var pkgCandidates []string
	var err error
	if len(args) == 0 {
		e.logger.Debug("no packages provided; processing all packages")
		pkgCandidates, err = e.getAllPackages()
		if err != nil {
			return nil, err
		}
	} else {
		pkgCandidates, err = e.getPackagesFromArgs(args)
		if err != nil {
			return nil, err
		}
	}
	packages, err := e.filterPackages(pkgCandidates)
	if err != nil {
		return nil, err
	}
	e.logger.Debug("package list populated", "package_list", packages)
	return packages, nil
}

func (e *Engine) getAllPackages() ([]string, error) {
	dirs, err := e.fileSystem.ListDirs(e.source)
	if err != nil {
		return nil, err
	}
	candidates := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		candidate, err := filepath.Rel(e.source, dir)
		if err != nil {
			return nil, fmt.Errorf("rel %s %s: %w", e.source, dir, err)
		}
		candidates = append(candidates, candidate)
	}
	return candidates, nil
}

func (e *Engine) getPackagesFromArgs(candidates []string) ([]string, error) {
	result := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate == "." {
			return nil, &HintedError{
				Op:   fmt.Sprintf("read package %s", candidate),
				Hint: "move root files to suitable directory (`zsh/`, `bash/`, etc.)",
				Err:  ErrRootIsNotPkg,
			}
		}
		pkgPath := filepath.Clean(candidate)
		isDir, err := e.fileSystem.IsDir(filepath.Join(e.source, pkgPath))
		if err != nil {
			return nil, fmt.Errorf("read package %s: %w", candidate, err)
		}
		if !isDir {
			return nil, &HintedError{
				Op:   fmt.Sprintf("read package %s", candidate),
				Hint: fmt.Sprintf("make sure the %s is a directory", candidate),
				Err:  ErrPkgIsNotDir,
			}
		}
		result = append(result, pkgPath)
	}
	return result, nil
}

func (e *Engine) filterPackages(candidates []string) ([]string, error) {
	e.logger.Debug("filtering packages", "candidates", candidates, "filter", e.ignore.items)
	result := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		shouldIgnore, err := e.ignore.shouldIgnorePkg(candidate)
		if err != nil {
			return nil, err
		}
		if shouldIgnore {
			e.logger.Debug("ignoring package candidate", "candidate", candidate)
			continue
		}
		e.logger.Debug("adding package to process", "package", candidate)
		result = append(result, candidate)
	}
	return result, nil
}
