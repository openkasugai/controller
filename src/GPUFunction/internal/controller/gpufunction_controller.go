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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplecomv1 "GPUFunction/api/v1"

	/* Additional files */
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"strings"
	/* Additional files end here */)

// GPUFunctionReconciler reconciles a GPUFunction object
type GPUFunctionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

/* Event type */
const (
	CREATE = iota
	UPDATE
	DELETE
)

/* Overall Status type */
const (
	PENDING = "Pending"
	RUNNING = "Running"
)

var FunctionIndexEnable map[string]map[int32]bool = make(map[string]map[int32]bool)

//+kubebuilder:rbac:groups=example.com,resources=gpufunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=gpufunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=gpufunctions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GPUFunction object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *GPUFunctionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var crData examplecomv1.GPUFunction
	var podEnvData []corev1.EnvVar
	var eventKind int // 0:Add, 1:Upd,  2:Del
	var imageInfo string
	var podGPUID string
	var podData *corev1.Pod
	var podCRData corev1.PodSpec
	var podArgsData []string

	var err error

	// Get CR information
	err = r.Get(ctx, req.NamespacedName, &crData)
	if errors.IsNotFound(err) {
		// If CR does not exist
		logger.Info("NotFound to fetch CR")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else if err != nil {
		logger.Error(err, "unable to fetch CR")
		return ctrl.Result{}, err
	}

	// Supports daemonization
	myNodeName := os.Getenv("K8S_NODENAME")
	if myNodeName == crData.Spec.NodeName {
		// do nothing
	} else {
		// Do nothing except the target worker node
		return ctrl.Result{}, nil
	}

	// Get the event type
	eventKind = r.GetEventKind(&crData)
	if eventKind == CREATE {
		// In case of creation
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Create", "Create Start")

		var fromFunctionData examplecomv1.FunctionData
		var fromFunctionKind string
		var previousListKey []string
		var nextListKey []string

		var previousConnectionCRName string
		var nextConnectionCRName string
		var fromConnectionKind string
		var toConnectionKind string

		var rxProtocol string
		var txProtocol string

		if len(crData.Spec.PreviousFunctions) == 0 {
			fromConnectionKind = "wb-start-of-chain"
		} else {
			for key, _ := range crData.Spec.PreviousFunctions {
				previousListKey = append(previousListKey, key)
			}
			// From side ConnectionCR name generation
			MakeConnectionCRName(ctx,
				crData.Spec.DataFlowRef.Name,
				crData.Spec.PreviousFunctions[previousListKey[0]].WBFunctionRef.Name,
				req.Name,
				&previousConnectionCRName)
			logger.Info("previousConnectionCRName : " + previousConnectionCRName)

			// Get Kind of ConnectionCR on the From side
			err = r.GetConnectionData(ctx,
				req.Namespace,
				previousConnectionCRName,
				&fromConnectionKind)
			if nil != err {
				if errors.IsNotFound(err) {
					// do nothing
				} else {
					logger.Error(err, "GetConnectionData() error")
				}
				// Requeue due to CR information acquisition error on source device
				return ctrl.Result{Requeue: true}, nil
			}
		}

		if "high-infer" != crData.Spec.FunctionName && "low-infer" != crData.Spec.FunctionName {
			if len(crData.Spec.NextFunctions) == 0 {
				toConnectionKind = "wb-end-of-chain"
			} else {
				for key, _ := range crData.Spec.NextFunctions {
					nextListKey = append(nextListKey, key)
				}
				// Generate the To side ConnectionCR name
				MakeConnectionCRName(ctx,
					crData.Spec.DataFlowRef.Name,
					req.Name,
					crData.Spec.NextFunctions[nextListKey[0]].WBFunctionRef.Name,
					&nextConnectionCRName)
				logger.Info("nextConnectionCRName : " + nextConnectionCRName)

				// Get Kind of To side ConnectionCR
				err = r.GetConnectionData(ctx,
					req.Namespace,
					nextConnectionCRName,
					&toConnectionKind)
				if nil != err {
					if errors.IsNotFound(err) {
						// do nothing
					} else {
						logger.Error(err, "GetConnectionData() error")
					}
					// Requeue due to CR information acquisition error on source device
					return ctrl.Result{Requeue: true}, nil
				}
			}
		} else {
			toConnectionKind = "wb-end-of-chain"
		}

		switch fromConnectionKind {
		case examplecomv1.ConnectionCRKindPCIe:
			rxProtocol = "DMA"
		case examplecomv1.ConnectionCRKindEth:
			rxProtocol = "TCP"
		default:
			rxProtocol = ""
		}

		switch toConnectionKind {
		case examplecomv1.ConnectionCRKindPCIe:
			txProtocol = "DMA"
		case examplecomv1.ConnectionCRKindEth:
			txProtocol = "TCP"
		case "wb-end-of-chain":
			txProtocol = "RTP"
		default:
			txProtocol = ""
		}

		logger.Info("rxProtocol : " + rxProtocol)
		logger.Info("txProtocol : " + txProtocol)
		if "" == rxProtocol || "" == txProtocol {
			// Requeue until CR information for the source/destination device can be obtained
			return ctrl.Result{Requeue: true}, nil
		}

		switch rxProtocol {
		case "DMA":
			err = r.GetFunctionSpecData(ctx,
				crData.Spec.PreviousFunctions[previousListKey[0]].WBFunctionRef,
				&fromFunctionData, &fromFunctionKind)
			if nil != err {
				if errors.IsNotFound(err) {
					// do nothing
				} else {
					logger.Error(err, "GetFunctionSpecData() error")
				}
				// Requeue until CR information of the source device is available
				return ctrl.Result{Requeue: true}, nil
			}
		case "TCP":
			// do nothing
		default:
			// do nothing
		}

		var configData examplecomv1.GPUFuncConfig
		// Get config information
		err = r.GetConfigData(ctx, crData.Spec.ConfigName, rxProtocol, txProtocol, &configData)
		if nil != err {
			// If there is no Config information for GPUFunc
			return ctrl.Result{}, err
		}

		// Get NetworkAttachmentDefinition information
		var cniResourceName string
		if "sriov" == configData.VirtualNetworkDeviceDriverType && 0 == len(configData.IPAM) {
			err = r.GetNetworkAttachmentDefinitionData(ctx,
				req.Namespace,
				crData.Spec.NodeName,
				configData.VirtualNetworkDeviceDriverType,
				&cniResourceName)
			if nil != err {
				return ctrl.Result{}, err
			}
		}

		var containerName string
		containerName = "gpu-container"

		// Get container image information
		imageInfo = configData.ImageURI

		// env information sorting
		/* In the October FY23 version, there are no Envs on the Config template
		if len(configData.Template.Spec.Containers[0].Env) != 0 {
			podEnvData = configData.Template.Spec.Containers[0].Env
		}
		*/
		var frameSizeWidth string
		var frameSizeHeight string
		if len(configData.Envs) != 0 {
			for key, val := range configData.Envs {
				var configEnvData corev1.EnvVar
				configEnvData.Name = key
				configEnvData.Value = val
				podEnvData = append(podEnvData, configEnvData)
				if key == examplecomv1.Width {
					frameSizeWidth = val
				}
				if key == examplecomv1.Height {
					frameSizeHeight = val
				}
			}
		}
		if len(crData.Spec.Envs) != 0 {
			for listCount := 0; listCount < len(crData.Spec.Envs); listCount++ {
				if *crData.Spec.AcceleratorIDs[listCount].PartitionName ==
					req.Name {
					envDatalist := crData.Spec.Envs[listCount].EachEnv
					for dataIndex := 0; dataIndex < len(envDatalist); dataIndex++ {
						var crEnvData corev1.EnvVar
						crEnvData.Name = envDatalist[dataIndex].EnvKey
						crEnvData.Value = envDatalist[dataIndex].EnvValue
						podEnvData = append(podEnvData, crEnvData)
					}
					break
				}
			}
		}
		if len(crData.Spec.AcceleratorIDs) != 0 {
			var gpuFunctionEnvData corev1.EnvVar
			gpuFunctionEnvData.Name = "CUDA_VISIBLE_DEVICES"
			gpuFunctionEnvData.Value = crData.Spec.AcceleratorIDs[0].ID
			podEnvData = append(podEnvData, gpuFunctionEnvData)
			podGPUID = crData.Spec.AcceleratorIDs[0].ID
		}
		switch rxProtocol {
		case "DMA":
			if len(fromFunctionData.AcceleratorIDs) != 0 {
				var fpgaFunctionEnvData corev1.EnvVar
				fpgaFunctionEnvData.Name = "FPGA_DEV"
				fpgaFunctionEnvData.Value = fromFunctionData.AcceleratorIDs[0].ID
				podEnvData = append(podEnvData, fpgaFunctionEnvData)
			}

			if len(crData.Spec.SharedMemory.FilePrefix) != 0 {
				var gpuFunctionEnvData corev1.EnvVar
				gpuFunctionEnvData.Name = "FILE_PREFIX"
				gpuFunctionEnvData.Value = crData.Spec.SharedMemory.FilePrefix
				podEnvData = append(podEnvData, gpuFunctionEnvData)
			}

			if len(crData.Spec.SharedMemory.CommandQueueID) != 0 {
				var gpuFunctionEnvData corev1.EnvVar
				gpuFunctionEnvData.Name = "CONNECTOR_ID"
				gpuFunctionEnvData.Value = crData.Spec.SharedMemory.CommandQueueID
				podEnvData = append(podEnvData, gpuFunctionEnvData)
			}
		case "TCP":
			// do nothing
		default:
			// do nothing
		}

		var inputIP string
		var inputPort string
		if crData.Spec.Params[examplecomv1.InputIP].StrVal != "" {
			inputIP = crData.Spec.Params[examplecomv1.InputIP].StrVal
		}
		if crData.Spec.Params[examplecomv1.InputPort].IntVal >= 0 {
			inputPort = strconv.Itoa(int(crData.Spec.Params[examplecomv1.InputPort].IntVal))
		}

		var outputIP string
		var outputPort string
		var outputMAC string
		if crData.Spec.Params[examplecomv1.OutputIP].StrVal != "" {
			var outputEnvDataIPAddress corev1.EnvVar
			outputEnvDataIPAddress.Name = "RECEIVING_SERVER_IP"
			outputEnvDataIPAddress.Value =
				crData.Spec.Params[examplecomv1.OutputIP].StrVal
			outputIP = crData.Spec.Params[examplecomv1.OutputIP].StrVal
			podEnvData = append(podEnvData, outputEnvDataIPAddress)
		}
		if crData.Spec.Params[examplecomv1.OutputPort].IntVal >= 0 {
			var outputEnvDataPort corev1.EnvVar
			outputEnvDataPort.Name = "RECEIVING_SERVER_PORT"
			outputEnvDataPort.Value =
				strconv.Itoa(int(crData.Spec.Params[examplecomv1.OutputPort].IntVal))
			outputPort =
				strconv.Itoa(int(crData.Spec.Params[examplecomv1.OutputPort].IntVal))
			podEnvData = append(podEnvData, outputEnvDataPort)
		}
		if crData.Spec.Params[examplecomv1.OutputMAC].StrVal != "" {
			var outputEnvDataMACAddress corev1.EnvVar
			outputEnvDataMACAddress.Name = "RECEIVING_SERVER_MAC"
			outputEnvDataMACAddress.Value =
				crData.Spec.Params[examplecomv1.OutputMAC].StrVal
			outputMAC = crData.Spec.Params[examplecomv1.OutputMAC].StrVal
			podEnvData = append(podEnvData, outputEnvDataMACAddress)
		}

		// Args conversion process
		if len(configData.Template.Spec.Containers[0].Args) != 0 {
			var argsTempData []string
			// Iterate through the list of Args
			for argsIndex := 0; argsIndex < len(configData.Template.Spec.Containers[0].Args); argsIndex++ {
				strData := configData.Template.Spec.Containers[0].Args[argsIndex]
				replaceStr := r.ReplaceArgsData(ctx, strData, frameSizeWidth, frameSizeHeight,
					inputIP, inputPort, outputIP, outputPort, outputMAC)
				argsTempData = append(argsTempData, replaceStr)
			}
			var stringJoin = make([]byte, 0, 1000)
			for cntData := 0; cntData < len(argsTempData); cntData++ {
				stringJoin = append(stringJoin, argsTempData[cntData]...)
				stringJoin = append(stringJoin, ' ')
			}
			podArgsData = append(podArgsData, string(stringJoin))
		}

		// Generate Pod creation information: containers and later
		if len(configData.Template.Spec.Containers) != 0 {
			for containerIndex := 0; containerIndex < len(configData.Template.Spec.Containers); containerIndex++ {
				cfgData := configData.Template.Spec.Containers[containerIndex]
				var tempContainer corev1.Container
				tempContainer.Name =
					containerName + strconv.Itoa(containerIndex)
				switch {
				case cfgData.Image != "":
					tempContainer.Image = cfgData.Image
				default:
					tempContainer.Image = imageInfo
				}
				tempContainer.ImagePullPolicy = corev1.PullIfNotPresent
				tempContainer.Env = podEnvData
				tempContainer.Env = append(tempContainer.Env, corev1.EnvVar{
					Name: "K8S_POD_NAMESPACE",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.namespace",
						},
					},
				})
				tempContainer.Env = append(tempContainer.Env, corev1.EnvVar{
					Name: "K8S_POD_NAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				})
				tempContainer.Command =
					append(tempContainer.Command, cfgData.Command...)
				tempContainer.Args =
					append(tempContainer.Args, podArgsData...)
				tempContainer.WorkingDir = cfgData.WorkingDir
				tempContainer.SecurityContext = cfgData.SecurityContext
				tempContainer.VolumeMounts =
					append(tempContainer.VolumeMounts, cfgData.VolumeMounts...)
				if "sriov" == configData.VirtualNetworkDeviceDriverType {
					if nil == cfgData.Resources.Requests {
						cfgData.Resources.Requests = make(corev1.ResourceList)
					}
					if nil == cfgData.Resources.Limits {
						cfgData.Resources.Limits = make(corev1.ResourceList)
					}
					if 0 == len(configData.IPAM) {
						cfgData.Resources.Requests[corev1.ResourceName(cniResourceName)] =
							resource.Quantity(resource.MustParse("1"))
						cfgData.Resources.Limits[corev1.ResourceName(cniResourceName)] =
							resource.Quantity(resource.MustParse("1"))
					} else if cfgData.Name != nil && *cfgData.Name != "sidecar" {
						cfgData.Resources.Requests[corev1.ResourceName("nvidia.com/mlnx_sriov_netdevice")] =
							resource.Quantity(resource.MustParse("1"))
						cfgData.Resources.Limits[corev1.ResourceName("nvidia.com/mlnx_sriov_netdevice")] =
							resource.Quantity(resource.MustParse("1"))
					}
					tempContainer.Resources = cfgData.Resources
				} else {
					tempContainer.Resources = cfgData.Resources
				}

				tempContainer.Ports = cfgData.Ports

				podCRData.Containers =
					append(podCRData.Containers, tempContainer)
			}
		}
		if len(configData.Template.Spec.Volumes) != 0 {
			podCRData.Volumes = configData.Template.Spec.Volumes
		}
		podCRData.RestartPolicy = configData.Template.Spec.RestartPolicy
		podCRData.HostNetwork = configData.Template.Spec.HostNetwork
		podCRData.HostIPC = configData.Template.Spec.HostIPC
		if configData.Template.Spec.ShareProcessNamespace {
			b := true
			podCRData.ShareProcessNamespace = &b
		}
		// Add affinity

		// Create Pod creation information: affinity
		nodeSelectorRequirement := corev1.NodeSelectorRequirement{}
		nodeSelectorRequirement.Key = "kubernetes.io/hostname"
		nodeSelectorRequirement.Operator = corev1.NodeSelectorOpIn
		nodeSelectorRequirement.Values = make([]string, 1)
		nodeSelectorRequirement.Values[0] = crData.Spec.NodeName

		nodeSelectorTerm := corev1.NodeSelectorTerm{}
		nodeSelectorTerm.MatchExpressions =
			make([]corev1.NodeSelectorRequirement, 1)
		nodeSelectorTerm.MatchExpressions[0] = nodeSelectorRequirement

		nodeSelector := corev1.NodeSelector{}
		nodeSelector.NodeSelectorTerms =
			make([]corev1.NodeSelectorTerm, 1)
		nodeSelector.NodeSelectorTerms[0] = nodeSelectorTerm

		nodeAffinity := corev1.NodeAffinity{}
		nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution =
			&nodeSelector

		affinity := corev1.Affinity{}
		affinity.NodeAffinity = &nodeAffinity

		podCRData.Affinity = &affinity

		// Annotation settings
		if true == configData.AdditionalNetwork {
			podData = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: strings.ToLower(req.Name +
						"-mps-dgpu-" + podGPUID + "-pod"),
					Namespace: req.Namespace,
					Annotations: map[string]string{
						"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"" +
							strings.ToLower(crData.Spec.NodeName+"-config-net-"+
								configData.VirtualNetworkDeviceDriverType) +
							"\",\"ips\": " + strings.ToLower("[\""+
							crData.Spec.Params[examplecomv1.MyIP].StrVal+"\"]") +
							" } ]",
					},
				},
				Spec: podCRData,
			}

			if "sriov" == configData.VirtualNetworkDeviceDriverType && 0 < len(configData.IPAM) {
				podData.Annotations["k8s.v1.cni.cncf.io/networks"] = configData.IPAM[0]
				podData.Annotations["ethernet.swb.example.com/network"] = "sriov"
			}
		} else if mults := os.Getenv("K8S_POD_ANNOTATION"); mults != "" {
			podData = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: strings.ToLower(req.Name +
						"-mps-dgpu-" + podGPUID + "-pod"),
					Namespace: req.Namespace,
					Annotations: map[string]string{
						"k8s.v1.cni.cncf.io/networks": "ipvlan",
					},
				},
				Spec: podCRData,
			}
		} else {
			podData = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: strings.ToLower(req.Name +
						"-mps-dgpu-" + podGPUID + "-pod"),
					Namespace: req.Namespace,
				},
				Spec: podCRData,
			}
		}
		for k, v := range configData.Annotations {
			if podData.Annotations == nil {
				podData.Annotations = map[string]string{}
			}
			podData.Annotations[k] = v
		}
		podData.Labels = map[string]string{}
		podData.Labels["swb/func-type"] = "gpufunc"
		podData.Labels["swb/func-name"] = strings.ToLower(req.Name + "-gpu-pod")
		for k, v := range configData.Labels {
			podData.Labels[k] = v
		}

		// Set the parent-child relationship between CR and pod
		if err = ctrl.SetControllerReference(&crData,
			podData,
			r.Scheme); err != nil {
			logger.Error(err, "ctrl.SetControllerReference err")
			return ctrl.Result{}, err
		}
		// Pass Pod creation information to api-server
		err = r.Create(ctx, podData)
		if err != nil {
			logger.Error(err, "Failed to create Pod.")
		} else {
			logger.Info("Success to create Pod.")

			// crData.Status.IPAddress = crIP // @TODO to confirm

			/* // @TODO Needs consideration
			if _, ok := status["0"]; !ok {
			    status["0"] = make(map[string]string)
			}

			status["contena_id"][
			    strconv.Itoa(crData.Spec.FuncCHID)] = STATUS_OK
			crData.Status.AccStatus = status
			*/
			r.UpdCustomResource(ctx, &crData, RUNNING,
				podCRData.Containers,
				*configData.RxProtocol, *configData.TxProtocol,
				configData.VirtualNetworkDeviceDriverType,
				configData.AdditionalNetwork)
		}
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Create", "Create End")

	} else if eventKind == UPDATE {
		// In case of update
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Update", "Update Start")
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Update", "Update End")
	} else if eventKind == DELETE {
		// In case of deletion
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Delete", "Delete Start")

		FunctionIndexEnable[crData.Spec.DeviceType][*crData.Status.FunctionIndex] = true

		// Delete the Finalizer statement.
		err = r.DelCustomResource(ctx, &crData)
		if err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Delete", "Delete End")
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GPUFunctionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.GPUFunction{}).
		Complete(r)
}

