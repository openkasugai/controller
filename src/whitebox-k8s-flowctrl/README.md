# whitebox-k8s-flowctrl

## Controllers for the schedule function

| Name | File Name |
| -- | -- |
| Scheduler | wbscheduler_controller.go | 
| DataFlow Controller | dataflow_controller.go |
| FunctionTarget Controller | functiontarget_controller.go | 
| FunctionType Controller | functiontype_controller.go | 
| FunctionChain Controller | functionchain_controller.go | 

## whitebox-k8s-flowctrl operation check
  Start by running make install/make run, then create ComputeResource.
  Apply the YAML files for FunctionInfo (configmap), FunctionType, and FunctionChain for preliminary setup, and then apply the YAML for DataFlow to perform operation checks.

### Start the console
#### Install
    make install
    kubectl get crd
      ⇒ Verify that various CRDs are generated.

#### Boot
    make run

### Preliminary preparation

#### Namespace Registration
    kubectl create namespace test01
    kubectl create namespace wbfunc-imgproc
    kubectl create namespace chain-imgproc
    kubectl create namespace cluster01

#### FunctionTarget CR generation confirmation
    kubectl get FunctionTarget
      ⇒ Verify that FunctionTarget CR is generated.

#### Generate ConfigMap for FunctionInfo
    kubectl apply -f config/samples/ntt-hpc_v1_configmap_funcinfo.yaml
    kubectl get configmap -n wbfunc-imgproc
      ⇒ Verify that funcinfo-decode, funcinfo-filter-resize, and funcinfo-high-infer are generated.
    kubectl get configmap -n wbfunc-imgproc funcinfo-decode -o yaml
    kubectl get configmap -n wbfunc-imgproc funcinfo-filter-resize -o yaml
    kubectl get configmap -n wbfunc-imgproc funcinfo-high-infer -o yaml

#### FunctionType CR generation
    kubectl apply -f config/samples/ntt-hpc_v1_functionkind.yaml
    kubectl describe FunctionType -n wbfunc-imgproc
      ⇒ Verify that data is entered into Status.FunctionTargetKindCandidates.
        Verify that Status.Status is Ready.

#### FunctionChain CR Generation
    kubectl apply -f config/samples/ntt-hpc_v1_functionchain.yaml
    kubectl describe -n chain-imgproc FunctionChain
      ⇒ Verify that Status.Status is Ready.

### DataFlow CR Generation
    kubectl apply -f config/samples/ntt-hpc_v1_dataflow.yaml

### DataFlow/WBFunction/WBConnection CR generation confirmation
    kubectl get DataFlow -n test01
    kubectl get WBFunction -n test01
    kubectl get WBConnection -n test01
      ⇒ Verify that CRs for DataFlow, WBConnection, and WBFunction are generated.

    kubectl describe DataFlow -n test01
    kubectl describe WBFunction -n test01
    kubectl describe WBConnection -n test01
      ⇒ If WBFunction controller and WBConnection controller are not operating, confirm that DataFlow status is updated to "WBFunction/WBConnection created" by DataFlow controller.
        "" → "Scheduling in progress" → "WBFunction/WBConnection creation in progress" → "WBFunction/WBConnection created"      
      ⇒ If WBFunction controller and WBConnection controller are operating, confirm that DataFlow status is updated to "Deployed" by DataFlow controller.
        "" → "Scheduling in progress" → "WBFunction/WBConnection creation in progress" → "WBFunction/WBConnection created" → "Deployed"      

    kubectl get WBFunction -n test01 -o yaml
    kubectl get WBConnection -n test01 -o yaml
      ⇒ Output in YAML format

### Stop
    Stop with ctrl+c

### Uninstall
    make uninstall
    kubectl get crd
      ⇒ Verify that various CRDs are deleted.

### Delete ConfigMap
    kubectl delete configmap -n wbfunc-imgproc funcinfo-decode funcinfo-filter-resize funcinfo-high-infer
    kubectl get configmap -n wbfunc-imgproc
      ⇒ Verify that funcinfo-decode, funcinfo-filter-resize, and funcinfo-high-infer are deleted.

### Delete Namespace
    kubectl delete namespace test01 wbfunc-imgproc chain-imgproc cluster01
    kubectl get namespace
      ⇒ Verify that test01, wbfunc-imgproc, chain-imgproc, and cluster01 are deleted.

## Testing Scheduler and Filter

There are two types of tests: envtest, which performs validation by starting a local virtual k8s cluster, and deploytest, which performs validation in an actual k8s environment on a server.

The execution procedures for each are shown below.

### envtest 

Executable with the following command.

```bash
make test 
```

### deploytest

Executable with the following procedure.

1. Create docker image
```bash
make docker-build IMG="IMAGE_NAME"
```

2. Deploy the controller

```bash
make deploy IMG="IMAGE_NAME"
```

3. Execute the test

```bash
make deploytest
```

4. Undeploy the controller

```bash
make undeploy
```
