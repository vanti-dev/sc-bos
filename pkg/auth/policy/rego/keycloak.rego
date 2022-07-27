package keycloak

issuer := "http://localhost:8888/realms/smart-core"
client_id := "sc-apps"

oidc_metadata_endpoint := concat("/", [issuer, ".well-known/openid-configuration"])

fetch_oidc_metadata(url) := http.send({
  "url": url,
  "method": "GET",
  "force_cache": true,
  "force_cache_duration_seconds": 86400 # 24 hours
}).body

oidc_metadata := fetch_oidc_metadata(oidc_metadata_endpoint)

jwks_endpoint := oidc_metadata.jwks_uri

jwks_request(url, force_cache) := res.body {
  res := http.send({
    "url": url,
    "method": "GET",
    "force_cache": force_cache,
    "force_cache_duration": 3600 # 1 hour
  })
}

jwt_header := io.jwt.decode(input.access_token)[0]

# Fetch the JWKS to verify the token against
# First, try with an agressive caching strategy to reduce load on the authorization server.
# If the cached JWKS lacks the key we need, make another request.
jwks := cached_jwks {
  cached_jwks := jwks_request(jwks_endpoint, true)
  # check that the JWKS has the key we want
  jwt_header.kid == cached_jwks.keys[_].kid
} else = fresh_jwks {
  fresh_jwks := jwks_request(jwks_endpoint, false)
}

now := time.now_ns()
verified_claims := payload {
  [valid, _, payload] := io.jwt.decode_verify(input.access_token, {
     "cert": jwks,
     "alg": "RS256",
     "iss": issuer,
     "aud": client_id,
     "time": now
  })
  valid
}