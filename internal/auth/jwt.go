package auth

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrWrongAudience = errors.New("token audience does not match expected audience")
	ErrWrongIssuer   = errors.New("token issuer does not match expected issuer")
	ErrExpired       = errors.New("token expired")
	ErrNotYetValid   = errors.New("token not yet valid")
)

type JWTCommonClaims struct {
	Issuer     string  `json:"iss"`
	Subject    string  `json:"sub"`
	Audience   string  `json:"aud"`
	Expiration JWTTime `json:"exp"`
	NotBefore  JWTTime `json:"nbf"`
	IssuedAt   JWTTime `json:"iat"`
}

// JWTTime is a wrapper for time.Time that marshals and unmarshals to JSON in the format expected by the JWT standard
// for standard claims - integer UNIX timestamps, in seconds, UTC.
type JWTTime time.Time

func (t *JWTTime) UnmarshalJSON(data []byte) error {
	var unix int64
	err := json.Unmarshal(data, &unix)
	if err != nil {
		return err
	}

	*t = JWTTime(time.Unix(unix, 0))
	return nil
}

func (t *JWTTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*t).Unix())
}

// JWTClaimVerifier checks that JWTCommonClaims are correct.
type JWTClaimVerifier struct {
	// Now is used to obtain the current time, used for the purpose of checking the time-based claim exp and nbf.
	// Optional - if nil, time.Now is used.
	Now func() time.Time
	// Audience is used to check the aud claim of the token for an exact match.
	// Optional - if Audience is empty, the aud claim is not checked.
	Audience string
	// Issuer is used to check the iss claim of the token for an exact match.
	// Optional - if Issuer is empty, then the iss claim is not checked.
	Issuer string
}

func (v *JWTClaimVerifier) Verify(claims JWTCommonClaims) error {
	if v.Audience != "" && claims.Audience != v.Audience {
		return ErrWrongAudience
	}
	if v.Issuer != "" && claims.Issuer != v.Issuer {
		return ErrWrongIssuer
	}

	var now time.Time
	if v.Now != nil {
		now = v.Now()
	} else {
		now = time.Now()
	}

	if now.Before(time.Time(claims.NotBefore)) {
		return ErrNotYetValid
	}
	if now.After(time.Time(claims.Expiration)) {
		return ErrExpired
	}
	return nil
}

type JWTScopes []string

func (s *JWTScopes) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.Join(*s, " "))
}

func (s *JWTScopes) UnmarshalJSON(data []byte) error {
	var joined string
	err := json.Unmarshal(data, &joined)
	if err != nil {
		return err
	}
	*s = strings.Fields(joined)
	return nil
}

func (s *JWTScopes) HasScopes(required ...string) bool {
	return RequireAll(required, *s)
}
