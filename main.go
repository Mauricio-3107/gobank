package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	err = store.Init()
	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8080", store)
	server.Run()
}
