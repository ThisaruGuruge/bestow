/*
All Rights Reversed (ɔ)
*/

package engine

import (
	"fmt"
	"log/slog"
)

const (
	actionLink    = "link"
	actionBackup  = "backup"
	actionSkip    = "skip"
	actionAdopt   = "adopt"
	actionRemove  = "remove"
	actionCreated = "created"
)

type EventType int

const (
	EventSuccess EventType = iota
	EventStep
	EventSkip
	EventWarn
	EventIgnore
)

type ActionType int

const (
	UpToDate ActionType = iota
	Skip
	Link
	Replace
	Backup
	Adopt
	Remove
)

func (a ActionType) String() string {
	switch a {
	case UpToDate:
		return "up-to-date"
	case Skip:
		return "skip"
	case Link:
		return "link"
	case Replace:
		return "replace"
	case Backup:
		return "backup"
	case Adopt:
		return "adopt"
	case Remove:
		return "remove"
	default:
		return fmt.Sprintf("Unknown %d", a)
	}
}

type ActionEvent struct {
	Action    string
	Msg       string
	EventType EventType
}

const backupExtension = "bestow.backup"
const tmpExtension = "bestow.tmp"

type fileAction interface {
	Execute(fs FileSystem) ([]ActionEvent, error)
	Type() ActionType
	Source() string
	Destination() string
}

type fileActionBase struct {
	source      string
	destination string
	logger      *slog.Logger
}

func (fab fileActionBase) Source() string {
	return fab.source
}

func (fab fileActionBase) Destination() string {
	return fab.destination
}

type fileActionUpToDate struct {
	fileActionBase
	reason string
}

func newFileActionUpToDate(source, destination, reason string, l *slog.Logger) *fileActionUpToDate {
	return &fileActionUpToDate{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
			logger:      l,
		},
		reason: reason,
	}
}

func (f *fileActionUpToDate) Execute(fs FileSystem) ([]ActionEvent, error) {
	f.logger.Debug(f.reason, "source", f.source, "destination", f.destination)
	return []ActionEvent{
		{EventType: EventIgnore},
	}, nil
}

func (f *fileActionUpToDate) Type() ActionType {
	return UpToDate
}

type fileActionSkip struct {
	fileActionBase
	reason string
}

func newFileActionSkip(source, destination, reason string, l *slog.Logger) *fileActionSkip {
	return &fileActionSkip{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
			logger:      l,
		},
		reason: reason,
	}
}

func (f *fileActionSkip) Execute(fs FileSystem) ([]ActionEvent, error) {
	return []ActionEvent{
		{
			Action:    actionSkip,
			Msg:       fmt.Sprintf("%s -> %s [%s]", f.source, f.destination, f.reason),
			EventType: EventSkip,
		},
	}, nil
}

func (f *fileActionSkip) Type() ActionType {
	return Skip
}

type fileActionLink struct {
	fileActionBase
}

func newFileActionLink(source, destination string, l *slog.Logger) *fileActionLink {
	return &fileActionLink{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
			logger:      l,
		},
	}
}

func (f *fileActionLink) Execute(fs FileSystem) ([]ActionEvent, error) {
	if err := fs.Link(f.source, f.destination); err != nil {
		return nil, err
	}
	return []ActionEvent{
		{
			Action:    actionLink,
			Msg:       fmt.Sprintf("%s -> %s", f.destination, f.source),
			EventType: EventSuccess,
		},
	}, nil
}

func (f *fileActionLink) Type() ActionType {
	return Link
}

type fileActionReplace struct {
	fileActionBase
}

func newFileActionReplace(source, destination string, l *slog.Logger) *fileActionReplace {
	return &fileActionReplace{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
			logger:      l,
		},
	}
}

