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
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
)

// TODO: support other types
const (
	undefined = iota
	device
	node
	rack
	floor
	other
)

func StringToConst(s string) int {
	switch s {
	case "device":
		return device
	case "node":
		return node
	case "rack":
		return rack
	case "floor":
		return floor
	default:
		// TODO: fix this to support OTHER pattern
		return other
	}
}

func ConstToString(i int) string {
	switch i {
	case device:
		return "device"
	case node:
		return "node"
	case rack:
		return "rack"
	case floor:
		return "floor"
	default:
		// TODO: fix this to support OTHER pattern
		return "other"
	}
}

// ConnectionTypeReconciler reconciles a ConnectionType object
type ConnectionTypeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type ConnectionInfoSpecJson struct {
	ConnectionInfoItem map[string]ConnectionInfoItem `json:"items"`
}

type ConnectionInfoItem struct {
	Latency string `json:"latency"`
	Perf    string `json:"perf"`
	Power   string `json:"power"`
	Type    string `json:"type"`
}

type Topology struct {
	// stores interface names from ConnectionInfo
	ConnectionEdges []string
	// stores connection of edges.
	ConnectionSides []ConnectionSide
	// first index corresponds ConnectionEdges
	// second index stores type
	ConnectionTopology [][]int
}

type ConnectionSide struct {
	// P1 and P2 are name of interface
	P1   string
	P2   string
	Type string
}

func (c *ConnectionSide) IsSame(p ConnectionSide) bool {
	if (c.P1 == p.P1 && c.P2 == p.P2) || (c.P1 == c.P2 && c.P2 == p.P1) {
		return true
	}
	return false
}

func (t *Topology) AppendSide(p1 string, p2 string, conType string) error {
	s := ConnectionSide{p1, p2, conType}
	found := false
	// ignore if duplicate side is defined
	for _, i := range t.ConnectionSides {
		if s.IsSame(i) == true {
			found = true
			break
		}
	}
	if found == false {
		t.ConnectionSides = append(t.ConnectionSides, s)
	}
	return nil
}

func (t *Topology) GetIndexFromEdgeName(n string) int {
	// get index of connection edge from interface name
	for i, v := range t.ConnectionEdges {
		if n == v {
			return i
		}
	}
	return -1
}

func (t *Topology) Generate() error {
	// generate ConnectionEdges and ConnectionTopology from ConnectionSides
	for _, i := range t.ConnectionSides {
		for _, j := range [2]string{i.P1, i.P2} {
			found := false
			for _, k := range t.ConnectionEdges {
				if j == k {
					found = true
					break
				}
			}
			if found == false {
				t.ConnectionEdges = append(t.ConnectionEdges, j)
			}
		}
	}

	// generate Topology
	t.ConnectionTopology = make([][]int, len(t.ConnectionEdges))
	for i := range t.ConnectionTopology {
		t.ConnectionTopology[i] = make([]int, len(t.ConnectionSides))
	}

	for i, v := range t.ConnectionSides {
		// convert type of the connection to const value to calculate reachability to the connection type
		t.ConnectionTopology[t.GetIndexFromEdgeName(v.P1)][i] = StringToConst(v.Type)
		t.ConnectionTopology[t.GetIndexFromEdgeName(v.P2)][i] = StringToConst(v.Type)
	}
	// fmt.Println(t.ConnectionEdges)
	// fmt.Println(t.ConnectionTopology)

	return nil
}

