package smartcore.bos.HubApi

import data.scutil.token.token_has_role
import data.scutil.rpc.read_request
import data.scutil.rpc.verb_match

# Allow anybody to request information about nodes.
# This is useful for status monitoring.
allow { read_request }

allow {
  token_has_role("operator")
  verb_match({"Inspect", "Test"})
}