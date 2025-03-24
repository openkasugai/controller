package testutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1" //nolint:stylecheck // ST1019: intentional import as another name
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	goccy "github.com/goccy/go-yaml"
	. "github.com/onsi/gomega"
)

var (
	ErrElementsNotEqual      = errors.New("Elements doesn't equal")
	ErrInvalidPathIsInputeed = errors.New("Invalid path is inputted")
	ErrFailedToConvert       = errors.New("Failed to convert")
)

func decodeManifest(fileName string) (obj runtime.Object, err error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return obj, err
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err = decode(bytes, nil, nil)
	return obj, err
}

func GenerateExpectYaml(resource client.Object, apiVersion string, kind string, path string) error {

	resource.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: kind, Version: apiVersion})

	newFile, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer newFile.Close()
	y := printers.YAMLPrinter{}
	return y.PrintObj(resource, newFile)
}

func Deploy(ctx context.Context, k8sClient client.Client, paths ...string) (map[string]any, error) {

	ret := make(map[string]any)
	for _, p := range paths {

		if f, err := os.Stat(p); !os.IsNotExist(err) {

			if f.IsDir() {
				objs, err := deployFromDirectory(ctx, k8sClient, p)
				if err != nil {
					return nil, err
				}
				for k, v := range objs {
					ret[k] = v
				}
			} else {
				obj, err := decodeManifest(p)
				if err != nil {
					return nil, err
				}
				ret[p] = obj
				if err := deploy(ctx, k8sClient, obj); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}

	}

	return ret, nil
}

func DeployWithObjectMeta(ctx context.Context, k8sClient client.Client, paths ...string) (map[string]any, error) {

	ret := make(map[string]any)
	for _, p := range paths {

		if f, err := os.Stat(p); !os.IsNotExist(err) {

			if f.IsDir() {
				objs, err := deployFromDirectoryWithObjectMeta(ctx, k8sClient, p)
				if err != nil {
					return nil, err
				}
				for k, v := range objs {
					ret[k] = v
				}
			} else {
				obj, err := decodeManifest(p)
				if err != nil {
					return nil, err
				}
				ret[p] = obj
				if err := deployWithObjectMeta(ctx, k8sClient, obj); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}

	}

	return ret, nil
}

func deployFromDirectory(ctx context.Context, k8sClient client.Client, dirPath string) (map[string]any, error) {

	ret := make(map[string]any)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return ret, err
	}

	paths := make([]string, 0)
	for _, f := range files {
		fName := f.Name()
		ext := filepath.Ext(fName)
		if ext == ".yaml" {
			paths = append(paths, filepath.Join(dirPath, fName))
		}
	}
	return Deploy(ctx, k8sClient, paths...)
}

func deployFromDirectoryWithObjectMeta(ctx context.Context, k8sClient client.Client, dirPath string) (map[string]any, error) {

	ret := make(map[string]any)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return ret, err
	}

	paths := make([]string, 0)
	for _, f := range files {
		fName := f.Name()
		ext := filepath.Ext(fName)
		if ext == ".yaml" {
			paths = append(paths, filepath.Join(dirPath, fName))
		}
	}
	return DeployWithObjectMeta(ctx, k8sClient, paths...)
}

