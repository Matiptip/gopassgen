package password

import "testing"

func TestRandomLen(t *testing.T) {
	s, err := Random(16)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 16 {
		t.Fatalf("expected len 16, got %d", len(s))
	}
}

func TestRandomZero(t *testing.T) {
	s, err := Random(0)
	if err != nil {
		t.Fatal(err)
	}
	if s != "" {
		t.Fatalf("expected empty string for len 0, got %q", s)
	}
}
