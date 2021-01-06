apiVersion: v1
kind: Secret
metadata:
  name: inwx-credentials
data:
  username: $INWX_USER_OTP_BASE64
  password: $INWX_PASSWORD_OTP_BASE64
  otpKey: $INWX_OTPKEY_BASE64
