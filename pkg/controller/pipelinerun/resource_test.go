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

	_, err := findGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestFindGitResourceWithRepository(t *testing.T) {
	pipelineRun := makePipelineRunWithResources(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"))

	want := &pipelinev1.PipelineResourceSpec{
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

	r, err := findGitResource(pipelineRun)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, want) {
		t.Fatalf("got %+v, want %+v", r, want)
	}
}

func TestFindGitResourceWithMultipleRepositories(t *testing.T) {
	pipelineRun := makePipelineRunWithResources(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"),
		makeGitResourceBinding("https://github.com/tektoncd/pipeline", "master"))

	_, err := findGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestFindGitResourceWithNonGitResource(t *testing.T) {
	pipelineRun := makePipelineRunWithResources(
		makeImageResourceBinding("example.com/project/myimage"))

	_, err := findGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestGetRepoAndSHA(t *testing.T) {
	repoURL := "https://example.com/test/repo"
	resourceTests := []struct {
		name     string
		resType  pipelinev1.PipelineResourceType
		url      string
		revision string
		repo     string
		sha      string
		wantErr  string
	}{
		{"non-git resource", pipelinev1.PipelineResourceTypeImage, "", "", "", "", "non-git resource"},
		{"git resource with no url", pipelinev1.PipelineResourceTypeGit, "", "master", "", "", "failed to find param url"},
		{"git resource with no revision", pipelinev1.PipelineResourceTypeGit, repoURL, "", "", "", "failed to find param revision"},
		{"git resource", pipelinev1.PipelineResourceTypeGit, repoURL, "master", "test/repo", "master", ""},
	}

	for _, tt := range resourceTests {
		res := makePipelineResource(tt.resType, tt.url, tt.revision)

		repo, sha, err := getRepoAndSha(res)
		if !matchError(t, tt.wantErr, err) {
			t.Errorf("getRepoAndSha() %s: got error %v, want %s", tt.name, err, tt.wantErr)
			continue
		}

		if tt.repo != repo {
			t.Errorf("getRepoAndSha() %s: got repo %s, want %s", tt.name, repo, tt.repo)
		}

		if tt.sha != sha {
			t.Errorf("getRepoAndSha() %s: got SHA %s, want %s", tt.name, sha, tt.sha)
		}
	}
}

func TestExtractRepoFromGitHubURL(t *testing.T) {
	repoURLTests := []struct {
		name    string
		url     string
		repo    string
		wantErr string
	}{
		{"standard URL", "https://github.com/tektoncd/triggers", "tektoncd/triggers", ""},
		{"invalid URL", "http://192.168.0.%31/test/repo", "", "failed to parse repo URL.*invalid URL escape"},
		{"url with no repo path", "https://github.com/", "", "could not determine repo from URL"},
	}

	for _, tt := range repoURLTests {
		repo, err := extractRepoFromGitHubURL(tt.url)
		if !matchError(t, tt.wantErr, err) {
			t.Errorf("extractRepoFromGitHubURL() %s: got error %v, want %s", tt.name, err, tt.wantErr)
			continue
		}

		if tt.repo != repo {
			t.Errorf("getRepoAndSha() %s: got repo %s, want %s", tt.name, repo, tt.repo)
		}
	}
}

func makePipelineRunWithResources(opts ...tb.PipelineRunSpecOp) *pipelinev1.PipelineRun {
	return tb.PipelineRun(pipelineRunName, testNamespace, tb.PipelineRunSpec(
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

func makePipelineResource(resType pipelinev1.PipelineResourceType, url, rev string) *pipelinev1.PipelineResourceSpec {
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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomSuffix() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