func (r *GPUFunctionReconciler) GetfinalizerName(pCRData *examplecomv1.GPUFunction) string {
	/* Value to set in finalizer */
	gvks, _, _ := r.Client.Scheme().ObjectKinds(pCRData)
	return strings.ToLower(gvks[0].Kind) + ".finalizers." +
		strings.ToLower(gvks[0].Group+"."+gvks[0].Version)
}

func (r *GPUFunctionReconciler) GetEventKind(pCRData *examplecomv1.GPUFunction) int {
	var eventKind int
	eventKind = UPDATE
	/* Whether DeletionTimestamp is enabled or not */
	if !pCRData.ObjectMeta.DeletionTimestamp.IsZero() {
		eventKind = DELETE
	} else if !controllerutil.ContainsFinalizer(pCRData, r.GetfinalizerName(pCRData)) {
		/* Whether or not Finalizer is written*/
		eventKind = CREATE
	}
	return eventKind
}

func setFunctionIndex(devicetype string) (functionIndex int32) {

	functionIndex = -1
	var i int32

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		if _, ok := FunctionIndexEnable[devicetype][int32(0)]; !ok {
			FunctionIndexEnable[devicetype] = make(map[int32]bool)
			FunctionIndexEnable[devicetype][int32(0)] = false
			functionIndex = int32(0)
			break
		}
		for {
			if _, ok := FunctionIndexEnable[devicetype][i]; ok {
				if true != FunctionIndexEnable[devicetype][i] {
					i += 1
					continue
				} else {
					functionIndex = i
					FunctionIndexEnable[devicetype][i] = false
					break
				}
			} else {
				functionIndex = i
				FunctionIndexEnable[devicetype][i] = false
				break
			}
		}
	}
	return functionIndex
}

