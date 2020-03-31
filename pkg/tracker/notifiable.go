package tracker

// IsNotifiable returns true if this TaskRun should report its
// completion status as a GitHub status.
func IsNotifiable(ag annotationsGetter) bool {
	for k, v := range ag.Annotations() {
		if k == NotifiableName && v == "true" {
			return true
		}
	}
	return false
}
