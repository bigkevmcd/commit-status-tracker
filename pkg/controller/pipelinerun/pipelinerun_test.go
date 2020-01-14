package pipelinerun

import (
	"testing"

	tb "github.com/tektoncd/pipeline/test/builder"
)

func TestIsNotifiablePipelineRun(t *testing.T) {
	nt := []struct {
		name   string
		opts   []tb.PipelineRunOp
		wanted bool
	}{
		{"no labels", nil, false},
		{"no notifiable label", []tb.PipelineRunOp{tb.PipelineRunLabel("testing", "app")}, false},
		{"notifiable label", []tb.PipelineRunOp{tb.PipelineRunLabel(notifiable, "true")}, true},
		{"notifiable label is false", []tb.PipelineRunOp{tb.PipelineRunLabel(notifiable, "false")}, false},
	}

	for _, tt := range nt {
		r := tb.PipelineRun("test-pipeline-run-with-labels", "foo", tt.opts...)
		if b := IsNotifiablePipelineRun(r); b != tt.wanted {
			t.Errorf("IsNotifiablePipelineRun() %s got %v, wanted %v", tt.name, b, tt.wanted)
		}
	}
}
