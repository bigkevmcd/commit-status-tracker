package taskrun

import (
	"github.com/bigkevmcd/commit-status-tracker/pkg/tracker"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

type taskRunWrapper struct {
	*pipelinev1.TaskRun
}

func wrap(tr *pipelinev1.TaskRun) taskRunWrapper {
	return taskRunWrapper{tr}
}

// RunState returns whether or not a TaskRun was successful or
// not.
func (t taskRunWrapper) RunState() tracker.State {
	return tracker.ConditionsToState(t.Status.Conditions)
}

// Annotations returns the set of Annotations on the underlying TaskRun.
func (t taskRunWrapper) Annotations() map[string]string {
	return t.TaskRun.Annotations
}

// FindCommit attempts to find a GitCommit that can be tracked.
func (t taskRunWrapper) FindCommit() (*tracker.Commit, error) {
	return tracker.FindCommit(extractPipelineResources(t.Spec.Inputs.Resources))
}

func extractPipelineResources(bindings []pipelinev1.TaskResourceBinding) []*pipelinev1.PipelineResourceSpec {
	resources := make([]*pipelinev1.PipelineResourceSpec, len(bindings))
	for i, b := range bindings {
		resources[i] = b.ResourceSpec
	}
	return resources
}
