package password

import "testing"

func TestSymbols(t *testing.T) {
	s, err := Random(32, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 32 {
		t.Fatalf("esperado 32, obtenido %d", len(s))
	}
}

func TestNoAmbiguous(t *testing.T) {
	s, err := Random(64, false, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range "O0Il" {
		if containsRune(s, c) {
			t.Fatalf("se encontró carácter ambiguo %q en %q", c, s)
		}
	}
}

func containsRune(s string, r rune) bool {
	for _, x := range s {
		if x == r {
			return true
		}
	}
	return false
}
