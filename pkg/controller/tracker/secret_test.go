package tracker

import (
	"regexp"
	"testing"

	"github.com/bigkevmcd/commit-status-tracker/test"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	testToken = "abcdefghijklmnopqrstuvwxyz12345678901234"
)

func TestGetAuthSecretWithExistingToken(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	secret := test.MakeSecret(SecretName, map[string][]byte{"token": []byte(testToken)})
	objs := []runtime.Object{
		secret,
	}

	cl := fake.NewFakeClient(objs...)
	sec, err := GetAuthSecret(cl, secret.Namespace)
	if err != nil {
		t.Fatal(err)
	}
	if sec != testToken {
		t.Fatalf("got %s, want %s", sec, testToken)
	}
}

func TestGetAuthSecretWithNoSecret(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	objs := []runtime.Object{}

	cl := fake.NewFakeClient(objs...)
	_, err := GetAuthSecret(cl, "testing")

	wantErr := "error getting secret 'commit-status-tracker-git-secret' in namespace 'testing':.* not found"
	if !matchError(t, wantErr, err) {
		t.Fatalf("failed to match error when no secret: got %s, want %s", err, wantErr)
	}
}

func TestGetAuthSecretWithNoToken(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	secret := test.MakeSecret(SecretName, map[string][]byte{})
	objs := []runtime.Object{
		secret,
	}

	cl := fake.NewFakeClient(objs...)
	_, err := GetAuthSecret(cl, secret.Namespace)

	wantErr := "secret .* does not have a 'token' key"
	if !matchError(t, wantErr, err) {
		t.Fatalf("failed to match error when no secret: got %s, want %s", err, wantErr)
	}
}

func matchError(t *testing.T, s string, e error) bool {
	t.Helper()
	if s == "" && e == nil {
		return true
	}
	if s != "" && e == nil {
		return false
	}
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	return match
}
