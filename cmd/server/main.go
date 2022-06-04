package main

import (
	"log"

	"github.com/mingdinghan/golog/internal/server"
)

func main() {
	srv := server.NewHttpServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
