# Setting up to Push GitHub commit-statuses

## Prerequisites

You'll need this operator, and Tekton installed see the installation
[instructions](../README.md#installing).

And you'll need a GitHub auth token.

## Create a secret

Create a secret from your GitHub auth token, this command assumes a token is in `~/Downloads/token`.

```shell
$ kubectl create secret generic github-auth --from-file=$HOME/Downloads/token
```

## Annotating a PipelineRun

The operator watches for PipelineRuns with specific annotations.

This is an alpha operator, and the annotation names will change, but for now
you'll need...

```yaml
apiVersion: tekton.dev/v1alpha1
kind: PipelineRun
metadata:
  name: demo-pipeline-run
  annotations:
    "app.example.com/git-status": "true"
    "app.example.com/status-context": "demo-pipeline"
    "app.example.com/status-description": "this is a test"
spec:
  pipelineRef:
    name: demo-pipeline
  serviceAccountName: 'default'
  resources:
  - name: source
    resourceSpec:
      type: git
      params:
        - name: revision
          value: insert revision
        - name: url
          value: https://github.com/this/repo
```

The annotations are:

<table style="width=100%" border="1">
  <tr>
    <th>Name</th>
    <th>Description</th>
    <th>Required</th>
    <th>Default</th>
  </tr>
  <tr>
    <th>
      app.example.com/git-status
    </th>
    <td>
      This indicates that this `PipelineRun` should trigger commit-status notifications.
    </td>
    <td><b>Yes</b></td>
    <td></td>
  </tr>
  <tr>
    <th>
      app.example.com/status-context
    </th>
    <td>
      This is the [context](https://developer.github.com/v3/repos/statuses/#create-a-status) that will be reported, you can require named contexts in your branch protection rules.
    </td>
    <td>No</td>
    <td>"default"</td>
  </tr>
  <tr>
    <th>
      app.example.com/status-description
    </th>
    <td>
      This is used as the description of the context, not the commit.
    </td>
    <td>No</td>
    <td>""</td>
  </tr>
  <tr>
    <th>
     app.example.com/status-target-url
    </th>
    <td>
      If provided, then this will be linked in the GitHub web UI, this could be used to link to logs or output.
    </td>
    <td>No</td>
    <td>""</td>
  </tr>
</table>

## Detecting the Git Repository

Currently, this uses a simple mechanism to find the Git repository and SHA to
update the status of.

It looks for a single `PipelineResource` of type `git` and pulls the *url*
and *revision* from there.

If no suitable `PipelineResource` is found, then this will be logged as an
error, and _not_ retried.

## FAQ

 1. Does this work with `resourceRef`
    *Not yet, this is definitely on my TODO list*.
 1. Can this pull the repository details from a `pullrequest`
    `PipelineResource`?
    *Not yet, again, this is on my TODO list*.
