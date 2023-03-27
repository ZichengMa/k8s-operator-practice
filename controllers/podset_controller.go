/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "tutorial.kubebuilder.io/operator-practice/api/v1"
)

// PodSetReconciler reconciles a PodSet object
type PodSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=podsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=podsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=podsets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PodSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	// listen to the creation of PodSet objects

	// fetch the PodSet object
	podset := &batchv1.PodSet{}
	err := r.Get(context.TODO(), req.NamespacedName, podset) // get the PodSet object from the API server by its name and namespace and store it in the podset variable
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		// If the error is not a NotFound error, the second if statement is triggered, and the controller
		// returns a ctrl.Result{} along with the original error to requeue the request for reconciliation.
		// This is done to allow the controller to retry the operation at a later time, in case the error was temporary or transient.
		return ctrl.Result{}, err
	}

	podList := &corev1.PodList{} // create a new list of pods not podsets!
	// The lbs variable is a map of labels that will be used to filter the list of pods. In this case, the labels are set to app=podset.Name and version=v0.1.
	// This means that the query will only return pods that have the app label set to the name of the podset object and the version label set to v0.1.
	lbs := map[string]string{
		"app":     podset.Name,
		"version": "v0.1",
	}
	labelSelector := labels.SelectorFromSet(lbs) // create a label selector from the lbs map
	// the listOps variable is a client.ListOptions object that contains the options for querying the Kubernetes API server.
	listOps := &client.ListOptions{Namespace: podset.Namespace, LabelSelector: labelSelector} // create a list options object with the namespace and label selector
	if err = r.List(context.TODO(), podList, listOps); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.PodSet{}).
		Complete(r)
}
