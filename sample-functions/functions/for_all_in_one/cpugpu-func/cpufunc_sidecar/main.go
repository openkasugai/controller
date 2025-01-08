package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(time.Local).Format("2006-01-02T15:04:05.000Z07:00"))
	}
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l)
}

type NotificationAdded struct {
	Object *corev1.ConfigMap
}
type NotificationModified struct {
	OldObject *corev1.ConfigMap
	Object    *corev1.ConfigMap
}
type NotificationDeleted struct {
	Object *corev1.ConfigMap
}

func startWatcher(clientset *kubernetes.Clientset, ctx context.Context, namespace string, name string) (<-chan interface{}, error) {
	logger := zap.L()
	logger.Debug("watch setup",
		zap.String("status", "start"),
		zap.String("namespace", namespace), zap.String("configmap-name", name))

	out := make(chan interface{})

	// watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
	// 	timeOut := int64(10)
	// 	log.WithField("namespace", namespace).WithField("name", name).Info("watch ConfigMap")
	// 	opts := metav1.SingleObject(metav1.ObjectMeta{Namespace: namespace, Name: name})
	// 	opts.TimeoutSeconds = &timeOut
	// 	return clientset.CoreV1().ConfigMaps(namespace).Watch(ctx, opts)
	// }
	selector := fields.OneTermEqualSelector("metadata.name", name)
	listwatch := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		string(corev1.ResourceConfigMaps),
		namespace,
		selector,
	)
	//listwatch.WatchFunc = watchFunc
	/*
		watcher, _ := toolswatch.NewRetryWatcher("1", listwatch)


		go func() {

			for event := range watcher.ResultChan() {
				item := event.Object.(*corev1.ConfigMap)
				log.WithFields(log.Fields{"type": event.Type, "item": item}).
					Info("event")

				switch event.Type {
				case watch.Modified:
					//notify <- struct{}{}
				case watch.Bookmark:
				case watch.Error:
				case watch.Deleted:
				case watch.Added:
					//notify <- struct{}{}
				}
			}
		}()
	*/

	_, controller := cache.NewInformer(listwatch, &corev1.ConfigMap{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			logger.Info("watch", zap.String("status", "event"), zap.String("type", "Add"))
			out <- &NotificationAdded{obj.(*corev1.ConfigMap)}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			logger.Info("watch", zap.String("status", "event"), zap.String("type", "Update"))
			out <- &NotificationModified{oldObj.(*corev1.ConfigMap), newObj.(*corev1.ConfigMap)}
		},
		DeleteFunc: func(obj interface{}) {
			logger.Info("watch", zap.String("status", "event"), zap.String("type", "Delete"))
			out <- &NotificationDeleted{obj.(*corev1.ConfigMap)}
		},
	})

	go func() {
		defer close(out)
		controller.Run(ctx.Done())
	}()

	logger.Debug("watch setup",
		zap.String("status", "end"),
		zap.String("namespace", namespace), zap.String("configmap-name", name))

	return out, nil
}

func IsInsideKubernetesClusterEnvironment() bool {
	fi, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token")
	return os.Getenv("KUBERNETES_SERVICE_HOST") != "" &&
		os.Getenv("KUBERNETES_SERVICE_PORT") != "" &&
		err == nil && !fi.IsDir()

}

func k8sRestConfig() (*rest.Config, error) {
	logger := zap.L()

	if kubeconfig, ok := os.LookupEnv("KUBECONFIG"); IsInsideKubernetesClusterEnvironment() && !ok {
		logger.Info("Environment: inside kubernetes cluster")

		cnf, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return cnf, nil
	} else {
		logger.Info("Environment: cmd")

		kubeconfigPath := kubeconfig
		if !ok {
			if home, err := os.UserHomeDir(); err != nil {
				return nil, err
			} else {
				kubeconfigPath = path.Join(home, "/.kube/config")
			}
		}

		logger.Info("Load kubeconfig", zap.String("path", kubeconfigPath))
		cnf, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, err
		}
		return cnf, nil
	}
}

