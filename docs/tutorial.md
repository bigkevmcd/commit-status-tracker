# Setting up to Push GitHub commit-statuses

## Auth Secret

You'll need a secret created from a GitHub auth token, this command assumes a token is in `~/Downloads/token`.

```shell
$ kubectl create secret generic github-auth --from-file=$HOME/Downloads/token
```

## Annotating a PipelineRun
