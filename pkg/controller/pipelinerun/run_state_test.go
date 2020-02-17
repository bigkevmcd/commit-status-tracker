package pipelinerun

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"

	"github.com/bigkevmcd/commit-status-tracker/pkg/controller/tracker"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
)

func TestGetPipelineRunStatus(t *testing.T) {
	statusTests := []struct {
		conditionType   apis.ConditionType
		conditionStatus corev1.ConditionStatus
		want            tracker.State
	}{
		{apis.ConditionSucceeded, corev1.ConditionTrue, tracker.Successful},
		{apis.ConditionSucceeded, corev1.ConditionUnknown, tracker.Pending},
		{apis.ConditionSucceeded, corev1.ConditionFalse, tracker.Failed},
	}

	for _, tt := range statusTests {
		w := pipelineRunWrapper{makePipelineRunWithCondition(tt.conditionType, tt.conditionStatus)}
		s := w.RunState()
		if s != tt.want {
			t.Errorf("RunState(%s) got %v, want %v", tt.conditionStatus, s, tt.want)
		}
	}
}

func makePipelineRunWithCondition(s apis.ConditionType, c corev1.ConditionStatus) *pipelinev1.PipelineRun {
	return tb.PipelineRun(pipelineRunName, testNamespace, tb.PipelineRunSpec(
		"tomatoes",
	), tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
		apis.Condition{Type: s, Status: c}),
		tb.PipelineRunTaskRunsStatus("trname", &pipelinev1.PipelineRunTaskRunStatus{
			PipelineTaskName: "task-1",
		})))
}