func (f *fileActionReplace) Execute(fs FileSystem) ([]ActionEvent, error) {
	var events []ActionEvent
	tmp := fmt.Sprintf("%s.%s", f.destination, tmpExtension)
	if err := fs.Move(f.destination, tmp); err != nil {
		return nil, err
	}
	removeStep := ActionEvent{
		Action:    actionRemove,
		Msg:       f.destination,
		EventType: EventStep,
	}
	events = append(events, removeStep)
	if err := fs.Link(f.source, f.destination); err != nil {
		if err := fs.Move(tmp, f.destination); err != nil {
			f.logger.Warn("failed to restore the tmp", "tmp_file", tmp, "original_file", f.destination)
			return nil, fmt.Errorf("recover %s %s: %w", tmp, f.destination, err)
		}
		return nil, err
	}
	if err := fs.Remove(tmp); err != nil {
		f.logger.Warn("failed to remove the tmp", "tmp_file", tmp)
		return nil, fmt.Errorf("remove %s: %w", tmp, err)
	}
	linkStep := ActionEvent{
		Action:    actionLink,
		Msg:       fmt.Sprintf("%s -> %s", f.destination, f.source),
		EventType: EventSuccess,
	}
	events = append(events, linkStep)
	return events, nil
}

func (f *fileActionReplace) Type() ActionType {
	return Replace
}

type fileActionBackup struct {
	fileActionBase
	backup string
}

func newFileActionBackup(source, destination, backup string, l *slog.Logger) *fileActionBackup {
	return &fileActionBackup{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
			logger:      l,
		},
		backup: backup,
	}
}

func (f *fileActionBackup) Execute(fs FileSystem) ([]ActionEvent, error) {
	if err := fs.Move(f.destination, f.backup); err != nil {
		return nil, err
	}
	var events []ActionEvent
	moveStep := ActionEvent{
		Action:    actionBackup,
		Msg:       fmt.Sprintf("%s -> %s", f.destination, f.backup),
		EventType: EventStep,
	}
	events = append(events, moveStep)
	if err := fs.Link(f.source, f.destination); err != nil {
		if err := fs.Move(f.backup, f.destination); err != nil {
			f.logger.Warn("failed to restore the backup", "backup_file", f.backup, "original_file", f.destination)
			return nil, fmt.Errorf("recover %s %s: %w", f.backup, f.destination, err)
		}
		return nil, err
	}
	linkStep := ActionEvent{
		Action:    actionLink,
		Msg:       fmt.Sprintf("%s -> %s", f.destination, f.source),
		EventType: EventSuccess,
	}
	events = append(events, linkStep)
	return events, nil
}

func (f *fileActionBackup) Type() ActionType {
	return Backup
}

type fileActionAdopt struct {
	fileActionBase
}

func newFileActionAdopt(source, destination string, l *slog.Logger) *fileActionAdopt {
	return &fileActionAdopt{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
			logger:      l,
		},
	}
}

func (f *fileActionAdopt) Execute(fs FileSystem) ([]ActionEvent, error) {
	if err := fs.Move(f.destination, f.source); err != nil {
		return nil, err
	}
	var events []ActionEvent
	moveStep := ActionEvent{
		Action:    actionAdopt,
		Msg:       fmt.Sprintf("%s -> %s", f.destination, f.source),
		EventType: EventStep,
	}
	events = append(events, moveStep)
	if err := fs.Link(f.source, f.destination); err != nil {
		if err := fs.Move(f.source, f.destination); err != nil {
			f.logger.Warn("failed to restore the original", "new_file", f.source, "original_file", f.destination)
			return nil, fmt.Errorf("recover %s %s: %w", f.source, f.destination, err)
		}
		return nil, err
	}
	linkStep := ActionEvent{
		Action:    actionLink,
		Msg:       fmt.Sprintf("%s -> %s", f.destination, f.source),
		EventType: EventSuccess,
	}
	events = append(events, linkStep)
	return events, nil
}

func (f *fileActionAdopt) Type() ActionType {
	return Adopt
}

type fileActionRemove struct {
	fileActionBase
}

func newFileActionRemove(source, destination string, l *slog.Logger) *fileActionRemove {
	return &fileActionRemove{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
			logger:      l,
		},
	}
}

func (f *fileActionRemove) Execute(fs FileSystem) ([]ActionEvent, error) {
	if err := fs.Remove(f.destination); err != nil {
		return nil, err
	}
	return []ActionEvent{
		{
			Action:    actionRemove,
			Msg:       f.destination,
			EventType: EventSuccess,
		},
	}, nil
}

func (f *fileActionRemove) Type() ActionType {
	return Remove
}
