package pipelinerun

// TODO: Determine a base domain.
const (
	notifiableName      = "app.example.com/git-status"
	statusContextName   = "app.example.com/status-context"
	statusTargetURLName = "app.example.com/status-target-url"

	// TODO: This could also come from a ConfigMap based on the context.
	statusDescriptionName = "app.example.com/status-description"
)
