package pipelinerun

import (
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"

	"github.com/bigkevmcd/commit-status-tracker/pkg/controller/tracker"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

type pipelineRunWrapper struct {
	*pipelinev1.PipelineRun
}

// RunState returns whether or not a PipelineRun was successful or
// not.
//
// It can return a Pending result if the task has not yet completed.
// TODO: will likely need to work out if a task was killed OOM.
func (p pipelineRunWrapper) RunState() tracker.State {
	for _, c := range p.Status.Conditions {
		if c.Type == apis.ConditionSucceeded {
			switch c.Status {
			case
				corev1.ConditionFalse:
				return tracker.Failed
			case corev1.ConditionTrue:
				return tracker.Successful
			case corev1.ConditionUnknown:
				return tracker.Pending
			}
		}
	}
	return tracker.Pending
}
