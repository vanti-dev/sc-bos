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

allow {
  pprof_permission
  startswith(input.path, "/__/debug/pprof/")
}

log_level_permission { token_has_role("admin") }
log_level_permission { token_has_role("super-admin") }

pprof_permission { token_has_role("admin") }
pprof_permission { token_has_role("super-admin") }
