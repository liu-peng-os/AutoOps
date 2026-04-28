package db

import "testing"

func TestShouldRetryDBConnect(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "unexpected eof", err: errString("failed to receive message: unexpected EOF"), want: true},
		{name: "connection refused", err: errString("dial tcp 127.0.0.1:5432: connectex: connection refused"), want: true},
		{name: "invalid password", err: errString("password authentication failed for user"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRetryDBConnect(tt.err); got != tt.want {
				t.Fatalf("shouldRetryDBConnect(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

type errString string

func (e errString) Error() string {
	return string(e)
}
