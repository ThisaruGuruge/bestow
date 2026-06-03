/*
All Rights Reversed (ɔ)
*/

package engine

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ThisaruGuruge/bestow/internal/config"
	"github.com/ThisaruGuruge/bestow/internal/constant"
)

type InitContext struct {
	Force      bool
	IgnoreList []string
}

func (e *Engine) Init(ctx *InitContext) (*ExecuteSummary, error) {
	e.logger.Debug("initializing bestow")
	appConfigDir := config.AppConfigHome()
	configFile := filepath.Join(appConfigDir, constant.ConfigFile)
	ignoreFile := filepath.Join(appConfigDir, constant.IgnoreFile)
	if err := e.checkExistingFiles(configFile, ignoreFile, ctx.Force); err != nil {
		return nil, err
	}
	if err := e.fileSystem.CreateDir(appConfigDir); err != nil {
		return nil, err
	}
	configAction, err := e.createConfigFile(e.source, e.destination, ctx.Force, appConfigDir)
	if err != nil {
		return nil, err
	}
	ignoreAction, err := e.createIgnoreFile(appConfigDir, ctx.Force, ctx.IgnoreList)
	if err != nil {
		return nil, err
	}
	actions := []ActionEvent{*configAction, *ignoreAction}
	return &ExecuteSummary{Actions: actions, OperationSummary: &Summary{}}, nil
}

func (e *Engine) checkExistingFiles(configFile, ignoreFile string, force bool) error {
	if force {
		return nil
	}
	configExists, err := e.fileSystem.Exists(configFile)
	if err != nil {
		return err
	}
	ignoreExists, err := e.fileSystem.Exists(ignoreFile)
	if err != nil {
		return err
	}
	existing := make([]string, 0, 2)
	if configExists {
		existing = append(existing, configFile)
	}
	if ignoreExists {
		existing = append(existing, ignoreFile)
	}
	if len(existing) > 0 {
		fileString := strings.Join(existing, ", ")
		return &HintedError{
			Op:   fmt.Sprintf("exists %s", fileString),
			Hint: "remove the existing files or use --force",
			Err:  ErrFileExists,
		}
	}
	return nil
}

func (e *Engine) createIgnoreFile(appConfigDir string, force bool, ignoreList []string) (*ActionEvent, error) {
	e.logger.Debug("creating ignore file")
	ignoreFile := filepath.Join(appConfigDir, constant.IgnoreFile)
	exists, err := e.fileSystem.Exists(ignoreFile)
	if err != nil {
		return nil, err
	}
	if exists {
		if !force {
			return nil, &HintedError{
				Op:   fmt.Sprintf("create ignorefile %s", ignoreFile),
				Hint: "use --force to overwrite",
				Err:  ErrFileExists,
			}
		}
		e.logger.Warn("ignore file exists; overwriting", "ignore-file", ignoreFile)
	}
	e.logger.Debug("initializing ignore list", "ignore-list", ignoreList)
	if err := e.fileSystem.CreateFile(ignoreFile, getIgnoreFileContent(ignoreList)); err != nil {
		return nil, err
	}
	return &ActionEvent{
		Label:     "[init]",
		Action:    actionCreated,
		Msg:       ignoreFile,
		EventType: EventSuccess,
	}, nil
}

func getIgnoreFileContent(ignoreList []string) string {
	result := []string{"# Global Ignore List for Bestow"}
	for _, item := range ignoreList {
		result = append(result, strings.TrimSpace(item))
	}
	return strings.Join(result, "\n")
}

func (e *Engine) createConfigFile(source, destination string, force bool, appConfigDir string) (*ActionEvent, error) {
	configFile := filepath.Join(appConfigDir, constant.ConfigFile)
	e.logger.Debug("creating the config file", "path", configFile)
	exists, err := e.fileSystem.Exists(configFile)
	if err != nil {
		return nil, err
	}
	if exists {
		if !force {
			return nil, &HintedError{
				Op:   fmt.Sprintf("create configfile %s", configFile),
				Hint: "use --force to overwrite",
				Err:  ErrFileExists,
			}
		}
		e.logger.Warn("config file exists; overwriting", "config-file", configFile)
	}
	config, err := config.GetDefaultConfigTemplate(source, destination)
	if err != nil {
		return nil, fmt.Errorf("load config %s %s: %w", source, destination, err)
	}
	if err := e.fileSystem.CreateFile(configFile, config); err != nil {
		return nil, err
	}
	return &ActionEvent{
		Label:     "[init]",
		Action:    actionCreated,
		Msg:       configFile,
		EventType: EventSuccess,
	}, nil
}