func (r *GPUFunctionReconciler) UpdCustomResource(ctx context.Context,
	pCRData *examplecomv1.GPUFunction, status string,
	podContainers []corev1.Container,
	rxProtocol string,
	txProtocol string,
	virtualNetworkDeviceDriverType string,
	additionalNetwork bool) error {
	logger := log.FromContext(ctx)
	var err error
	var functionIndex int32

	if status == RUNNING {
		if nil != pCRData.Spec.FunctionIndex {
			if _, ok := FunctionIndexEnable[pCRData.Spec.DeviceType][int32(0)]; !ok {
				FunctionIndexEnable[pCRData.Spec.DeviceType] = make(map[int32]bool)
			}
			functionIndex = *pCRData.Spec.FunctionIndex
			FunctionIndexEnable[pCRData.Spec.DeviceType][functionIndex] = false
		} else {
			functionIndex = setFunctionIndex(pCRData.Spec.DeviceType)
		}

		// Write a Finalizer
		controllerutil.AddFinalizer(pCRData, r.GetfinalizerName(pCRData))
		// status update
		pCRData.Status.StartTime = metav1.Now()
		pCRData.Status.Status = status
		// addition
		pCRData.Status.DataFlowRef = pCRData.Spec.DataFlowRef
		pCRData.Status.FunctionName = pCRData.Spec.FunctionName
		pCRData.Status.ImageURI = podContainers[0].Image
		pCRData.Status.SharedMemory = pCRData.Spec.SharedMemory
		pCRData.Status.RxProtocol = &rxProtocol
		pCRData.Status.TxProtocol = &txProtocol
		pCRData.Status.ConfigName = pCRData.Spec.ConfigName
		pCRData.Status.VirtualNetworkDeviceDriverType = virtualNetworkDeviceDriverType
		pCRData.Status.AdditionalNetwork = &additionalNetwork
		pCRData.Status.FunctionIndex = &functionIndex
	}
	err = r.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "Status Update Error.")
	} else {
		logger.Info("Status Update.")
	}
	return err
}

