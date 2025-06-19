package smartcore.traits

import future.keywords.in

import data.scutil.token.token_has_permission
import data.scutil.rpc.read_request

allow {
  token_has_permission("trait:read")
  read_request
}

allow {
  token_has_permission("trait:write")
}

allow {
  token_has_permission("trait:*")
}