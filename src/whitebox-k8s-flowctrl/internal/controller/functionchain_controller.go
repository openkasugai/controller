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
	"errors"
	"fmt"
	"strings"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// FunctionChainReconciler reconciles a FunctionChain object
type FunctionChainReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// in/out params for funcMaxInOutPairs
type maxNum struct {
	in  int
	out int
}

// from/to params for funcCount
type funcCount struct {
	from int
	to   int
}

//+kubebuilder:rbac:groups=example.com,resources=functionchains,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=functionchains/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=functionchains/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=functiontypes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=connectiontypes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

func (r *FunctionChainReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var functionChain ntthpcv1.FunctionChain
	l.Info("fetching FunctionChain Resource")
	if err := r.Get(ctx, req.NamespacedName, &functionChain); err != nil {
		l.Error(err, "unable to fetch FunctionChain")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// check if all functions are "Ready"
	var functionTypeList ntthpcv1.FunctionTypeList

	l.Info("fetching FunctionTypeList")
	err := r.List(ctx, &functionTypeList, &client.ListOptions{
		Namespace: functionChain.Spec.FunctionTypeNamespace,
	})
	if err != nil {
		l.Error(err, "unable to fetch FunctionTypeList")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	funcMaxInOutPairs := map[string]map[string]maxNum{}

	for fKey, f := range functionChain.Spec.Functions {
		found := false
		for _, ft := range functionTypeList.Items {
			if f.FunctionName == ft.Spec.FunctionName && f.Version == ft.Spec.Version {
				if ft.Status.Status == "Ready" {
					found = true

					// get FunctionInfo
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

					// parse functioninfo cm json
					for k, j := range cm.Data {
						if k == "spec" {
							funcInfoJson := j
							var result []map[string]interface{}
							json.Unmarshal([]byte(funcInfoJson), &result)

							// get maxInputnum/maxOutputNum per spec
							maxnum := maxNum{}
							specName := ""
							for _, v := range result {
								specName = v["name"].(string)
								if v["maxInputNum"] != nil {
									maxnum.in = int(v["maxInputNum"].(float64))
								} else {
									maxnum.in = 1
								}
								if v["maxOutputNum"] != nil {
									maxnum.out = int(v["maxOutputNum"].(float64))
								} else {
									maxnum.out = 1
								}
								// set to map
								if _, ok := funcMaxInOutPairs[fKey]; ok {
									funcMaxInOutPairs[fKey][specName] = maxnum
								} else {
									// init
									tmp := map[string]maxNum{}
									funcMaxInOutPairs[fKey] = tmp
									funcMaxInOutPairs[fKey][specName] = maxnum
								}
							}
						}
					}
					break
				} else {
					l.Error(err, fmt.Sprintf("Function was found but not Ready state: %s(%s)",
						ft.Spec.FunctionName, ft.Status.Status))
					r.Recorder.Eventf(&functionChain,
						corev1.EventTypeNormal,
						"Reconcile",
						"Function was found but not Ready state: %s(%s)",
						ft.Spec.FunctionName, ft.Status.Status)
					functionChain.Status.Status = "Error"
					r.Status().Update(ctx, &functionChain)
					return ctrl.Result{}, nil
				}
			}
		}
		if found == false {
			l.Error(err, fmt.Sprintf("Function was not found: %s(%s)", f.FunctionName, f.Version))
			r.Recorder.Eventf(&functionChain,
				corev1.EventTypeNormal,
				"Reconcile",
				"Reconcile Function was not found: %s(%s)",
				f.FunctionName, f.Version)
			functionChain.Status.Status = "Error"
			r.Status().Update(ctx, &functionChain)
			return ctrl.Result{}, nil
		}
	}

	funcFromToCount := map[string]funcCount{}
	funcFromPorts := map[string][]int{}
	funcToPorts := map[string][]int{}

	// check if each connectiontype has AvailabeDeviceKind for functions(from and to)
	for _, c := range functionChain.Spec.Connections {
		if c.ConnectionTypeName == "auto" {
			// save functionKey-parametors(from/to) pairs to map
			// count from function num
			if _, ok := funcFromToCount[c.From.FunctionKey]; ok {
				tmp := funcFromToCount[c.From.FunctionKey]
				tmp.from++
				funcFromToCount[c.From.FunctionKey] = tmp
			} else {
				tmp := funcCount{}
				tmp.from = 1
				tmp.to = 0
				funcFromToCount[c.From.FunctionKey] = tmp
			}
			// count to function num
			if _, ok := funcFromToCount[c.To.FunctionKey]; ok {
				tmp := funcFromToCount[c.To.FunctionKey]
				tmp.to++
				funcFromToCount[c.To.FunctionKey] = tmp
			} else {
				tmp := funcCount{}
				tmp.from = 0
				tmp.to = 1
				funcFromToCount[c.To.FunctionKey] = tmp
			}

			// push port number to slice
			funcFromPorts[c.From.FunctionKey] = append(funcFromPorts[c.From.FunctionKey], int(c.From.Port))
			funcToPorts[c.To.FunctionKey] = append(funcToPorts[c.To.FunctionKey], int(c.To.Port))

			// if auto, just continue
			// continue
		}

		/* The code below does not work because connectiontype is not implemented
		// get the connectiontype from spec
		var ck ntthpcv1.ConnectionType
		l.Info("fetching ConnectionType Resource")
		if err := r.Get(ctx, client.ObjectKey{Namespace: functionChain.Spec.ConnectionTypeNamespace,
			Name: c.ConnectionTypeName}, &ck); err != nil {
			l.Error(err, "unable to fetch ConnectionType")
			r.Recorder.Eventf(&functionChain,
				corev1.EventTypeNormal,
				"Reconcile",
				"unable to fetch ConnectionType")
			functionChain.Status.Status = "Error"
			r.Status().Update(ctx, &functionChain)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		// check if the connection kind status is ready
		if ck.Status.Status != "Ready" {
			l.Error(err, fmt.Sprintf("connectiontype %s is not Ready (%s)", ck.Name, ck.Status.Status))
			r.Recorder.Eventf(&functionChain,
				corev1.EventTypeNormal,
				"Reconcile",
				"connectiontype %s is not Ready (%s)",
				ck.Name, ck.Status.Status)
			functionChain.Status.Status = "Error"
			r.Status().Update(ctx, &functionChain)
			return ctrl.Result{}, nil
		}

		var fromto [2]string = [2]string{c.From.FunctionKey, c.To.FunctionKey}
		for _, ft := range fromto {
			found := false
			// if c.From or c.To is reserved word(start or end), skip check and continue
			if ft == "wb-start-of-chain" || ft == "wb-end-of-chain" {
				continue
			}
			fs := functionChain.Spec.Functions[ft]

			for _, i := range functionTypeList.Items {
				if fs.FunctionName == i.Spec.FunctionName && fs.Version == i.Spec.Version {
					for _, rkc := range i.Status.RegionTypeCandidates {
						for aif, _ := range ck.Status.AvailableInterfaces {
							if strings.HasPrefix(aif, rkc) {
								found = true
								l.Info(fmt.Sprintf("ConnectionType '%s' has AvailableDeviceKinds(%s) for function %s(%s)",
									ck.Name, aif, i.Spec.FunctionName, i.Spec.Version))
								break
							}
						}
						if found {
							break
						}
					}
					if found {
						break
					}
				}
			}
			if !found {
				l.Error(err, fmt.Sprintf("connectiontype %s does not have AvailableDeviceKinds for %s(%s)",
					ck.Name, fs.FunctionName, fs.Version))
				r.Recorder.Eventf(&functionChain,
					corev1.EventTypeNormal,
					"Reconcile",
					"connectiontype %s does not have AvailableDeviceKinds for %s(%s)",
					ck.Name, fs.FunctionName, fs.Version)
				functionChain.Status.Status = "Error"
				r.Status().Update(ctx, &functionChain)
				return ctrl.Result{}, nil
			}
		}
		*/
	}

	// check function's maxInputNum/maxOutputNum, and duplicate connections
	errmsg := ""

	// check from/to branch/integration num
	errmsg += r.checkFunctionMaxInOut(l, ctx, functionChain, funcFromToCount, funcMaxInOutPairs)
	// check 'from' duplicate connection
	errmsg += r.checkDupulicatePort(l, ctx, "from", functionChain, funcFromPorts)
	// check 'to' duplicate connection
	errmsg += r.checkDupulicatePort(l, ctx, "to", functionChain, funcToPorts)

	if errmsg != "" {
		l.Error(errors.New("validation check"), "Validation check... "+errmsg)
		r.Recorder.Eventf(&functionChain, corev1.EventTypeNormal, "Reconcile",
			"functionchain validation error %s", errmsg)
		functionChain.Status.Status = "Error"
		r.Status().Update(ctx, &functionChain)
		return ctrl.Result{}, nil
	}

	// FunctionChain Status Update
	l.Info("Setting FunctionChain Status")
	functionChain.Status.Status = "Ready"
	if err := r.Status().Update(ctx, &functionChain); err != nil {
		l.Error(err, "failed to update FunctionChain resource status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FunctionChainReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.FunctionChain{}).
		Complete(r)
}

// check from/to branch/integration num
func (r *FunctionChainReconciler) checkFunctionMaxInOut(l logr.Logger, ctx context.Context,
	functionChain ntthpcv1.FunctionChain, funcFromToCount map[string]funcCount,
	funcMaxInOutPairs map[string]map[string]maxNum) string {
	errmsg := ""

	for fKey, count := range funcFromToCount {
		// if start/end-of-chain, no check
		if strings.Contains(fKey, "wb-start-of-chain") || strings.Contains(fKey, "wb-end-of-chain") {
			continue
		}

		// compare specs(max in/out num) and function from/to count
		ok := false
		for _, maxnum := range funcMaxInOutPairs[fKey] {
			// found
			if count.from <= maxnum.out && count.to <= maxnum.in {
				ok = true
			}
		}
		// not found in specs, error
		if !ok {
			errmsg += fmt.Sprintf("(functionKey:%s to/from=%d/%d unmatch specs(max in/out))", fKey, count.to, count.from)
		}
	}

	return errmsg
}

// check duplicate connection
func (r *FunctionChainReconciler) checkDupulicatePort(l logr.Logger, ctx context.Context,
	t string, functionChain ntthpcv1.FunctionChain, funcPorts map[string][]int) string {
	errmsg := ""

	for fKey, s := range funcPorts {
		cm := map[int]int{}
		for _, port := range s {
			cm[port]++
		}
		for k, v := range cm {
			// from/to port k of fKey function is defined twice
			if v > 1 {
				errmsg += fmt.Sprintf("(%s.functionKey:%s duplicate port:%d)", t, fKey, k)
			}
		}
	}

	return errmsg
}
