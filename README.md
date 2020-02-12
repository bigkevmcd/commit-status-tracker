# commit-status-tracker Go Operator [![Docker Repository on Quay](https://quay.io/repository/kmcdermo/commit-status-tracker/status "Docker Repository on Quay")](https://quay.io/repository/kmcdermo/commit-status-tracker)

## Overview

This operator tracks completed [Tekton](https://github.com/tektoncd/pipeline) [PipelineRuns](https://github.com/tektoncd/pipeline/blob/master/docs/pipelineruns.md) and attempts to create a [GitHub Commit Status](https://developer.github.com/v3/repos/statuses/) with the success or failure of the PipelineRun.

## Why?

If you're running tasks that are important parts of your deployment flow, you
can define policies that require specific checks are carried out before code can
be merged.

These can be enforced by GitHub, using their [branch
protection](https://help.github.com/en/github/administering-a-repository/configuring-protected-branches)
mechanism.

If you want your Tekton Pipelines to be a part of this, then you'll want to
report the success or failure of your PipelineRuns to Github (you might also
want Tasks, but that's implemented yet).

## Prerequisites

- [go][go_tool] version v1.12+.
- [docker][docker_tool] version 17.03+
- [kubectl][kubectl_tool] v1.11.3+
- [operator-sdk][operator_install]
- Access to a Kubernetes v1.11.3+ cluster

## Getting Started

### Cloning the repository

Checkout the Operator repository

```
$ git clone https://github.com/bigkevmcd/commit-status-tracker.git
$ cd commit-status-tracker
```
### Pulling the dependencies

Run the following command

```
$ go mod tidy
```

### Building the operator

Build the operator image and push it to a public registry, such as quay.io:

```
$ export IMAGE=quay.io/example-inc/commit-status-tracker:v0.0.1
$ operator-sdk build $IMAGE
$ docker push $IMAGE
```

### Using the image

```shell
# Update the operator manifest to use the built image name (if you are performing these steps on OSX, see note below)
$ sed -i 's|REPLACE_IMAGE|quay.io/example-inc/commit-status-tracker:v0.0.1|g' deploy/operator.yaml
# On OSX use:
$ sed -i "" 's|REPLACE_IMAGE|quay.io/example-inc/commit-status-tracker:v0.0.1|g' deploy/operator.yaml
```

**NOTE** The `quay.io/example-inc/commit-status-tracker:v0.0.1` is an example. You should build and push the image for your repository.

### Installing

You *must* have Tekton [Pipeline](https://github.com/tektoncd/pipeline/) installed before installing this operator:

```shell
$ kubectl apply -f https://github.com/tektoncd/pipeline/releases/download/v0.9.1/release.yaml
```

And then you can install the statuses operator with:

```shell
$ kubectl apply -f https://github.com/bigkevmcd/operator-statuses/releases/download/v0.0.1/release.yaml
```

### Uninstalling

```shell
$ kubectl delete -f https://github.com/bigkevmcd/operator-statuses/releases/download/v0.0.1/release.yaml
```

### Troubleshooting

Use the following command to check the operator logs.

```shell
$ kubectl logs commit-status-tracker
```

### Running Tests

```shell
$ go test -v ./...
```

[dep_tool]: https://golang.github.io/dep/docs/installation.html
[go_tool]: https://golang.org/dl/
[kubectl_tool]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[docker_tool]: https://docs.docker.com/install/
[operator_sdk]: https://github.com/operator-framework/operator-sdk
[operator_install]: https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md
