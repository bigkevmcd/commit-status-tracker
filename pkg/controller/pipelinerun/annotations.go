package pipelinerun

const (
	notifiableName      = "tekton.dev/git-status"
	statusContextName   = "tekton.dev/status-context"
	statusTargetURLName = "tekton.dev/status-target-url"

	// TODO: This could also come from a ConfigMap based on the context.
	statusDescriptionName = "tekton.dev/status-description"
)
