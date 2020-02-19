package tracker

import (
	"reflect"
	"testing"

	"github.com/bigkevmcd/commit-status-tracker/test"
	tb "github.com/bigkevmcd/commit-status-tracker/test/builder"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

const (
	rtImage = pipelinev1.PipelineResourceTypeImage
	rtGit   = pipelinev1.PipelineResourceTypeGit
)

func TestFindCommit(t *testing.T) {
	repoURL := "https://example.com/test/repo.git"
	resourceTests := []struct {
		name    string
		res     []*pipelinev1.PipelineResourceSpec
		want    *Commit
		wantErr string
	}{
		{"non-git resource", specs(spec(rtImage, "", "")), nil, "failed to find a git resource"},
		{"git resource with no url", specs(spec(rtGit, "", "master")), nil, "failed to find param url"},
		{"git resource with no revision", specs(spec(rtGit, repoURL, "")), nil, "failed to find param revision"},
		{"git resource", specs(spec(rtGit, repoURL, "master")), &Commit{repoURL, "master"}, ""},
		{"specs git resources", specs(spec(rtGit, repoURL, "master"), spec(rtGit, repoURL, "master")), nil, "multiple git resources"},
	}

	for _, tt := range resourceTests {
		t.Run(tt.name, func(t *testing.T) {
			gr, err := FindCommit(tt.res)
			if !test.MatchError(t, tt.wantErr, err) {
				t.Errorf("FindCommit() %s: got error %v, want %s", tt.name, err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.want, gr) {
				t.Errorf("FindCommit() %s: got res  %#v, want %#v", tt.name, gr, tt.want)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	resourceTests := []struct {
		name    string
		repoURL string
		want    string
	}{
		{"url with .git", "https://github.com/test/test.git", "test/test"},
		{"url without .git", "https://github.com/org/repo", "org/repo"},
	}

	for _, tt := range resourceTests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commit{RepoURL: tt.repoURL}
			v, err := c.Repo()
			if err != nil {
				t.Errorf("Repo() got an error: %s", err)
			}
			if v != tt.want {
				t.Errorf("Repo() got %#v, want %#v", v, tt.want)
			}
		})
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
		{"url with .git", "https://github.com/tektoncd/triggers.git", "tektoncd/triggers", ""},
		{"invalid URL", "http://192.168.0.%31/test/repo", "", "failed to parse repo URL.*invalid URL escape"},
		{"url with no repo path", "https://github.com/", "", "could not determine repo from URL"},
	}

	for _, tt := range repoURLTests {
		repo, err := extractRepoFromGitHubURL(tt.url)
		if !test.MatchError(t, tt.wantErr, err) {
			t.Errorf("extractRepoFromGitHubURL() %s: got error %v, want %s", tt.name, err, tt.wantErr)
			continue
		}

		if tt.repo != repo {
			t.Errorf("FindCommit() %s: got repo %s, want %s", tt.name, repo, tt.repo)
		}
	}
}

func specs(v ...*pipelinev1.PipelineResourceSpec) []*pipelinev1.PipelineResourceSpec {
	return v
}
func spec(t pipelinev1.PipelineResourceType, url, rev string) *pipelinev1.PipelineResourceSpec {
	return tb.MakePipelineResource(t, url, rev)
}
