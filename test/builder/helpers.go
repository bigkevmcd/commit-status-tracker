package builder

import (
	"math/rand"

	tb "github.com/tektoncd/pipeline/test/builder"
	"knative.dev/pkg/apis"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pipelineRunName = "test-pipeline-run"
	taskRunName     = "test-task-run"
	testNamespace   = "test-namespace"
)

func MakePipelineRunWithResources(res ...*pipelinev1.PipelineResourceSpec) *pipelinev1.PipelineRun {
	bound := make([]tb.PipelineRunSpecOp, 0)
	for _, r := range res {
		bound = append(bound, tb.PipelineRunResourceBinding(
			"testing"+MakeRandomString(),
			tb.PipelineResourceBindingResourceSpec(r)))
	}

	return tb.PipelineRun(pipelineRunName, testNamespace, tb.PipelineRunSpec(
		pipelineRunName,
		bound...,
	),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded}),
			tb.PipelineRunTaskRunsStatus("trname", &pipelinev1.PipelineRunTaskRunStatus{
				PipelineTaskName: "task-1",
			}),
		), tb.PipelineRunLabel("label-key", "label-value"))
}

func MakeTaskRunWithInputResources(res ...*pipelinev1.PipelineResourceSpec) *pipelinev1.TaskRun {
	inputs := make([]tb.TaskRunInputsOp, 0)
	for _, r := range res {
		inputs = append(inputs, tb.TaskRunInputsResource(
			"testing"+MakeRandomString(),
			tb.TaskResourceBindingResourceSpec(r)))
	}
	return tb.TaskRun(taskRunName, testNamespace, tb.TaskRunSpec(
		tb.TaskRunInputs(inputs...),
	), tb.TaskRunStatus(tb.StatusCondition(
		apis.Condition{Type: apis.ConditionSucceeded}),
	), tb.TaskRunLabel("label-key", "label-value"))
}

func MakeGitResource(url, rev string) *pipelinev1.PipelineResourceSpec {
	return MakePipelineResource(pipelinev1.PipelineResourceTypeGit, url, rev)
}

func MakePipelineResource(resType pipelinev1.PipelineResourceType, url, rev string) *pipelinev1.PipelineResourceSpec {
	spec := &pipelinev1.PipelineResourceSpec{
		Type: resType,
	}
	if url != "" {
		spec.Params = append(spec.Params,
			pipelinev1.ResourceParam{
				Name:  "url",
				Value: url,
			})
	}
	if rev != "" {
		spec.Params = append(spec.Params,
			pipelinev1.ResourceParam{
				Name:  "revision",
				Value: rev,
			})
	}
	return spec
}

func MakeSecret(name string, data map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		Type: corev1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Data: data,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func MakeRandomString() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func MakeGitPipelineResourceSpec(url, rev string) *pipelinev1.PipelineResourceSpec {
	return &pipelinev1.PipelineResourceSpec{
		Type: pipelinev1.PipelineResourceTypeGit,
		Params: []pipelinev1.ResourceParam{
			{
				Name:  "url",
				Value: url,
			}, {
				Name:  "revision",
				Value: rev,
			},
		},
	}
}
