package tracker

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

var (
	ErrNoGitResource        = errors.New("failed to find a git resource")
	ErrMultipleGitResources = errors.New("found multiple git resources")
)

// Commit represents the repo/ref that the tracker sends statuses
// notifications for.
type Commit struct {
	RepoURL string
	Ref     string
}

// Repo extracts the "org/repo" from the Commit's RepoURL.
func (c Commit) Repo() (string, error) {
	return extractRepoFromGitHubURL(c.RepoURL)
}

// FindCommit extracts the details of commit/ref from a "git" PipelineResource.
//
// An error is returned if:
//
//   no "git" resource is found,
//   multiple "git" resources are found
//   the found "git" resource has no url or revision
func FindCommit(res []*pipelinev1.PipelineResourceSpec) (*Commit, error) {
	gits := make([]*pipelinev1.PipelineResourceSpec, 0)
	for _, r := range res {
		if r.Type == pipelinev1.PipelineResourceTypeGit {
			gits = append(gits, r)
		}
	}
	if len(gits) == 0 {
		return nil, ErrNoGitResource
	}
	if len(gits) > 1 {
		return nil, ErrMultipleGitResources
	}
	found := gits[0]
	u, err := getResourceParamByName(found.Params, "url")
	if err != nil {
		return nil, fmt.Errorf("failed to find param url in FindCommit: %w", err)
	}
	rev, err := getResourceParamByName(found.Params, "revision")
	if err != nil {
		return nil, fmt.Errorf("failed to find param revision in FindCommit: %w", err)
	}
	return &Commit{RepoURL: u, Ref: rev}, nil
}

func getResourceParamByName(params []pipelinev1.ResourceParam, name string) (string, error) {
	for _, p := range params {
		if p.Name == name {
			return p.Value, nil
		}
	}
	return "", fmt.Errorf("no resource parameter with name %s", name)
}

// TODO This only parses GitHub repo paths, would need work to parse GitLab repo
// paths too (can have more components).
func extractRepoFromGitHubURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("failed to parse repo URL %s: %w", s, err)
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("could not determine repo from URL: %v", u)
	}
	repo := strings.TrimSuffix(parts[2], ".git")
	return fmt.Sprintf("%s/%s", parts[1], repo), nil
}
