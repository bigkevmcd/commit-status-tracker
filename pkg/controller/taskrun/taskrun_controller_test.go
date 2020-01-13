package taskrun

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

var (
	namespace   = "test-namespace"
	taskRunName = "test-task-run"
)

// TestTaskRunController runs ReconcileTaskRun.Reconcile() against a
// fake client that tracks a Memcached object.
func TestTaskRunController(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	taskRun := &pipelinev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskRunName,
			Namespace: namespace,
		},
		Spec: pipelinev1.TaskRunSpec{
			Inputs: pipelinev1.TaskRunInputs{
				Resources: []pipelinev1.TaskResourceBinding{
					pipelinev1.TaskResourceBinding{
						PipelineResourceBinding: pipelinev1.PipelineResourceBinding{
							Name: "source",
							ResourceSpec: &pipelinev1.PipelineResourceSpec{
								Type: "git",
								Params: []pipelinev1.ResourceParam{
									pipelinev1.ResourceParam{
										Name:  "revision",
										Value: "master",
									},
									pipelinev1.ResourceParam{
										Name:  "url",
										Value: "https://github.com/GoogleContainerTools/skaffold",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	objs := []runtime.Object{
		taskRun,
	}

	s := scheme.Scheme
	s.AddKnownTypes(pipelinev1.SchemeGroupVersion, taskRun)
	cl := fake.NewFakeClient(objs...)
	r := &ReconcileTaskRun{client: cl, scheme: s}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      taskRunName,
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
