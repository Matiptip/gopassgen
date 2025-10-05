package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Matip/gopassgen/internal/password"
)

func main() {
	length := flag.Int("len", 16, "Longitud de la contraseña")
	count := flag.Int("n", 1, "Cantidad a generar")
	flag.Parse()

	if *length < 1 || *length > 256 {
		log.Fatal("len debe estar entre 1 y 256")
	}
	if *count < 1 || *count > 1000 {
		log.Fatal("n debe estar entre 1 y 1000")
	}

	for i := 0; i < *count; i++ {
		s, err := password.Random(*length)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
	}
}
