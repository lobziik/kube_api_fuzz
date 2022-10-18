set -x

docker run --rm --network=host \
 -v "$1":/kube_configs \
 -v /Users/dmoiseev/work/kube_api_fuzz/results:/results \
 schemathesis/schemathesis:stable \
  run \
  --show-errors-tracebacks \
  --request-cert /kube_configs/clientCert.pem \
  --request-cert-key /kube_configs/clientKey.pem \
  --request-tls-verify false \
  --junit-xml /results/junit.xml \
  --workers 1 \
  --checks all \
  https://host.docker.internal:"$2"/openapi/v2