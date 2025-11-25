package policy

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/smart-core-os/sc-bos/pkg/auth/token"
)

func BenchmarkPolicy(b *testing.B) {
	attrs := Attributes{
		Service:    "smartcore.traits.OnOff",
		Method:     "UpdateOnOff",
		Request:    json.RawMessage(`{"name":"test","onOff":{"state":"ON"}}`),
		TokenValid: true,
		TokenClaims: token.Claims{
			SystemRoles: []string{"admin"},
		},
	}
	run := func(b *testing.B, policy Policy) {
		for i := 0; i < b.N; i++ {
			result, err := policy.EvalPolicy(context.Background(), "data.smartcore.allow", attrs)
			if err != nil {
				b.Error(err)
			}

			if !result.Allowed() {
				b.Errorf("expected interation %d to suceed", i)
			}
		}
	}

	b.Run("static", func(b *testing.B) {
		policy := Default(false)
		run(b, policy)
	})

	b.Run("cachedStatic", func(b *testing.B) {
		policy := Default(true)
		run(b, policy)
	})
}
