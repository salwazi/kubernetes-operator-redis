# Redis Operator Deployment Guide

This guide provides instructions on how to deploy the Redis Operator, create a Redis instance using the Custom Resource Definition (CRD), and confirm the deployment and secret in use.

## Prerequisites

- Access to a Kubernetes cluster (minikube)
- kubectl configured to interact with your cluster
- go version v1.21.0+
- docker version 17.03+.
- Operator SDK (Optional for development)

## Deployment Steps

### 1. Clone the Repository

First, clone the repository containing the Redis Operator to your local machine:

```sh
git clone https://github.com/salwazi/kubernetes-operator-redis.git
```


**Install the CRDs into the cluster:**

```sh
make install
```

**Run the operator locally**

```sh
make run
```

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:
In a
```sh
kubectl apply -f config/crd/bases/cache.example.com_redis.yaml
```


**Confirming the Deployment and Secret**
You can check the status of the Deployment and the Secret as well as the CRD from the cluster

```sh
kubectl get crd
kubectl get redis
kubectl get deployments
kubectl get secret
```

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/crd/bases/cache.example.com_redis.yaml
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

Since we only run the operator locally, exiting the shell is sufficient.


## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

