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
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
)

// FunctionTargetReconciler reconciles a FunctionTarget object
type FunctionTargetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// json parsing
type ComputeResourceValue struct {
	Available        bool                 `json:"available"`
	MaxDataFlows     int                  `json:"total"`
	CurrentDataFlows int                  `json:"used"`
	Functions        map[string]Functions `json:"functions"`
}

type Functions struct {
	Available        bool `json:"available"`
	MaxDataFlows     int  `json:"total"`
	CurrentDataFlows int  `json:"used"`
}

//+kubebuilder:rbac:groups=example.com,resources=functiontargets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=functiontargets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=functiontargets/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=computeresources,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FunctionTarget object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *FunctionTargetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var computeresource ntthpcv1.ComputeResource

	l.Info("fetching ComputeResource Resource")
	if err := r.Get(ctx, req.NamespacedName, &computeresource); err != nil {
		l.Error(err, "unable to fetch ComputeResource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, ri := range computeresource.Status.Regions {
		// create or update functiontarget resource
		ftname := computeresource.Status.NodeName + "." + ri.DeviceType + "-" +
			strconv.Itoa(int(ri.DeviceIndex)) + "." + ri.Name
		functiontarget := &ntthpcv1.FunctionTarget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ftname,
				Namespace: req.Namespace,
			},
		}

		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, functiontarget, func() error {

			l.Info(fmt.Sprintf("CreateOrUpdate for FunctionTarget Resouce"))
			functiontarget.Spec.ComputeResourceRef = ntthpcv1.WBNamespacedName(req.NamespacedName)

			// set the owner so that garbage collection can kicks in
			if err := ctrl.SetControllerReference(&computeresource, functiontarget, r.Scheme); err != nil {
				l.Error(err, "unable to set ownerReference from ComputeResource to FunctionTarget")
				return err
			}

			// end of ctrl.CreateOrUpdate
			return nil

		}); err != nil {

			// error handling of ctrl.CreateOrUpdate
			l.Error(err, "unable to ensure functiontarget is correct")
			return ctrl.Result{}, err
		}

		// Update FunctionTarget status
		functiontarget.Status.RegionName = ri.Name
		functiontarget.Status.RegionType = ri.Type
		functiontarget.Status.NodeName = computeresource.Status.NodeName
		functiontarget.Status.DeviceType = ri.DeviceType
		functiontarget.Status.DeviceIndex = ri.DeviceIndex
		functiontarget.Status.Available = ri.Available
		functiontarget.Status.Status = ri.Status
		functiontarget.Status.MaxFunctions = ri.MaxFunctions
		functiontarget.Status.CurrentFunctions = ri.CurrentFunctions
		functiontarget.Status.MaxCapacity = ri.MaxCapacity
		functiontarget.Status.CurrentCapacity = ri.CurrentCapacity
		if ri.Functions != nil {
			functiontarget.Status.Functions = []ntthpcv1.FunctionCapStruct{}
			for _, fis := range ri.Functions {
				var fts ntthpcv1.FunctionCapStruct
				fts.FunctionIndex = fis.FunctionIndex
				fts.FunctionName = fis.FunctionName
				fts.Available = fis.Available
				fts.MaxDataFlows = fis.MaxDataFlows
				fts.CurrentDataFlows = fis.CurrentDataFlows
				fts.MaxCapacity = fis.MaxCapacity
				fts.CurrentCapacity = fis.CurrentCapacity
				functiontarget.Status.Functions = append(functiontarget.Status.Functions, fts)
			}
		} else {
			functiontarget.Status.Functions = nil
		}

		l.Info("Update FunctionTarget status")
		if err := r.Status().Update(ctx, functiontarget); err != nil {
			l.Error(err, "unable to update df status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FunctionTargetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.ComputeResource{}).
		Complete(r)
}
