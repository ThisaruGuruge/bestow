/*
All Rights Reversed (ɔ)
*/

package output

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

var successStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Green)

var StepStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Cyan)

var warnStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Yellow)

var hintStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Magenta)

func Success(message string, args ...any) {
	text := fmt.Sprintf(message, args...)
	fmt.Println(successStyle.Render(text))
}

func PrintAction(label, action, msg string, t Type) {
	message := fmt.Sprintf("%s %s %s", label, action, msg)
	var text string
	switch t {
	case TypeSuccess:
		text = successStyle.Render(message)
	case TypeStep:
		text = StepStyle.Render(message)
	case TypeWarn:
		text = warnStyle.Render(message)
	case TypeHint:
		text = hintStyle.Render(message)
	}
	lipgloss.Println(text)
}
