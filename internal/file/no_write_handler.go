/*
All Rights Reversed (ɔ)
*/

package file

import (
	"log/slog"
	"path/filepath"
)

// NoWriteHandler is the implementation of the System using io, os, and bufio go modules.
type NoWriteHandler struct {
	baseHandler
	createdDirs map[string]bool
}

// NewNoWriteHandler returns a new NoWriteHandler with the provided logger l.
func NewNoWriteHandler(l *slog.Logger) *NoWriteHandler {
	return &NoWriteHandler{
		createdDirs: make(map[string]bool),
		baseHandler: baseHandler{logger: l},
	}
}

// CreateFile creates a file in the provided path and writes the provided content to the file.
func (h *NoWriteHandler) CreateFile(path, content string) error {
	h.logger.Debug("writing to file", "file", path)
	h.logger.Debug("successfully written to file", "path", path)
	return nil
}

// CreateDir creates a directory on the provided path, including all the parent directories.
func (h *NoWriteHandler) CreateDir(path string) error {
	h.logger.Debug("creating directory", "path", path)
	if h.createdDirs[path] {
		h.logger.Debug("directory already created", "path", path)
		return nil
	}
	exists, err := h.IsDir(path)
	if err != nil {
		return err
	}
	if exists {
		h.logger.Debug("directory already exists", "path", path)
		h.createdDirs[path] = true
		return nil
	}
	h.createdDirs[path] = true
	h.logger.Debug("created directory", "path", path)
	return nil
}

// Link creates a symlink of a provided src in the provided target.
// If the target directory does not exist, link will create all the parent directories.
func (h *NoWriteHandler) Link(src, target string) error {
	h.logger.Debug("creating symlink", "source", src, "target", target)
	destParent := filepath.Dir(target)
	if err := h.CreateDir(destParent); err != nil {
		return err
	}
	h.logger.Debug("link created", "source", src, "target", target)
	return nil
}

// Move moves a file from src to target
// If the target directory does not exist, move will create all the parent directories.
func (h *NoWriteHandler) Move(src, target string) error {
	h.logger.Debug("moving file", "source", src, "target", target)
	destParent := filepath.Dir(target)
	if err := h.CreateDir(destParent); err != nil {
		return err
	}
	h.logger.Debug("moved file", "from", src, "to", target)
	return nil
}

// Remove removes the file in the provided path.
func (h *NoWriteHandler) Remove(path string) error {
	h.logger.Debug("removing the file", "path", path)
	exists, err := h.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		h.logger.Warn("file does not exist", "operation", "remove", "file", path)
		return nil
	}
	h.logger.Debug("successfully removed the file", "file_name", path)
	return nil
}
