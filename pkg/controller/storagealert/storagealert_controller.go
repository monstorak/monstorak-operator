package storagealert

import (
	"context"

	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monitoring "github.com/coreos/prometheus-operator/pkg/client/versioned"
	alertsv1alpha1 "github.com/monstorak/monstorak/pkg/apis/alerts/v1alpha1"
	manifests "github.com/monstorak/monstorak/pkg/manifests"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_storagealert")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

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

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner StorageAlerts
	// err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &alertsv1alpha1.StorageAlerts{},
	// })
	// if err != nil {
	// 	return err
	// }

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
func (r *ReconcileStorageAlert) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling StorageAlert")

	// Fetch the StorageAlert instance
	instance := &alertsv1alpha1.StorageAlert{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	f := manifests.NewFactory(instance.Namespace)
	prometheusK8sRules, err := f.PrometheusK8sRules()
	prometheusK8sRules.Namespace = instance.Spec.StorageAlert.PrometheusNamespace

	// Set StorageAlerts instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, prometheusK8sRules, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	CreateOrUpdatePrometheusRule(prometheusK8sRules)

	return reconcile.Result{}, nil
}

// CreateOrUpdatePrometheusRule Creates or Updates prometheusRule object
func CreateOrUpdatePrometheusRule(p *monv1.PrometheusRule) {
	reqLogger := log.WithValues("Prometheus Namespace: ", p.ObjectMeta.Namespace)
	cfg, err := config.GetConfig()
	mclient, err := monitoring.NewForConfig(cfg)
	pclient := mclient.MonitoringV1().PrometheusRules(p.GetNamespace())
	oldRule, err := pclient.Get(p.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err := pclient.Create(p)
		if err != nil {
			reqLogger.Error(err, "creating PrometheusRule object failed", "Prometheus Namespace: ", p.ObjectMeta.Namespace)
			return
		}
		reqLogger.Info("PrometheusRule Created.", "Prometheus Namespace: ", p.ObjectMeta.Namespace)
		return
	}
	if err != nil {
		reqLogger.Error(err, "retrieving PrometheusRule object failed", "Prometheus Namespace: ", p.ObjectMeta.Namespace)
		return
	}

	p.ResourceVersion = oldRule.ResourceVersion

	_, err = pclient.Update(p)
	if err != nil {
		reqLogger.Error(err, "updating PrometheusRule object failed", "Prometheus Namespace: ", p.ObjectMeta.Namespace)
		return
	} else {
		reqLogger.Info("PrometheusRule Updated.", "Prometheus Namespace: ", p.ObjectMeta.Namespace)
		return
	}
}
