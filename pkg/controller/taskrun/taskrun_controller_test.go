package taskrun

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

	"github.com/bigkevmcd/commit-status-tracker/pkg/tracker"
	"github.com/bigkevmcd/commit-status-tracker/test"
	ctb "github.com/bigkevmcd/commit-status-tracker/test/builder"
)

var (
	testNamespace   = "test-namespace"
	pipelineRunName = "test-task-run"
	testToken       = "abcdefghijklmnopqrstuvwxyz12345678901234"
)

var _ reconcile.Reconciler = &ReconcileTaskRun{}

// TestTaskRunControllerPendingState runs ReconcileTaskRun.Reconcile() against a
// fake client that tracks TaskRun objects.
func TestTaskRunControllerPendingState(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		taskRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// TestTaskRunReconcileWithPreviousPending tests a TaskRun that
// we've already sent a pending notification.
func TestTaskRunReconcileWithPreviousPending(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		taskRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}

	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// TestTaskRunControllerSuccessState runs ReconcileTaskRun.Reconcile() against a
// fake client that tracks TaskRun objects.
func TestTaskRunControllerSuccessState(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue})))
	objs := []runtime.Object{
		taskRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// TestTaskRunControllerFailedState runs ReconcileTaskRun.Reconcile() against a
// fake client that tracks TaskRun objects.
func TestTaskRunControllerFailedState(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionFalse})))
	objs := []runtime.Object{
		taskRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// TestTaskRunReconcileWithNoGitCredentials tests a non-notifable
// TaskRun.
func TestTaskRunReconcileNonNotifiable(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionFalse})))
	objs := []runtime.Object{
		taskRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// TestTaskRunReconcileWithNoGitCredentials tests a notifable TaskRun
// with no "git" resource.
func TestTaskRunReconcileWithNoGitRepository(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources()
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		taskRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// TestTaskRunReconcileWithNoGitCredentials tests a notifable TaskRun
// with multiple "git" resources.
func TestTaskRunReconcileWithGitRepositories(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"),
		ctb.MakeGitResource("https://github.com/tektoncd/task", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{
		taskRun,
		ctb.MakeSecret(tracker.SecretName, map[string][]byte{"token": []byte(testToken)}),
	}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// TestTaskRunReconcileWithNoGitCredentials tests a notifable TaskRun
// with a "git" resource, but with no Git credentials.
func TestTaskRunReconcileWithNoGitCredentials(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"),
		ctb.MakeGitResource("https://github.com/tektoncd/task", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{taskRun}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
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

// If the TaskRun can't be loaded then this isn't an error.
func TestTaskRunControllerMissingTaskRun(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
	r, data := makeReconciler(taskRun)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "unknown-task-run",
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
	_, ok := data.Statuses["master"]
	if ok {
		t.Fatal("status incorrectly queued for unknown pipeline")
	}
}

// If the TaskRun has a bad git repository then this should fail.
func TestTaskRunControllerBadGitRepo(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	taskRun := ctb.MakeTaskRunWithInputResources(
		ctb.MakeGitResource("http://192.168.0.%31/test/repo", "master"))
	applyOpts(
		taskRun,
		tb.TaskRunAnnotation(tracker.NotifiableName, "true"),
		tb.TaskRunAnnotation(tracker.StatusContextName, "test-context"),
		tb.TaskRunAnnotation(tracker.StatusDescriptionName, "testing"),
		tb.TaskRunStatus(
			tb.StatusCondition(
				apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
	objs := []runtime.Object{taskRun}
	r, data := makeReconciler(taskRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRun.Name,
			Namespace: testNamespace,
		},
	}
	_, err := r.Reconcile(req)
	if !test.MatchError(t, "failed to parse repo URL", err) {
		t.Errorf("unexpected error returned: %s", err)
	}

	_, ok := data.Statuses["master"]
	if ok {
		t.Fatal("status incorrectly queued for unknown pipeline")
	}

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

func applyOpts(pr *pipelinev1.TaskRun, opts ...tb.TaskRunOp) {
	for _, o := range opts {
		o(pr)
	}
}

func makeReconciler(pr *pipelinev1.TaskRun, objs ...runtime.Object) (*ReconcileTaskRun, *fakescm.Data) {
	s := scheme.Scheme
	s.AddKnownTypes(pipelinev1.SchemeGroupVersion, pr)
	cl := fake.NewFakeClient(objs...)
	client, data := fakescm.NewDefault()
	fakeClientFactory := func(s string) *scm.Client {
		return client
	}
	return &ReconcileTaskRun{
		client:     cl,
		scheme:     s,
		scmFactory: fakeClientFactory,
		taskRuns:   make(taskRunTracker),
	}, data
}

func fatalIfError(t *testing.T, err error, format string, a ...interface{}) {
	t.Helper()
	if err != nil {
		t.Fatalf(format, a...)
	}
}

func assertNoStatusesRecorded(t *testing.T, d *fakescm.Data) {
	t.Helper()
	if l := len(d.Statuses["master"]); l != 0 {
		t.Fatalf("too many statuses recorded, got %v, wanted 0", l)
	}
}
