package vanti.bsp.ew.TestApi

import future.keywords.in

test_roles := {"Test.User", "Test.Admin"}
test_scopes := {"Test.Read", "Test.Write"}

token_has_role[role] {
  # the role is a valid test role
  role in test_roles
  # we have a valid token
  input.token_valid
  # the access token contains that role
  input.token_claims.Roles[_] = role
}

token_has_scope[scope] {
  # the scope is a valid test scope
  scope in test_scopes
  # we have a valid token
  input.token_valid
  # the access token contains that scope
  input.token_claims.Scopes[_] = scope
}

# service tokens don't use scopes
token_has_scope[scope] {
  test_scopes[scope]
  input.token_valid
  input.token_claims
}

# Admin user can write any data they want
valid_write_data {
  token_has_role["Test.Admin"]
}

# Other users can only write polite messages
valid_write_data {
  startswith(input.request.test.data, "please")
}

default allow := false

allow {
  input.method == "GetTest"
  token_has_role[_]  # any test_role will do
  token_has_scope["Test.Read"]
}

# allow any validated client certificate to access GetTest
allow {
  input.method == "GetTest"
  input.certificate_valid
}

allow {
  input.method == "UpdateTest"
  token_has_role[_]
  token_has_scope["Test.Write"]

  # check that data being written is OK
  valid_write_data
}