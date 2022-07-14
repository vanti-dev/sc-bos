package microsoft

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/vanti-dev/ew-auth-poc/pkg/auth"
)

var testJWK jose.JSONWebKey

func init() {
	err := json.Unmarshal([]byte(`
{
    "p": "2XF4loLMoILQBaCSUOgiNWPWSXcfWHJuuVIkjZtZLPjVpzJyAM2RiTJlSHbNJgQ8x3IPCZsz0hZm7S7fBuIP5xXncgqI0bDD95_LnXLVzEt9ai2_OUZqRJaU4R2z8ekU_KxHdEqI2Y-honzSwTJdSU66PVL491TwVzhFfO-Eccc",
    "kty": "RSA",
    "q": "w1-Ni5LXupWaprt3TSbysW_0PSQvQ3l_nddqjwHuExzyvMaIJSBdkE5435SRF4otv8siEaNoEc3rv-RmB7-a9PQX3w9HUuH7xVpSc6z77DOQbsAyFMSsPm_TU-wkjuET4AHdH6dwLLoRLzSNn5zg8s3yIHW9sEacIkd34kR2C98",
    "d": "QblHFUjpsi7aObwyzfCLxeSch-3wA3P0AwAL8QiA8CoHHfDnLUyaobCjGJIX2dzRQ9bTRLk1Eip_XR0eUFYyYOdlnyQ4o1w_s4Yh0I1G7wBIWQ_R6d91Pko0VBbPUAzjUgp_wSf8Mwd5WKi3sTeNpbzmRmo_Qiw-oMrDkUSSDhuExsxIMv4Nsuvc536csDHrdUsHMBoH4znJwo8cU2rf8bGEc7RBClCtB9QB7vlpG_Cmfy_56xOlxtoa7cSx6aVTWqQA-dcnELjlNOEkxab0VvQYDBHpcrJNVxRk0YSzqwUwRSP_XX_bXeHDdtDjjsiu7zXcOZYwf-1ZWq8dIcYWUQ",
    "e": "AQAB",
    "use": "sig",
    "kid": "test-microsoft",
    "qi": "XR2b4iAvOg70i5jfstaX9xUmXqwJWn-xMC06EUlX-ASTGeivYNLNtDBxdiEzZdYDi6Ff-3CTvujZNhLUzHe3ZufnIkk-c9Pg8l3LT1XIMqiBMth56AFwKcNgpi2VTXFMmluteq9Zm8fgQdzF0pVqhHhVB6etKXLydMTX4xI2FIk",
    "dp": "DVPwMUGJK1l5SK8D6OOhnSYmb1BW4oP5F3DafreD6sbBycLEzBnNjtxA6wUlI-xkkVEDbPJPZdQrsOQLlY9rrB1il1Hf-wJbsKThxR_RzPjfkg-Fxgvz3YphS78XVX-U4rvokU80PimYna4K-P6OSz0BG1snmXliVeZEFBlWQ-c",
    "alg": "RS256",
    "dq": "eDXe3bYoToGmelh2e14vjcWYvdX5zsZ_IxtcUUmatt_k6wos0ssrRpNCBC9hZt56b7wI5llouyu8loFt1e6q5XUbCwBXnkO1qiR8_ve8ugSwJbTUG7s3T-N2X_i4NcF_fSEocUWQ27RQxn7LdR4NutfE1vwbDs8jWmQm-22sRHk",
    "n": "pfKYMHI9u4SPh6hiQUkey89DEpsfTHTZhksGrY0O4Te82ntgIBJiN5RsgNILFdzXGGvZ7Mq4CNU12OiPa4yrwy742EpnM1_bi3BzGyxMf7xfuOkHeRb9ztaP-w9hULgxOYm-Esa63xgl3-gswccZk8VBzFLnYvBssuLjcAD6x4cGr_lhZZe8-96KToA_4G-3H3GlhyjQhc7UG7hd7mwDQbTiHFUhYN9uVC3742W8rpYbTGett9vFOS6ESAM1rAkeZVxCrgh9BaGJluGGevL_ct4uJWQfk1m8wYrD_VrIgnIJjzc0NiMWmTRFuaATe0r-mzOjudMepz6cN24SovypWQ"
}
`), &testJWK)
	if err != nil {
		panic(err)
	}
}

func TestTokenVerifier_VerifyAccessToken(t *testing.T) {
	now := time.Date(2022, 05, 31, 15, 32, 0, 0, time.UTC)

	validAudience := "c1ba833e-c813-4d3e-aabe-5ddcb59eb66f"
	invalidAudience := "c96119b5-3da0-4e3f-8377-8528faeba174"

	keySet := auth.NewLocalKeySet(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{testJWK.Public()}})
	verifier := NewTokenVerifier(keySet, &ClientConfig{
		Tenant:   "64a32c38-9418-4c14-8e1a-5ecc5da2c005",
		ClientId: validAudience,
	}, WithClock(func() time.Time { return now }))

	type testCase struct {
		expect error
		iat    time.Time
		nbf    time.Time
		exp    time.Time
		aud    string
	}

	cases := map[string]testCase{
		"Valid_Now": {
			iat: now.Add(-time.Hour),
			nbf: now.Add(-time.Hour),
			exp: now.Add(time.Hour),
			aud: validAudience,
		},
		"Expired": {
			expect: auth.ErrExpired,
			iat:    now.Add(-2 * time.Hour),
			nbf:    now.Add(-2 * time.Hour),
			exp:    now.Add(-time.Hour),
			aud:    validAudience,
		},
		"Not_Valid_Yet": {
			expect: auth.ErrNotYetValid,
			iat:    now.Add(-time.Hour),
			nbf:    now.Add(time.Hour),
			exp:    now.Add(2 * time.Hour),
			aud:    validAudience,
		},
		"Wrong_Audience": {
			expect: auth.ErrWrongAudience,
			iat:    now.Add(-time.Hour),
			nbf:    now.Add(-time.Hour),
			exp:    now.Add(time.Hour),
			aud:    invalidAudience,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			token := createAccessToken(t, testJWK, c.iat, c.nbf, c.exp, c.aud, "foo@example.com", "Foo Bar",
				nil, nil)

			_, err := verifier.VerifyAccessToken(context.Background(), token)
			if !errors.Is(err, c.expect) {
				t.Error(err)
			}
		})
	}
}

func createAccessToken(t *testing.T, key jose.JSONWebKey, iat time.Time, nbf time.Time, exp time.Time, aud string,
	email string, name string, roles []string, scopes []string) string {

	payload := AccessTokenPayload{
		JWTCommonClaims: auth.JWTCommonClaims{
			Audience:   aud,
			IssuedAt:   auth.JWTTime(iat),
			NotBefore:  auth.JWTTime(nbf),
			Expiration: auth.JWTTime(exp),
		},
		PreferredUsername: email,
		Name:              name,
		Roles:             roles,
		Scp:               scopes,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, nil)
	if err != nil {
		t.Fatal(err)
	}

	signed, err := signer.Sign(payloadBytes)
	if err != nil {
		t.Fatal(err)
	}

	jws, err := signed.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return jws
}
