Security Definitions
====================
*Definitions of terms used or referred to in the Smart Core security architecture*

## Access Control
An overarching term for how a system decides who can access a system, and which operations they may perform on
which resources.

  - **Subject** - the entity attempting to perform an operation e.g. a user
  - **Object**  - the targeted resource of the operation e.g. a device
  - **Action**  - what the subject is attempting to do to the object e.g. turn on the light

### Authentication
Determining the identity of an subject, such as a person or service.
Usernames & passwords, certificates and API secrets are examples of ways for a subject to authenticate with a server.

### Authorization
Determining if the subject making the request has permission to do so. This could be achieved by looking the subject up
in an Access Control List, checking the roles encoded in an access token etc.

As the permissions available to a subject are very often linked to who they are, Authorization may depend on
Authentication.

#### Role-based Access Control
A scheme where permissions (such as "turn on the lights in room 4") are assigned to roles, and roles are assigned to
subjects. A subject has permission to perform an operation if one of their assigned roles has permission.

#### Attribute-based Access Control
A scheme where the attributes of the subject, object and action involved in a request are evaluated by a policy,
which decides whether to allow or deny the request.

Policies may be expressed as a list of rules, or code (e.g. Rego in Open Policy Agent). If the policy language
is expressive enough, Role-based Access Control can be implemented within a policy.

  - **Policy** - a set of rules that take attributes as input and outputs a decision
  - **Policy Decision Point** - the component of the system which evaluates requests against the policies
  - **Policy Enforcement Point** - the component of the system which protects the resources based on the PDP's decisions
  - **Attributes** - pieces of data describing the subject, object, action and request context
    - Subject attributes: Job title, list of roles, clearance level, place of work, employer
    - Object attributes: Object type, owner, department, location
    - Action attributes: Type of operation e.g. read / update / delete
    - Contextual: Time of day