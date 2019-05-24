package storagealert

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/monstorak/monstorak/pkg/tasks"

	alertsv1alpha1 "github.com/monstorak/monstorak/pkg/apis/alerts/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_storagealert")

// Add creates a new StorageAlert Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStorageAlert{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("storagealert-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource StorageAlert
	err = c.Watch(&source.Kind{Type: &alertsv1alpha1.StorageAlert{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources and requeue the owner StorageAlerts
	err = c.Watch(&source.Kind{Type: &corev1.Namespace{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &corev1.Namespace{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileStorageAlert{}

// ReconcileStorageAlert reconciles a StorageAlert object
type ReconcileStorageAlert struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a StorageAlerts object and makes changes based on the state read
// and what is in the StorageAlerts.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
var (
	reconcilePeriod = 10 * time.Second
	reconcileResult = reconcile.Result{RequeueAfter: reconcilePeriod}
)

const (
	FailedPrerequisites              string = "Some prerequisites are not met"
	FailedRetrievePrometheusRule     string = "Failed to retrieve Prometheus Rules"
	FailedCreateUpdatePrometheusRule string = "Failed to create/update Prometheus Rules"
)

func (r *ReconcileStorageAlert) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling StorageAlerts")

	// Fetch the StorageAlert instance
	instance := &alertsv1alpha1.StorageAlert{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	for index, storageSpec := range instance.Spec.Storage {
		storageNamespace := storageSpec.Namespace
		serviceMonitor := storageSpec.ServiceMonitor
		storageProvider := storageSpec.Provider
		storageVersion := storageSpec.Version
		reqLogger.WithValues("Storage Provider", storageProvider, "Storage Namespace", storageNamespace,
			"Storage Version", storageVersion, "Service Monitor", serviceMonitor)
		reqLogger.Info("Reconciling StorageAlert", "Index", index)
		// Check prerequisites
		err = tasks.Prerequisites(storageNamespace, serviceMonitor)
		if err != nil {
			// Prerequisites not met, requeue
			reqLogger.Error(err, FailedPrerequisites)
			return reconcileResult, err
		}
		// Get prometheusRule
		prometheusRule, err := tasks.GetPrometheusRule(storageNamespace, storageProvider, storageVersion)
		if err != nil {
			// Failed to retrieve prometheusRule, requeue
			reqLogger.Error(err, FailedRetrievePrometheusRule)
			return reconcileResult, err
		}
		// Deploy prometheusRule
		err = tasks.DeployPrometheusRule(storageNamespace, prometheusRule)
		if err != nil {
			// Failed to create/update prometheusRule, requeue
			reqLogger.Error(err, FailedCreateUpdatePrometheusRule)
			return reconcileResult, err
		}
	}

	return reconcile.Result{}, nil
}
