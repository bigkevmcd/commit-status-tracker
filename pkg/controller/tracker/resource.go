package tracker

import (
	"fmt"
	"net/url"
	"strings"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

// FindGitResource locates a Git PipelineResource in a PipelineRun.
//
// If no Git resources are found, an error is returned.
// If more than one Git resource is found, an error is returned.
func FindGitResource(p *pipelinev1.PipelineRun) (*pipelinev1.PipelineResourceSpec, error) {
	var spec *pipelinev1.PipelineResourceSpec
	for _, r := range p.Spec.Resources {
		if r.ResourceSpec == nil {
			continue
		}
		if r.ResourceSpec.Type == pipelinev1.PipelineResourceTypeGit {
			if spec != nil {
				return nil, fmt.Errorf("found multiple git PipelineResources in the PipelineRun %s", p.ObjectMeta.Name)
			}
			spec = r.ResourceSpec
		}
	}
	if spec == nil {
		return nil, fmt.Errorf("failed to find a git PipelineResource in the PipelineRun %s", p.ObjectMeta.Name)
	}

	return spec, nil
}

// TODO This only parses GitHub repo paths, would need work to parse GitLab repo
// paths too (can have more components).
func GetRepoAndSHA(p *pipelinev1.PipelineResourceSpec) (string, string, error) {
	if p.Type != pipelinev1.PipelineResourceTypeGit {
		return "", "", fmt.Errorf("failed to get repo and SHA from non-git resource: %s", p)
	}
	u, err := getResourceParamByName(p.Params, "url")
	if err != nil {
		return "", "", fmt.Errorf("failed to find param url in GetRepoAndSHA: %w", err)
	}

	rev, err := getResourceParamByName(p.Params, "revision")
	if err != nil {
		return "", "", fmt.Errorf("failed to find param revision in GetRepoAndSHA: %w", err)
	}
	repo, err := extractRepoFromGitHubURL(u)
	if err != nil {
		return "", "", fmt.Errorf("GetRepoAndSHA failed: %w", err)
	}

	return strings.TrimSuffix(repo, ".git"), rev, nil
}

func getResourceParamByName(params []pipelinev1.ResourceParam, name string) (string, error) {
	for _, p := range params {
		if p.Name == name {
			return p.Value, nil
		}
	}
	return "", fmt.Errorf("no resource parameter with name %s", name)
}

func extractRepoFromGitHubURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("failed to parse repo URL %s: %w", s, err)
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("could not determine repo from URL: %v", u)
	}
	return fmt.Sprintf("%s/%s", parts[1], parts[2]), nil
}
