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
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrl "sigs.k8s.io/controller-runtime"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TopologyInfoReconciler reconciles a TopologyInfo object
type TopologyInfoReconciler struct {
	client.Client
	Scheme                     *runtime.Scheme
	TopologyinfoNamespacedName types.NamespacedName
	TopologydataNamespacedName types.NamespacedName
}

// Controller Name
const (
	ControllerName = "topologyinfo"
)

// Logger Setting
const (
	loggerKeyControllerTi      = "topologyinfo"
	loggerKeyControllerGroupTi = "example.com"
	loggerKeyControllerKindTi  = "TopologyInfo"
)

// Resource Kind
const (
	TopologyinfoKind = "TopologyInfo"
	WbconnKind       = "WBConnection"
)

// Resource Name
const (
	WbconnName = "-wbconnection-"
)

// TopologyData ConfigMap KEY
const (
	TopologyCmEntitiesKey  = "entities"
	TopologyCmRelationsKey = "relations"
)

// EntityType
const (
	EntityTypeNode      = "node"
	EntityTypeDevice    = "device"
	EntityTypeInterface = "interface"
	EntityTypeNetwork   = "network"
)

// Info for each EntityType
const (
	EntityInfoNode      = "nodeInfo"
	EntityInfoDevice    = "deviceInfo"
	EntityInfoInterface = "interfaceInfo"
	EntityInfoNetwork   = "networkInfo"
)

// UsedType
const (
	UsedTypeIn    = "Incoming"
	UsedTypeOut   = "Outgoing"
	UsedTypeInout = "IncomingAndOutgoing"
)

// WBConnection Status
const (
	WbconnStatusDeployed = "Deployed"
)

// Finalizer Name
const (
	DomainName    = "example.com"
	FinalizerName = "topologyinfo-finalizer"
)

//+kubebuilder:rbac:groups=example.com,resources=topologyinfos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=topologyinfos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=topologyinfos/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=configmaps/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=wbconnections,verbs=get;list;watch
//+kubebuilder:rbac:groups=example.com,resources=wbconnections/status,verbs=get

func (r *TopologyInfoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// assign log entry
	l := log.FromContext(ctx)

	// check resource name and namespace
	if req.NamespacedName.Name == r.TopologydataNamespacedName.Name &&
		req.NamespacedName.Namespace == r.TopologydataNamespacedName.Namespace {
		// check event triggerd by TopologyData Configmap
		l.Info("execute checkTopologyDataEvent function")
		return r.checkTopologyDataEvent(ctx, req)
	} else if strings.Contains(req.NamespacedName.Name, WbconnName) {
		// update TopologyData resource status CapacityUsed
		l.Info("execute updateCapacityUsed function")
		return r.updateCapacityUsed(ctx, req)
	} else {
		return ctrl.Result{}, nil
	}
}

