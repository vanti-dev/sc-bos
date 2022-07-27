package vanti.ew_auth_poc.TestApi

import future.keywords.in

test_roles := {"Test.User", "Test.Admin"}
test_scopes := {"Test.Read", "Test.Write"}

user_has_role[role] {
  # the role is a valid test role
  role in test_roles
  # the access token contains that role
  input.Authorization.Roles[_] = role
}

user_has_scope[scope] {
  # the scope is a valid test scope
  scope in test_scopes
  # the access token contains that scope
  input.Authorization.Scopes[_] = scope
}

# Admin user can write any data they want
valid_write_data {
  user_has_role["Test.Admin"]
}

# Other users can only write polite messages
valid_write_data {
  startswith(input.Request.test.data, "please")
}

default allow := false

allow {
  input.Method == "GetTest"
  user_has_role[_]  # any test_role will do
  user_has_scope["Test.Read"]
}

allow {
  input.Method == "UpdateTest"
  user_has_role[_]
  user_has_scope["Test.Write"]

  # check that data being written is OK
  valid_write_data
}