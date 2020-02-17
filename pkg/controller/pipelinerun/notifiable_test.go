package pipelinerun

import (
	"testing"

	"github.com/bigkevmcd/commit-status-tracker/pkg/controller/tracker"
	tb "github.com/tektoncd/pipeline/test/builder"
)

func TestIsNotifiablePipelineRun(t *testing.T) {
	nt := []struct {
		name string
		opts []tb.PipelineRunOp
		want bool
	}{
		{"no labels", nil, false},
		{"no notifiable label", []tb.PipelineRunOp{tb.PipelineRunAnnotation("testing", "app")}, false},
		{"notifiable label", []tb.PipelineRunOp{tb.PipelineRunAnnotation(tracker.NotifiableName, "true")}, true},
		{"notifiable label is false", []tb.PipelineRunOp{tb.PipelineRunAnnotation(tracker.NotifiableName, "false")}, false},
	}

	for _, tt := range nt {
		r := tb.PipelineRun("test-pipeline-run-with-labels", "foo", tt.opts...)
		if b := isNotifiablePipelineRun(r); b != tt.want {
			t.Errorf("IsNotifiablePipelineRun() %s got %v, want %v", tt.name, b, tt.want)
		}
	}
}
