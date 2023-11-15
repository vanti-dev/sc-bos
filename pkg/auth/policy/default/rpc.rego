package scutil.rpc

import future.keywords.in

read_verbs = {
  "Get",
  "List",
  "Pull",
  "Describe"
}

write_verbs = {
  "Create",
  "Update",
  "Delete"
}

verb_match(verbs) {
  regex.match(concat("", ["^(", concat("|", verbs), ")[A-Z]"]), input.method)
}

read_request {
  verb_match(read_verbs)
}

write_request {
  verb_match(write_verbs)
}

rpc_match(service, method) {
  input.service == service
  input.method == method
}

rpc_match_methods(service, methods) {
  input.service == service
  input.method in methods
}

rpc_match_verbs(service, verbs) {
  input.service == service
  verb_match(verbs)
}
