package tracker

type stateGetter interface {
	RunState() State
}

type annotationsGetter interface {
	Annotations() map[string]string
}

// FindCommit locates a Git PipelineResource and extracts the details.
//
// If no Git resources are found, an error should be returned.
// If more than one Git resource is found, an error should be returned.
type gitRefFinder interface {
	FindCommit() (*Commit, error)
}

type trackableResource interface {
	stateGetter
	annotationsGetter
	gitRefFinder
}
