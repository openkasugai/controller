---
weight: 2
title: "argo.go改修"
---
# argo.go改修
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* `argo.go`を編集し、Custom Kubeflow SDKでデプロイする各リソースに合わせて編集する。

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [go.mod改修](../modify-go.mod)

## 手順
1. `argo.go`をテキストエディタで開く。
```
$ vi backend/src/v2/compiler/argocompiler/argo.go
```

2. importするパッケージを以下のように編集する。
```go
 package argocompiler
 
 import (
+       "encoding/json"
        "fmt"
+       "log"
+       "sort"
        "strings"
 
        wfapi "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
~~OMITTED~~
        k8score "k8s.io/api/core/v1"
        k8sres "k8s.io/apimachinery/pkg/api/resource"
        k8smeta "k8s.io/apimachinery/pkg/apis/meta/v1"
+       "gopkg.in/yaml.v2"
 )
```

3. 各種関数を定義する。
* componentとconnectionを判定する関数を定義する。(※1)
 ```go
        // TODO(Bobgy): add an option -- dev mode, ImagePullPolicy should only be Always in dev mode.
 }
 
+// connectionとcomponentを判定しキーをスライスで返す
+func JudgeComponentType(comps map[string]*pipelinespec.ComponentSpec) ([]string, []string) {
+       connect_keys := []string{}
+       component_keys := []string{}
+       for k, v := range comps {
+               if v.ConnectionKind != "" {
+                       connect_keys = append(connect_keys, k)
+               } else if v.StartPointIP != "" && v.StartPointProtocol != "" {
+                       connect_keys = append(connect_keys, k)
+               } else if v.EndPointIP != "" && v.EndPointProtocol != "" {
+                       connect_keys = append(connect_keys, k)
+               } else {
+                       component_keys = append(component_keys, k)
+               }
+       }
+       return connect_keys, component_keys
+}
+
```

* 特定のコンポーネント名が存在するかを確認する関数を定義する。
```go
       return connect_keys, component_keys
}

+// componentに特定のフィールド名が存在するかを判定
+func JudgeComponentName(comps []string, comp_name string) bool {
+       for _, k := range comps {
+               if k == comp_name {
+                       return true
+               }
+       }
+       return false
+}
```

