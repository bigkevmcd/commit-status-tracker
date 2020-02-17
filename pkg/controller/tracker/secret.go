package tracker

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// TODO: what should these be called?
	secretName = "commit-status-tracker-git-secret"
	secretID   = "token"
)

// GetAuthSecret attempts to find a Secret in the provided namespace, using the
// client.
//
// Returns the string of the secret if found, otherwise returns an error.
func GetAuthSecret(c client.Client, ns string) (string, error) {
	secret := &corev1.Secret{}
	err := c.Get(context.TODO(), getNamespaceSecretName(ns), secret)
	if err != nil {
		return "", fmt.Errorf("failed to GetAuthSecret, error getting secret '%s' in namespace '%s': '%q'", secretName, ns, err)
	}

	tokenData, ok := secret.Data[secretID]
	if !ok {
		return "", fmt.Errorf("failed to GetAuthSecret, secret %s does not have a 'token' key", ns)
	}
	return string(tokenData), nil
}

func getNamespaceSecretName(s string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: s,
		Name:      secretName,
	}

}
