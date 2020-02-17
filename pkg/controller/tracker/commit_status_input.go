package tracker

import (
	"github.com/jenkins-x/go-scm/scm"
)

type stateHelper interface {
	RunState() State
	Annotations() map[string]string
}

// getCommitStatusInput extracts the various bits from a PipelineRun and
// returns a status record for submitting to the upstream Git Hosting
// Service.
//
// See https://developer.github.com/v3/repos/statuses/#create-a-status and
// https://github.com/jenkins-x/go-scm/blob/b48d209334ed7b167bad3326a481ae3964c7c1a1/scm/repo.go#L88
func getCommitStatusInput(r stateHelper) *scm.StatusInput {
	return &scm.StatusInput{
		State:  convertState(r.RunState()),
		Label:  getAnnotationByName(r, StatusContextName, "default"),
		Desc:   getAnnotationByName(r, StatusDescriptionName, ""),
		Target: getAnnotationByName(r, StatusTargetURLName, ""),
	}
}

func getAnnotationByName(r stateHelper, name, def string) string {
	for k, v := range r.Annotations() {
		if k == name {
			return v
		}
	}
	return def
}

// convertState converts between pipeline run state, and the commit status.
func convertState(s State) scm.State {
	switch s {
	case Failed:
		return scm.StateFailure
	case Pending:
		return scm.StatePending
	case Successful:
		return scm.StateSuccess
	default:
		return scm.StateUnknown
	}
}
