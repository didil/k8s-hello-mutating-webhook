package main

import (
	"log"

	"github.com/didil/k8s-hello-mutating-webhook/webhook/api"
)

func main() {
	err := api.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
