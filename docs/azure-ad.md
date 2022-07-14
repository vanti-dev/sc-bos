Microsoft Identity Platform
===========================
Azure AD issues access tokens in the JSON Web Token format, signed using the `RS256` algorithm (asymmetric).
The signing public keys can be downloaded from Microsoft.

[Reference Documentation](https://docs.microsoft.com/en-us/azure/active-directory/develop/access-tokens)

### Important Claims

| Claim   | Title      | Description                                                       |
|---------|------------|-------------------------------------------------------------------|
| `aud`   | Audience   | The UUID of the application this token was issued for             |
| `nbf`   | Not Before | Unix timestamp - token is not valid before this time              |
| `exp`   | Expiry     | Unix timestamp - token is not valid after this time               |
| `scp`   | Scopes     | Scopes enabled on this token - User tokens only. Space seperated. |
| `roles` | Roles      | App roles issued to token's subject (user or app)                 |
| `idtyp` | ID Type    | For app tokens, has the value `"app"`. Absent for user tokens.    |

### Validating Tokens

Things that need to be checked:
  - Token needs a valid signature from Microsoft's public key
  - `aud` claim must match the resource server's App ID
  - Current time must be after `nbf` and before `exp`
  - Roles must allow the current operation
  - (User Tokens Only) Scopes must allow the current operation