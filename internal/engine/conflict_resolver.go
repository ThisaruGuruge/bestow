package engine

type ExistingType int

const (
	ExistingManagedSymlink ExistingType = iota
	ExistingForeignSymlink
	ExistingRegularFile
	ExistingDir
)

type ConflictResolver interface {
	Resolve(src, dest string, existing ExistingType) (ResolveStrategy, error)
}

type StaticResolver struct {
	strategy ResolveStrategy
}

// TODO: Make sure to have a prune method to clear the history
func (sr StaticResolver) Resolve(src, dest string, existing ExistingType) (ResolveStrategy, error) {
	return sr.strategy, nil
}
