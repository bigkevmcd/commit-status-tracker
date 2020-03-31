package tracker

const (
	NotifiableName      = "tekton.dev/git-status"
	StatusContextName   = "tekton.dev/status-context"
	StatusTargetURLName = "tekton.dev/status-target-url"

	// TODO: This could also come from a ConfigMap based on the context.
	StatusDescriptionName = "tekton.dev/status-description"
)
