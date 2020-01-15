package pipelinerun

import (
	"testing"

	tb "github.com/tektoncd/pipeline/test/builder"
)

func TestLabelByName(t *testing.T) {
	nt := []struct {
		name string
		opts []tb.PipelineRunOp
		want string
	}{
		{"no labels", nil, "default"},
		{"no matching label",
			[]tb.PipelineRunOp{tb.PipelineRunAnnotation("testing", "app")},
			"default"},
		{"with matching label",
			[]tb.PipelineRunOp{tb.PipelineRunAnnotation(statusContextLabel, "test-lint")},
			"test-lint"},
	}

	for _, tt := range nt {
		r := tb.PipelineRun("test-pipeline-run-with-labels", "foo", tt.opts...)
		if b := getAnnotationByName(r, statusContextLabel, "default"); b != tt.want {
			t.Errorf("Context() %s got %v, want %v", tt.name, b, tt.want)
		}
	}
}