* OpenKasugaiコントローラで作成する各リソースのTemplateを成形する関数を定義する。(※2)
```go
       }
       return false
}
+// dataflowTemplateを成形
+func CreateDataflowTemp(spec *pipelinespec.PipelineSpec) wfapi.Template {
+       type DataflowManifestMeta struct {
+               Name      string `yaml:"name"`
+               Namespace string `yaml:"namespace"`
+       }
+
+       type DataFlowManifestSpec struct {
+               FunctionChain          string `yaml:"functionChain"`
+               FunctionChainNamespace string `yaml:"functionChainNamespace"`
+               StartPointIP           string `yaml:"startPointIP"`
+               StartPointPort         int64  `yaml:"startPointPort"`
+               StartPointProtocol     string `yaml:"startPointProtocol"`
+               EndPointIP             string `yaml:"endPointIP"`
+               EndPointPort           int64  `yaml:"endPointPort"`
+               EndPointProtocol       string `yaml:"endPointProtocol"`
+       }
+
+       type DataFlowManifest struct {
+               ApiVersion string               `yaml:"apiVersion"`
+               Kind       string               `yaml:"kind"`
+               Metadata   DataflowManifestMeta `yaml:"metadata"`
+               Spec       DataFlowManifestSpec `yaml:"spec"`
+       }
+
+       dfm := DataFlowManifest{
+               ApiVersion: "example.com/v1",
+               Kind:       "DataFlow",
+               Metadata: DataflowManifestMeta{
+                       Name:      spec.Meta.DataflowName,
+                       Namespace: spec.Meta.DataflowNamespace,
+               },
+               Spec: DataFlowManifestSpec{
+                       FunctionChain:          spec.Meta.FunctionChainName,
+                       FunctionChainNamespace: spec.Meta.FunctionChainNamespace,
+                       StartPointIP:           spec.Components["start"].StartPointIP,
+                       StartPointPort:         spec.Components["start"].StartPointPort,
+                       StartPointProtocol:     spec.Components["start"].StartPointProtocol,
+                       EndPointIP:             spec.Components["end"].EndPointIP,
+                       EndPointPort:           spec.Components["end"].EndPointPort,
+                       EndPointProtocol:       spec.Components["end"].EndPointProtocol,
+               },
+       }
+
+       yaml_dfm, err := yaml.Marshal(dfm)
+       if err != nil {
+               log.Fatal(err)
+       }
+
+       dftemp := wfapi.Template{
+               Name: "dataflow",
+               Resource: &wfapi.ResourceTemplate{
+                       Action:   "create",
+                       Manifest: string(yaml_dfm),
+               },
+       }
+
+       return dftemp
+}
+
+// FunctionKindTemplateを成形
+func CreateFunctionKindTemp(spec *pipelinespec.PipelineSpec, comp_key string) wfapi.Template {
+
+       type FunctionKindManifestSpec struct {
+               Name                    string `yaml:"name"`
+               Function_info_name      string `yaml:"function_info_name"`
+               Function_info_namespace string `yaml:"function_info_namespace"`
+               Version                 string `yaml:"version"`
+       }
+
+       type FunctionKindManifestMeta struct {
+               Name      string `yaml:"name"`
+               Namespace string `yaml:"namespace"`
+       }
+
+       type FunctionKindManifest struct {
+               ApiVersion string                   `yaml:"apiVersion"`
+               Kind       string                   `yaml:"kind"`
+               Metadata   FunctionKindManifestMeta `yaml:"metadata"`
+               Spec       FunctionKindManifestSpec `yaml:"spec"`
+       }
+
+       fkmf := FunctionKindManifest{
+               ApiVersion: "example.com/v1",
+               Kind:       "FunctionKind",
+               Metadata: FunctionKindManifestMeta{
+                       Name:      "fk-" + comp_key,
+                       Namespace: spec.Meta.FunctionKindNamespace,
+               },
+               Spec: FunctionKindManifestSpec{
+                       Name:                    comp_key,
+                       Function_info_name:      "funcinfo-" + comp_key,
+                       Function_info_namespace: spec.Meta.FunctionKindNamespace,
+                       Version:                 spec.Components[comp_key].Version,
+               },
+       }
+
+       yaml_fkmf, err := yaml.Marshal(fkmf)
+       if err != nil {
+               log.Fatal(err)
+       }
+
+       fktemp := wfapi.Template{
+               Name: "fk-" + comp_key,
+               Resource: &wfapi.ResourceTemplate{
+                       Action:   "create",
+                       Manifest: string(yaml_fkmf),
+               },
+       }
+
+       return fktemp
+}
+
+// ConfigMapのTemplateを成形
+func CreateConfigMapTemp(spec *pipelinespec.PipelineSpec, comp_key string) wfapi.Template {
+
+       manifestJsonformat := make(map[string]interface{})
+       manifestJsonformat["items"] = make(map[string]interface{})
+
+       manifestYamlformat := make(map[string]interface{})
+       manifestYamlformat["apiVersion"] = "v1"
+       manifestYamlformat["kind"] = "ConfigMap"
+       manifestYamlformat["metadata"] = map[string]interface{}{
+               "name":      "funcinfo-" + comp_key,
+               "namespace": spec.Components[comp_key].Namespace,
+       }
+
+       innerData := make(map[string]interface{})
+
+       for inf := range spec.Components[comp_key].Info {
+               itemsMap := manifestJsonformat["items"].(map[string]interface{})
+               for eth := range spec.Components[comp_key].Info[inf].Items {
+                       item := spec.Components[comp_key].Info[inf].Items[eth]
+                       itemsMap[eth] = map[string]interface{}{
+                               "configName":        item.ConfigName,
+                               "coreMin":           item.CoreMin,
+                               "coreMax":           item.CoreMax,
+                               "totalBase":         item.TotalBase,
+                               "capacityTotalBase": item.CapacityTotalBase,
+                       }
+               }
+
+               // すべてのエントリーを追加した後でJSONに変換
+               jsonData, err := json.Marshal(manifestJsonformat)
+               if err != nil {
+                       log.Fatal(err)
+               }
+               innerData[inf] = string(jsonData)
+       }
+
+       manifestYamlformat["data"] = innerData
+
+       // mapをyamlに変換
+       yaml_cm, err := yaml.Marshal(manifestYamlformat)
+       if err != nil {
+               log.Fatal(err)
+       }
+
+       cmaptemp := wfapi.Template{
+               Name: "funcinfo-" + comp_key,
+               Resource: &wfapi.ResourceTemplate{
+                       Action:   "create",
+                       Manifest: string(yaml_cm),
+               },
+       }
+
+       return cmaptemp
+}
```

