# Secrets From GCP

This plugin generates Kubernetes secrets from a GCP Maneged Secret.

Here is a simple secret-from-gcp config object and the secret it creates:

```yaml
apiVersion: pollination.cloud/v1
kind: SecretsFromGCP
metadata:
  name: secret-name
  namespace: default
source:
  projectId: pollination-staging-1d6a8      # GCP Project ID
  name: kustomization-plugin-example-secret # GCP secret name
  version: latest                           # GCP secret version (default: latest)
keys:                                       # Keys from the secret to be used
  - username
  - password
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: secret-name-g7gctm97b4
  namespace: default
type: Opaque
data:
  password: cGFzc3dvcmQ=
  username: YWRtaW4=
```

The plugin's API also allows for more configuration:

```yaml
apiVersion: pollination.cloud/v1
kind: SecretsFromGCP
metadata:
  name: secret-name
  namespace: default
source:
  projectId: pollination-staging-1d6a8      # GCP Project ID
  name: kustomization-plugin-example-secret # GCP secret name
  version: latest                           # GCP secret version (default: latest)
type: Opaque                                # Kubernetes secret type (default: Opaque)
disableNameSuffixHash: true                 # disable creating name suffix hash used to trigger service rollout
behavior: replace                           # secret update behavior (create | replace | merge) 
keys:                                       # Keys from the secret to be used
  - username
  - password
```