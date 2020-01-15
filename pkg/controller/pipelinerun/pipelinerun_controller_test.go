package pipelinerun

import (
	"reflect"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	fakescm "github.com/jenkins-x/go-scm/scm/driver/fake"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
)

var (
	testNamespace   = "test-namespace"
	pipelineRunName = "test-pipeline-run"
	testToken       = "abcdefghijklmnopqrstuvwxyz12345678901234"
)

// TestPipelineRunController runs ReconcilePipelineRun.Reconcile() against a
// fake client that tracks PipelineRun objects.
func TestPipelineRunController(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := makePipelineRunWithResources(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(notifiableName, "true"),
		tb.PipelineRunAnnotation(statusContextName, "test-context"),
		tb.PipelineRunAnnotation(statusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue})))
	objs := []runtime.Object{
		pipelineRun,
		makeSecret(map[string][]byte{"token": []byte(testToken)}),
	}

	s := scheme.Scheme
	s.AddKnownTypes(pipelinev1.SchemeGroupVersion, pipelineRun)
	cl := fake.NewFakeClient(objs...)
	client, data := fakescm.NewDefault()
	fakeClientFactory := func(s string) *scm.Client {
		return client
	}
	r := &ReconcilePipelineRun{client: cl, scheme: s, scmFactory: fakeClientFactory}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	wanted := &scm.Status{State: scm.StateSuccess, Label: "test-context", Desc: "testing", Target: ""}
	status := data.Statuses["master"][0]
	if !reflect.DeepEqual(status, wanted) {
		t.Fatalf("commit-status notification got %#v, wanted %#v\n", status, wanted)
	}
}

// TestPipelineRunReconcileWithNoGitCredentials tests a non-notifable
// PipelineRun.
func TestPipelineRunReconcileNonNotifiable(t *testing.T) {
	t.Skip()
}

// TestPipelineRunReconcileWithNoGitCredentials tests a notifable PipelineRun
// with no "git" resource.
func TestPipelineRunReconcileWithNoGitRepository(t *testing.T) {
	t.Skip()
}

// TestPipelineRunReconcileWithNoGitCredentials tests a notifable PipelineRun
// with multiple "git" resources.
func TestPipelineRunReconcileWithGitRepositories(t *testing.T) {
	t.Skip()
}

// TestPipelineRunReconcileWithNoGitCredentials tests a notifable PipelineRun
// with a "git" resource, but with no Git credentials.
func TestPipelineRunReconcileWithNoGitCredentials(t *testing.T) {
	t.Skip()
}

func applyOpts(pr *pipelinev1.PipelineRun, opts ...tb.PipelineRunOp) {
	for _, o := range opts {
		o(pr)
	}
}
