apiVersion: v1
kind: Secret
metadata:
  name: inwx-credentials
data:
  username: $INWX_USER_BASE64
  password: $INWX_PASSWORD_BASE64
