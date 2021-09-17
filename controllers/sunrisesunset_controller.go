/*
Copyright 2021.

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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	harleyb123v1alpha1 "github.com/HarleyB123/harleykubeoperator/api/v1alpha1"
)

// SunriseSunsetReconciler reconciles a SunriseSunset object
type SunriseSunsetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=harleyb123.com,resources=sunrisesunsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=harleyb123.com,resources=sunrisesunsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=harleyb123.com,resources=sunrisesunsets/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *SunriseSunsetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	v := &harleyb123v1alpha1.SunriseSunset{}
	err := r.Get(ctx, req.NamespacedName, v)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new deployment.
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: v.Name, Namespace: v.Namespace}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define and create a new deployment.
			dep := r.deploymentForSunrise(v)
			if err = r.Create(ctx, dep); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	// Ensure the deployment size is the same as the spec.
	size := v.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		if err = r.Update(ctx, found); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, nil
}

// deploymentForMemcached returns a Deployment object for data from m.
func (r *SunriseSunsetReconciler) deploymentForSunrise(m *harleyb123v1alpha1.SunriseSunset) *appsv1.Deployment {
	lbls := labelsForApp(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: lbls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: lbls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "harleyb123/sunrisesunset",
						Name:  "sunrisesunset",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5000,
							Name:          "sunset",
						}},
					}},
				},
			},
		},
	}

	// Set Memcached instance as the owner and controller.memcac
	// NOTE: calling SetControllerReference, and setting owner references in
	// general, is important as it allows deleted objects to be garbage collected.
	controllerutil.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForApp creates a simple set of labels for Memcached.
func labelsForApp(name string) map[string]string {
	return map[string]string{"cr_name": name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SunriseSunsetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&harleyb123v1alpha1.SunriseSunset{}).
		Complete(r)
}
