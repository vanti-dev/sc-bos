# Authorization and Access Control

*TODO*: This document has not been updated to reflect the new authentication model.

Authorization (authz) is process of taking some confirmed identity and determining what that identity is allowed to do.
Access control is the process of enforcing those permissions.

See [the Authenentication Schedule](./authn-schedule.md) for a available mechanisms for identifying and validating
external entities.

In Smart Core BOS, authz is configured and enforced using a collection of policy rules that say things like

> "If the user has the admin role, allow access to write apis"

We use [Open Policy Agent](https://www.openpolicyagent.org/) (OPA) to implement our authz policies. The default policies
can be found in [the `auth/policy/default` package](../pkg/auth/policy/default).

In general our authz policy is:

1. Everybody can see the status of the system, even unauthenticated users. This includes the reflect api and read-only
   hub and enrollment apis needed to understand if things are working or not.
2. Your `roles` define what you can do across all devices and APIs in the system.
    1. An `admin` can do anything with any API.
    2. A `viewer` can only read from any API, they can't update or stop anything.
3. Your `zones` define which devices you can access
    1. Having `zones=["FOO"]` grants you access to `UpdateLight({name:"FOO"})` but not `UpdateLight({name:"BAR"})`
    2. Zones are implicit prefix wildcards, zone `FOO` grants access to `FOO` and `FOO/BAR` but not `FOOBAR`