func (t *Topology) GetPathToTargetType(targetType int, types []int, path []int,
	pathMap *map[string]ntthpcv1.DestinationStruct) (bool, error) {
	// targetType -> target type to find (const value)
	// types      -> stores connection type of each hops
	// path       -> stores current finding path
	// pathMap    -> stores map of paths. the key is the interface name of the destination.

	curIndex := path[len(path)-1]
	prevType := types[len(types)-1]

	// TODO: remove these hardcoded strings from the code. related to issue #31
	FPGA_STATIC_DELIMITER := "-static-"
	IONAMES := [4]string{"pcie", "host-mem", "host-100gether", "rdma"}

	// the first interface is just I/O, do nothing.
	var f bool
	if len(path) == 1 {
		for _, a := range IONAMES {
			if t.ConnectionEdges[curIndex] == a {
				f = true
				break
			}
		}
		if f || strings.Contains(t.ConnectionEdges[curIndex], FPGA_STATIC_DELIMITER) {
			// fmt.Printf("first interface (%s) must be acc's interface. return\n", t.ConnectionEdges[curIndex])
			return false, nil
		}
	}

	// fmt.Printf("current index is %d\n", curIndex)
	// fmt.Printf("interface name of current index is %s\n", t.ConnectionEdges[curIndex])
	// fmt.Printf("current path is %d\n", path)
	// fmt.Printf("interfaces of path are [")
	// for _, m := range path {
	// 	fmt.Printf("%s(%d) -> ", t.ConnectionEdges[m], m)
	// }
	// fmt.Printf("end]\n")

	// next hop variables
	var nexthopIndexes []int
	var nexthopTypes []int
	var nexthopIsDest []bool

	for i, ty := range t.ConnectionTopology[curIndex] {
		// ty is connection type of nexthop

		if ty < prevType {
			// do nothing if the nexthop is smaller than previous type
			continue
		} else if ty > targetType {
			// do nothing if the nexthop is larger than target type
			continue
		} else if ty >= rack && prevType >= rack {
			// do nothing if the current and nexthop are both larger than RACK
			// TODO: fix this to support FLOOR pattern
			continue
		} else {
			// else, append index to nexthops
			for x, z := range t.ConnectionTopology {
				// x is the index of the other edge of the side
				if x == curIndex {
					// continue because it makes loop
					continue
				}
				if z[i] <= 0 {
					// continue because route does not exist
					continue
				} else if z[i] == ty {
					found := false
					// skip if the nexthopIndexes has already have the index of nexthop
					for _, j := range nexthopIndexes {
						if x == j {
							found = true
							break
						}
					}
					if found {
						continue
					}

					// skip if the path has already have the index of nexthop
					for _, k := range path {
						if x == k {
							found = true
							break
						}
					}
					if found {
						continue
					}

					// // skip if the nexthop's interface name has entire of current one
					// // example: pci -> gv100-pci
					// if strings.Contains(t.ConnectionEdges[x], t.ConnectionEdges[curIndex]) {
					// 	//fmt.Printf("%s has %s. continue\n", t.ConnectionEdges[x], t.ConnectionEdges[curIndex])
					// 	continue
					// }

					// skip if the nexthop's interface name has prefix of functiontarget when type is not DEVICE
					if targetType != device {
						var p bool
						for _, a := range IONAMES {
							if t.ConnectionEdges[x] == a {
								// fmt.Printf("%s has acc name (%s). continue\n", t.ConnectionEdges[x], a)
								p = true
								break
							}
						}
						if !p && strings.Index(t.ConnectionEdges[x], FPGA_STATIC_DELIMITER) < 0 {
							continue
						}
					}

					// append it for nexthop
					var d bool
					if ty == targetType {
						d = true
					}
					// fmt.Printf("appending next hop: %s(%d) -> %s(%d)\n", t.ConnectionEdges[curIndex], prevType, t.ConnectionEdges[x], ty)
					nexthopIndexes = append(nexthopIndexes, x)
					nexthopTypes = append(nexthopTypes, ty)
					nexthopIsDest = append(nexthopIsDest, d)
				}
			}
		}
	}

	ret := false
	// update lists of path and type
	for i, c := range nexthopIndexes {
		newpath := append(path, c) //nolint:gocritic // intentionally assign to another silce(variable) 'newpath'
		var newtypes []int
		if types[0] == undefined {
			// override the first types
			newtypes = []int{nexthopTypes[i]}
		} else {
			newtypes = append(types, nexthopTypes[i]) //nolint:gocritic // intentionally assign to another silce(variable) 'newtype'
		}
		if nexthopIsDest[i] {
			_, ok := (*pathMap)[t.ConnectionEdges[c]]
			if ok {
				// fmt.Printf("path updating. current=%d, new=%d\n", len((*pathMap)[t.ConnectionEdges[c]].Route), len(newpath))
				if len((*pathMap)[t.ConnectionEdges[c]].Route) > len(newpath) {
					// fmt.Printf("path updated. current=%d, new=%d\n", len((*pathMap)[t.ConnectionEdges[c]].Route), len(newpath))
					(*pathMap)[t.ConnectionEdges[c]] = t.GetDestinationStruct(newtypes, newpath)
				}
			} else {
				// fmt.Printf("path added: the key is %d\n", c)
				(*pathMap)[t.ConnectionEdges[c]] = t.GetDestinationStruct(newtypes, newpath)
			}
		}
		// check nexthops recursively if targetType != DEVICE
		if targetType != device {
			ret, _ = t.GetPathToTargetType(targetType, newtypes, newpath, pathMap)
		}
	}

	if len(*pathMap) > 0 {
		ret = true
	}

	return ret, nil
}

func (t *Topology) GetDestinationStruct(types []int, path []int) ntthpcv1.DestinationStruct {
	// generate map[string]ntthpcv1.DestinationStruct from path and types
	var routes []ntthpcv1.RouteStruct
	for i, p := range path {
		if i == len(path)-1 {
			break
		}
		routes = append(routes, ntthpcv1.RouteStruct{Type: ConstToString(types[i]),
			From: t.ConnectionEdges[p],
			To:   t.ConnectionEdges[path[i+1]]})
	}
	return ntthpcv1.DestinationStruct{Route: routes}
}

