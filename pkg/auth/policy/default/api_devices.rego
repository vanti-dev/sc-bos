package smartcore.bos.DevicesApi

# Only allow requests if they have valid authentication.
allow {
  input.certificate_valid
}
allow {
  input.token_valid
}