func (r *GPUFunctionReconciler) DelCustomResource(ctx context.Context,
	pCRData *examplecomv1.GPUFunction) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	/* Delete the Finalizer description */
	if controllerutil.ContainsFinalizer(pCRData,
		r.GetfinalizerName(pCRData)) {
		controllerutil.RemoveFinalizer(pCRData, r.GetfinalizerName(pCRData))
		err := r.Update(ctx, pCRData)
		if err != nil {
			logger.Error(err, "RemoveFinalizer Update Error.")
		}
	}
	return err
}

// Function Data Get
func (r *GPUFunctionReconciler) GetFunctionSpecData(
	ctx context.Context,
	functionCRName examplecomv1.WBNamespacedName,
	functionCRSpecData *examplecomv1.FunctionData,
	functionKind *string) error {
	logger := log.FromContext(ctx)

	var err error

	var existCount int = 0
	var notFoundErr error
	var elseErr error
	var strmapFunctionCRData map[string]interface{}

	var functionKindList []string = []string{examplecomv1.FunctionCRKindFPGA,
		examplecomv1.FunctionCRKindGPU,
		examplecomv1.FunctionCRKindGATE,
		examplecomv1.FunctionCRKindCPU}

	fcrData := &unstructured.Unstructured{}

	for n := 0; n < len(functionKindList); n++ {
		fcrData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    functionKindList[n],
		})
		err = r.Get(ctx, client.ObjectKey{
			Namespace: functionCRName.Namespace,
			Name:      functionCRName.Name}, fcrData)
		if errors.IsNotFound(err) {
			existCount += 1
			notFoundErr = err
		} else if err != nil {
			if examplecomv1.FunctionCRKindGATE != functionKindList[n] {
				logger.Info("unable to fetch FunctionCR")
				elseErr = err
			} else {
				existCount += 1
			}
		} else {
			*functionKind = functionKindList[n]
			break
		}
		if existCount == len(functionKindList) {
			logger.Info("FunctionCR does not exist")
			err = notFoundErr
		} else {
			err = elseErr
		}
	}

	if err != nil {
		return err
	}

	if len(fcrData.Object) != 0 {
		// Store spec information
		strmapFunctionCRData, _, _ =
			unstructured.NestedMap(fcrData.Object, "spec")

		// Convert the obtained mapdata to byte type
		bytes, err := json.Marshal(strmapFunctionCRData)
		if err != nil {
			logger.Error(err, "unable to json.marshal")
			return err
		}
		// Replace with a struct
		err = json.Unmarshal(bytes, &functionCRSpecData)
		if err != nil {
			logger.Error(err, "unable to json.unmarshal")
			return err
		}
	}
	return nil
}

