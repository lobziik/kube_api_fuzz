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
 --workers 4 \
 --hypothesis-verbosity verbose \
 --checks all \
 --data-generation-method all \
 --debug-output-file "$OUTPUT_DIR/st_events.json" \
 --report "$OUTPUT_DIR/schemathesis_io_report.tar.gz" \
 "$APISERVER_URL/openapi/v2" 2>&1 | tee "$OUTPUT_DIR/st_stdout.txt"

