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
	"crypto/rand"
	"encoding/hex"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	if custom.Spec.Type == "basic-auth" && custom.Spec.Username == "admin" && custom.Spec.PasswordLen == 40 && custom.Spec.RotationPeriod > 0 {
		l.Info("reconciling CustomSecret: CustomSecret", "name", custom.Name, "type", custom.Spec.Type, "username", custom.Spec.Username, "passwordLen", custom.Spec.PasswordLen, "rotationPeriod", custom.Spec.RotationPeriod)

		// Generate a random password
		password, err := generateRandomPassword(custom.Spec.PasswordLen)
		if err != nil {
			l.Error(err, "unable to generate random password")
			return ctrl.Result{}, err
		}

		exisitngSecret := &corev1.Secret{}
		if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: custom.Name}, exisitngSecret); err != nil {
			if err != nil && client.IgnoreNotFound(err) != nil {
				l.Error(err, "unable to fetch Secret")
				return ctrl.Result{}, err
			}

			// Create or update the secret
			if err == nil {
				// Secret exists, check rotation period
				lastRotationTime := custom.Status.LastRotationTime.Time
				if time.Since(lastRotationTime) > time.Duration(custom.Spec.RotationPeriod)*time.Second {
					exisitngSecret.StringData["password"] = password
					if err := r.Update(ctx, exisitngSecret); err != nil {
						l.Error(err, "unable to update exisitngSecret")
						return ctrl.Result{}, err
					}
					l.Info("Secret updated", "name", custom.Name)

					custom.Status.LastRotationTime = metav1.Now()
					if err := r.Status().Update(ctx, custom); err != nil {
						l.Error(err, "unable to update CustomSecret status")
						return ctrl.Result{}, err
					}
					l.Info("CustomSecret status updated", "name", custom.Name)
				} else {
					l.Info("Secret not rotated: rotation period not over", "name", custom.Name)
				}
			} else { // secret does not exist
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      custom.Name,
						Namespace: custom.Namespace,
					},
					Type: corev1.SecretTypeBasicAuth,
					StringData: map[string]string{
						"username": custom.Spec.Username,
						"password": password,
					},
				}

				if err := r.Create(ctx, secret); err != nil {
					l.Error(err, "unable to create Secret")
					return ctrl.Result{}, err
				}
				l.Info("Secret created", "name", custom.Name)

				custom.Status.LastRotationTime = metav1.Now()
				if err := r.Status().Update(ctx, custom); err != nil {
					l.Error(err, "unable to update CustomSecret status")
					return ctrl.Result{}, err
				}
				l.Info("CustomSecret status updated", "name", custom.Name)
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.CustomSecret{}).
		Named("customsecret").
		Complete(r)
}

func generateRandomPassword(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