// Get Connection Data
func (r *GPUFunctionReconciler) GetConnectionData(
	ctx context.Context,
	connectionCRNamespace string,
	connectionCRName string,
	connectionKind *string) error {
	logger := log.FromContext(ctx)

	var err error

	var existCount int = 0
	var notFoundErr error
	var elseErr error
	var connectionKindList []string = []string{examplecomv1.ConnectionCRKindPCIe,
		examplecomv1.ConnectionCRKindEth}

	connectionCRData := &unstructured.Unstructured{}

	for n := 0; n < len(connectionKindList); n++ {
		connectionCRData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    connectionKindList[n],
		})
		err = r.Get(ctx, client.ObjectKey{
			Namespace: connectionCRNamespace,
			Name:      connectionCRName}, connectionCRData)
		if errors.IsNotFound(err) {
			existCount += 1
			notFoundErr = err
		} else if err != nil {
			logger.Info("unable to fetch ConnectionCR")
			elseErr = err
		} else {
			*connectionKind = connectionKindList[n]
			break
		}
		if existCount == len(connectionKindList) {
			logger.Info("ConnectionCR does not exist")
			err = notFoundErr
		} else {
			err = elseErr
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// Get Config information for GPUFunc
func (r *GPUFunctionReconciler) GetConfigData(
	ctx context.Context,
	configName string,
	rxProtocol string,
	txProtocol string,
	pConfigData *examplecomv1.GPUFuncConfig) error {

	// logger := log.FromContext(ctx)
	var err error
	var mapData map[string]interface{}
	var configSliceData []examplecomv1.GPUFuncConfig

	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ {
		tmpData := &unstructured.Unstructured{}
		tmpData.SetGroupVersionKind(schema.GroupVersionKind{
			Kind:    "ConfigMap",
			Version: "v1",
		})

		err = r.Get(ctx,
			client.ObjectKey{
				Namespace: "default",
				Name:      configName,
			}, tmpData)

		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			break
		} else {
			mapData, _, _ =
				unstructured.NestedMap(tmpData.Object, "data")
			for _, jsonRecord := range mapData {
				err = json.Unmarshal([]byte(jsonRecord.(string)), &configSliceData)
			}
			for _, configData := range configSliceData {
				cmRxProtocol := *configData.RxProtocol
				cmTxProtocol := *configData.TxProtocol
				if cmRxProtocol == rxProtocol && cmTxProtocol == txProtocol {
					*pConfigData = configData
				}
			}
		}
	}
	return err
}

