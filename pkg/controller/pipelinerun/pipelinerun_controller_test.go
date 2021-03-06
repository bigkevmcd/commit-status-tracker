package pipelinerun

import (
	"reflect"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	fakescm "github.com/jenkins-x/go-scm/scm/driver/fake"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/bigkevmcd/commit-status-tracker/pkg/tracker"
	ctb "github.com/bigkevmcd/commit-status-tracker/test/builder"
)

var (
	testNamespace   = "test-namespace"
	pipelineRunName = "test-pipeline-run"
	testToken       = "abcdefghijklmnopqrstuvwxyz12345678901234"
)

var _ reconcile.Reconciler = &ReconcilePipelineRun{}

// TestPipelineRunControllerPendingState runs ReconcilePipelineRun.Reconcile() against a
// fake client that tracks PipelineRun objects.
func TestPipelineRunControllerPendingState(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(tracker.NotifiableName, "true"),
		tb.PipelineRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.PipelineRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		pipelineRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	wanted := &scm.Status{State: scm.StatePending, Label: "test-context", Desc: "testing", Target: ""}
	status := data.Statuses["master"][0]
	if !reflect.DeepEqual(status, wanted) {
		t.Fatalf("commit-status notification got %#v, wanted %#v\n", status, wanted)
	}
}

// TestPipelineRunReconcileWithPreviousPending tests a PipelineRun that
// we've already sent a pending notification.
func TestPipelineRunReconcileWithPreviousPending(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(tracker.NotifiableName, "true"),
		tb.PipelineRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.PipelineRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		pipelineRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}

	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	// This runs Reconcile twice.
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	// This cleans out the existing date for the data, because the fake scm
	// client updates in-place, so there's no way to know if it received multiple
	// pending notifications.
	delete(data.Statuses, "master")
	res, err = r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	// There should be no recorded statuses, because the state is still pending
	// and the fake client's state was deleted above.
	assertNoStatusesRecorded(t, data)
}

// TestPipelineRunControllerSuccessState runs ReconcilePipelineRun.Reconcile() against a
// fake client that tracks PipelineRun objects.
func TestPipelineRunControllerSuccessState(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(tracker.NotifiableName, "true"),
		tb.PipelineRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.PipelineRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue})))
	objs := []runtime.Object{
		pipelineRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	wanted := &scm.Status{State: scm.StateSuccess, Label: "test-context", Desc: "testing", Target: ""}
	status := data.Statuses["master"][0]
	if !reflect.DeepEqual(status, wanted) {
		t.Fatalf("commit-status notification got %#v, wanted %#v\n", status, wanted)
	}
}

// TestPipelineRunControllerFailedState runs ReconcilePipelineRun.Reconcile() against a
// fake client that tracks PipelineRun objects.
func TestPipelineRunControllerFailedState(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(tracker.NotifiableName, "true"),
		tb.PipelineRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.PipelineRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionFalse})))
	objs := []runtime.Object{
		pipelineRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	wanted := &scm.Status{State: scm.StateFailure, Label: "test-context", Desc: "testing", Target: ""}
	status := data.Statuses["master"][0]
	if !reflect.DeepEqual(status, wanted) {
		t.Fatalf("commit-status notification got %#v, wanted %#v\n", status, wanted)
	}
}

// TestPipelineRunReconcileWithNoGitCredentials tests a non-notifable
// PipelineRun.
func TestPipelineRunReconcileNonNotifiable(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		pipelineRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	assertNoStatusesRecorded(t, data)
}

// TestPipelineRunReconcileWithNoGitCredentials tests a notifable PipelineRun
// with no "git" resource.
func TestPipelineRunReconcileWithNoGitRepository(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources()
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(tracker.NotifiableName, "true"),
		tb.PipelineRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.PipelineRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		pipelineRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	assertNoStatusesRecorded(t, data)
}

// TestPipelineRunReconcileWithNoGitCredentials tests a notifable PipelineRun
// with multiple "git" resources.
func TestPipelineRunReconcileWithGitRepositories(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"),
		ctb.MakeGitResource("https://github.com/tektoncd/pipeline", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(tracker.NotifiableName, "true"),
		tb.PipelineRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.PipelineRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		pipelineRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	assertNoStatusesRecorded(t, data)
}

// TestPipelineRunReconcileWithNoGitCredentials tests a notifable PipelineRun
// with a "git" resource, but with no Git credentials.
func TestPipelineRunReconcileWithNoGitCredentials(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"),
		ctb.MakeGitResource("https://github.com/tektoncd/pipeline", "master"))
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(tracker.NotifiableName, "true"),
		tb.PipelineRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.PipelineRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{pipelineRun}
	r, data := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	assertNoStatusesRecorded(t, data)

}

func TestKeyForCommit(t *testing.T) {
	inputTests := []struct {
		repo string
		ref  string
		want string
	}{
		{"tekton/triggers", "e1466db56110fa1b813277c1647e20283d3370c3",
			"7b2841ab8791fece7acdc0b3bb6e398c7a184273"},
	}

	for _, tt := range inputTests {
		if v := keyForCommit(tt.repo, tt.ref); v != tt.want {
			t.Errorf("keyForCommit(%#v, %#v) got %#v, want %#v", tt.repo, tt.ref, v, tt.want)
		}
	}
}

func applyOpts(pr *pipelinev1.PipelineRun, opts ...tb.PipelineRunOp) {
	for _, o := range opts {
		o(pr)
	}
}

func makeReconciler(pr *pipelinev1.PipelineRun, objs ...runtime.Object) (*ReconcilePipelineRun, *fakescm.Data) {
	s := scheme.Scheme
	s.AddKnownTypes(pipelinev1.SchemeGroupVersion, pr)
	cl := fake.NewFakeClient(objs...)
	client, data := fakescm.NewDefault()
	fakeClientFactory := func(s string) *scm.Client {
		return client
	}
	return &ReconcilePipelineRun{
		client:       cl,
		scheme:       s,
		scmFactory:   fakeClientFactory,
		pipelineRuns: make(pipelineRunTracker),
	}, data
}

func fatalIfError(t *testing.T, err error, format string, a ...interface{}) {
	if err != nil {
		t.Fatalf(format, a...)
	}
}

func assertNoStatusesRecorded(t *testing.T, d *fakescm.Data) {
	if l := len(d.Statuses["master"]); l != 0 {
		t.Fatalf("too many statuses recorded, got %v, wanted 0", l)
	}
}
