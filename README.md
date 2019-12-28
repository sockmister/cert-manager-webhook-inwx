# ACME Webhook for INWX

## Requirements
-   [go](https://golang.org/) >= 1.13.0
-   [helm](https://helm.sh/) >= v3.0.0
-   [kubernetes](https://kubernetes.io/) >= v1.14.0
-   [cert-manager](https://cert-manager.io/) >= 0.12.0

## Installation

### cert-manager

Follow the [instructions](https://cert-manager.io/docs/installation/) using the cert-manager documentation to install it within your cluster.

### Webhook

```bash
helm install --namespace cert-manager cert-manager-webhook-inwx deploy/cert-manager-webhook-inwx
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
You have two choices for the configuration for the credentials but you can also mix them.
When `username` or `password` are set, theses values are preferred and the secret will not be used.

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
    name: letsencrypt-staging
  secretName: example-cert
```

## Development

### Running the test suite

1. Download test binaries
    ```bash
    scripts/fetch-test-binaries.sh
    ```

1. Create a new test account at [https://ote.inwx.com/en/customer/signup](https://ote.inwx.com/en/customer/signup) or use an existing account

1. Go to [https://ote.inwx.de/en/nameserver2#tab=ns](https://ote.inwx.de/en/nameserver2#tab=ns) and add a new domain

1. Copy `testdata/config.json.tpl` to `testdata/config.json` and replace username and password placeholders

1. Download dependencies
    ```bash
    go mod download
    ```
   
1. Run tests with your created domain
    ```bash
    go TEST_ZONE_NAME="$YOUR_NEW_DOMAIN." test .
    ```
   
### Building the container image

```bash
docker build -t registry.gitlab.com/smueller18/cert-manager-webhook-inwx .
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