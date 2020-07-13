apiVersion: v1
kind: Secret
metadata:
  name: inwx-credentials
stringData:
  username: $INWX_USER
  password: $INWX_PASSWORD