* Dagの依存関係を探索し取得する関数を定義する。
```go

       return cmaptemp
}
+func flatten(arrays []string) []string {
+       var result []string
+
+       for _, array := range arrays {
+               result = append(result, array)
+       }
+
+       return result
+}
+
+// 再起呼び出しでDagの依存関係を探索して取得
+func RegressionTasks(depends []string, spec *pipelinespec.PipelineSpec, result *[]map[string]interface{}) {
+       var dependmap = make(map[string]interface{})
+       froms := []string{}
+
+       //コネクタ名もしくはタスク名を配列で指定されているのでループ
+       for _, depend := range depends {
+               //引数のコネクタ/タスクの依存関係を取得
+               for _, from := range spec.Root.Dag.Tasks[depend].DependentTasks {
+                       dependmap["from"] = from
+                       dependmap["to"] = depend
+
+                       deepCopy := make(map[string]interface{})
+                       for key, value := range dependmap {
+                               deepCopy[key] = value
+                       }
+
+                       *result = append(*result, deepCopy)
+               }
+               for _, flat := range flatten(spec.Root.Dag.Tasks[depend].DependentTasks) {
+                       froms = append(froms, flat)
+               }
+
+       }
+
+       if len(froms) != 0 {
+               RegressionTasks(froms, spec, result)
+       }
+
+}
+func CreateFunctionChainTemp(spec *pipelinespec.PipelineSpec, comp_keys []string, connect_keys []string) wfapi.Template {
+
+       data := make(map[string]interface{})
+
+       data["apiVersion"] = "example.com/v1"
+       data["kind"] = "FunctionChain"
+       data["metadata"] = map[string]interface{}{
+               "name":      spec.Meta.FunctionChainName,
+               "namespace": spec.Meta.FunctionChainNamespace,
+       }
+
+       //functions生成
+       innerfunc := make(map[string]map[string]interface{})
+
+       for _, comp_key := range comp_keys {
+               if comp_key != "start" && comp_key != "end" {
+                       innerfunc[comp_key] = map[string]interface{}{
+                               "functionName": comp_key,
+                               "version":      spec.Components[comp_key].Version,
+                       }
+               }
+       }
+
+       //connections生成
+       dagDep := GetDagDepends(spec, connect_keys)
+
+       var con_slice []map[string]interface{}
+       for _, connect_key := range connect_keys {
+               if connect_key != "start" && connect_key != "end" {
+                       con_slice = append(con_slice, map[string]interface{}{
+                               "to":             dagDep[connect_key]["to"],
+                               "from":           dagDep[connect_key]["from"],
+                               "fromPort":       spec.Components[connect_key].FromPort,
+                               "toPort":         spec.Components[connect_key].ToPort,
+                               "connectionKind": spec.Components[connect_key].ConnectionKind,
+                       })
+               }
+       }
+
+       data["spec"] = map[string]interface{}{
+               "functionKindNamespace":   spec.Meta.FunctionKindNamespace,
+               "connectionKindNamespace": spec.Meta.ConnectionKindNamespace,
+               "functions":               innerfunc,
+               "connections":             con_slice,
+       }
+
+       // オブジェクトをyamlに変換
+       yaml_cm, err := yaml.Marshal(data)
+       if err != nil {
+               log.Fatal(err)
+       }
+
+       fctemp := wfapi.Template{
+               Name: "functionchain",
+               Resource: &wfapi.ResourceTemplate{
+                       Action:   "create",
+                       Manifest: string(yaml_cm),
+               },
+       }
+
+       return fctemp
+}
```

