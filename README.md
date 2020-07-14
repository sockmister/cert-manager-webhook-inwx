# ACME Webhook for INWX

This project provides a cert-manager ACME Webhook for [INWX](https://inwx.de/) and a corresponding helm chart.

The helm chart is listed at Helm Hub in repository [smueller18](https://hub.helm.sh/charts/smueller18) at <https://hub.helm.sh/charts/smueller18/cert-manager-webhook-inwx>.

## Requirements

-   [helm](https://helm.sh/) >= v3.0.0
-   [kubernetes](https://kubernetes.io/) >= v1.14.0
-   [cert-manager](https://cert-manager.io/) >= 0.12.0

## Configuration

The following table lists the configurable parameters of the cert-manager chart and their default values.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `groupName` | Group name of the API service. | `cert-manager-webhook-inwx.smueller18.gitlab.com` |
| `credentialsSecretRef` | Name of secret where INWX credentials are stored. Used for RBAC to allow reading the secret by the service account name of webhook. | `inwx-credentials` |
| `deployment.loglevel` | Number for the log level verbosity of webhook deployment | 2 |
| `certManager.namespace` | Namespace where cert-manager is deployed to. | `cert-manager` |
| `certManager.serviceAccountName` | Service account of cert-manager installation. | `cert-manager` |
| `image.repository` | Image repository | `registry.gitlab.com/smueller18/cert-manager-webhook-inwx` |
| `image.tag` | Image tag | `v0.1.0` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `service.type` | API service type | `ClusterIP` |
| `service.port` | API service port | `443` |
| `resources` | CPU/memory resource requests/limits | `{}` |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `affinity` | Node affinity for pod assignment | `{}` |
| `tolerations` | Node tolerations for pod assignment | `[]` |

## Installation

### cert-manager

Follow the [instructions](https://cert-manager.io/docs/installation/) using the cert-manager documentation to install it within your cluster.

### Webhook

```bash
helm repo add smueller18 https://smueller18.gitlab.io/helm-charts
helm repo update
helm install --namespace cert-manager cert-manager-webhook-inwx smueller18/cert-manager-webhook-inwx
```

**Note**: The kubernetes resources used to install the Webhook should be deployed within the same namespace as the cert-manager.

To uninstall the webhook run
```bash
helm uninstall --namespace cert-manager cert-manager-webhook-inwx
```

## Issuer

Create a `ClusterIssuer` or `Issuer` resource as following:
```yaml
apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    # The ACME server URL
    server: https://acme-staging-v02.api.letsencrypt.org/directory

    # Email address used for ACME registration
    email: mail@example.com # REPLACE THIS WITH YOUR EMAIL!!!

    # Name of a secret used to store the ACME account private key
    privateKeySecretRef:
      name: letsencrypt-staging

    solvers:
      - dns01:
          webhook:
            groupName: cert-manager-webhook-inwx.smueller18.gitlab.com
            solverName: inwx
            config:
              ttl: 300 # default 300
              sandbox: false # default false

              # prefer using secrets!
              # username: USERNAME
              # password: PASSWORD

              usernameSecretKeyRef:
                name: inwx-credentials
                key: username
              passwordSecretKeyRef:
                name: inwx-credentials
                key: password
```

### Credentials
For accessing INWX DNS provider, you need the username and password of the account.
You have two choices for the configuration for the credentials, but you can also mix them.
When `username` or `password` are set, these values are preferred, and the secret will not be used.

If you choose another name for the secret than `inwx-credentials`, ensure you modify the value `credentialsSecretRef` in `values.yaml`.

The secret for the example above will look like this:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: inwx-credentials
stringData:
  username: USERNAME
  password: PASSWORD
```

### Create a certificate

Finally you can create certificates, for example:

```yaml
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: example-cert
  namespace: cert-manager
spec:
  commonName: example.com
  dnsNames:
    - example.com
  issuerRef:
    kind: ClusterIssuer
    name: letsencrypt-staging
  secretName: example-cert
```

## Development

### Requirements

-   [go](https://golang.org/) >= 1.13.0

### Running the test suite

1. Download test binaries
    ```bash
    scripts/fetch-test-binaries.sh
    ```

1. Create a new test account at [https://ote.inwx.com/en/customer/signup](https://ote.inwx.com/en/customer/signup) or use an existing account

1. Go to [https://ote.inwx.de/en/nameserver2#tab=ns](https://ote.inwx.de/en/nameserver2#tab=ns) and add a new domain

1. Copy `testdata/config.json.tpl` to `testdata/config.json` and replace username and password placeholders

1. Copy `testdata/secret-inwx-credentials.yaml.tpl` to `testdata/secret-inwx-credentials.yaml` and replace username and password placeholders

1. Download dependencies
    ```bash
    go mod download
    ```

1. Run tests with your created domain
    ```bash
    TEST_ZONE_NAME="$YOUR_NEW_DOMAIN." go test -cover .
    ```

### Building the container image

```bash
docker build -t registry.gitlab.com/smueller18/cert-manager-webhook-inwx:master .
```

### Running the full suite with microk8s

Tested with Ubuntu:

```bash
sudo snap install microk8s --classic
sudo microk8s.enable dns rbac
sudo microk8s.kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.12.0/cert-manager.yaml
sudo microk8s.config > /tmp/microk8s.config
export KUBECONFIG=/tmp/microk8s.config
helm install --namespace cert-manager cert-manager-webhook-inwx deploy/cert-manager-webhook-inwx
```
