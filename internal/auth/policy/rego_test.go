package policy

import (
	"context"
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/rego"
	"github.com/vanti-dev/bsp-ew/internal/auth"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
)

var exampleAttributes = Attributes{
	Service: "vanti.bsp.ew.TestApi",
	Method:  "UpdateTest",
	Request: &gen.UpdateTestRequest{
		Test: &gen.Test{Data: "please foobar"},
	},
	TokenValid: true,
	TokenClaims: auth.TokenClaims{
		Issuer:  "foobar",
		Subject: "barfoo",
		Roles:   []string{"Test.User"},
		Scopes:  []string{"Test.Read", "Test.Write"},
	},
}

func BenchmarkLoadRegoCached(b *testing.B) {
	attr := exampleAttributes

	for i := 0; i < b.N; i++ {
		partial, err := LoadRegoCached("data.vanti.bsp.ew.TestApi.allow")
		if err != nil {
			b.Fatal(err)
		}

		// make each iteration have different data
		attr.Request.(*gen.UpdateTestRequest).Test.Data = fmt.Sprintf("please %d", i)

		result, err := partial.Rego(rego.Input(attr)).Eval(context.Background())
		if err != nil {
			b.Error(err)
		}

		if !result.Allowed() {
			b.Errorf("expected iteration %d to succeed", i)
		}
	}
}

func BenchmarkSimple(b *testing.B) {
	attr := exampleAttributes

	for i := 0; i < b.N; i++ {
		// make each iteration have different data
		attr.Request.(*gen.UpdateTestRequest).Test.Data = fmt.Sprintf("please %d", i)

		r := rego.New(
			rego.Compiler(RegoCompiler),
			rego.Input(attr),
			rego.Query("data.vanti.bsp.ew.TestApi.allow"),
		)
		result, err := r.Eval(context.Background())
		if err != nil {
			b.Error(err)
		}

		if !result.Allowed() {
			b.Errorf("expected iteration %d to succeed", i)
		}
	}
}
