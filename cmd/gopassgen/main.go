package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Matiptip/gopassgen/internal/password"
	"github.com/Matiptip/gopassgen/pkg/version"
)

func main() {
	length := flag.Int("len", 16, "Longitud de la contraseña")
	count := flag.Int("n", 1, "Cantidad de contraseñas a generar")
	useSymbols := flag.Bool("symbols", false, "Incluir símbolos especiales")
	noAmbiguous := flag.Bool("no-ambiguous", false, "Excluir caracteres ambiguos (O, 0, I, l)")
	showVersion := flag.Bool("v", false, "Mostrar versión y salir")
	outputFile := flag.String("o", "", "Archivo de salida (opcional). Ej: -o passwords.txt")
	flag.Parse()

	if *showVersion {
		fmt.Println("gopassgen", version.Version)
		return
	}

	if *length < 1 || *length > 256 {
		log.Fatal("len debe estar entre 1 y 256")
	}
	if *count < 1 || *count > 1000 {
		log.Fatal("n debe estar entre 1 y 1000")
	}

	// Generar contraseñas usando password.Random (...)
	passwords := make([]string, *count)
	for i := 0; i < *count; i++ {
		p, err := password.Random(*length, *useSymbols, *noAmbiguous)
		if err != nil {
			log.Fatalf("error generando contraseña: %v", err)
		}
		passwords[i] = p
	}

	// Si no se pasa -o, crear nombre por timestamp
	if *outputFile == "" {
		timestamp := time.Now().Format("20060102_150405")
		*outputFile = fmt.Sprintf("passwords_%s.txt", timestamp)
	}

	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("error creando archivo: %v", err)
	}
	defer f.Close()

	for _, p := range passwords {
		if _, err := fmt.Fprintln(f, p); err != nil {
			log.Fatalf("error escribiendo archivo: %v", err)
		}
	}

	fmt.Printf("✅ %d contraseñas guardadas en %s\n", *count, *outputFile)
}
