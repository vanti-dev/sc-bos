package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// SASKey is a secret for performing HMAC operations.
// Standard textual representation is base64-encoded.
type SASKey []byte

func ParseSASKey(secret string) (SASKey, error) {
	var sas SASKey
	err := sas.UnmarshalText([]byte(secret))
	if err != nil {
		return nil, err
	}
	return sas, nil
}

//goland:noinspection GoMixedReceiverTypes
func (k *SASKey) UnmarshalText(text []byte) error {
	maxLen := base64.StdEncoding.DecodedLen(len(text))
	newSas := make([]byte, maxLen)

	actualLen, err := base64.StdEncoding.Decode(newSas, text)
	if err != nil {
		return err
	}
	newSas = newSas[:actualLen]
	*k = newSas
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (k SASKey) MarshalText() ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(k)))
	base64.StdEncoding.Encode(buf, k[:])
	return buf, nil
}

// GenerateSASToken will generate a token for accessing Azure IoT Hub, from a base64-encoded key.
// The token is scoped to a specific resource tree identified by uri - only resources which have uri as a prefix can be
// accessed with this token.
// Device SAS tokens must be device-scoped e.g. hubname.azure-devices.net/devices/device1
// FormatSASURI can be used to generate an appropriate URI.
//
// policyName is an optional shared access policy to referred to - omit when using device-registry credentials.
//
// Algorithm specified at https://learn.microsoft.com/en-us/azure/iot-hub/iot-hub-dev-guide-sas?tabs=node
func GenerateSASToken(ctx context.Context, signer Signer, uri, policyName string, expires time.Time) (string, error) {
	expiresUnix := expires.Unix()

	// HMAC-256 of "{urlencoded uri}\n{expiry timestamp}"
	signPayload := fmt.Sprintf("%s\n%d", url.QueryEscape(uri), expiresUnix)
	rawSignature, err := signer.Sign(ctx, []byte(signPayload))
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(rawSignature)

	const (
		paramSig        = "sig"
		paramExpiry     = "se"
		paramPolicyName = "skn"
		paramURI        = "sr"
	)

	values := make(url.Values)
	values.Set(paramSig, signature)
	values.Set(paramExpiry, strconv.FormatInt(expiresUnix, 10))
	if policyName != "" {
		values.Set(paramPolicyName, policyName)
	}
	values.Set(paramURI, uri)

	return "SharedAccessSignature " + values.Encode(), nil
}

func FormatSASURI(hostname, deviceID, moduleID string) string {
	if moduleID != "" {
		return fmt.Sprintf("%s/devices/%s/modules/%s", hostname, deviceID, moduleID)
	} else {
		return fmt.Sprintf("%s/devices/%s", hostname, deviceID)
	}
}
