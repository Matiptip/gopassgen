package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/Matiptip/gopassgen/internal/password"
	"github.com/Matiptip/gopassgen/pkg/version"
)

func main() {
	go func() {
		log.Println("pprof disponible en http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("Error iniciando servidor pprof: %v", err)
		}
	}()

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
	fs := flag.NewFlagSet("export", flag.ExitOnError)

	length := fs.Int("len", 16, "Longitud de la contraseña")
	count := fs.Int("n", 5, "Cantidad de contraseñas")
	useSymbols := fs.Bool("symbols", false, "Incluir símbolos especiales")
	noAmbiguous := fs.Bool("no-ambiguous", false, "Excluir caracteres ambiguos")
	format := fs.String("format", "text", "Formato de salida: text|json")
	output := fs.String("o", "", "Archivo de salida (opcional)")
	parallel := fs.Int("parallel", 1, "Número de workers concurrentes (1 = secuencial)")
	cpuprofile := fs.String("cpuprofile", "", "Archivo de perfil de CPU (ej: cpu.pprof)")
	memprofile := fs.String("memprofile", "", "Archivo de perfil de memoria (ej: mem.pprof)")

	// Parsear SOLO los flags del subcomando export
	fs.Parse(os.Args[2:])

	// Si el usuario pasó nombre relativo, convertir a absoluto para evitar dudas de carpeta
	if *cpuprofile != "" && !filepath.IsAbs(*cpuprofile) {
		abs, _ := filepath.Abs(*cpuprofile)
		*cpuprofile = abs
	}
	if *memprofile != "" && !filepath.IsAbs(*memprofile) {
		abs, _ := filepath.Abs(*memprofile)
		*memprofile = abs
	}

	// Log de control (verás esto en la terminal)
	log.Printf("cpuprofile=%q memprofile=%q\n", *cpuprofile, *memprofile)

	// Iniciar CPU profile (si corresponde)
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("No se pudo crear cpuprofile: %v", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("Error iniciando cpuprofile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	passwords := make([]string, *count)
	if *parallel > 1 {
		generateParallel(passwords, *length, *useSymbols, *noAmbiguous, *parallel)
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

	// Heap profile al final (si se pidió)
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatalf("No se pudo crear memprofile: %v", err)
		}
		defer f.Close()
		runtime.GC() // mejora representatividad del heap profile
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("Error escribiendo heap profile: %v", err)
		}
	}

	fmt.Printf("✅ %d contraseñas exportadas en %s (%s)\n", *count, *output, *format)
}
func generateParallel(passwords []string, length int, useSymbols bool, noAmbiguous bool, workers int) {
	var wg sync.WaitGroup
	tasks := make(chan int, len(passwords))
	results := make(chan struct {
		index int
		value string
		err   error
	}, len(passwords))

	// Cargar tareas
	for i := range passwords {
		tasks <- i
	}
	close(tasks)

	// Lanzar workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range tasks {
				p, err := password.Random(length, useSymbols, noAmbiguous)
				results <- struct {
					index int
					value string
					err   error
				}{i, p, err}
			}
		}()
	}

	// Esperar finalización
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
