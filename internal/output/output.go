/*
All Rights Reversed (ɔ)
*/

package output

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ThisaruGuruge/bestow/internal/engine"
)

var successStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Green)

var stepStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Cyan)

var warnStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Yellow)

var hintStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Magenta)

func PrintAction(action engine.ActionEvent) {
	var message string
	if action.Label == "" {
		message = fmt.Sprintf("%s %s", action.Action, action.Msg)
	} else {
		message = fmt.Sprintf("%s %s %s", action.Label, action.Action, action.Msg)
	}
	var text string
	switch action.EventType {
	case engine.EventSuccess:
		text = successStyle.Render(message)
	case engine.EventStep:
		text = stepStyle.Render(message)
	case engine.EventWarn:
		text = warnStyle.Render(message)
	case engine.EventIgnore:
		return
	}
	lipgloss.Println(text)
}

func PrintSummary(summary *engine.ExecuteSummary) {
	for _, action := range *summary.Actions {
		PrintAction(action)
	}
	summaryFields := 7
	parts := make([]string, 0, summaryFields)
	if summary.OperationSummary.Stowed > 0 {
		parts = append(parts, fmt.Sprintf("stowed: %d", summary.OperationSummary.Stowed))
	}
	if summary.OperationSummary.Unstowed > 0 {
		parts = append(parts, fmt.Sprintf("unstowed: %d", summary.OperationSummary.Unstowed))
	}
	if summary.OperationSummary.Replaced > 0 {
		parts = append(parts, fmt.Sprintf("replaced: %d", summary.OperationSummary.Replaced))
	}
	if summary.OperationSummary.Backed > 0 {
		parts = append(parts, fmt.Sprintf("backed up: %d", summary.OperationSummary.Backed))
	}
	if summary.OperationSummary.Adopted > 0 {
		parts = append(parts, fmt.Sprintf("adopted: %d", summary.OperationSummary.Adopted))
	}
	if summary.OperationSummary.Skipped > 0 {
		parts = append(parts, fmt.Sprintf("skipped: %d", summary.OperationSummary.Skipped))
	}
	if summary.OperationSummary.UpToDate > 0 {
		parts = append(parts, fmt.Sprintf("up to date: %d", summary.OperationSummary.UpToDate))
	}
	if len(parts) == 0 {
		return
	}
	lipgloss.Println(strings.Join(parts, "   "))
}

func PrintHint(hint string) {
	message := "[hint] " + hint
	lipgloss.Fprintln(os.Stderr, hintStyle.Render(message))
}
