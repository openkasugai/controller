/* Copyright 2025 NTT Corporation , FUJITSU LIMITED */

package main

import (
	"context"
	examplecomv1 "example.com/FPGACheckPerDF/api/v1"
	"example.com/FPGACheckPerDF/pkg/dfcheck"
	paramcheck "example.com/FPGACheckPerDF/pkg/paramcheck"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(examplecomv1.AddToScheme(scheme))
}

func main() {
	mainRet := 1
	ctx := context.Background()

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		mng, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme})
		if err != nil {
			fmt.Println("new manager function error.", err)
			break
		}

		cCommand := ""
		cCommandMode := ""
		cNodeName := ""
		cNamespace := ""

		intRet := paramcheck.ParamMain(ctx, &cCommand, &cCommandMode, &cNodeName, &cNamespace)
		if 0 != intRet {
			break
		}

		if "get" == cCommandMode {
			err = dfcheck.ShellDeleteAllMain(ctx, cCommand)
			if nil != err {
				break
			}
			err = dfcheck.LogDeleteAllMain(ctx, cCommand)
			if nil != err {
				break
			}
			intRet = dfcheck.CollectMain(ctx, mng, cCommand, cCommandMode, cNodeName, cNamespace)
			if 0 != intRet {
				break
			}
		} else if "check" == cCommandMode {
			err = dfcheck.LogDeleteAllMain(ctx, cCommand)
			if nil != err {
				break
			}
			intRet = dfcheck.CollectMain(ctx, mng, cCommand, cCommandMode, cNodeName, cNamespace)
			if 0 != intRet {
				break
			}
			err = dfcheck.ShellExecMain(ctx, cCommand)
			if nil != err {
				break
			}
		} else {
			intRet = dfcheck.CollectMain(ctx, mng, cCommand, cCommandMode, cNodeName, cNamespace)
			if 0 != intRet {
				break
			}
		}
		mainRet = 0
	}
	os.Exit(mainRet)
}
