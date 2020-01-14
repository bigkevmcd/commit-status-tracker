package pipelinerun

import (
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

const (
	notifiable = "app.example.com/git-status"
)

// IsNotifiablePipelineRun returns true if this PipelineRun should report its
// completion status as a GitHub status.
func IsNotifiablePipelineRun(p *pipelinev1.PipelineRun) bool {
	for k, v := range p.Labels {
		if k == notifiable && v == "true" {
			return true
		}
	}
	return false
}
