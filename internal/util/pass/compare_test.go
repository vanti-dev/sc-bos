package pass

import "testing"

func TestCompare(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{name: "empty", secret: ""},
		{name: "pass1", secret: "hello"},
		{name: "pass2", secret: "0928u3j@£F@*DF@N 23d9n2f3oAS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range 5 {
				hash, err := Hash([]byte(tt.secret))
				if err != nil {
					t.Errorf("hash failed %v", err)
				}
				if err := Compare(hash, []byte(tt.secret)); (err != nil) != tt.wantErr {
					t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
