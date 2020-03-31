package tracker

import (
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1beta1"
)

// State represents the state of a Pipeline.
type State int

const (
	Pending State = iota
	Failed
	Successful
	Error
)

func (s State) String() string {
	names := [...]string{
		"Pending",
		"Failed",
		"Successful",
		"Error"}
	return names[s]
}

// ConditionToState processes a set of conditions looking for a
// ConditionSucceeded and returns a commit-status compatible state for the run.
//
// It can return a Pending result if the task has not yet completed.
// TODO: will likely need to work out if a task was killed OOM.
func ConditionsToState(conditions duckv1.Conditions) State {
	for _, c := range conditions {
		if c.Type == apis.ConditionSucceeded {
			switch c.Status {
			case
				corev1.ConditionFalse:
				return Failed
			case corev1.ConditionTrue:
				return Successful
			case corev1.ConditionUnknown:
				return Pending
			}
		}
	}
	return Pending
}
