package main

import "fmt"

import (
	"github.com/c0rner/pak"
	"log"
)

func main() {
	fmt.Println("PAK")

	p, err := pak.OpenReader("data/res.pak")
	if err != nil {
		log.Fatal("Unable to open file:", err)
	}

	for _, f := range p.File {
		fmt.Printf("File: %s Size: %d Cksum: 0x%X\n", f.Path(), f.Size(), f.Cksum())
	}
	p.Close()
}
