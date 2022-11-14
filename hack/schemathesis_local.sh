REPO_ROOT=$(dirname "${BASH_SOURCE}")/..
OUTPUT_DIR="$REPO_ROOT/results"
KUBE_CONFIGS="$1"
APISERVER_URL="$2"

test -d "$REPO_ROOT/python_venv" || (echo "No python virtualenv found, exit"; exit 1)

source "$REPO_ROOT/python_venv/bin/activate"

echo "KUBE CONFIGS IN: $KUBE_CONFIGS"
echo "APISERVER URL IS: $APISERVER_URL"


set -x
st run \
 --show-errors-tracebacks \
 --request-cert "$KUBE_CONFIGS/clientCert.pem" \
 --request-cert-key "$KUBE_CONFIGS/clientKey.pem" \
 --request-tls-verify false \
 --junit-xml "$OUTPUT_DIR/junit.xml" \
 --workers 1 \
 --checks all \
 "$APISERVER_URL/openapi/v2"

