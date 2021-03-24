# Kustomize Plugins

This repository contains Plugins for Kustomize used to deploy Pollination infrastructure.

## Requirements
* Go 1.16
* [gcloud](https://cloud.google.com/sdk/docs/install)
* Unix filesystem (sorry Windows... 😅)

## Usage

Kustomize plugins are not very mature yet but they do enable us to quickly pack quite a bit of functionality into our deployment tools. Due to the lack of maturity some of the installation and development process is a bit strange.

You can install kustomize and the plugins in this repo using the following command:

```console
> make install-plugins
```

You will have to set the `XDG_CONFIG_HOME` environment variable before running any kustomize that requires plugins so it will know where to find them. We recommend you set this in value permanently in your environment variables for ease of development.

To run the examples you must also ensure that your `gcloud` profile is pointing towards the `pollination-staging-1d6a8` project on GCP and that you have access to the required secrets in this project.

```console
> export XDG_CONFIG_HOME=$HOME/.config/

> kustomize build --enable-alpha-plugins example

2021/03/24 16:52:30 Attempting plugin load from '/root/.config/kustomize/plugin/pollination.cloud/v1/secretsfromgcp/SecretsFromGCP.so'
apiVersion: v1
data:
  host: ZXhhbXBsZS5jb20=
kind: Secret
metadata:
  name: host-name-g7gctm97b4
  namespace: default
type: Opaque
---
apiVersion: v1
data:
  password: cGFzc3dvcmQ=
  username: YWRtaW4=
kind: Secret
metadata:
  name: secret-name-fd8288cb7g
  namespace: default
type: Opaque
```

## Plugins

- [SecretsFromGCP](kustomize/plugin/pollination.cloud/v1/secretsfromgcp/): Generate Kubernetes secrets from a GCP sealed secret

## References and Other Docs

* [Kustomize Go Plugin Example](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/gopluginguidedexample/): An helpful tutorial to understand the plugin development flow
* [Kustomize Go Plugin Caveats](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/goplugincaveats/): Explains why we need to compile kustomize from "scratch" to use plugins
* [Built In Plugins](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/builtins/): List of built in kustomize plugins
* [Kustomize Secrets from GCP/AWS](https://github.com/ForgeCloud/ksecrets): Helpful repo we drew heavy inspiration from when setting up this codebase