package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Matiptip/gopassgen/internal/password"
	"github.com/Matiptip/gopassgen/pkg/version"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: gopassgen [generate|export|version]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		handleGenerate()
	case "export":
		handleExport()
	case "version":
		fmt.Println("gopassgen", version.Version)
	default:
		fmt.Println("Comando desconocido:", os.Args[1])
	}
}

// ✅ Genera contraseñas simples (modo tradicional)
func handleGenerate() {
	length := flag.Int("len", 16, "Longitud de la contraseña")
	count := flag.Int("n", 1, "Cantidad de contraseñas")
	useSymbols := flag.Bool("symbols", false, "Incluir símbolos especiales")
	noAmbiguous := flag.Bool("no-ambiguous", false, "Excluir caracteres ambiguos")
	flag.CommandLine.Parse(os.Args[2:])

	for i := 0; i < *count; i++ {
		p, err := password.Random(*length, *useSymbols, *noAmbiguous)
		if err != nil {
			log.Fatalf("error generando contraseña: %v", err)
		}
		fmt.Println(p)
	}
}

// ✅ Exporta contraseñas a archivo (modo paralelo opcional)
func handleExport() {
	length := flag.Int("len", 16, "Longitud de la contraseña")
	count := flag.Int("n", 5, "Cantidad de contraseñas")
	useSymbols := flag.Bool("symbols", false, "Incluir símbolos especiales")
	noAmbiguous := flag.Bool("no-ambiguous", false, "Excluir caracteres ambiguos")
	format := flag.String("format", "text", "Formato de salida: text|json")
	output := flag.String("o", "", "Archivo de salida (opcional)")
	parallel := flag.Bool("parallel", false, "Generar contraseñas en modo concurrente")
	flag.CommandLine.Parse(os.Args[2:])

	passwords := make([]string, *count)

	if *parallel {
		generateParallel(passwords, *length, *useSymbols, *noAmbiguous)
	} else {
		for i := 0; i < *count; i++ {
			p, err := password.Random(*length, *useSymbols, *noAmbiguous)
			if err != nil {
				log.Fatalf("error generando contraseña: %v", err)
			}
			passwords[i] = p
		}
	}

	if *output == "" {
		timestamp := time.Now().Format("20060102_150405")
		*output = fmt.Sprintf("passwords_%s.%s", timestamp, *format)
	}

	savePasswords(*output, *format, passwords)
	fmt.Printf("✅ %d contraseñas exportadas en %s (%s)\n", *count, *output, *format)
}

// ✅ Generación concurrente
func generateParallel(passwords []string, length int, useSymbols bool, noAmbiguous bool) {
	var wg sync.WaitGroup
	wg.Add(len(passwords))
	results := make(chan struct {
		index int
		value string
		err   error
	}, len(passwords))

	for i := range passwords {
		go func(i int) {
			defer wg.Done()
			p, err := password.Random(length, useSymbols, noAmbiguous)
			results <- struct {
				index int
				value string
				err   error
			}{i, p, err}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.err != nil {
			log.Fatalf("error generando contraseña: %v", r.err)
		}
		passwords[r.index] = r.value
	}
}

// ✅ Guardado de contraseñas (text o JSON)
func savePasswords(output string, format string, passwords []string) {
	file, err := os.Create(output)
	if err != nil {
		log.Fatalf("error creando archivo: %v", err)
	}
	defer file.Close()

	switch format {
	case "json":
		data, err := json.MarshalIndent(passwords, "", "  ")
		if err != nil {
			log.Fatalf("error codificando JSON: %v", err)
		}
		file.Write(data)
	default:
		for _, p := range passwords {
			fmt.Fprintln(file, p)
		}
	}
}
