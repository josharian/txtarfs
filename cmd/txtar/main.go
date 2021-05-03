package main

import (
	"log"
	"os"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

func main() {
	src := os.DirFS(".")
	ar, err := txtarfs.From(src)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(txtar.Format(ar))
}