* Dagの依存関係情報を取得する関数を定義する。
```go

       return fctemp
}
+func GetDagDepends(spec *pipelinespec.PipelineSpec, connect_keys []string) map[string]map[string]interface{} {
+       endDepend := []string{"end"}
+
+       result_map := make(map[string]map[string]interface{})
+
+       depends := []map[string]interface{}{}
+
+       RegressionTasks(endDepend, spec, &depends)
+       for _, connect_key := range connect_keys {
+               inner_map := make(map[string]interface{})
+               for _, depend := range depends {
+                       if depend["from"] == connect_key {
+                               //toのstart・endをwb-start-of-chain・wb-end-of-chainに変換
+                               if depend["to"] == "start" {
+                                       depend["to"] = "wb-start-of-chain"
+                               } else if depend["to"] == "end" {
+                                       depend["to"] = "wb-end-of-chain"
+                               }
+
+                               inner_map["to"] = depend["to"]
+                               result_map[connect_key] = inner_map
+                       } else if depend["to"] == connect_key {
+                               //fromのstart・endをwb-start-of-chain・wb-end-of-chainに変換
+                               if depend["from"] == "start" {
+                                       depend["from"] = "wb-start-of-chain"
+                               } else if depend["from"] == "end" {
+                                       depend["from"] = "wb-end-of-chain"
+                               }
+
+                               inner_map["from"] = depend["from"]
+                               result_map[connect_key] = inner_map
+                       }
+               }
+
+       }
+
+       return result_map
+}
```

* 各リソースのデプロイ順番を決定するためにソートする処理を定義する。
```go

       return result_map
}
+// slice内に指定の値があるかチェックし要素数を返す
+// deploy順番をソートするために使用
+func index(slice []string, item string) int {
+       for i := range slice {
+               //文字列が前方一致したものを判定
+               if strings.HasPrefix(item, slice[i]) {
+                       return i
+               }
+       }
+       return -1
+}
+
```

4. Compile関数の中を以下のように編集する。(※3)
```go
        if spec.GetPipelineInfo().GetName() == "" {
                return nil, fmt.Errorf("pipelineInfo.name is empty")
        }
-       deploy, err := compiler.GetDeploymentConfig(spec)
-       if err != nil {
-               return nil, err
~~OMITTED~~
-                       return nil, fmt.Errorf("bug: cloned Kubernetes spec message does not have expected type")
-               }
-       }
+
+       //Components内のconnectionとcomponentを仕分ける
+       connect_keys, component_keys := JudgeComponentType(spec.Components)
 
-       // initialization
        wf := &wfapi.Workflow{
                TypeMeta: k8smeta.TypeMeta{
                        APIVersion: "argoproj.io/v1alpha1",
                        Kind:       "Workflow",
                },
                ObjectMeta: k8smeta.ObjectMeta{
-                       GenerateName: retrieveLastValidString(spec.GetPipelineInfo().GetName()) + "-",
-                       // Note, uncomment the following during development to view argo inputs/outputs in KFP UI.
-                       // TODO(Bobgy): figure out what annotations we should use for v2 engine.
-                       // For now, comment this annotation, so that in KFP UI, it shows argo input/output params/artifacts
-                       // suitable for debugging.
-                       //
-                       // Annotations: map[string]string{
-                       //      "pipelines.kubeflow.org/v2_pipeline": "true",
-                       // },
+                       GenerateName: "swb-dataflow-",
                },
                Spec: wfapi.WorkflowSpec{
-                       PodMetadata: &wfapi.Metadata{
-                               Annotations: map[string]string{
-                                       "pipelines.kubeflow.org/v2_component": "true",
-                               },
-                               Labels: map[string]string{
-                                       "pipelines.kubeflow.org/v2_component": "true",
+                       Entrypoint: "deploy",
+                       Templates: []wfapi.Template{
+                               {
+                                       Name: "deploy",
+                                       Steps: []wfapi.ParallelSteps{
+                                               {
+                                                       Steps: []wfapi.WorkflowStep{},
+                                               },
+                                       },
                                },
                        },
-                       ServiceAccountName: "pipeline-runner",
-                       Entrypoint:         tmplEntrypoint,
                },
        }
-       c := &workflowCompiler{
-               wf:        wf,
-               templates: make(map[string]*wfapi.Template),
-               // TODO(chensun): release process and update the images.
-               driverImage:   "gcr.io/ml-pipeline/kfp-driver@sha256:8e60086b04d92b657898a310ca9757631d58547e76bbbb8bfc376d654bef1707",
-               launcherImage: "gcr.io/ml-pipeline/kfp-launcher@sha256:50151a8615c8d6907aa627902dce50a2619fd231f25d1e5c2a72737a2ea4001e",
-               job:           job,
-               spec:          spec,
-               executors:     deploy.GetExecutors(),
-       }
-       if opts != nil {
-               if opts.DriverImage != "" {
-                       c.driverImage = opts.DriverImage
-               }
-               if opts.LauncherImage != "" {
-                       c.launcherImage = opts.LauncherImage
+
+       for _, v := range component_keys {
+               if v != "start" && v != "end" {
+                       //FunctionKind作成
+                       wf.Spec.Templates = append(wf.Spec.Templates, CreateFunctionKindTemp(spec, v))
+
+                       //ConfigMap作成
+                       wf.Spec.Templates = append(wf.Spec.Templates, CreateConfigMapTemp(spec, v))
                }
-               if opts.PipelineRoot != "" {
-                       job.RuntimeConfig.GcsOutputDirectory = opts.PipelineRoot
+       }
+
+       wf.Spec.Templates = append(wf.Spec.Templates, CreateDataflowTemp(spec))
+
+       //functionchainを生成
+       wf.Spec.Templates = append(wf.Spec.Templates, CreateFunctionChainTemp(spec, component_keys, connect_keys))
+
+       //デプロイ順序に合わせてソート
+       //ソートの順番を定義
+       orderList := []string{
+               "funcinfo-",
+               "fk-",
+               "functionchain",
+               "dataflow",
+       }
+
+       // Templates内をソート
+       sort.Slice(wf.Spec.Templates, func(i, j int) bool {
+               return index(orderList, wf.Spec.Templates[i].Name) < index(orderList, wf.Spec.Templates[j].Name)
+       })
+
+       //wfの情報からDeployを生成
+       for i := 0; len(wf.Spec.Templates) > i; i++ {
+               // Templates内でdeployの要素が何番目かを確認
+               if wf.Spec.Templates[i].Name == "deploy" {
+                       // Templatesの要素数分ループ
+                       for k, v := range wf.Spec.Templates {
+                               //deployが格納されていないものを選別
+                               if k != i {
+                                       //workflowstep生成
+                                       tmpWorkflowStep := []wfapi.WorkflowStep{}
+                                       tmpWorkflowStep = append(tmpWorkflowStep, wfapi.WorkflowStep{
+                                               Name:     v.Name,
+                                               Template: v.Name,
+                                       })
+                                       //ParalelSteps生成
+                                       tmpParallelSteps := wfapi.ParallelSteps{
+                                               Steps: tmpWorkflowStep,
+                                       }
+                                       wf.Spec.Templates[i].Steps = append(wf.Spec.Templates[i].Steps, tmpParallelSteps)
+                               }
+                       }
                }
        }
 
-       // compile
-       err = compiler.Accept(job, kubernetesSpec, c)
+       // Debug
+       got, err := yaml.Marshal(wf)
+       if err != nil {
+               log.Fatal(err)
+       }
+       println(string(got), "\n")
+
```

