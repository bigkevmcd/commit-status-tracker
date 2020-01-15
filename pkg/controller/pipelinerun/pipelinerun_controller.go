package pipelinerun

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	pipelinesv1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

var log = logf.Log.WithName("controller_pipelinerun")

// Add creates a new PipelineRun Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePipelineRun{client: mgr.GetClient(), scheme: mgr.GetScheme(), scmFactory: createClient}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("pipelinerun-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &pipelinesv1alpha1.PipelineRun{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcilePipelineRun implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePipelineRun{}

// ReconcilePipelineRun reconciles a PipelineRun object
type ReconcilePipelineRun struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	scmFactory scmClientFactory
}

// Reconcile reads that state of the cluster for a PipelineRun object and makes changes based on the state read
// and what is in the PipelineRun.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePipelineRun) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PipelineRun")

	// Fetch the PipelineRun instance
	pipelineRun := &pipelinesv1alpha1.PipelineRun{}
	err := r.client.Get(context.TODO(), request.NamespacedName, pipelineRun)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !isNotifiablePipelineRun(pipelineRun) {
		reqLogger.Info("not a notifable pipeline run")
		return reconcile.Result{}, nil
	}

	status := getPipelineRunState(pipelineRun)
	if status == Pending {
		reqLogger.Info("pipelineRun still pending")
		return reconcile.Result{}, nil
	}

	res, err := findGitResource(pipelineRun)
	if err != nil {
		reqLogger.Error(err, "failed to find a git resource")
		return reconcile.Result{}, nil
	} else {
		reqLogger.Info("found a git resource", "resource", res)
	}

	repo, sha, err := getRepoAndSHA(res)
	if err != nil {
		reqLogger.Error(err, "failed to parse the URL and SHA correctly")
		return reconcile.Result{}, nil
	} else {
		reqLogger.Info("found a git resource", "resource", res)
	}

	// TODO: Create a GitHub status.
	secret, err := getAuthSecret(r.client, request.Namespace)
	if err != nil {
		reqLogger.Error(err, "failed to get an authSecret")
		return reconcile.Result{}, nil
	}

	client := r.scmFactory(secret)
	commitStatusInput := getCommitStatusInput(pipelineRun)
	reqLogger.Info("creating a github status for", "resource", res, "status", commitStatusInput)
	s, _, err := client.Repositories.CreateStatus(context.Background(), repo, sha, commitStatusInput)
	reqLogger.Info("created a github status", "status", s)
	return reconcile.Result{}, nil
}
