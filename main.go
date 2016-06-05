package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	domain := os.Args[1]
	if len(os.Args) != 2 {
		log.Fatal("Usage:\n\n\tgo run main.go <url>\n\n")
	}
	resp, err := http.Get(domain)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Request complete: ", resp.StatusCode)
}
