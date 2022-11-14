PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
ENVTEST = go run ${PROJECT_DIR}/vendor/sigs.k8s.io/controller-runtime/tools/setup-envtest
ENVTEST_K8S_VERSION = 1.25

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: run_envtest
run_envtest: ## Run kube apiserver and etcd (envtest)
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --bin-dir $(PROJECT_DIR)/bin)" go run ${PROJECT_DIR}/cmd/envtest_runner --apiserverStdout="${PROJECT_DIR}/results/apiserverStdout.log" --apiserverStderr="${PROJECT_DIR}/results/apiserverStderr.log"


.PHONY: deps
deps: ## download dependencies
	go mod tidy && go mod vendor


.PHONY: check-and-reinit-submodules
check-and-reinit-submodules:
	@if git submodule status | egrep -q '^[-]|^[+]' ; then \
		echo "Need to reinitialize git submodules"; \
		git submodule update --init; \
	fi


.PHONY: python_venv
python_venv: python_venv/touchfile ## create python virtualenv, install local copy on schemathesis into it


python_venv/touchfile: check-and-reinit-submodules ./schemathesis/poetry.lock
	test -d python_venv || virtualenv python_venv
	. python_venv/bin/activate; pip install -e ./schemathesis;
	touch python_venv/touchfile