/*
All Rights Reversed (ɔ)
*/

package file

type ExistingType string

const (
	ExistingManagedSymlink ExistingType = "managed_symlink"
	ExistingForeignSymlink ExistingType = "foreign_symlink"
	ExistingRegularFile    ExistingType = "regular_file"
	ExistingDir            ExistingType = "directory"
	ExistingUnknown        ExistingType = "unknown_type"
)

type LabelledHandler struct {
	label string
}

// Label returns the label that this filesystem operation should print.
// Helpful for printing dry run or file system label per each operation.
func (l *LabelledHandler) Label() string {
	return l.label
}
