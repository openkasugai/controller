
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.28.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: fmtcheck
fmtcheck: ## Run go fmt against code, check only.
	gofmt -d .
	[ $(shell gofmt -l . | wc -l) -eq 0 ]

.PHONY: staticcheck
staticcheck: staticcheck_bin ## Run staticcheck against code , instead go vet.
	$(STATICCHECK) ./... || echo


.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt staticcheck envtest ## Run tests.
	ACK_GINKGO_DEPRECATIONS=1.16.5 \
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" \
	go test ./lib/... ./internal/controller/... ./test/envtest/...  -coverprofile cover.out --ginkgo.noColor --ginkgo.failFast --test.v 
	

# In deploytest, it is necessary to execute other tests serially, but due to the current specification of go test,
# serial execution is difficult, so they need to be run one at a time.
# If the use of ginkgo becomes a prerequisite, this method will be improved.
.PHONY: deploytest
deploytest: manifests generate fmt staticcheck envtest ## Run tests.
	ACK_GINKGO_DEPRECATIONS=1.16.5 \
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" \
	go test ./test/deploytest/splitfilter/...  -coverprofile cover.out --ginkgo.noColor --ginkgo.failFast --test.v 

##test: manifests generate fmt vet envtest ## Run tests.
##	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -coverprofile cover.out


.PHONY: debug
ifdef FUNC
FOCUS = --ginkgo.focus $(FUNC)
endif
debug: fmt envtest ## Debug test .
	ACK_GINKGO_DEPRECATIONS=1.16.5 \
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" \
	go test ./... -coverprofile cover.out \
		--ginkgo.noColor \
		--ginkgo.failFast \
		--test.v \
		--ginkgo.reportPassed \
		$(FOCUS)

.PHONY: report
report: test docker-build ## Run lcov for output coverage report html.
	$(RUNBUILDER) gcov2lcov -infile=cover.out -outfile=coverage.lcov
	$(RUNBUILDER) genhtml coverage.lcov -o site

.PHONY: report-http
report-http:
	cd site ; python3 -m http.server 8888

##@ Build

.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./cmd/main.go

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	sudo buildah bud -t ${IMG} -f ./Dockerfile ..  

.PHONY: docker-build-wo-test
docker-build-wo-test: 
	sudo buildah bud -t ${IMG} -f ./Dockerfile ..  

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

# PLATFORMS defines the target platforms for the manager image be built to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
STATICCHECK ?= $(LOCALBIN)/staticcheck

## Tool Versions
KUSTOMIZE_VERSION ?= v5.1.1
CONTROLLER_TOOLS_VERSION ?= v0.14.0

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@v0.0.0-20240317073005-bd9ea79e8d18
.PHONY: staticcheck_bin
staticcheck_bin: $(STATICCHECK) ## Download staticcheck locally if necessary.
$(STATICCHECK): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install honnef.co/go/tools/cmd/staticcheck@v0.5.1

test_single_cpu:
	@make -C config/samples/single-function/cpu test

test_single_cpu_clean:
	@make -C config/samples/single-function/cpu test_delete

test_single_fpga:
	@make -C config/samples/single-function/fpga test

test_single_fpga_clean:
	@make -C config/samples/single-function/fpga test_delete

test_multi_cpu:
	@make -C config/samples/IT_multi-function/cpu-cpu test

test_multi_cpu_clean:
	@make -C config/samples/IT_multi-function/cpu-cpu test_delete

test_2node_cpu:
	@make -C config/samples/multi-node/cpu-cpu test

test_2node_cpu_clean:
	@make -C config/samples/multi-node/cpu-cpu test_delete

test_2node_3func_cpu:
	@make -C config/samples/multi-node/cpu-cpu test_3func

test_2node_3func_cpu_clean:
	@make -C config/samples/multi-node/cpu-cpu test_3func_delete

test_2node_fpga:
	@make -C config/samples/multi-node/fpga-cpu test

test_2node_fpga_clean:
	@make -C config/samples/multi-node/fpga-cpu test_delete

test_2node_3func_fpga:
	@make -C config/samples/multi-node/fpga-cpu test_3func

test_2node_3func_fpga_clean:
	@make -C config/samples/multi-node/fpga-cpu test_3func_delete

test_autosche_01:
	@make -C config/samples/auto-schedule/pattern01 test_3func

test_autosche_01_clean:
	@make -C config/samples/auto-schedule/pattern01 test_3func_delete

test_autosche_02:
	@make -C config/samples/auto-schedule/pattern02 test_3func

test_autosche_02_clean:
	@make -C config/samples/auto-schedule/pattern02 test_3func_delete

test_alreadysche_01:
	@make -C config/samples/already-schedule/pattern01 test_3func

test_alreadysche_01_clean:
	@make -C config/samples/already-schedule/pattern01 test_3func_delete

test_alreadysche_02:
	@make -C config/samples/already-schedule/pattern02 test_3func

test_alreadysche_02_clean:
	@make -C config/samples/already-schedule/pattern02 test_3func_delete