func Check(expResource client.Object, actResource client.Object, keys ...string) error {

	actResource.GetObjectKind().SetGroupVersionKind(expResource.GetObjectKind().GroupVersionKind())

	y := printers.YAMLPrinter{}

	var expByte bytes.Buffer
	err := y.PrintObj(expResource, &expByte)
	if err != nil {
		return err
	}

	var actByte bytes.Buffer
	err = y.PrintObj(actResource, &actByte)
	if err != nil {
		return err
	}

	for _, key := range keys {

		if key[0] != '.' {
			return ErrInvalidPathIsInputeed
		}

		key = key[1:]

		var expMap map[interface{}]interface{}
		var actMap map[interface{}]interface{}

		if key == "" {
			err = parseKey(expByte.String(), "", &expMap)
			if err != nil {
				return err
			}
			err = parseKey(actByte.String(), "", &actMap)
			if err != nil {
				return err
			}
			if msg, res := recursiveCompare(key, expMap, actMap); !res {
				return fmt.Errorf("%w : %v", ErrElementsNotEqual, msg)
			}
		} else {

			spKey := strings.Split(key, ".")
			if len(spKey) == 1 {

				err = parseKey(expByte.String(), "", &expMap)
				if err != nil {
					return err
				}
				err = parseKey(actByte.String(), "", &actMap)
				if err != nil {
					return err
				}

				if msg, res := recursiveCompare(key, expMap[key], actMap[key]); !res {
					return fmt.Errorf("%w : %v", ErrElementsNotEqual, msg)
				}

			} else {
				keyEnd := spKey[len(spKey)-1]
				keyBase := "." + strings.Join(spKey[:len(spKey)-1], ".")

				err = parseKey(expByte.String(), keyBase, &expMap)
				if err != nil {
					return err
				}
				err = parseKey(actByte.String(), keyBase, &actMap)
				if err != nil {
					return err
				}

				if msg, res := recursiveCompare(keyBase, expMap[keyEnd], actMap[keyEnd]); !res {
					return fmt.Errorf("%w : %v", ErrElementsNotEqual, msg)
				}
			}

		}

	}

	return nil
}

func parseKey(source string, key string, dest interface{}) error {

	source = editInputStringForReadYAML_Format(source)

	path, err := goccy.PathString("$.TmpKey" + key)
	if err != nil {
		return err
	}

	err = path.Read(strings.NewReader(source), dest)
	return err
}

func recursiveCompare(key string, exp interface{}, act interface{}) (string, bool) {

	var result bool = true
	var retErrMessage string = ""

	switch expElem := exp.(type) {
	case []interface{}:
		actElem, ok := act.([]interface{})
		if !ok {
			retErrMessage += fmt.Sprintf("act doesn't have key : %v \n", key)
			return retErrMessage, false
		}

		if len(expElem) != len(actElem) {
			retErrMessage += fmt.Sprintf("Length of arrays are not same at key %v \n", key)
			retErrMessage += fmt.Sprintf("Expect : %v Actual : %v\n", exp, act)
			return retErrMessage, false
		}

		for i, expValue := range expElem {
			if msgs, res := recursiveCompare(fmt.Sprintf("%v[%v]", key, i), expValue, actElem[i]); !res {
				result = false
				retErrMessage += msgs
			}
		}

	case map[interface{}]interface{}:
		actElem, ok := act.(map[interface{}]interface{})
		if !ok {
			retErrMessage += fmt.Sprintf("act doesn't have key : %v \n", key)
			return retErrMessage, false
		}

		for k, v := range expElem {

			actValue, ok := actElem[k]
			if ok {
				if msgs, res := recursiveCompare(fmt.Sprintf("%v.%v", key, k), v, actValue); !res {
					result = false
					retErrMessage += msgs
				}
			} else {
				retErrMessage += fmt.Sprintf("Not existing key %v for actual result. %v\n", k, key)
				retErrMessage += fmt.Sprintf("Expect : %v Actual : %v\n", exp, act)
			}
		}
	case map[string]interface{}:
		actElem, ok := act.(map[string]interface{})

		if !ok {
			retErrMessage += fmt.Sprintf("act doesn't have key : %v \n", key)
			return retErrMessage, false
		}

		for k, v := range expElem {
			actValue, ok := actElem[k]
			if ok {
				if msgs, res := recursiveCompare(fmt.Sprintf("%v.%v", key, k), v, actValue); !res {
					result = false
					retErrMessage += msgs
				}
			} else {
				retErrMessage += fmt.Sprintf("Not existing key %v for actual result. %v\n", k, key)
				retErrMessage += fmt.Sprintf("Expect : %v Actual : %v\n", exp, act)
			}
		}
	default:
		if exp != act {
			retErrMessage += fmt.Sprintf("Values are defferent on Key %v \n", key)
			retErrMessage += fmt.Sprintf("Expect : %v Actual : %v\n", exp, act)
			return retErrMessage, false
		}
	}

	return retErrMessage, result
}