// Generate ConnectionCR name
func MakeConnectionCRName(
	ctx context.Context,
	dataFlowName string,
	fromCRName string,
	toCRName string,
	connectionCRName *string) {

	// logger := log.FromContext(ctx)

	var concatenationFlag bool = false
	var fromFunctionName string
	var toFunctionName string

	fromCRNameSplitList := strings.Split(fromCRName, "-")
	for splitIndex := 0; splitIndex < len(fromCRNameSplitList); splitIndex++ {
		if true == concatenationFlag {
			fromFunctionName = fromFunctionName +
				fromCRNameSplitList[splitIndex] + "-"
		}

		if "wbfunction" == fromCRNameSplitList[splitIndex] {
			concatenationFlag = true
		}
	}
	concatenationFlag = false
	toCRNameSplitList := strings.Split(toCRName, "-")
	for splitIndex := 0; splitIndex < len(toCRNameSplitList); splitIndex++ {
		if "" != toFunctionName {
			toFunctionName = toFunctionName + "-"
		}

		if true == concatenationFlag {
			toFunctionName = toFunctionName +
				toCRNameSplitList[splitIndex]
		}

		if "wbfunction" == toCRNameSplitList[splitIndex] {
			concatenationFlag = true
		}
	}

	*connectionCRName = dataFlowName + "-wbconnection-" + fromFunctionName + toFunctionName
}

