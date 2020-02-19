package pipelinerun

import (
	"github.com/bigkevmcd/commit-status-tracker/pkg/tracker"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

type pipelineRunWrapper struct {
	*pipelinev1.PipelineRun
}

func wrap(pr *pipelinev1.PipelineRun) pipelineRunWrapper {
	return pipelineRunWrapper{pr}
}

// RunState returns whether or not a PipelineRun was successful or
// not.
func (p pipelineRunWrapper) RunState() tracker.State {
	return tracker.ConditionsToState(p.Status.Conditions)
}

func (p pipelineRunWrapper) Annotations() map[string]string {
	return p.PipelineRun.Annotations
}

func (p pipelineRunWrapper) FindCommit() (*tracker.Commit, error) {
	return tracker.FindCommit(extractPipelineResources(p.Spec.Resources))
}

func extractPipelineResources(bindings []pipelinev1.PipelineResourceBinding) []*pipelinev1.PipelineResourceSpec {
	resources := make([]*pipelinev1.PipelineResourceSpec, len(bindings))
	for i, b := range bindings {
		resources[i] = b.ResourceSpec
	}
	return resources
}