* wf変数に対し成形したリソースを格納するよう変更したため戻り値を書き換える。
```go
       }
       println(string(got), "\n")

-       return c.wf, err
+       return wf, err
 }
 
 func retrieveLastValidString(s string) string {
```

## 次の手順について
[visitor.go改修](../modify-visitor.go)

## 補足事項
(※1) Custom Kubeflow SDKのIR YAMLの書式はcomponentとconnectionの2種類を定義する仕様で実装している。将来的に当該2種類以外の書式を追加する場合は別途判定処理を追記する必要がある。   
(※2) OpenKasugaiコントローラのバージョンアップなどでデプロイするリソースの仕様変更があった場合は適宜テンプレート成形処理を追加および変更必要がある。   
(※3)   
  - OpenKasugaiコントローラはリソースのデプロイ順序によってはデプロイに失敗する。
  - そのため以下処理でデプロイ順序を指定している。
  - OpenKasugaiコントローラのバージョンに合わせてリソースのデプロイ順序を指定する必要がある。
```go
+       //デプロイ順序に合わせてソート
+       //ソートの順番を定義
+       orderList := []string{
+               "funcinfo-",
+               "fk-",
+               "functionchain",
+               "dataflow",
+       }
+
+       // Templates内をソート
+       sort.Slice(wf.Spec.Templates, func(i, j int) bool {
+               return index(orderList, wf.Spec.Templates[i].Name) < index(orderList, wf.Spec.Templates[j].Name)
+       })
```