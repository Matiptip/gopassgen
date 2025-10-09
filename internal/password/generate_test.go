package password

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ✅ Test básico de longitud
func TestRandomLength(t *testing.T) {
	p, err := Random(16, false, false)
	if err != nil {
		t.Fatalf("error generando contraseña: %v", err)
	}
	if len(p) != 16 {
		t.Errorf("esperaba longitud 16, obtuvo %d", len(p))
	}
}

// ✅ Test sin símbolos y sin ambiguos
func TestRandomNoSymbolsNoAmbiguous(t *testing.T) {
	p, err := Random(32, false, true)
	if err != nil {
		t.Fatalf("error generando contraseña: %v", err)
	}
	for _, r := range p {
		if strings.ContainsRune("!@#$%^&*()-_=+[]{}<>?/|~O0Il", r) {
			t.Errorf("carácter prohibido encontrado: %c", r)
		}
	}
}

// ✅ Test de creación de archivo de salida (integración simple)
func TestExportFileCreation(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "passwords_test.txt")

	file, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("error creando archivo: %v", err)
	}
	defer file.Close()

	p, err := Random(8, false, false)
	if err != nil {
		t.Fatalf("error generando contraseña: %v", err)
	}

	if _, err := file.WriteString(p); err != nil {
		t.Fatalf("error escribiendo archivo: %v", err)
	}

	stat, err := os.Stat(outFile)
	if err != nil {
		t.Fatalf("error accediendo al archivo: %v", err)
	}

	if stat.Size() == 0 {
		t.Error("el archivo fue creado pero está vacío")
	}
}

// ✅ Benchmark: performance general
func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Random(32, true, true)
	}
}
