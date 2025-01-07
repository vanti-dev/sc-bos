package policy

import (
	"context"
	"embed"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/open-policy-agent/opa/rego"
)

func TestValidate(t *testing.T) {
	allow := rego.ResultSet{{
		Expressions: []*rego.ExpressionValue{
			{
				Value: true,
				Text:  "allow",
			},
		},
	}}

	deny := rego.ResultSet{{
		Expressions: []*rego.ExpressionValue{
			{
				Value: false,
				Text:  "allow",
			},
		},
	}}

	empty := rego.ResultSet{}

	type testCase struct {
		attr          Attributes
		responses     map[string]rego.ResultSet
		expectErr     error
		expectQueries []string
	}

	cases := map[string]testCase{
		"Hierarchy": {
			attr: Attributes{
				Protocol: ProtocolGRPC,
				Service:  "foo.bar.baz",
			},
			expectErr: ErrPermissionDenied,
			expectQueries: []string{
				"data.foo.bar.baz.allow",
				"data.foo.bar.allow",
				"data.foo.allow",
				"data.grpc_default.allow",
			},
		},
		"ShortCircuit_Positive": {
			attr: Attributes{
				Protocol: ProtocolGRPC,
				Service:  "foo.bar.baz",
			},
			responses: map[string]rego.ResultSet{
				"data.foo.bar.baz.allow": empty,
				"data.foo.bar.allow":     allow,
				"data.foo.allow":         deny,
			},
			expectErr: nil,
			expectQueries: []string{
				"data.foo.bar.baz.allow",
				"data.foo.bar.allow",
			},
		},
		"ShortCircuit_Negative": {
			attr: Attributes{
				Protocol: ProtocolGRPC,
				Service:  "foo.bar.baz",
			},
			responses: map[string]rego.ResultSet{
				"data.foo.bar.baz.allow": empty,
				"data.foo.bar.allow":     deny,
				"data.foo.allow":         allow,
			},
			expectErr: ErrPermissionDenied,
			expectQueries: []string{
				"data.foo.bar.baz.allow",
				"data.foo.bar.allow",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			policy := &mockPolicy{responses: c.responses}
			_, err := Validate(context.Background(), policy, c.attr)
			if !errors.Is(err, c.expectErr) {
				t.Errorf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(c.expectQueries, policy.queries); diff != "" {
				t.Errorf("wrong query sequence (-want +got):\n%s", diff)
			}
		})
	}
}

type mockPolicy struct {
	responses map[string]rego.ResultSet
	queries   []string
}

func (p *mockPolicy) EvalPolicy(ctx context.Context, query string, attributes Attributes) (rego.ResultSet, error) {
	p.queries = append(p.queries, query)
	return p.responses[query], nil
}

//go:embed testdata
var testdata embed.FS

func TestValidate_Integration(t *testing.T) {
	policy, err := FromFS(testdata)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		attr      Attributes
		expectErr error
	}
	cases := map[string]testCase{
		"foo.bar": {
			attr: Attributes{
				Protocol: ProtocolGRPC,
				Service:  "foo.bar",
			},
			expectErr: ErrPermissionDenied,
		},
		"foo.baz": {
			attr: Attributes{
				Protocol: ProtocolGRPC,
				Service:  "foo.baz",
			},
			expectErr: nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Validate(context.Background(), policy, c.attr)
			if !errors.Is(err, c.expectErr) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