// Get the Annotations information for NetworkAttachmentDefinition
func (r *GPUFunctionReconciler) GetNetworkAttachmentDefinitionData(
	ctx context.Context,
	namespace string,
	nodeName string,
	networkAttachmentDeviceDriverType string,
	pCNIResourceName *string) error {

	// logger := log.FromContext(ctx)
	var err error
	var strmapMetadata map[string]interface{}
	var strmapAnnotations map[string]string

	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ {
		tmpData := &unstructured.Unstructured{}
		tmpData.SetGroupVersionKind(schema.GroupVersionKind{
			Kind:    "NetworkAttachmentDefinition",
			Version: "k8s.cni.cncf.io/v1",
		})

		err = r.Get(ctx,
			client.ObjectKey{
				Namespace: namespace,
				Name:      nodeName + "-config-net-" + networkAttachmentDeviceDriverType,
			}, tmpData)

		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			break
		} else if len(tmpData.Object) != 0 {
			// Store metadata information
			strmapMetadata, _, _ =
				unstructured.NestedMap(tmpData.Object, "metadata")
			strmapAnnotations, _, _ =
				unstructured.NestedStringMap(strmapMetadata, "annotations")

			*pCNIResourceName = strmapAnnotations["k8s.v1.cni.cncf.io/resourceName"]
		}
	}
	return err
}