// Add key and shift content for being able to parse in yaml format.
func editInputStringForReadYAML_Format(source string) string {
	// Split the input data into lines
	lines := strings.Split(source, "\n")

	// Initialize a slice to store the modified lines
	var modifiedLines []string

	// Iterate through the lines
	for _, line := range lines {
		line = strings.Replace(line, "---", "", -1)
		// if line != "" && (strings.Contains(line, "-") || (strings.Contains(line, ":"))) {
		if line != "" {
			modifiedLine := "  " + line
			modifiedLines = append(modifiedLines, modifiedLine)
		}
	}

	// Join the modified lines with newlines
	return "TmpKey :\n" + strings.Join(modifiedLines, "\n")
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func CreateDataFlow(ctx context.Context, df ntthpcv1.DataFlow, k8sClient client.Client) error {
	tmp := &ntthpcv1.DataFlow{}
	*tmp = df
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	tmp.Status = df.Status
	err = k8sClient.Status().Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func CreateSchedulingData(ctx context.Context, sd ntthpcv1.SchedulingData, k8sClient client.Client) error {
	tmp := &ntthpcv1.SchedulingData{}
	*tmp = sd
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	tmp.Status = sd.Status
	err = k8sClient.Status().Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllOf(ctx context.Context, k8sClient client.Client, nameSpace string, objs ...client.Object) error {
	for _, obj := range objs {
		err := k8sClient.DeleteAllOf(ctx, obj, client.InNamespace(nameSpace))
		if err != nil {
			return err
		}
	}
	return nil
}

// If you convert a runtime.Object to a client.Object,
// There is a problem where the .Status field and below disappear.
// Currently, the only way to avoid this is to convert it directly to the original resource.

func deploy(ctx context.Context, cli client.Client, resource runtime.Object) error {

	const updateChallengeNum = 5

	switch r := resource.(type) {
	case *ntthpcv1.ComputeResource:
		dep := &ntthpcv1.ComputeResource{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.ComputeResource{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.ConnectionTarget:
		dep := &ntthpcv1.ConnectionTarget{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.ConnectionTarget{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.DataFlow:
		dep := &ntthpcv1.DataFlow{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.DataFlow{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.FunctionChain:
		dep := &ntthpcv1.FunctionChain{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.FunctionChain{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.FunctionType:
		dep := &ntthpcv1.FunctionType{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.FunctionType{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.FunctionTarget:
		dep := &ntthpcv1.FunctionTarget{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.FunctionTarget{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.SchedulingData:
		dep := &ntthpcv1.SchedulingData{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.SchedulingData{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.TopologyInfo:
		dep := &ntthpcv1.TopologyInfo{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.TopologyInfo{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.WBConnection:
		dep := &ntthpcv1.WBConnection{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.WBConnection{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.WBFunction:
		dep := &ntthpcv1.WBFunction{}
		dep.Name = r.Name
		dep.Namespace = r.Namespace
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.WBFunction{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	default:
		dep := r.(client.Object)
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
	}

	return nil
}

func deployWithObjectMeta(ctx context.Context, cli client.Client, resource runtime.Object) error {

	const updateChallengeNum = 5

	switch r := resource.(type) {
	case *ntthpcv1.ComputeResource:
		dep := &ntthpcv1.ComputeResource{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.ComputeResource{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.ConnectionTarget:
		dep := &ntthpcv1.ConnectionTarget{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.ConnectionTarget{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.DataFlow:
		dep := &ntthpcv1.DataFlow{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.DataFlow{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.FunctionChain:
		dep := &ntthpcv1.FunctionChain{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.FunctionChain{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.FunctionType:
		dep := &ntthpcv1.FunctionType{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.FunctionType{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.FunctionTarget:
		dep := &ntthpcv1.FunctionTarget{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.FunctionTarget{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.SchedulingData:
		dep := &ntthpcv1.SchedulingData{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.SchedulingData{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.TopologyInfo:
		dep := &ntthpcv1.TopologyInfo{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.TopologyInfo{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Status().Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.WBConnection:
		dep := &ntthpcv1.WBConnection{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.WBConnection{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	case *ntthpcv1.WBFunction:
		dep := &ntthpcv1.WBFunction{}
		dep.ObjectMeta = r.ObjectMeta
		dep.Spec = r.Spec
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
		var err error
		for i := 0; i < updateChallengeNum; i++ {
			dep = &ntthpcv1.WBFunction{}
			srcInfo := types.NamespacedName{Namespace: r.Namespace, Name: r.Name}
			if err := cli.Get(ctx, srcInfo, dep); err != nil {
				continue
			}
			dep.Status = r.Status
			err = cli.Update(ctx, dep)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return err

	default:
		dep := r.(client.Object)
		if err := cli.Create(ctx, dep); err != nil {
			return err
		}
	}

	return nil
}

func GetResourceFromYaml[T any](path string) (T, error) {

	obj, err := decodeManifest(path)
	if err != nil {
		var tmp T
		return tmp, err
	}

	ret, ok := obj.(T)
	if !ok {
		var tmp T
		return tmp, ErrFailedToConvert
	} else {
		return ret, nil
	}
}

func CreateExpectYaml(ctx context.Context, k8sClient client.Client, dest, name, nameSpace, apiVersion, expStatus string, resource any) {
	Eventually(func(g Gomega) {
		switch v := resource.(type) {
		case ntthpcv1.DataFlow:
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(
					ctx, client.ObjectKey{Namespace: nameSpace, Name: name}, &v)).To(Succeed())
				g.Expect(v.Status.Status).Should(Equal(expStatus))
				g.Expect(GenerateExpectYaml(&v, apiVersion, "DataFlow", dest)).Should(Succeed())
			}).Should(Succeed())
		case ntthpcv1.SchedulingData:
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(
					ctx, client.ObjectKey{Namespace: nameSpace, Name: name}, &v)).To(Succeed())
				g.Expect(v.Status.Status).Should(Equal(expStatus))
				g.Expect(GenerateExpectYaml(&v, apiVersion, "SchedulingData", dest)).Should(Succeed())
			}).Should(Succeed())
		}
	}).Should(Succeed())
}

func RemoveFinalizers(ctx context.Context, k8sClient client.Client, namespace string, resources ...client.Object) error {

	for _, resource := range resources {
		switch resource.(type) {
		case *ntthpcv1.ComputeResource:
			list := &ntthpcv1.ComputeResourceList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.ConnectionTarget:
			list := &ntthpcv1.ConnectionTargetList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.DataFlow:
			list := &ntthpcv1.DataFlowList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.FunctionChain:
			list := &ntthpcv1.FunctionChainList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.FunctionType:
			list := &ntthpcv1.FunctionTypeList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.FunctionTarget:
			list := &ntthpcv1.FunctionTargetList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.SchedulingData:
			list := &ntthpcv1.SchedulingDataList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.TopologyInfo:
			list := &ntthpcv1.TopologyInfoList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.WBConnection:
			list := &ntthpcv1.WBConnectionList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *ntthpcv1.WBFunction:
			list := &ntthpcv1.WBFunctionList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		case *corev1.ConfigMap:
			list := &corev1.ConfigMapList{}
			if err := k8sClient.List(ctx, list, client.InNamespace(namespace)); err != nil {
				return err
			}
			for _, item := range list.Items {
				item.SetFinalizers(nil)
				if err := k8sClient.Update(ctx, &item); err != nil {
					return err
				}
			}

		default:
			return fmt.Errorf("unsupported resource type: %T", resource)
		}
	}

	return nil
}

func ValToAddr[T any](in T) *T {
	return &in
}
