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

	examplecomv1 "CPUFunction/api/v1"
	/* Additional files */
	"encoding/json"
	"os"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	/* Additional files end here */)

// CPUFunctionReconciler reconciles a CPUFunction object
type CPUFunctionReconciler struct {
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

var FunctionIndexEnable = map[int32]bool{}

// +kubebuilder:rbac:groups=example.com,resources=cpufunctions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=cpufunctions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.com,resources=cpufunctions/finalizers,verbs=update
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CPUFunction object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *CPUFunctionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var crData examplecomv1.CPUFunction
	var podEnvData []corev1.EnvVar
	var eventKind int // 0:Add, 1:Upd,  2:Del
	var imageInfo string
	var podData *corev1.Pod
	var podCRData corev1.PodSpec
	var podArgsData []string

	// Get CR information
	err := r.Get(ctx, req.NamespacedName, &crData)
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
		var toFunctionData examplecomv1.FunctionData
		var toFunctionKind string
		var previousListKey []string
		var nextListKey []string

		var previousConnectionCRName string
		var nextConnectionCRName string
		var fromConnectionKind string
		var toConnectionKind string

		var rxProtocol string
		var txProtocol string

		if "cpu-decode" != crData.Spec.FunctionName {
			for key, _ := range crData.Spec.PreviousFunctions {
				previousListKey = append(previousListKey, key)
			}
			if len(crData.Spec.PreviousFunctions) == 0 {
				fromConnectionKind = "wb-start-of-chain"
			} else {
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
		} else {
			fromConnectionKind = "wb-start-of-chain"
		}

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
		switch fromConnectionKind {
		case examplecomv1.ConnectionCRKindPCIe:
			rxProtocol = "DMA"
		case examplecomv1.ConnectionCRKindEth:
			rxProtocol = "TCP"
		case "wb-start-of-chain":
			rxProtocol = "RTP"
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

		// Get destination/source FunctionCR
		switch crData.Spec.FunctionName {
		case "cpu-decode":
			switch txProtocol {
			case "DMA":
				err = r.GetFunctionSpecData(ctx,
					crData.Spec.NextFunctions[nextListKey[0]].WBFunctionRef,
					&toFunctionData,
					&toFunctionKind)
				if nil != err {
					if errors.IsNotFound(err) {
						// do nothing
					} else {
						logger.Error(err, "GetFunctionSpecData() error")
					}
					// Requeue until the CR information of the connected device can be obtained
					return ctrl.Result{Requeue: true}, nil
				}
			case "TCP":
				// do nothing
			default:
				// do nothing
			}
		case "copy-branch":
			// do nothing
		case "glue-fdma-to-tcp":
			err = r.GetFunctionSpecData(ctx,
				crData.Spec.PreviousFunctions[previousListKey[0]].WBFunctionRef,
				&fromFunctionData,
				&fromFunctionKind)
			if nil != err {
				if errors.IsNotFound(err) {
					// do nothing
				} else {
					logger.Error(err, "GetFunctionSpecData() error")
				}
				// Requeue until CR information of the source device is available
				return ctrl.Result{Requeue: true}, nil
			}
		case "cpu-filter-resize-high-infer", "cpu-filter-resize-low-infer":
			// do nothing
		default:
			// do nothing
		}

		// Get config information
		var configData examplecomv1.CPUFuncConfig
		err = r.GetConfigData(ctx, crData.Spec.ConfigName, rxProtocol, txProtocol, &configData)
		if nil != err {
			// If there is no Config information for CPUFunc
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
		containerName = "cpu-container"

		// Get container image information
		imageInfo = configData.ImageURI

		// env information sorting
		/* In the October FY23 version, there are no Envs on the Config template
		if len(configData.Template.Spec.Containers[0].Env) != 0 {
		   podEnvData = configData.Template.Spec.Containers[0].Env
		}
		*/
		if len(configData.Envs) != 0 {
			for key, val := range configData.Envs {
				var configEnvData corev1.EnvVar
				configEnvData.Name = key
				configEnvData.Value = val
				podEnvData = append(podEnvData, configEnvData)
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

		switch crData.Spec.FunctionName {
		case "cpu-decode":
			if crData.Spec.Params[examplecomv1.InputIP].StrVal != "" {
				var inputEnvDataIPAddress corev1.EnvVar
				inputEnvDataIPAddress.Name = "DECENV_VIDEOSRC_IPA"
				inputEnvDataIPAddress.Value =
					crData.Spec.Params[examplecomv1.InputIP].StrVal
				podEnvData = append(podEnvData, inputEnvDataIPAddress)
			}

			if crData.Spec.Params[examplecomv1.InputPort].IntVal >= 0 {
				var inputEnvDataPort corev1.EnvVar
				inputEnvDataPort.Name = "DECENV_VIDEOSRC_PORT"
				inputEnvDataPort.Value =
					strconv.Itoa(int(crData.Spec.Params[examplecomv1.InputPort].IntVal))
				podEnvData = append(podEnvData, inputEnvDataPort)
			}

			// This passes from Params[examplecomv1.FPS] in CPUFunction CR
			// to the environment variables of the decoding module in Pod
			if crData.Spec.Params[examplecomv1.FPS].IntVal >= 0 {
				var inputEnvDataFPS corev1.EnvVar
				var decEnvFrameFPS int32
				decEnvFrameFPS = crData.Spec.Params[examplecomv1.FPS].IntVal
				inputEnvDataFPS.Name = "DECENV_FRAME_FPS"
				inputEnvDataFPS.Value =
					strconv.Itoa(int(decEnvFrameFPS))
				podEnvData = append(podEnvData, inputEnvDataFPS)
			}
			switch txProtocol {
			case "DMA":
				if len(toFunctionData.AcceleratorIDs) != 0 {
					var fpgaFunctionEnvData corev1.EnvVar
					fpgaFunctionEnvData.Name = "DECENV_FPGA_DEV_NAME"
					fpgaFunctionEnvData.Value = toFunctionData.AcceleratorIDs[0].ID
					podEnvData = append(podEnvData, fpgaFunctionEnvData)
				}

				if len(crData.Spec.SharedMemory.FilePrefix) != 0 {
					var fpgaFunctionEnvData corev1.EnvVar
					fpgaFunctionEnvData.Name = "DECENV_DPDK_FILE_PREFIX"
					fpgaFunctionEnvData.Value = crData.Spec.SharedMemory.FilePrefix
					podEnvData = append(podEnvData, fpgaFunctionEnvData)
				}

				if len(crData.Spec.SharedMemory.CommandQueueID) != 0 {
					var fpgaFunctionEnvData corev1.EnvVar
					fpgaFunctionEnvData.Name = "DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID"
					fpgaFunctionEnvData.Value = crData.Spec.SharedMemory.CommandQueueID
					podEnvData = append(podEnvData, fpgaFunctionEnvData)
				}
			case "TCP":
				if crData.Spec.Params[examplecomv1.OutputIP].StrVal != "" {
					var outputEnvDataIPAddress corev1.EnvVar
					outputEnvDataIPAddress.Name = "DECENV_OUTDST_IPA"
					outputEnvDataIPAddress.Value =
						crData.Spec.Params[examplecomv1.OutputIP].StrVal
					podEnvData = append(podEnvData, outputEnvDataIPAddress)
				}

				if crData.Spec.Params[examplecomv1.OutputPort].IntVal >= 0 {
					var outputEnvDataPort corev1.EnvVar
					outputEnvDataPort.Name = "DECENV_OUTDST_PORT"
					outputEnvDataPort.Value =
						strconv.Itoa(int(crData.Spec.Params[examplecomv1.OutputPort].IntVal))
					podEnvData = append(podEnvData, outputEnvDataPort)
				}
			default:
				// do nothing
			}
		case "copy-branch":
			// do nothing
		case "glue-fdma-to-tcp":
			if len(fromFunctionData.AcceleratorIDs) != 0 {
				var fpgaFunctionEnvData corev1.EnvVar
				fpgaFunctionEnvData.Name = "GLUEENV_FPGA_DEV_NAME"
				fpgaFunctionEnvData.Value = fromFunctionData.AcceleratorIDs[0].ID
				podEnvData = append(podEnvData, fpgaFunctionEnvData)
			}

			if len(crData.Spec.SharedMemory.FilePrefix) != 0 {
				var fpgaFunctionEnvData corev1.EnvVar
				fpgaFunctionEnvData.Name = "GLUEENV_DPDK_FILE_PREFIX"
				fpgaFunctionEnvData.Value = crData.Spec.SharedMemory.FilePrefix
				podEnvData = append(podEnvData, fpgaFunctionEnvData)
			}

			if len(crData.Spec.SharedMemory.CommandQueueID) != 0 {
				var fpgaFunctionEnvData corev1.EnvVar
				fpgaFunctionEnvData.Name = "GLUEENV_FPGA_DMA_DEV_TO_HOST_CONNECTOR_ID"
				fpgaFunctionEnvData.Value = crData.Spec.SharedMemory.CommandQueueID
				podEnvData = append(podEnvData, fpgaFunctionEnvData)
			}
		case "cpu-filter-resize-high-infer", "cpu-filter-resize-low-infer":
			if crData.Spec.Params[examplecomv1.InputPort].IntVal >= 0 {
				var inputEnvDataPort corev1.EnvVar
				inputEnvDataPort.Name = "FRENV_INPUT_PORT"
				inputEnvDataPort.Value =
					strconv.Itoa(int(crData.Spec.Params[examplecomv1.InputPort].IntVal))
				podEnvData = append(podEnvData, inputEnvDataPort)
			}

			if crData.Spec.Params[examplecomv1.OutputIP].StrVal != "" {
				var outputEnvDataIPAddress corev1.EnvVar
				outputEnvDataIPAddress.Name = "FRENV_OUTPUT_IP"
				outputEnvDataIPAddress.Value =
					crData.Spec.Params[examplecomv1.OutputIP].StrVal
				podEnvData = append(podEnvData, outputEnvDataIPAddress)
			}

			if crData.Spec.Params[examplecomv1.OutputPort].IntVal >= 0 {
				var outputEnvDataPort corev1.EnvVar
				outputEnvDataPort.Name = "FRENV_OUTPUT_PORT"
				outputEnvDataPort.Value =
					strconv.Itoa(int(crData.Spec.Params[examplecomv1.OutputPort].IntVal))
				podEnvData = append(podEnvData, outputEnvDataPort)
			}
		default:
			// do nothing
		}

		// Args conversion process
		if len(configData.Template.Spec.Containers[0].Args) != 0 {
			var argsTempData []string
			// Iterate through the list of Args
			for argsIndex := 0; argsIndex < len(configData.Template.Spec.Containers[0].Args); argsIndex++ {
				strData := configData.Template.Spec.Containers[0].Args[argsIndex]
				replaceStr := r.ReplaceArgsData(ctx, strData, crData.Spec.Params,
					crData.Spec.FunctionName, fromFunctionData.FunctionName, "", configData)
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
					Name:      strings.ToLower(req.Name + "-cpu-pod"),
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
		} else if mults := os.Getenv("K8S_POD_ANNOTATION_CPU"); mults != "" {
			podData = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.ToLower(req.Name + "-cpu-pod"),
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
					Name:      strings.ToLower(req.Name + "-cpu-pod"),
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
		podData.Labels["swb/func-type"] = "cpufunc"
		podData.Labels["swb/func-name"] = strings.ToLower(req.Name + "-cpu-pod")
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
			r.UpdCustomResource(ctx, &crData,
				RUNNING, podCRData.Containers,
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

		FunctionIndexEnable[*crData.Status.FunctionIndex] = true

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
func (r *CPUFunctionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.CPUFunction{}).
		Complete(r)
}

func (r *CPUFunctionReconciler) GetfinalizerName(pCRData *examplecomv1.CPUFunction) string {
	// Value to set in the finalizer
	gvks, _, _ := r.Client.Scheme().ObjectKinds(pCRData)
	return strings.ToLower(gvks[0].Kind) + ".finalizers." +
		strings.ToLower(gvks[0].Group+"."+gvks[0].Version)
}

func (r *CPUFunctionReconciler) GetEventKind(pCRData *examplecomv1.CPUFunction) int {
	var eventKind int
	eventKind = UPDATE
	// Whether or not there is a deletion timestamp
	if !pCRData.ObjectMeta.DeletionTimestamp.IsZero() {
		eventKind = DELETE
	} else if !controllerutil.ContainsFinalizer(pCRData, r.GetfinalizerName(pCRData)) {
		// Whether or not Finalizer is written
		eventKind = CREATE
	}
	return eventKind
}

func setFunctionIndex() (functionIndex int32) {

	functionIndex = -1
	var i int32

	if 0 != len(FunctionIndexEnable) {
		for i = 0; i < int32(len(FunctionIndexEnable)); i++ {
			if true != FunctionIndexEnable[i] {
				continue
			} else {
				functionIndex = i
				FunctionIndexEnable[i] = false
				break
			}
		}
		if functionIndex < 0 {
			functionIndex = int32(len(FunctionIndexEnable))
			FunctionIndexEnable[functionIndex] = false
		}
	} else {
		FunctionIndexEnable[int32(0)] = false
		functionIndex = int32(0)
	}
	return functionIndex
}

func (r *CPUFunctionReconciler) UpdCustomResource(ctx context.Context,
	pCRData *examplecomv1.CPUFunction, status string,
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
			functionIndex = *pCRData.Spec.FunctionIndex
			FunctionIndexEnable[functionIndex] = false
		} else {
			functionIndex = setFunctionIndex()
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

func (r *CPUFunctionReconciler) DelCustomResource(ctx context.Context,
	pCRData *examplecomv1.CPUFunction) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	// Delete the Finalizer statement.
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
func (r *CPUFunctionReconciler) GetFunctionSpecData(
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
func (r *CPUFunctionReconciler) GetConnectionData(
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

// Get Config information for CPUFunc
func (r *CPUFunctionReconciler) GetConfigData(
	ctx context.Context,
	configName string,
	rxProtocol string,
	txProtocol string,
	pConfigData *examplecomv1.CPUFuncConfig) error {
	logger := log.FromContext(ctx)
	logger.Info("debug::GetConfigData", "configName", configName, "rxProtocol", rxProtocol, "txProtocol", txProtocol)
	var err error
	var mapData map[string]interface{}
	var configSliceData []examplecomv1.CPUFuncConfig

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

// String concatenation of IP address and port number
func MakeForwardingInfo(ctx context.Context,
	branchIP intstr.IntOrString,
	branchPort intstr.IntOrString,
	forwarding *string) {

	// logger := log.FromContext(ctx)

	var ip string
	var port string
	var num int

	if "" != branchPort.StrVal {
		port = branchPort.StrVal
	} else if 0 <= branchPort.IntVal {
		port =
			strconv.Itoa(int(branchPort.IntVal))
	}

	if "" != branchIP.StrVal {
		ip = branchIP.StrVal
	}

	if "" != port && "" != ip {
		ipList := strings.Split(ip, ",")
		portList := strings.Split(port, ",")
		if len(ipList) > len(portList) {
			num = len(ipList)
		} else {
			num = len(portList)
		}
		for branchIndex := 0; branchIndex < num; branchIndex++ {
			if "" == ipList[branchIndex] || "" == portList[branchIndex] {
				break
			}

			if 0 != branchIndex {
				*forwarding = *forwarding + ","
			}
			*forwarding = *forwarding + ipList[branchIndex] + ":" +
				portList[branchIndex]
		}
	}
}

// Get the Annotations information for NetworkAttachmentDefinition
func (r *CPUFunctionReconciler) GetNetworkAttachmentDefinitionData(
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
func (r *CPUFunctionReconciler) ReplaceArgsData(
	ctx context.Context, strData string, params map[string]intstr.IntOrString,
	myFunctionName string, fromFunctionName string, toFunctionName string,
	configData examplecomv1.CPUFuncConfig) string {

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
		case examplecomv1.ArgsReceiving:
			if "" != params[examplecomv1.InputIP].StrVal &&
				0 <= params[examplecomv1.InputPort].IntVal {
				strData =
					params[examplecomv1.InputIP].StrVal +
						":" +
						strconv.Itoa(int(params[examplecomv1.InputPort].IntVal))
			}
		case examplecomv1.ArgsNum:
			if "" != params[examplecomv1.BranchOutputIP].StrVal {
				countNum := strings.Count(params[examplecomv1.BranchOutputIP].StrVal, ",")
				forwardingNum := strconv.Itoa(countNum + 1)
				strData = forwardingNum
			}
		case examplecomv1.ArgsForwarding:
			var forwarding string

			switch myFunctionName {
			case "copy-branch":
				MakeForwardingInfo(ctx,
					params[examplecomv1.BranchOutputIP],
					params[examplecomv1.BranchOutputPort],
					&forwarding)
			case "glue-fdma-to-tcp":
				MakeForwardingInfo(ctx,
					params[examplecomv1.GlueOutputIP],
					params[examplecomv1.GlueOutputPort],
					&forwarding)
			default:
				// do nothing
			}
			strData = forwarding
		case examplecomv1.ArgsMemSize:
			if "" != *configData.CopyMemorySize {
				strData = *configData.CopyMemorySize
			}
		case examplecomv1.ArgsWidth:
			switch fromFunctionName {
			case "filter-resize-high-infer":
				strData = "1280"
			case "filter-resize-low-infer":
				strData = "416"
			default:
				// do nothing
			}
		case examplecomv1.ArgsHeight:
			switch fromFunctionName {
			case "filter-resize-high-infer":
				strData = "1280"
			case "filter-resize-low-infer":
				strData = "416"
			default:
				// do nothing
			}
		default:
			// do nothing
		}
		strData = strings.ReplaceAll(strData, changeStr, replaceStr)
		index = indexAfter
	}
	return strData
}
