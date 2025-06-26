package smartcore

import data.scutil.token.token_has_role
import data.scutil.rpc.rpc_match_verbs
import data.scutil.rpc.read_request
import data.scutil.rpc.write_request

# Common rules for services that follow Smart Core conventions.

# admin based access is unrestricted
allow {token_has_role("admin")}
allow {token_has_role("super-admin")}
# certificate based access is unrestricted, this may change in future
allow {input.certificate_valid}

# Commissioners can do anything apart from create new users
allow { token_has_role("commissioner") }
# Operators can read or update anything by default.
# SerivceApi has special rules for this role.
allow { token_has_role("operator"); read_request }
allow { token_has_role("operator"); write_request }
# Viewers have read-only access, no write or mutation rpcs allowed
allow {
  token_has_role("viewer")
  read_request
}