func (r *TopologyInfoReconciler) checkTopologyDataEvent(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// assign log entry
	l := log.FromContext(ctx)

	// fetch TopologyData Configmap
	var topologyDataCm corev1.ConfigMap

	l.Info("fetching TopologyData ConfigMap")
	if err := r.Get(ctx, req.NamespacedName, &topologyDataCm); err != nil {
		l.Error(err, "unable to fetch TopologyData ConfigMap")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// get finalizer name for topologyData ConfigMap
	finalizerName := getFinalizerName(&topologyDataCm)
	// set flag which explain TopologyData resource is already exists or not
	isNewTopologyInfo := false

	// check deletion timestamp
	if !topologyDataCm.ObjectMeta.DeletionTimestamp.IsZero() {
		// check if TopologyData ConfigMap has finalizer or not
		if controllerutil.ContainsFinalizer(&topologyDataCm, finalizerName) {
			// execute deleteTopologyData function
			l.Info("execute deleteTopologyInfo function")
			return r.deleteTopologyInfo(ctx, req, &topologyDataCm, finalizerName)
		}
	} else {
		// check if TopologyInfo resource has already created or not
		l.Info("check TopologyInfo resource existence")
		var topologyInfoCr ntthpcv1.TopologyInfo
		namespacedName := client.ObjectKey{Name: r.TopologyinfoNamespacedName.Name,
			Namespace: r.TopologyinfoNamespacedName.Namespace}
		err := r.Get(ctx, namespacedName, &topologyInfoCr)
		if errors.IsNotFound(err) {
			isNewTopologyInfo = true
		}
		// execute CreateOrUpdateToplogyData function
		l.Info("execute createOrUpdateTopologyInfo function")
		return r.createOrUpdateTopologyInfo(ctx, req, &topologyDataCm, isNewTopologyInfo, finalizerName)
	}
	return ctrl.Result{}, nil
}

func (r *TopologyInfoReconciler) createOrUpdateTopologyInfo(ctx context.Context, req ctrl.Request,
	topologyDataCm *corev1.ConfigMap, isNewTopologyInfo bool, finalizerName string) (ctrl.Result, error) {

	// assign log entry
	l := log.FromContext(ctx)

	var err error = nil

	// convert TopologyData ConfigMap JSON data to struct
	l.Info("converting TopologyData ConfigMap entities JSON data to struct")
	var topologyDataCmEntities []ntthpcv1.EntityInfo
	if entitiesData, keyExists := topologyDataCm.Data[TopologyCmEntitiesKey]; keyExists {
		if len(entitiesData) == 0 {
			l.Error(err, "TopologyData ConfigMap entities data is empty")
			return ctrl.Result{Requeue: false}, nil
		}
		if err := json.Unmarshal([]byte(entitiesData), &topologyDataCmEntities); err != nil {
			l.Error(err, "JSON parsing failed for entities")
			return ctrl.Result{Requeue: false}, nil
		}
	} else {
		l.Error(err, "TopologyData ConfigMap entities key does not exist")
		return ctrl.Result{Requeue: false}, nil
	}

	// validate entities data
	if err := r.validateEntitiesData(topologyDataCmEntities, l); err != nil {
		l.Error(err, "failed to validate entities data")
		return ctrl.Result{Requeue: false}, nil
	}

	l.Info("converting TopologyData ConfigMap relations data to struct")
	var topologyDataCmRelations []ntthpcv1.RelationInfo
	if relationsData, keyExists := topologyDataCm.Data[TopologyCmRelationsKey]; keyExists {
		if len(relationsData) == 0 {
			l.Error(err, "TopologyData ConfigMap relations data is empty")
			return ctrl.Result{Requeue: false}, nil
		}
		if err := json.Unmarshal([]byte(relationsData), &topologyDataCmRelations); err != nil {
			l.Error(err, "JSON parsing failed for entities")
			return ctrl.Result{Requeue: false}, nil
		}
	} else {
		l.Error(err, "TopologyData ConfigMap relations key does not exist")
		return ctrl.Result{Requeue: false}, nil
	}

	// validate relations data
	if err := r.validateRelationsData(topologyDataCmRelations, l); err != nil {
		l.Error(err, "failed to validate relations data")
		return ctrl.Result{Requeue: false}, nil
	}

	var topologyInfoCr ntthpcv1.TopologyInfo

	// if TopologyInfo resource has not created yet
	if isNewTopologyInfo {
		l.Info("TopologyInfo does not exists yet, so create new")
		// create TopologyInfo Spec
		topologyInfoCr = ntthpcv1.TopologyInfo{
			ObjectMeta: metav1.ObjectMeta{
				Name:      r.TopologyinfoNamespacedName.Name,
				Namespace: r.TopologyinfoNamespacedName.Namespace,
			},
			Spec: ntthpcv1.TopologyInfoSpec{
				TopologyDataCMRef: []ntthpcv1.WBNamespacedName{{
					Name:      req.NamespacedName.Name,
					Namespace: req.NamespacedName.Namespace,
				}},
			},
		}

		// create TopologyInfo resource
		l.Info("create for TopologyInfo resource")
		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, &topologyInfoCr, func() error {
			return nil
		}); err != nil {
			l.Error(err, "unable to Create TopologyInfo resource")
			return ctrl.Result{}, err
		}

		// set entities data
		topologyInfoCr.Status.Entities = topologyDataCmEntities

		// set relations data
		topologyInfoCr.Status.Relations = topologyDataCmRelations

		// if TopologyInfo resource is already exists
	} else {
		l.Info("TopologyInfo resource already exists, so update it")
		// fetch TopologyInfo resource
		l.Info("fetching TopologyInfo resource")
		namespacedName := client.ObjectKey{Name: r.TopologyinfoNamespacedName.Name,
			Namespace: r.TopologyinfoNamespacedName.Namespace}
		if err := r.Get(ctx, namespacedName, &topologyInfoCr); err != nil {
			l.Error(err, "unable to fetch TopologyInfo resource")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		// convert ToplogyInfo resource entities list to Map
		crEntitiesMap := convertEntitiesToMap(topologyInfoCr.Status.Entities)

		// convert TopologyData ConfigMap entities list to Map
		cmEntitiesMap := convertEntitiesToMap(topologyDataCmEntities)

		// update and delete entities in ToplogyInfo resource
		for tmpCrEntityKey, crEntity := range crEntitiesMap {
			if cmEntity, exists := cmEntitiesMap[tmpCrEntityKey]; exists {
				// update entity properties
				updateEntityProperties(&crEntity, &cmEntity, ctx)
				crEntitiesMap[tmpCrEntityKey] = crEntity
			} else {
				// delete entity from TopologyInfo resource which doesn't exists in TopologyData ConfigMap
				delete(crEntitiesMap, tmpCrEntityKey)
			}
		}

		// add new entities into ToplogyInfo resource which exists only in TopologyData ConfigMap
		for tmpCmEntityKey, cmEntity := range cmEntitiesMap {
			if _, exists := crEntitiesMap[tmpCmEntityKey]; !exists {
				crEntitiesMap[tmpCmEntityKey] = cmEntity
			}
		}

		// convert Entities map to list to update TopologyInfo Resource
		topologyInfoCr.Status.Entities = mapToEntitiesList(crEntitiesMap)

		// convert ToplogyInfo resource relations list to map
		crRelationsMap := convertRelationsToMap(topologyInfoCr.Status.Relations)
		// convert TopologyData ConfigMap relations list to map
		cmRelationsMap := convertRelationsToMap(topologyDataCmRelations)

		// update and delete relations in TopologyInfo resource
		for tmpCrRelationKey, crRelation := range crRelationsMap {
			if cmRelation, exists := cmRelationsMap[tmpCrRelationKey]; exists {
				// update relation property
				crRelation.Available = cmRelation.Available
				crRelationsMap[tmpCrRelationKey] = crRelation
			} else {
				// delete relation from TopologyInfo resource which doesn't exists in TopologyData ConfigMap
				delete(crRelationsMap, tmpCrRelationKey)
			}
		}

		// add new relations into ToplogyInfo Resource which exists only in TopologyData ConfigMap
		for tmpCmRelationKey, cmRelation := range cmRelationsMap {
			if _, exists := crRelationsMap[tmpCmRelationKey]; !exists {
				crRelationsMap[tmpCmRelationKey] = cmRelation
			}
		}
		// convert relations map to List to update TopologyInfo Resource
		topologyInfoCr.Status.Relations = mapToRelationsList(crRelationsMap)
	}

	// update TopologyInfo status
	l.Info("update TopologyInfo status")
	if err := r.Status().Update(ctx, &topologyInfoCr); err != nil {
		l.Error(err, "unable to update TopologyInfo status")
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(topologyDataCm, finalizerName) {
		// add finalizer to TopologyData ConfigMap
		controllerutil.AddFinalizer(topologyDataCm, finalizerName)
		l.Info("update TopologyData ConfigMap to add Finalizer")
		if err := r.Update(ctx, topologyDataCm); err != nil {
			l.Error(err, "unable to update TopologyData ConfigMap")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *TopologyInfoReconciler) deleteTopologyInfo(ctx context.Context, req ctrl.Request,
	topologyDataCm *corev1.ConfigMap, finalizerName string) (ctrl.Result, error) {

	// assign log entry
	l := log.FromContext(ctx)

	var err error

	// fetch TopologyInfo resource
	var topologyInfoCr ntthpcv1.TopologyInfo
	namespacedName := client.ObjectKey{Name: r.TopologyinfoNamespacedName.Name,
		Namespace: r.TopologyinfoNamespacedName.Namespace}
	l.Info("fetching TopologyInfo resource")
	if err := r.Get(ctx, namespacedName, &topologyInfoCr); err != nil {
		l.Error(err, "unable to fetch TopologyInfo resource")
		return ctrl.Result{}, err
	}

	// delete TopologyInfo resource
	l.Info("delete TopologyInfo resource")
	if err = r.Delete(ctx, &topologyInfoCr); err != nil {
		l.Error(err, "unable to delete TopologyInfo resource")
		return ctrl.Result{}, err
	}

	// delete finalizer from TopologyInfo ConfigMap
	controllerutil.RemoveFinalizer(topologyDataCm, finalizerName)
	l.Info("update TopologyInfo ConfigMap to remove Finalizer")
	if err := r.Update(ctx, topologyDataCm); err != nil {
		l.Error(err, "unable to update TopologyInfo ConfigMap")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// update TopologyInfo resource CapacityUsed
func (r *TopologyInfoReconciler) updateCapacityUsed(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// assign log entry
	l := log.FromContext(ctx)

	// fetch TopologyInfo resource
	var topologyInfoCr ntthpcv1.TopologyInfo
	l.Info("fetching TopologyInfo resource")
	namespacedName := client.ObjectKey{Name: r.TopologyinfoNamespacedName.Name,
		Namespace: r.TopologyinfoNamespacedName.Namespace}
	if err := r.Get(ctx, namespacedName, &topologyInfoCr); err != nil {
		l.Error(err, "unable to fetch TopologyInfo resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// fetch WBConnection resource
	var wbConnection ntthpcv1.WBConnection
	l.Info("fetching WBConnection resource")
	if err := r.Get(ctx, req.NamespacedName, &wbConnection); err != nil {
		l.Error(err, "unable to fetch WBConnction resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// status check
	l.Info("check WBConnection Status")
	if wbConnection.Status.Status != WbconnStatusDeployed {

		l.Info("WBConnection is not yet in the Deployed state, abort Reconcile")
		return ctrl.Result{}, nil

	} else if wbConnection.Status.Status == WbconnStatusDeployed {

		l.Info("WBConnection is now in the Deployed state, updating TopologyInfo Status")
		crEntitiesIndexMap := make(map[string]int)
		for idx, topologyEntity := range topologyInfoCr.Status.Entities {
			crEntitiesIndexMap[topologyEntity.ID] = idx
		}

		// get WBConnection connectionPath
		if wbConnection.Spec.ConnectionPath != nil {
			l.Info("get WBConnection ConnectionPath")
			for _, wbConnPathValue := range wbConnection.Spec.ConnectionPath {
				if idx, exists := crEntitiesIndexMap[wbConnPathValue.EntityID]; exists {
					topologyEntity := &topologyInfoCr.Status.Entities[idx]

					// set CapacityUsed
					l.Info("update CapacityUsed")
					switch wbConnPathValue.UsedType {
					case UsedTypeIn:
						if topologyEntity.CapacityInfo == nil {
							topologyEntity.CapacityInfo = &ntthpcv1.CapacityInfo{}
						}
						topologyEntity.CapacityInfo.CurrentIncomingCapacity += wbConnection.Spec.Requirements.Capacity
					case UsedTypeOut:
						if topologyEntity.CapacityInfo == nil {
							topologyEntity.CapacityInfo = &ntthpcv1.CapacityInfo{}
						}
						topologyEntity.CapacityInfo.CurrentOutgoingCapacity += wbConnection.Spec.Requirements.Capacity
					case UsedTypeInout:
						if topologyEntity.CapacityInfo == nil {
							topologyEntity.CapacityInfo = &ntthpcv1.CapacityInfo{}
						}
						topologyEntity.CapacityInfo.CurrentIncomingCapacity += wbConnection.Spec.Requirements.Capacity
						topologyEntity.CapacityInfo.CurrentOutgoingCapacity += wbConnection.Spec.Requirements.Capacity
					}
				}
			}
			// update TopologyInfo Resource status
			l.Info("update TopologyInfo resource status")
			if err := r.Status().Update(ctx, &topologyInfoCr); err != nil {
				l.Error(err, "unable to update TopologyInfo status")
				return ctrl.Result{}, err
			}
		} else {
			l.Info("WBConnection has no connectionpath information, abort reconcile")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func getFinalizerName(topologyDataCm *corev1.ConfigMap) string {
	// string for finalizer
	return strings.ToLower(DomainName) + "/" + strings.ToLower(FinalizerName)
}

// generate a unique key for an entity
func generateTmpEntityKey(entity ntthpcv1.EntityInfo) string {
	return fmt.Sprint(entity.ID)
}

// generate a unique key for a relation
func generateTmpRelationKey(relation ntthpcv1.RelationInfo) string {
	return fmt.Sprintf("%s-%s-%s", relation.Type, relation.From, relation.To)
}

// convert list of entities to a map
func convertEntitiesToMap(entities []ntthpcv1.EntityInfo) map[string]ntthpcv1.EntityInfo {
	entityMap := make(map[string]ntthpcv1.EntityInfo)
	for _, entity := range entities {
		key := generateTmpEntityKey(entity)
		entityMap[key] = entity
	}
	return entityMap
}

// convert list of relations to a map
func convertRelationsToMap(relations []ntthpcv1.RelationInfo) map[string]ntthpcv1.RelationInfo {
	relationMap := make(map[string]ntthpcv1.RelationInfo)
	for _, relation := range relations {
		key := generateTmpRelationKey(relation)
		relationMap[key] = relation
	}
	return relationMap
}

// convert map of entities back to a list
func mapToEntitiesList(entitiesMap map[string]ntthpcv1.EntityInfo) []ntthpcv1.EntityInfo {
	entities := []ntthpcv1.EntityInfo{}
	for _, entity := range entitiesMap {
		entities = append(entities, entity)
	}
	return entities
}

// convert map of relations back to a list
func mapToRelationsList(relationsMap map[string]ntthpcv1.RelationInfo) []ntthpcv1.RelationInfo {
	relations := []ntthpcv1.RelationInfo{}
	for _, relation := range relationsMap {
		relations = append(relations, relation)
	}
	return relations
}

// validate Entity Data
func (r *TopologyInfoReconciler) validateEntitiesData(topologyDataCmEntities []ntthpcv1.EntityInfo, l logr.Logger) error {
	checkedEntityIds := make(map[string]bool)
	for idx, cmEntity := range topologyDataCmEntities {
		if cmEntity.ID == "" || cmEntity.Type == "" {
			return fmt.Errorf("entity at index %d has one or more empty required fields", idx)
		} else {
			if _, exists := checkedEntityIds[cmEntity.ID]; exists {
				return fmt.Errorf("EntityId %s at index %d is duplicated", cmEntity.ID, idx)
			}
			checkedEntityIds[cmEntity.ID] = true
			// validate EntityType Data
			if err := validateEntityTypeData(cmEntity, l); err != nil {
				return err
			}
			if cmEntity.CapacityInfo != nil {
				// if TopologyData ConfigMap has IncomingCapacityUsed value, set it to 0
				if cmEntity.CapacityInfo.CurrentIncomingCapacity != 0 {
					cmEntity.CapacityInfo.CurrentIncomingCapacity = 0
				}
				// if TopologyData ConfigMap has OutgoingCapacityUsed value, set it to 0
				if cmEntity.CapacityInfo.MaxOutgoingCapacity != 0 {
					cmEntity.CapacityInfo.CurrentOutgoingCapacity = 0
				}
			}
		}
	}
	return nil
}

// validate EntityType Data
func validateEntityTypeData(cmEntity ntthpcv1.EntityInfo, l logr.Logger) error {
	switch cmEntity.Type {
	case EntityTypeNode:
		if cmEntity.NodeInfo == nil {
			return fmt.Errorf("EntityType is %s but corresponding field %s is nil", cmEntity.Type, EntityInfoNode)
		}
		if cmEntity.DeviceInfo != nil || cmEntity.InterfaceInfo != nil || cmEntity.NetworkInfo != nil {
			return fmt.Errorf("EntityType is %s but unnecessary field is set", cmEntity.Type)
		}
	case EntityTypeDevice:
		if cmEntity.DeviceInfo == nil {
			return fmt.Errorf("EntityType is %s but corresponding field %s is nil", cmEntity.Type, EntityInfoDevice)
		}
		if cmEntity.NodeInfo != nil || cmEntity.InterfaceInfo != nil || cmEntity.NetworkInfo != nil {
			return fmt.Errorf("EntityType is %s but unnecessary field is set", cmEntity.Type)
		}
	case EntityTypeInterface:
		if cmEntity.InterfaceInfo == nil {
			return fmt.Errorf("EntityType is %s but corresponding field %s is nil", cmEntity.Type, EntityInfoInterface)
		}
		if cmEntity.NodeInfo != nil || cmEntity.DeviceInfo != nil || cmEntity.NetworkInfo != nil {
			return fmt.Errorf("EntityType is %s but unnecessary field is set", cmEntity.Type)
		}
	case EntityTypeNetwork:
		if cmEntity.NetworkInfo == nil {
			return fmt.Errorf("EntityType is %s but corresponding field %s is nil", cmEntity.Type, EntityInfoNetwork)
		}
		if cmEntity.NodeInfo != nil || cmEntity.DeviceInfo != nil || cmEntity.InterfaceInfo != nil {
			return fmt.Errorf("EntityType is %s but unnecessary field is set", cmEntity.Type)
		}
	default:
		return fmt.Errorf("EntityType %s is not one of the expected values (%s, %s, %s, %s)",
			cmEntity.Type, EntityTypeNode, EntityTypeDevice, EntityInfoInterface, EntityInfoNetwork)
	}
	return nil
}

// validate Relation Data
func (r *TopologyInfoReconciler) validateRelationsData(topologyDataCmRelations []ntthpcv1.RelationInfo, l logr.Logger) error {
	checkedRelationTmpKeys := make(map[string]bool)
	for idx, cmRelation := range topologyDataCmRelations {
		tmpCmRelationKey := generateTmpRelationKey(cmRelation)
		if cmRelation.Type == "" || cmRelation.From == "" || cmRelation.To == "" {
			return fmt.Errorf("Relation at index %d has one or more empty fields", idx)
		} else {
			if _, exists := checkedRelationTmpKeys[tmpCmRelationKey]; exists {
				return fmt.Errorf("Relation at index %d is duplicated", idx)
			}
			checkedRelationTmpKeys[tmpCmRelationKey] = true
		}
	}
	return nil
}

// update entity properties
func updateEntityProperties(crEntity, cmEntity *ntthpcv1.EntityInfo, ctx context.Context) {
	var crCurrentIncomingCapacity, crCurrentOutgoingCapacity int32
	switch cmEntity.Type {
	case EntityTypeNode:
		if cmEntity.NodeInfo != nil {
			crEntity.NodeInfo = cmEntity.NodeInfo
		}
	case EntityTypeDevice:
		if cmEntity.DeviceInfo != nil {
			crEntity.DeviceInfo = cmEntity.DeviceInfo
		}
	case EntityTypeInterface:
		if cmEntity.InterfaceInfo != nil {
			crEntity.InterfaceInfo = cmEntity.InterfaceInfo
		}
	case EntityTypeNetwork:
		if cmEntity.NetworkInfo != nil {
			crEntity.NetworkInfo = cmEntity.NetworkInfo
		}
	}

	if crEntity.CapacityInfo != nil {
		crCurrentIncomingCapacity = crEntity.CapacityInfo.CurrentIncomingCapacity
		crCurrentOutgoingCapacity = crEntity.CapacityInfo.CurrentOutgoingCapacity
	}

	if cmEntity.CapacityInfo != nil {
		if cmEntity.CapacityInfo.MaxIncomingCapacity != 0 {
			if crEntity.CapacityInfo == nil {
				crEntity.CapacityInfo = &ntthpcv1.CapacityInfo{}
			}
			crEntity.CapacityInfo.MaxIncomingCapacity = cmEntity.CapacityInfo.MaxIncomingCapacity
		}
		if cmEntity.CapacityInfo.MaxOutgoingCapacity != 0 {
			if crEntity.CapacityInfo == nil {
				crEntity.CapacityInfo = &ntthpcv1.CapacityInfo{}
			}
			crEntity.CapacityInfo.MaxOutgoingCapacity = cmEntity.CapacityInfo.MaxOutgoingCapacity
		}
	}

	if crEntity.CapacityInfo != nil {
		crEntity.CapacityInfo.CurrentIncomingCapacity = crCurrentIncomingCapacity
		crEntity.CapacityInfo.CurrentOutgoingCapacity = crCurrentOutgoingCapacity
	}
	crEntity.Available = cmEntity.Available
}

// ignore finalizer update conducted by TopologyInfo Controller
func ignoreFinelizerUpdate(r client.Client) predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectNew.GetDeletionTimestamp() != nil {
				return true
			} else {
				oldObject, oldIsCm := e.ObjectOld.(*corev1.ConfigMap)
				newObject, newIsCm := e.ObjectNew.(*corev1.ConfigMap)
				if !oldIsCm || !newIsCm {
					return true
				} else {
					if !reflect.DeepEqual(oldObject.Data, newObject.Data) {
						return true
					} else {
						return false
					}
				}
			}
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *TopologyInfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named(ControllerName).
		For(&corev1.ConfigMap{}).WithEventFilter(ignoreFinelizerUpdate(r.Client)).
		WithLogConstructor(func(req *reconcile.Request) logr.Logger {
			return mgr.GetLogger().WithValues("controller", loggerKeyControllerTi, "controllerGroup",
				loggerKeyControllerGroupTi, "controllerKind", loggerKeyControllerKindTi)
		}).
		Watches(&ntthpcv1.WBConnection{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
