package pipelinerun

import (
	"math/rand"
	"reflect"
	"testing"

	tb "github.com/tektoncd/pipeline/test/builder"
	"knative.dev/pkg/apis"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

func TestFindGitResourceWithNoRepository(t *testing.T) {
	pipelineRun := makePipelineRunWithResources()

	_, err := FindGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestFindGitResourceWithRepository(t *testing.T) {
	pipelineRun := makePipelineRunWithResources(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"))

	wanted := &pipelinev1.PipelineResourceSpec{
		Type: "git",
		Params: []pipelinev1.ResourceParam{
			pipelinev1.ResourceParam{
				Name:  "url",
				Value: "https://github.com/tektoncd/triggers",
			},
			pipelinev1.ResourceParam{
				Name:  "revision",
				Value: "master",
			},
		},
	}

	r, err := FindGitResource(pipelineRun)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, wanted) {
		t.Fatalf("got %+v, wanted %+v", r, wanted)
	}
}

func TestFindGitResourceWithMultipleRepositories(t *testing.T) {
	pipelineRun := makePipelineRunWithResources(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"),
		makeGitResourceBinding("https://github.com/tektoncd/pipeline", "master"))

	_, err := FindGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestFindGitResourceWithNonGitResource(t *testing.T) {
	pipelineRun := makePipelineRunWithResources(
		makeImageResourceBinding("example.com/project/myimage"))

	_, err := FindGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func makePipelineRunWithResources(opts ...tb.PipelineRunSpecOp) *pipelinev1.PipelineRun {
	return tb.PipelineRun("pear", "foo", tb.PipelineRunSpec(
		"tomatoes", opts...,
	), tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
		apis.Condition{Type: apis.ConditionSucceeded}),
		tb.PipelineRunTaskRunsStatus("trname", &pipelinev1.PipelineRunTaskRunStatus{
			PipelineTaskName: "task-1",
		}),
	), tb.PipelineRunLabel("label-key", "label-value"))
}

func makeGitResourceBinding(url, rev string) tb.PipelineRunSpecOp {
	return tb.PipelineRunResourceBinding("some-resource"+randomSuffix(),
		tb.PipelineResourceBindingResourceSpec(&pipelinev1.PipelineResourceSpec{
			Type: pipelinev1.PipelineResourceTypeGit,
			Params: []pipelinev1.ResourceParam{{
				Name:  "url",
				Value: url,
			}, {
				Name:  "revision",
				Value: rev,
			}}}))
}

func makeImageResourceBinding(url string) tb.PipelineRunSpecOp {
	return tb.PipelineRunResourceBinding("some-resource"+randomSuffix(),
		tb.PipelineResourceBindingResourceSpec(&pipelinev1.PipelineResourceSpec{
			Type: pipelinev1.PipelineResourceTypeImage,
			Params: []pipelinev1.ResourceParam{{
				Name:  "url",
				Value: url,
			},
			}}))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomSuffix() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
