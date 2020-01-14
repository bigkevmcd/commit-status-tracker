package pipelinerun

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// TODO: what should these be called?
	secretName = "github-auth"
	secretID   = "token"
)

func getAuthSecret(c client.Client, ns types.NamespacedName) (string, error) {
	secret := &corev1.Secret{}
	err := c.Get(context.TODO(), ns, secret)
	if err != nil {
		return "", fmt.Errorf("failed to getAuthSecret, error getting secret %s: '%q'", ns, err)
	}

	tokenData, ok := secret.Data[secretID]
	if !ok {
		return "", fmt.Errorf("failed to getAuthSecret, secret %s does not have a 'token' key", ns)
	}
	return string(tokenData), nil
}

func getNamespaceSecretName(s string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: s,
		Name:      secretName,
	}

}