func (t *Topology) GetAvailableInterfaces(targetType string) (map[string]ntthpcv1.AvailableInterfaceStruct,
	error) {
	// edges := []string{}
	ifs := map[string]ntthpcv1.AvailableInterfaceStruct{}
	for i, _ := range t.ConnectionEdges {
		p := map[string]ntthpcv1.DestinationStruct{}
		// fmt.Println("=============================")
		// fmt.Printf("finding path from: %s\n", t.ConnectionEdges[i])
		// fmt.Println("=============================")
		ret, _ := t.GetPathToTargetType(StringToConst(targetType), []int{undefined}, []int{i}, &p)
		if ret {
			// for _, v := range p {
			// 	fmt.Printf("[end: %s]\n", j)
			// 	fmt.Printf("path: ")
			// 	for j, k := range v.Route {
			// 		fmt.Printf("[%d]%s-(%s)->%s, ", j, k.From, k.Type, k.To)
			// 	}
			// 	fmt.Printf("->END\n")
			// }

			// append it to ntthpcv1.AvailableInterfaces
			ifs[t.ConnectionEdges[i]] = ntthpcv1.AvailableInterfaceStruct{Destinations: p}
		}
	}
	return ifs, nil
}

//+kubebuilder:rbac:groups=example.com,resources=connectiontypes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=connectiontypes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=connectiontypes/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConnectionType object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ConnectionTypeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// get connectiontype resource
	var Ct ntthpcv1.ConnectionType
	l.Info("fetching ConnectionType Resource")
	if err := r.Get(ctx, req.NamespacedName, &Ct); err != nil {
		l.Error(err, "unable to fetch ConnectionType")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var topology Topology

	// do it for each namespaces
	for _, n := range Ct.Spec.ConnectionInfoNameSpaces {
		// do it for per connectioninfo resources
		var configMaps v1.ConfigMapList

		// get configmaps in the namespace
		l.Info("fetching configmap")
		err := r.List(ctx, &configMaps, &client.ListOptions{
			Namespace: n,
		})
		if err != nil {
			l.Error(err, "unable to fetch ConfigMap")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		// for each configmap
		for _, cm := range configMaps.Items {
			// check it is connectioninfo
			// TODO: fix this literal
			if !strings.HasPrefix(cm.Name, "conninfo-") {
				continue
			}

			// for each device or io
			for src, jsBody := range cm.Data {
				// parse json body and get targets
				var mapData ConnectionInfoSpecJson
				err := json.Unmarshal([]byte(jsBody), &mapData)
				if err != nil {
					l.Error(err, fmt.Sprintf("json parse failed with key: %v, val: %v", src, jsBody))
				}
				// for each dst
				for dst, item := range mapData.ConnectionInfoItem {
					topology.AppendSide(src, dst, item.Type)
				}
			}
		}
	}

	// generate topology from connectioninfos
	topology.Generate()

	// check reachiablity to "spec.Name" for each edges
	ifs, _ := topology.GetAvailableInterfaces(Ct.Name)

	// update status.AvailableInterfaces if AvailableDevicesList is updated
	updated := false

	if len(ifs) != len(Ct.Status.AvailableInterfaces) {
		updated = true
	} else {
		for k, _ := range Ct.Status.AvailableInterfaces {
			found := false
			for e, _ := range ifs {
				// Todo: check deep diff
				if k == e {
					found = true
					break
				}
			}
			if found == false {
				updated = true
				break
			}
		}
	}

	if updated {
		Ct.Status.AvailableInterfaces = ifs
		ifl := []string{}
		for i, _ := range ifs {
			ifl = append(ifl, i)
		}
		Ct.Status.Interfaces = ifl
	} else if Ct.Status.AvailableInterfaces == nil {
		// put empty list into the status when the connectiontype is created
		Ct.Status.AvailableInterfaces = map[string]ntthpcv1.AvailableInterfaceStruct{}
		Ct.Status.Interfaces = []string{}
	}

	var newStat string

	if len(Ct.Status.AvailableInterfaces) > 0 {
		newStat = "Ready"
	} else {
		newStat = "Not Ready"
	}

	if newStat != Ct.Status.Status {
		Ct.Status.Status = newStat
	}

	// Update Spec.Name to the resource name to avoid issue that
	// can't fetch metadata when the resource is saved by other resource (dataflow)
	// refer: issue 27
	if Ct.Spec.ConnectionTypeName != Ct.Name {
		Ct.Spec.ConnectionTypeName = Ct.Name
		if err := r.Update(ctx, &Ct); err != nil {
			l.Error(err, "failed to update ConnectionType resource spec")
			return ctrl.Result{}, err
		}
	}

	if err := r.Status().Update(ctx, &Ct); err != nil {
		l.Error(err, "failed to update ConnectionType resource status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConnectionTypeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.ConnectionType{}).
		Complete(r)
}
