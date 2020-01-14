package pipelinerun

import (
	"testing"

	tb "github.com/tektoncd/pipeline/test/builder"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

var (
	namespace       = "test-namespace"
	pipelineRunName = "test-pipeline-run"
)

// TestPipelineRunController runs ReconcilePipelineRun.Reconcile() against a
// fake client that tracks PipelineRun objects.
func TestPipelineRunController(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := makePipelineRunWithResources(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"))
	applyOpts(pipelineRun, tb.PipelineRunLabel(notifiable, "true"))
	objs := []runtime.Object{
		pipelineRun,
	}

	s := scheme.Scheme
	s.AddKnownTypes(pipelinev1.SchemeGroupVersion, pipelineRun)
	cl := fake.NewFakeClient(objs...)
	r := &ReconcilePipelineRun{client: cl, scheme: s}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: namespace,
		},
	}
	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Fatal("reconcile requeued request")
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
