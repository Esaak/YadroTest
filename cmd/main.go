package main

import (
	"github.com/Esaak/YadroTest/internal/root"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <input_file>")
		return
	}

	root.Execute(os.Args[1])
}
