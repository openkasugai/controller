/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED

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
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
)

// FunctionTypeReconciler reconciles a FunctionType object
type FunctionTypeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.com,resources=functiontypes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=functiontypes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=functiontypes/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FunctionType object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *FunctionTypeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// get ConfigMap from manifest
	var ft ntthpcv1.FunctionType
	l.Info("fetching FunctionType Resource")
	if err := r.Get(ctx, req.NamespacedName, &ft); err != nil {
		l.Error(err, "unable to fetch FunctionType")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// update the status of function resource
	cm := &corev1.ConfigMap{}
	l.Info(fmt.Sprintf("Finding configmap with name %v in namespace %v",
		ft.Spec.FunctionInfoCMRef.Name, ft.Spec.FunctionInfoCMRef.Namespace))

	if err := r.Get(ctx, client.ObjectKey{
		Name:      ft.Spec.FunctionInfoCMRef.Name,
		Namespace: ft.Spec.FunctionInfoCMRef.Namespace,
	}, cm); err != nil {
		l.Error(err, "unable to fetch configmap")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for k, j := range cm.Data {

		switch k {

		case "deployableItems":
			funcInfoJson := j
			var result []map[string]string
			json.Unmarshal([]byte(funcInfoJson), &result)

			for _, v := range result {
				rt := v["regionType"]
				found := false

				for _, i := range ft.Status.RegionTypeCandidates {
					if i == rt {
						found = true
						break
					}
				}

				if !found && rt != "" {
					ft.Status.RegionTypeCandidates = append(ft.Status.RegionTypeCandidates, rt)
				}
			}

		case "recommend":
			l.Info("recommend key found")

		case "spec":
			l.Info("spec key found")

		default:
			l.Info("unsupported key found", "unsupported key", k)
		}
	}

	if len(ft.Status.RegionTypeCandidates) == 0 {
		ft.Status.Status = "Error"
	} else {
		ft.Status.Status = "Ready"
	}

	if err := r.Status().Update(ctx, &ft); err != nil {
		l.Error(err, "failed to update FunctionType resource status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FunctionTypeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.FunctionType{}).
		Complete(r)
}
