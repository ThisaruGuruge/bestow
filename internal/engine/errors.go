/*
All Rights Reversed (ɔ)
*/

package engine

import "fmt"

type EngineError struct {
	Message string
	Cause   error
	Hint    string
}

func (e *EngineError) Error() string {
	message := e.Message
	if e.Cause != nil {
		message = fmt.Sprintf("%s: %v", message, e.Cause)
	}
	if e.Hint != "" {
		message = fmt.Sprintf("%s: [Hint] %s", message, e.Hint)
	}
	return message
}

func (e *EngineError) Unwrap() error { return e.Cause }
