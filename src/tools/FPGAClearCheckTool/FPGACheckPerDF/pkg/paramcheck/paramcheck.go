/* Copyright 2025 NTT Corporation , FUJITSU LIMITED */

package paramcheck

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func ParamMain(ctx context.Context,
	pCommand *string,
	pCommandMode *string,
	pNodeName *string,
	pNamespace *string) int {

	intRet := 0

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		err := paramGet(ctx, pCommand, pCommandMode, pNodeName, pNamespace)
		if nil != err {
			break
		}
		if "get" == *pCommandMode ||
			"check" == *pCommandMode ||
			"debug" == *pCommandMode {
		} else {
			fmt.Println("mode parameter error. input mode get or check.")
			intRet = 1
			break
		}
		if "" == *pNodeName {
			cNodeName := os.Getenv("K8S_NODENAME")
			if "" == cNodeName {
				if "debug" != *pCommandMode {
					fmt.Println("NodeName get error. input node or export K8S_NODENAME")
					intRet = 1
					break
				}
			} else {
				*pNodeName = cNodeName
			}
		}
	}

	return intRet
}

func paramGet(ctx context.Context,
	pCommand *string,
	pCommandMode *string,
	pNodeName *string,
	pNamespace *string) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		for i, v := range os.Args {
			if 0 == i {
				*pCommand = v
			} else if 1 == i {
				*pCommandMode = v
			} else {
				break
			}
		}

		if len(os.Args) >= 2 {
			var prmNamespace *string
			var prmNodeName *string

			fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			prmNodeName = fs.String("node", "default", "NodeName Setting")
			prmNamespace = fs.String("n", "default", "NameSpace Setting")

			err = fs.Parse(os.Args[2:])
			if nil != err {
				fmt.Println("parameter parse error.", err)
				break
			}

			if nil != prmNodeName {
				if "default" != *prmNodeName {
					*pNodeName = *prmNodeName
				}
			}
			if nil != prmNamespace {
				*pNamespace = *prmNamespace
			}
		}
	}
	return err
}
