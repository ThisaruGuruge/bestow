/*
All Rights Reversed (ɔ)
*/

package file

type ExistingType int

const (
	ExistingManagedSymlink ExistingType = iota
	ExistingForeignSymlink
	ExistingRegularFile
	ExistingDir
	ExistingUnknown
)
