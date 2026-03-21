package log

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	charmlog "github.com/charmbracelet/log"
)

func NewCharmLogger() *CharmLogger {
	styles := charmlog.DefaultStyles()

	styles.Levels[charmlog.InfoLevel] = lipgloss.NewStyle().
		SetString("[Info]:").
		Foreground(lipgloss.AdaptiveColor{
			Light: "2",
			Dark:  "10",
		})
	styles.Levels[charmlog.WarnLevel] = lipgloss.NewStyle().
		SetString("[Warn]:").
		Foreground(lipgloss.AdaptiveColor{
			Light: "3",
			Dark:  "11",
		})
	styles.Levels[charmlog.ErrorLevel] = lipgloss.NewStyle().
		SetString("[Error]:").
		Foreground(lipgloss.AdaptiveColor{
			Light: "1",
			Dark:  "9",
		})
	styles.Levels[charmlog.DebugLevel] = lipgloss.NewStyle().
		SetString("[Debug]:").
		Foreground(lipgloss.AdaptiveColor{
			Light: "6",
			Dark:  "14",
		})

	l := charmlog.New(os.Stderr)
	l.SetStyles(styles)
	return &CharmLogger{logger: l}
}

type CharmLogger struct {
	logger *charmlog.Logger
}

func (l *CharmLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}
func (l *CharmLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}
func (l *CharmLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
func (l *CharmLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}
func (l *CharmLogger) SetLevel(level Level) {
	switch level {
	case LevelDebug:
		l.logger.SetLevel(charmlog.DebugLevel)
	case LevelInfo:
		l.logger.SetLevel(charmlog.InfoLevel)
	case LevelWarn:
		l.logger.SetLevel(charmlog.WarnLevel)
	case LevelError:
		l.logger.SetLevel(charmlog.ErrorLevel)
	}
}
