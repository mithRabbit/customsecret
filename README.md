# Customsecret Generator

A Kubernetes controller for managing CustomSecret resources, which automatically generates and rotates secrets.

## Description

The CustomSecret controller allows you to define custom secrets in your Kubernetes cluster. It supports automatic generation and rotation of secrets based on the specifications provided in the CustomSecret resource.

## Demo

You can watch a demo of the CustomSecret controller in action using Asciinema:

[![asciicast](https://asciinema.org/a/706702.svg)](https://asciinema.org/a/706702?t=30)

Or view the embedded player below:

<script src="https://asciinema.org/a/706702.js" id="asciicast-706702" async="true"></script>

## Getting Started

### Prerequisites
- go version v1.23.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/customsecret:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/customsecret:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples have default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Usage

To create a CustomSecret resource, you can use the provided example YAML file.

```sh
kubectl apply -f config/samples/customsecret_v1alpha1_customsecret.yaml
```

This will create a CustomSecret resource with the following specifications:
- Type: `basic-auth`
- Username: `admin`
- Password length: `40` characters
- Rotation period: `90` seconds

The controller will automatically generate a random password and create a Kubernetes Secret with the specified username and password. The password will be rotated according to the specified rotation period.

### Example CustomSecret Resource

```yaml
apiVersion: api.example.com/v1alpha1
kind: CustomSecret
metadata:
  name: customsecret-sample
  namespace: default
spec:
  type: basic-auth
  username: admin
  passwordLen: 40
  rotationPeriod: 86400 # 1 day in seconds
```

## Development

### Prerequisites

- go version v1.23.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### Running Locally

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run the controller locally:

```sh
make run
```

### Building and Deploying

1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/customsecret:tag
```

2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/customsecret:tag
```

### Testing

1. Create a CustomSecret resource:

```sh
kubectl apply -f config/samples/customsecret_v1alpha1_customsecret.yaml
```

2. Verify that the secret is created and rotated as expected.

### Using Helm

1. Build the Helm chart:

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. Package the Helm chart:

```sh
helm package ./dist/chart
```

3. Deploy the Helm chart:

```sh
helm install customsecret ./customsecret-0.1.0.tgz
```

4. Upgrade the Helm chart:

```sh
helm upgrade customsecret ./customsecret-0.1.0.tgz
```

5. Uninstall the Helm chart:

```sh
helm uninstall customsecret
```

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

