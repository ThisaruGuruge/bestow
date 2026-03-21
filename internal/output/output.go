package output

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var successStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("2"))

func Success(message string) {
	fmt.Println(successStyle.Render(fmt.Sprintf("[Success]: %s", message)))
}
