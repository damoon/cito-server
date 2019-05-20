package main

import (
	"flag"
	"log"

	"github.com/damoon/cito-server"
)

func main() {

	addr := flag.String("address", ":8080", "default server address, ':8080'")

	flag.Parse()

	log.Printf("server will listen on: %s\n", *addr)

	cito.RunServer(*addr)
}
