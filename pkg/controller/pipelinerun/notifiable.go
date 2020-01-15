package pipelinerun

import (
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

// isNotifiablePipelineRun returns true if this PipelineRun should report its
// completion status as a GitHub status.
func isNotifiablePipelineRun(pr *pipelinev1.PipelineRun) bool {
	for k, v := range pr.Labels {
		if k == notifiableLabel && v == "true" {
			return true
		}
	}
	return false
}
