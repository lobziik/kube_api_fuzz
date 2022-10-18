#### Quickstart 

Prerequisites:
- Docker
- go environment 1.18+

1. ```make deps && make run_envtest```
2. then grab path and port from apiserver url and start schemathesis via simple runner script (docker required):  ```hack/schemathesis_docker.sh {temp_folder_pwd} {apiserver_port}```

Kubeconfig might be found in the tempfolder, envtest's apiserver interactable via `kubectl` with this kubeconfig