# Policy files for the Smart Core BOS

TL;DR admins can do anything anywhere, zones are restricted by name, roles restrict permissions.

Test these rules by running `opa test .` from this directory.

We've split the policies into rules and utilities:

- token.rego contains token related utils like "does the token have this role"
- rpc.rego contains utils for working with requests, for example "is this a write request"

The bulk of the rules are in `smartcore.rego`, which applies these rules to the standard verbs and conventions of Smart
Core apis. Exceptions are encoded into specific rpc rego files.