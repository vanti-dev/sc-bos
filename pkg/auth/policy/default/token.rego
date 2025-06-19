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

token_permission_assignments = assignments {
  claims := valid_claims
  assignments := claims.permissions
}

token_has_role(role) {
  role in roles
  claims := valid_claims
  claims.system_roles[_] == role
}

token_has_permission(permission) {
  some assignment in token_permission_assignments
  assignment.permission == permission
  not assignment.scoped
}

token_has_permission(permission) {
  some assignment in token_permission_assignments
  assignment.permission == permission
  assignment.scoped
  assignment.resource_type == "NAMED_RESOURCE_PATH_PREFIX"
  startswith(input.request.name, concat("/", [assignment.resource, ""]))
}

token_has_permission(permission) {
  some assignment in token_permission_assignments
  assignment.permission == permission
  assignment.scoped
  assignment.resource_type in ["NAMED_RESOURCE", "NAMED_RESOURCE_PATH_PREFIX"]
  input.request.name == assignment.resource
}