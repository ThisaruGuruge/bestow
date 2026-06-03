/*
All Rights Reversed (ɔ)
*/

package engine

import (
	"fmt"

	"github.com/ThisaruGuruge/bestow/internal/file"
	"github.com/ThisaruGuruge/bestow/internal/output"
)

const (
	actionLink   = "[link  ]"
	actionBackup = "[backup]"
	actionSkip   = "[skip  ]"
	actionAdopt  = "[adopt ]"
	actionRemove = "[remove]"
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

const backupExtension = ".bestow.backup"

type FileAction interface {
	Execute(fs file.System, label string) error
	Type() ActionType
	Source() string
	Destination() string
}

type fileActionBase struct {
	source      string
	destination string
}

func (fab fileActionBase) Source() string {
	return fab.source
}

func (fab fileActionBase) Destination() string {
	return fab.destination
}

type FileActionUpToDate struct {
	fileActionBase
	reason string
}

func newFileActionUpToDate(source, destination, reason string) *FileActionUpToDate {
	return &FileActionUpToDate{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
		},
		reason: reason,
	}
}

func (f *FileActionUpToDate) Execute(fs file.System, label string) error {
	output.PrintAction(label, actionSkip, fmt.Sprintf("%s -> %s [%s]", f.source, f.destination, f.reason), output.TypeWarn)
	return nil
}

func (f *FileActionUpToDate) Type() ActionType {
	return UpToDate
}

type FileActionSkip struct {
	fileActionBase
	reason string
}

func newFileActionSkip(source, destination, reason string) *FileActionSkip {
	return &FileActionSkip{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
		},
		reason: reason,
	}
}

func (f *FileActionSkip) Execute(fs file.System, label string) error {
	output.PrintAction(label, actionSkip, fmt.Sprintf("%s -> %s [%s]", f.destination, f.source, f.reason), output.TypeWarn)
	return nil
}

func (f *FileActionSkip) Type() ActionType {
	return Skip
}

type FileActionLink struct {
	fileActionBase
}

func newFileActionLink(source, destination string) *FileActionLink {
	return &FileActionLink{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
		},
	}
}

func (f *FileActionLink) Execute(fs file.System, label string) error {
	if err := fs.Link(f.source, f.destination); err != nil {
		return err
	}
	output.PrintAction(label, actionLink, fmt.Sprintf("%s -> %s", f.destination, f.source), output.TypeSuccess)
	return nil
}

func (f *FileActionLink) Type() ActionType {
	return Link
}

type FileActionReplace struct {
	fileActionBase
}

func newFileActionReplace(source, destination string) *FileActionReplace {
	return &FileActionReplace{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
		},
	}
}

func (f *FileActionReplace) Execute(fs file.System, label string) error {
	if err := fs.Remove(f.destination); err != nil {
		return err
	}
	output.PrintAction(label, actionRemove, f.destination, output.TypeStep)
	if err := fs.Link(f.source, f.destination); err != nil {
		return err
	}
	output.PrintAction(label, actionLink, fmt.Sprintf("%s -> %s", f.destination, f.source), output.TypeSuccess)
	return nil
}

func (f *FileActionReplace) Type() ActionType {
	return Replace
}

type FileActionBackup struct {
	fileActionBase
	backup string
}

func newFileActionBackup(source, destination, backup string) *FileActionBackup {
	return &FileActionBackup{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
		},
		backup: backup,
	}
}

func (f *FileActionBackup) Execute(fs file.System, label string) error {
	if err := fs.Move(f.destination, f.backup); err != nil {
		return err
	}
	output.PrintAction(label, actionBackup, fmt.Sprintf("%s -> %s", f.destination, f.backup), output.TypeStep)
	if err := fs.Link(f.source, f.destination); err != nil {
		return err
	}
	output.PrintAction(label, actionLink, fmt.Sprintf("%s -> %s", f.destination, f.source), output.TypeSuccess)
	return nil
}

func (f *FileActionBackup) Type() ActionType {
	return Backup
}

type FileActionAdopt struct {
	fileActionBase
}

func newFileActionAdopt(source, destination string) *FileActionAdopt {
	return &FileActionAdopt{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
		},
	}
}

func (f *FileActionAdopt) Execute(fs file.System, label string) error {
	if err := fs.Move(f.destination, f.source); err != nil {
		return err
	}
	output.PrintAction(label, actionAdopt, fmt.Sprintf("%s -> %s", f.destination, f.source), output.TypeStep)
	if err := fs.Link(f.source, f.destination); err != nil {
		return err
	}
	output.PrintAction(label, actionLink, fmt.Sprintf("%s -> %s", f.destination, f.source), output.TypeSuccess)
	return nil
}

func (f *FileActionAdopt) Type() ActionType {
	return Adopt
}

type FileActionRemove struct {
	fileActionBase
}

func newFileActionRemove(source, destination string) *FileActionRemove {
	return &FileActionRemove{
		fileActionBase: fileActionBase{
			source:      source,
			destination: destination,
		},
	}
}

func (f *FileActionRemove) Execute(fs file.System, label string) error {
	if err := fs.Remove(f.destination); err != nil {
		return err
	}
	output.PrintAction(label, actionRemove, f.destination, output.TypeSuccess)
	return nil
}

func (f *FileActionRemove) Type() ActionType {
	return Remove
}
