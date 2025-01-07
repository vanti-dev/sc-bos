package http

import data.scutil.token.token_has_role

allow {
  log_level_permission
  input.path == "/__/log/level"
}

allow {
  # everyone can get the service version
  input.path == "/__/version"
}

log_level_permission { token_has_role("admin") }
log_level_permission { token_has_role("super-admin") }