package jwks

import (
	"encoding/json"
	"testing"

	"github.com/go-jose/go-jose/v3"
)

// keys for testing signature verification
var testJWK1, testJWK2 jose.JSONWebKey

func init() {
	err := json.Unmarshal([]byte(`
{
    "p": "-gmxNRWjTfO-CEzs-78nG8Jzf512kEMhNjXLGNlYVwCvREkc4ZJ5MCRavNy8qD9LdFi1j4Brd-y971F5xt7sjaMabnyxqEtqdrW-ttxItqSjHR1t0-7vth0wx8KCQVHost6H67TLcPtRiluzewHvECTvELYFsRxhA75Hn9Ddix8",
    "kty": "RSA",
    "q": "4ph81fDN-oumRpfiExI174Ug50S6twRlb5vCyh1HqgihVn4A5kp8UGQLzt1PbZlmu2lRbWJFvURt8VvCFZVDevwQ6a7NWJSyxiADrD1vPZqm8J0O3IR631M9b1TDAeNTYmsZIRVwODkC_G3vg0la9I6N-1tRMIpFd1OZKOj3ciE",
    "d": "AqeBjbvgnVVbiVc8OKW_t_A9so2ssNX3ldd-HIwE8ANZXKdAB9V00XncwChJ-GBxI4xSE8CjpBJs6skksf-dtH6k7L6UfUTmYJdq-tgyk14nMk1zUqgBJ-Lihd1VwDdP8_eZU1inzMvcT5xyKijZahhoJOdvIgTepHScrpXVCclhtVEYDgToazdj4iml7kK2xf-s3uyeW5EACEtVjuiGYmg_CzSXyBU-q20wZuNcZNM0WOc4cbbn1isUSlk1lsXwz03PJwhiTBzbNR0nD4yfoAzAwJB4b67L_WbHNUz4O2pjyOaqnOFZoJKYRQALxzeszNLXoiI-eTlap8lGPVe7gQ",
    "e": "AQAB",
    "use": "sig",
    "kid": "test-key-1",
    "qi": "DwfxiIUpDRZV2dXqYcWwdt8AxktueInhI4pNNIl1dUNFEA2i-LYmMLHiCcko6mXHCdm1Rs2AoSWEwG9pa0sOQkU3G7cZfLP904lumLCu6ftBYlT-WUCmOuYeSxJ7TaizvPKWNQWV_Z_bYuLOxe3YrGCOpABQrdHgenjHIzVjccY",
    "dp": "yMOrhCJBo7_IoEWUK3eK4WE6-AbpQmCEdFCxKNyrcABeuoeyJvVDVYJ7URY0bSuVXHA2KGlG4V44C8bx7trkOb3y5TA-PhGABJ1d6tnpkK2VQzV0EC3UT_gUSPFHQUeRfr3riTj7-VXyXRPQgz5ERERDqLlezJ0q0KSiQhKlMKU",
    "alg": "RS256",
    "dq": "ZUVVdaBbzoAfXil_Zpqa9GOBYxr6f9U9KHZqxj3zy3Bz-t3xtPrROHSeOmP6nbcTjOry83oaRQ6SPG6P_WlqcUq6nFX9fHtosteYDKCgWN4Hgj4PaErlR25CZMFzLiLVH4VSA9E7CEWiqgLQKtLcDbSwjAgx7wm9Jil8qCYGgUE",
    "n": "3VF-H7u3cAgZocxllIJmKvqrk6nE4419qGmC8bY5VW07RzOu7Zf9C_Jk-SxWM5hFC4yUivs-HVWnzdynyBS_jCom4oaXE1fS40WPZLLyLdXXwvqqs5EzbMwFdI3hPO0qGzZ8uaSGqUDZZTalAFSLDkXHGpHs7mm9_SbNlK_BV257JqzcfKjLeCoKWeh-LWu2Lat77HLV_wLvMmgOr6ZJpTRqZUm-9rKnLxzhvJNXssi3SLno_yRb-XMHkwPSBplssPVtSfX9x7yb7fp23v1ijZSTCwHmkSGFOAJIDY637xYKol3MJ18Ad38kS14y0UyhMwvaOlZth1O2t_ZJ52u8_w"
}
	`), &testJWK1)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(`
{
    "p": "5Pd8wH8cwazcnKvBHj1ta47XK_gJxjH0HqicAkToFkiVyZOJnstGTyieiu-u-frnoEu49SPBT-yoWhnc6ni5-VJ9a6akcM0aV9oGM5v552UqnVpPoW8gjfGwxmRZRQzYbackIFHhrNHpOHmgVerQYLr7n5FofzhaveRX4MbXx2c",
    "kty": "RSA",
    "q": "v_h6fCxa5I-ogMmisGWk6zIpFP4_zN8krX0XuNCJ4Lu7FRhwFLFBlF7nl4OrsBpQFRRhHIj0Dfu3YELEvwObFkCgh83h_898PcsUSyA0uQNHoQX7al-jJZarTWOCOJKaRW1W0fSv-kee_i_xxokr2m2doT33Pmmid1lHl5W_Qkc",
    "d": "oyXDG3_AAVhraiTKzIZw1tswO56-fICFTNfFxCDQ1HR_1S3I06U5xwDiN5-fpl0IU_47D9arpcZ3qHA8pt3HmsUMvCkQoNFjRMieutshvd2z6iaVsRcVqUZaDUliarZ2R6ty_4kPEBKzepwgq7LB6ewZpDE9ZAm6riPDgdr34jGqfCRd2olEejL26w0U7wv3USJW5CM5rBkftbiVeggdzP_nQskodRbKWY3kE2Pgweg1RsmqZIvafrQ9lTeSLRejX7wS0hTpd1HtiOgPiFg_jLD2K2jckfjBTPy9cNggl565UhpXqvlyaCbG18xYK1eC66-5i6HHlRWXE0cgeo225Q",
    "e": "AQAB",
    "use": "sig",
    "kid": "test-key-2",
    "qi": "Rt47hAJ1DkjU8mNJZBSNFkO_pmFNQbZoDwLZ6NalPPwnjDT_vVdq1Cw_WzORWPhz5pgyq2wzL0qi3G6rnrgxaDYAdr_38BE4aEqVJjKJOrPTk1uz6G41wuXh8c_7Hu_0Yj0htLyzQMf_B8afFkORVKgj-GjlWltKmEqX330jtVE",
    "dp": "WH0NSZfWlUMpP6NhTz6OOzNJFUUXAfHsVqzzHi1jRLloqi7K0QPeeFlKbIeVKCc_vUOGh7b5ztm3dproNfXSafjnX-NXSgD6XVl1bByryDHg9k8g11MLUdBGcWX22ijMvBQMcjEy9odpitn2jT3iqn-ZH2Ii8IfnCdxl2gj--6E",
    "alg": "RS256",
    "dq": "oQwPlYSQbBaowgJmXZ2oETfvhxEU7QZ2eqTq9bzdLo_Pjw8FWBascZB8sXtg2Uf5zvVd0taCCAkX-cWJ0MVxoeVtxwBNjJHAJQbta2kFUgESYl_mX4MEF1CjPTUx1cwHaB8mKtUfnNPg6lXGe0wwYfp7tv2JIe70wTNBAEY8QZk",
    "n": "q7LjYXpRlwy33ijf76CMv0YEirLXperNoEsxAYkoLFJAsDN-WNyycr88v9LXsTiOkkC13BHB0eRr6mixhdMtJPbNQ4VuhZr2PtbX7bjkPjROfDwAo9T-mYKr4GSRj2vx3f5GshPZm8p9gxme129iUuz-TarWE-80qAgyHFWhLpw1uTtzHEajjtjoneDgyYuhGK1xeiRSCH73M05tCk-kXfhQLnhr65awoAnWzF0k1ZrcXQsGFjh7JJCqac7CtX52xM_QDW8qgkizzmAl-1RJUVMKhsivLwCwJfb3cT0A-FxNjTeHwShKiEGz67MwBATGUYrEZwTzpEFCz3NogBnbkQ"
}
	`), &testJWK2)
	if err != nil {
		panic(err)
	}
}

func signJWS(t *testing.T, key jose.JSONWebKey, payload []byte) (jws string) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, nil)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := signer.Sign(payload)
	if err != nil {
		t.Fatal(err)
	}
	jws, err = signed.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return jws
}
