package policy

import (
	"context"
	"testing"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
)

func BenchmarkPolicy(b *testing.B) {
	attrs := Attributes{
		Service:    "smartcore.traits.OnOff",
		Method:     "UpdateOnOff",
		Request:    &traits.UpdateOnOffRequest{Name: "test", OnOff: &traits.OnOff{State: traits.OnOff_ON}},
		TokenValid: true,
		TokenClaims: token.Claims{
			Roles:  []string{"admin"},
			Scopes: []string{"Read", "Write"},
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
