REPO_ROOT=$(dirname "${BASH_SOURCE}")/..

go run "$REPO_ROOT/vendor/sigs.k8s.io/controller-runtime/tools/setup-envtest"  $@ --bin-dir "$REPO_ROOT/bin"