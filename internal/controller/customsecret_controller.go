/*
Copyright 2025.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/mithRabbit/customsecret/api/v1alpha1"
)

// CustomSecretReconciler reconciles a CustomSecret object
type CustomSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=api.example.com,resources=customsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.example.com,resources=customsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.example.com,resources=customsecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CustomSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.2/pkg/reconcile
func (r *CustomSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	custom := &apiv1alpha1.CustomSecret{}
	if err := r.Get(ctx, req.NamespacedName, custom); err != nil {
		l.Error(err, "unable to fetch CustomSecret")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	l.Info("reconciling CustomSecret", "name", custom.Name, "type", custom.Spec.Type, "username", custom.Spec.Username, "passwordLen", custom.Spec.PasswordLen, "rotationPeriod", custom.Spec.RotationPeriod)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.CustomSecret{}).
		Named("customsecret").
		Complete(r)
}