// Args conversion process
func (r *GPUFunctionReconciler) ReplaceArgsData(
	ctx context.Context, strData string,
	frameSizeWidth string, frameSizeHeight string,
	inputIP string, inputPort string,
	outputIP string, outputPort string, outputMAC string) string {

	// Repeat for the number of characters in Args
	for index := 0; index <= len(strData); index++ {
		if false == strings.Contains(strData[index:], "%") {
			break
		}

		indexBefore := strings.Index(strData[index:], "%") + index
		indexAfter := strings.Index(strData[indexBefore+1:], "%") + indexBefore + 2

		var changeStr string
		var replaceStr string

		switch strData[indexBefore:indexAfter] {
		case examplecomv1.ArgsWidth:
			changeStr = examplecomv1.ArgsWidth
			replaceStr = examplecomv1.ChangeArgsWidth + frameSizeWidth
		case examplecomv1.ArgsHeight:
			changeStr = examplecomv1.ArgsHeight
			replaceStr = examplecomv1.ChangeArgsHeight + frameSizeHeight
		case examplecomv1.ArgsArpIP:
			changeStr = examplecomv1.ArgsArpIP
			replaceStr = outputIP
		case examplecomv1.ArgsInputIP:
			changeStr = examplecomv1.ArgsInputIP
			replaceStr = examplecomv1.ChangeArgsIP + inputIP
		case examplecomv1.ArgsInputPort:
			changeStr = examplecomv1.ArgsInputPort
			replaceStr = examplecomv1.ChangeArgsPort + inputPort
		case examplecomv1.ArgsOutputIP:
			changeStr = examplecomv1.ArgsOutputIP
			replaceStr = examplecomv1.ChangeArgsIP + outputIP
		case examplecomv1.ArgsOutputPort:
			changeStr = examplecomv1.ArgsOutputPort
			replaceStr = examplecomv1.ChangeArgsPort + outputPort
		case examplecomv1.ArgsArpMAC:
			changeStr = examplecomv1.ArgsArpMAC
			replaceStr = outputMAC
		default:
			// do nothing
		}
		strData = strings.ReplaceAll(strData, changeStr, replaceStr)
		index = indexAfter
	}
	return strData
}
