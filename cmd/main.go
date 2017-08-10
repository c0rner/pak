package main

import "fmt"

import (
	"github.com/c0rner/pak"
	"log"
)

func main() {
	fmt.Println("PAK")

	r, err := pak.OpenReader("data/res.pak")
	if err != nil {
		log.Fatal("Unable to open file:", err)
	}
	r.Close()
}
