package media

import "testing"

func TestEndpoint(t *testing.T) {
	for _, test := range []struct {
		raw, want string
		secure    bool
	}{
		{"localhost:9000", "localhost:9000", false},
		{"https://minio.example", "minio.example", true},
	} {
		got, secure, err := endpoint(test.raw)
		if err != nil || got != test.want || secure != test.secure {
			t.Fatalf("endpoint(%q) = %q, %v, %v", test.raw, got, secure, err)
		}
	}
}
