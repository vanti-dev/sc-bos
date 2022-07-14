Smart Core Permissions
======================

## Attribute-based Access Control using Open Policy Agent
Smart Core APIs can be protected by [Open Policy Agent](https://www.openpolicyagent.org/).
OPA runs in-process, as a library. 

Rules are written in Rego, OPA's logic programming language. 
Each gRPC service has a package of the same name in Rego. That package should define an `allow` rule.
Access will be permitted only if the rule evaluates to `true`. This is enforced across all gRPC services
using middleware.

Rules have access to an `input` object, containing request metadata.

```json
{
  "Authorization": {
    "Issuer": "http://keycloak/realms/smart-core",
    "Subject": "(user or app identifier)",
    "Roles": ["Test.User"],
    "Scopes": ["Test.Read", "Test.Write"],
    "IsService": false
  },
  "Method": "UpdateTest",
  "Service": "vanti.ew_auth_poc.TestApi",
  "Request": {
    "test": {
      "data": "foobar"
    } 
  }
}
```

That input data can be checked by a Rego policy package, like so:

```rego
package vanti.ew_auch_poc.TestApi

import future.keywords.in

default allow := false

has_role[role] {
  some role in {"Test.Admin", "Test.User"}
  input.Authorization.Roles[_] == role
}

valid_scopes := {"Test.Read", "Test.Write"}

has_scope[scope] {
  some scope in valid_scopes
  input.Authorization.Scopes[_] == scope
}

# some service tokens may omit the list of scopes
# in this case we consider the service to have all scopes implicitly
has_scope[scope] {
  some scope in valid_scopes
  input.Authorization.IsService
  count(input.Authorization.Scopes) == 0
}

allow {
  input.Method == "GetTest"

  has_role[_] # has any valid role
  has_scope["Test.Read"]
}

allow {
  input.Method == "UpdateTest"
  
  has_role[_]
  has_scope["Test.Write"]
}
```

As can be seen above, the `Authorization` object contains data decoded and verified from the access token.
However, it should be noted that Rego is powerful enough to verify JWTs itself, so we could instead
move that logic into the policy. We would simply pass the access token in to the policy, which would
perform all authorization checking. This would allow us to add another authorization scheme without
rebuilding the app - it could even be done while the application is running (!).

There is also persistent `data`, which is shared between requests. The `data` hierarchy contains
both rules defined by packages, and any data imported into OPA. Data can be imported into the system while it is
running, if needed. OPA generally works with in-memory data, but it can be cached on the disk so it's still
available after a reboot, for example. For Smart Core, we might store the assignments of tenants to areas of the 
building in the data; the tenant data will be distributed to area controllers in the background