package policy

import (
	"context"
	"testing"

	"github.com/vanti-dev/sc-bos/pkg/auth/token"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func BenchmarkPolicy(b *testing.B) {
	attrs := Attributes{
		Service: "vanti.bsp.ew.TestApi",
		Method:  "UpdateTest",
		Request: &gen.UpdateTestRequest{
			Test: &gen.Test{Data: "please foobar"},
		},
		TokenValid: true,
		TokenClaims: token.Claims{
			Roles:  []string{"Test.User"},
			Scopes: []string{"Test.Read", "Test.Write"},
		},
	}
	run := func(b *testing.B, policy Policy) {
		for i := 0; i < b.N; i++ {
			result, err := policy.EvalPolicy(context.Background(), "data.vanti.bsp.ew.TestApi.allow", attrs)
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
