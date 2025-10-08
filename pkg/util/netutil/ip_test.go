package netutil

import (
	"testing"
)

func TestMergeHostPort(t *testing.T) {
	type testCase struct {
		in        []string
		expect    string
		expectErr bool
	}

	tests := map[string]testCase{
		"empty": {
			in:     []string{},
			expect: "",
		},
		"empty_strings": {
			in:     []string{"", "", ""},
			expect: "",
		},
		"only_hosts": {
			in:     []string{"host1", "host2", "host3"},
			expect: "host3",
		},
		"only_ports": {
			in:     []string{":8080", ":9090", ":7070"},
			expect: ":7070",
		},
		"hosts_and_ports": {
			in:     []string{"host1:8080", "host2:9090", "host3:7070"},
			expect: "host3:7070",
		},
		"mixed": {
			in:     []string{"host1", ":8080", "host2:9090", "host3", ":7070"},
			expect: "host3:7070",
		},
		"ipv4_and_ipv6": {
			in:     []string{"192.168.0.1", "[2001:db8::1]"},
			expect: "[2001:db8::1]",
		},
		// this doesn't work because IPv6 address are ambiguous without brackets in a context where a port might be present
		"unbracketed_ipv6": {
			in:        []string{"2001:db8::1"},
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := MergeHostPort(tc.in...)
			if result != tc.expect {
				t.Errorf("MergeHostPort(%v) = %q; want %q", tc.in, result, tc.expect)
			}
			if (err != nil) != tc.expectErr {
				t.Errorf("MergeHostPort(%v) error = %v; want error: %v", tc.in, err, tc.expectErr)
			}
		})
	}
}
