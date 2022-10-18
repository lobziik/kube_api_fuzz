PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
ENVTEST = go run ${PROJECT_DIR}/vendor/sigs.k8s.io/controller-runtime/tools/setup-envtest
ENVTEST_K8S_VERSION = 1.25



.PHONY: run_envtest
run_envtest: ## Run kube apiserver and etcd (envtest)
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --bin-dir $(PROJECT_DIR)/bin)" go run ${PROJECT_DIR}/cmd/envtest_runner


.PHONY: deps
deps: ## download dependencies
	go mod tidy && go mod vendor
