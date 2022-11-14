#### Quickstart 


##### Docker

Prerequisites:
- Docker
- go environment 1.18+

1. ```make deps && make run_envtest```
2. then grab path and port from apiserver url and start schemathesis via simple runner script (docker required):  ```hack/schemathesis_docker.sh {temp_folder_pwd} {apiserver_port}```

Kubeconfig might be found in the tempfolder, envtest's apiserver interactable via `kubectl` with this kubeconfig


##### Local

This method uses local copy of the latest schemathesis version from the `main` branch of [schemathesis](https://github.com/schemathesis/schemathesis) repository

Prerequisites:
- Python 3.8+
  - [virtualenv](https://virtualenv.pypa.io/en/latest/installation.html#via-pip)
- go environment 1.18+

1. ```make deps && make python_venv && make run_envtest```
2. then grab path and apiserver url runner output and start schemathesis via simple runner script:  ```hack/schemathesis_local.sh {temp_folder_pwd} {apiserver_url}```

Kubeconfig might be found in the tempfolder, envtest's apiserver interactable via `kubectl` with this kubeconfig


##### Notes

* `run_envtest` make target congigured in a way to store apiserver logs into `${PROJECT_DIR}/results` folder
