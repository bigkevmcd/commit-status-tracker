package tracker

import (
	"reflect"
	"testing"

	"github.com/bigkevmcd/commit-status-tracker/test"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

func TestFindGitResourceWithNoRepository(t *testing.T) {
	pipelineRun := test.MakePipelineRunWithResources()

	_, err := FindGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestFindGitResourceWithRepository(t *testing.T) {
	pipelineRun := test.MakePipelineRunWithResources(
		test.MakeGitResourceBinding("https://github.com/tektoncd/triggers", "master"))

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

	r, err := FindGitResource(pipelineRun)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, want) {
		t.Fatalf("got %+v, want %+v", r, want)
	}
}

func TestFindGitResourceWithMultipleRepositories(t *testing.T) {
	pipelineRun := test.MakePipelineRunWithResources(
		test.MakeGitResourceBinding("https://github.com/tektoncd/triggers", "master"),
		test.MakeGitResourceBinding("https://github.com/tektoncd/pipeline", "master"))

	_, err := FindGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestFindGitResourceWithNonGitResource(t *testing.T) {
	pipelineRun := test.MakePipelineRunWithResources(
		test.MakeImageResourceBinding("example.com/project/myimage"))

	_, err := FindGitResource(pipelineRun)
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
		{"git resource with .git", pipelinev1.PipelineResourceTypeGit, repoURL + ".git", "master", "test/repo", "master", ""},
	}

	for _, tt := range resourceTests {
		res := test.MakePipelineResource(tt.resType, tt.url, tt.revision)

		repo, sha, err := GetRepoAndSHA(res)
		if !matchError(t, tt.wantErr, err) {
			t.Errorf("GetRepoAndSHA() %s: got error %v, want %s", tt.name, err, tt.wantErr)
			continue
		}

		if tt.repo != repo {
			t.Errorf("GetRepoAndSHA() %s: got repo %s, want %s", tt.name, repo, tt.repo)
		}

		if tt.sha != sha {
			t.Errorf("GetRepoAndSHA() %s: got SHA %s, want %s", tt.name, sha, tt.sha)
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
			t.Errorf("GetRepoAndSHA() %s: got repo %s, want %s", tt.name, repo, tt.repo)
		}
	}
}