func k8sClientset(config *rest.Config) (*kubernetes.Clientset, error) {

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

const configTemplatePattern = "/config/*.tmpl"

func main() {
	logger := zap.L()
	defer logger.Sync()

	namespace := os.Getenv("K8S_POD_NAMESPACE")
	if namespace == "" {
		logger.Fatal("invalid env K8S_POD_NAMESPACE")
	}
	podname := os.Getenv("K8S_POD_NAME")
	if podname == "" {
		logger.Fatal("invalid env K8S_POD_NAME")
	}
	configmapName := "ethcrl." + podname

	prosessName := os.Getenv("SIDECAR_MNG_PROSESS_NAME")
	if prosessName == "" {
		logger.Fatal("invalid prosessName(SIDECAR_MNG_PROSESS_NAME)", zap.String("prosessName", prosessName))
	}

	waitForTemplate := make(chan struct{})
	go func() {
		defer close(waitForTemplate)
		for {
			files, err := filepath.Glob(configTemplatePattern)
			if err != nil {
				logger.Fatal("template file list faild")
			}
			if 0 < len(files) {
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	<-waitForTemplate

	config, err := k8sRestConfig()
	if err != nil {
		logger.Fatal("create k8s rest config failed", zap.Error(err))
	}

	clientset, err := k8sClientset(config)
	if err != nil {
		logger.Fatal("create k8s client failed", zap.Error(err))
	}

	ctx, stop := context.WithCancel(context.Background())
	notify, err := startWatcher(clientset, ctx, namespace, configmapName)
	if err != nil {
		logger.Fatal("start watch failed", zap.Error(err))
	}
	defer stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		logger.Fatal("Notify stop signal", zap.Any("Signal", s))
		stop()
	}()

	wkdone := make(chan struct{})
	go func() {
		logger.Info("worker start")
		defer close(wkdone)

		for n := range notify {
			switch o := n.(type) {
			case *NotificationAdded:
				logger.Info("notification added",
					zap.String("namespace", o.Object.ObjectMeta.Namespace),
					zap.String("name", o.Object.ObjectMeta.Name),
					zap.Any("uuid", o.Object.ObjectMeta.UID),
					zap.String("status", "notify"))
				runConfigReloadSequence(prosessName, o.Object)
			case *NotificationModified:
				logger.Info("notification modified",
					zap.String("namespace", o.Object.ObjectMeta.Namespace),
					zap.String("name", o.Object.ObjectMeta.Name),
					zap.Any("uuid", o.Object.ObjectMeta.UID),
					zap.String("status", "notify"))
				runConfigReloadSequence(prosessName, o.Object)
			case *NotificationDeleted:
				logger.Info("notification deleted",
					zap.String("namespace", o.Object.ObjectMeta.Namespace),
					zap.String("name", o.Object.ObjectMeta.Name),
					zap.Any("uuid", o.Object.ObjectMeta.UID),
					zap.String("status", "notify"))
			}
		}

		logger.Info("worker stop")
	}()

	<-wkdone

	if err := ctx.Err(); err != nil && err != context.Canceled {
		log.Fatal(ctx.Err())
	}
}

func findProsess(prosessName string) (*os.Process, error) {
	procPath := "/proc"
	entries, err := ioutil.ReadDir(procPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pid, err := strconv.Atoi(entry.Name())
			if err != nil {
				continue
			}

			cmdline, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
			if err != nil {
				continue
			}

			name := strings.TrimSpace(string(cmdline))
			if strings.HasPrefix(name, prosessName) {
				return os.FindProcess(pid)
			}
		}
	}
	return nil, nil
}

func runConfigReloadSequence(prosessName string, configmap *corev1.ConfigMap) {
	logger := zap.L()

	raw := []byte(configmap.Data["config"])
	temlCtx := map[string]any{}

	err := json.Unmarshal(raw, &temlCtx)
	if err != nil {
		logger.Fatal("configmap[config] load failed", zap.Error(err))
	}

	tmpl, err := template.ParseGlob(configTemplatePattern)
	if err != nil {
		logger.Fatal("template file list faild")
	}
	for _, t := range tmpl.Templates() {
		name, found := strings.CutSuffix(t.Name(), ".tmpl")
		if found {
			logger.Info("run template", zap.String("name", t.Name()))
			func() {
				f, err := os.Create("/config/" + name)
				if err != nil {
					logger.Fatal("run template create file failed",
						zap.String("template", t.Name()), zap.Error(err))
				}
				defer f.Close()
				err = t.Execute(f, temlCtx)
				if err != nil {
					logger.Fatal("run template failed",
						zap.String("template", t.Name()), zap.Error(err))
				}
			}()
		}
	}

	ps, err := findProsess(prosessName)
	if err != nil {
		logger.Error("find prosess failed", zap.Error(err))
	}
	if ps == nil {
		logger.Info("prosess not found", zap.String("prosessName", prosessName))
		return
	}
	err = ps.Signal(syscall.SIGHUP)
	if err != nil {
		logger.Error("send signals failed", zap.Error(err))
	}
}
