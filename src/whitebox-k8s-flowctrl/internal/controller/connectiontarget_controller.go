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
	"strconv"
	"strings"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ConnectionTargetReconciler reconciles a ConnectionTarget object
type ConnectionTargetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// json parsing
type IOResourceValue struct {
	Data map[string]Connections `json:"-"`
}

type Connections struct {
	Type          string `json:"type"`
	Available     bool   `json:"available"`
	IncomingTotal int    `json:"incomingTotal"`
	IncomingUsed  int    `json:"incomingUsed"`
	OutgoingTotal int    `json:"outgoingTotal"`
	OutgoingUsed  int    `json:"outgoingUsed"`
}

//+kubebuilder:rbac:groups=example.com,resources=connectiontargets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=connectiontargets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=connectiontargets/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConnectionTarget object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ConnectionTargetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO(user): your logic here
	l := log.FromContext(ctx)

	// Check that the updated resource is an IOResource
	if !strings.Contains(req.NamespacedName.Name, "ios-") {
		return ctrl.Result{}, nil
	}

	// fetch ConfigMap ConfigMap
	var io corev1.ConfigMap
	if err := r.Get(ctx, req.NamespacedName, &io); err != nil {
		l.Error(err, "unable to fetch ConfigMap")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// create ConnectionTarget Resource
	for k, v := range io.Data {
		// parse json
		var result map[string]map[string]interface{}
		json.Unmarshal([]byte(v), &result)

		nodeName := strings.ReplaceAll(req.NamespacedName.Name, "ios-", "")
		DeviceType := k[:strings.LastIndex(k, "-")]
		deviceindex, _ := strconv.Atoi(k[strings.LastIndex(k, "-")+1:])

		for _, str := range k[strings.LastIndex(k, "-")+1:] {
			if !('0' <= str && str <= '9') {
				DeviceType = k
				deviceindex = 0
				break
			}
		}

		var createResourceList []string

		l.Info("DeviceType:" + DeviceType)
		if strings.Contains(DeviceType, "host") {
			// get DeviceType
			var cm corev1.ConfigMap
			if err := r.Get(ctx, client.ObjectKey{Name: "compute-" + nodeName, Namespace: req.NamespacedName.Namespace}, &cm); err != nil {
				l.Error(err, "unable to fetch ComputeResource")
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			for k, v := range cm.Data {
				jsonBytes := []byte(v)

				// parse json
				var crv ComputeResourceValue
				if err := json.Unmarshal(jsonBytes, &crv); err != nil {
					l.Error(err, "jsonparsing")
				}
				funDeviceType := k[:strings.LastIndex(k, "-")]

				l.Info("funDeviceType:" + funDeviceType)
				if !strings.Contains(funDeviceType, "Func") {
					createResourceList = append(createResourceList, funDeviceType+"-"+DeviceType)
				}
			}
		} else if strings.Contains(DeviceType, "static") {
			// get DeviceType
			var cm corev1.ConfigMap
			if err := r.Get(ctx, client.ObjectKey{Name: "compute-" + nodeName, Namespace: req.NamespacedName.Namespace}, &cm); err != nil {
				l.Error(err, "unable to fetch ComputeResource")
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			for k, v := range cm.Data {
				jsonBytes := []byte(v)

				// parse json
				var crv ComputeResourceValue
				if err := json.Unmarshal(jsonBytes, &crv); err != nil {
					l.Error(err, "jsonparsing")
				}
				funDeviceType := k[:strings.LastIndex(k, "-")]
				if strings.Contains(funDeviceType, "Func") {
					createResourceList = append(createResourceList, funDeviceType)
				}
			}
		} else {
			createResourceList = append(createResourceList, DeviceType)
		}

		for _, createResource := range createResourceList {
			// Update ConnectionTargetStatus
			for ioName, _ := range result {
				// define objectName
				objectName := strings.ToLower(req.NamespacedName.Name + "-" + createResource)
				if ioName == "pcie" {
					objectName = objectName + "-" + strconv.Itoa(deviceindex) + "-" + ioName
				}

				newct := &ntthpcv1.ConnectionTarget{
					ObjectMeta: metav1.ObjectMeta{
						Name:      objectName,
						Namespace: req.NamespacedName.Namespace,
					},
				}

				// Create or Update ConnectionTarget object
				if _, err := ctrl.CreateOrUpdate(ctx, r.Client, newct, func() error {

					newct.Spec.IOResourceRef.Name = req.NamespacedName.Name
					newct.Spec.IOResourceRef.Namespace = req.NamespacedName.Namespace

					// end of ctrl.CreateOrUpdate
					return nil

				}); err != nil {

					// error handling of ctrl.CreateOrUpdate
					l.Error(err, "unable to ensure deployment is correct")
					return ctrl.Result{}, err
				}

				newct.Status.NodeName = nodeName
				if strings.Contains(createResource, "host") && strings.Contains(createResource, "-") {
					newct.Status.DeviceType = createResource[:strings.Index(createResource, "-")] //nolint:gocritic // Index() doesn't return '-1' because check createResource contains '-' before
				} else {
					newct.Status.DeviceType = createResource
				}
				newct.Status.DeviceIndex = deviceindex

				if err := r.Status().Update(ctx, newct); err != nil {
					l.Error(err, "unable to update connectiontarget status")
					return ctrl.Result{}, err
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConnectionTargetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&corev1.ConfigMap{}).
		Complete(r)
}
