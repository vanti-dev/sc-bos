package scutil.token

import future.keywords.in

roles := {
  "admin", # unrestricted user access
  "commissioner",
  "operator",
  "signage",
  "super-admin",
  "viewer"
}

valid_claims = claims {
  input.token_valid
  claims := input.token_claims
}

token_roles = roles {
  claims := valid_claims
  roles := claims.roles
}

token_zones = zones {
  claims := valid_claims
  zones := claims.zones
}

token_has_role(role) {
  role in roles
  claims := valid_claims
  claims.roles[_] == role
}

# match if the current request name is equal to or a sub-name of one of the tokens zones
token_matches_zone {
  some zone in token_zones
  startswith(input.request.name, concat("/", [zone, ""]))
}
token_matches_zone {
  token_zones[_] == input.request.name
}
