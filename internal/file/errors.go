package file

import "fmt"

type FileError struct {
	Message string
	Path    string
	Cause   error
}

func (e *FileError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Message, e.Path, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Path)
}

func (e *FileError) Unwrap() error { return e.Cause }
