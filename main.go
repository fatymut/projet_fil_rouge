package main

import (
	"log"
	"net/http"
	"os"

	"rouge/auth"
	"rouge/bdd"
	"rouge/route"
)

func main() {

	if len(os.Args) > 1 {
		if auth.HandleCli(os.Args) {
			return
		}
	}
	bdd.NewDB()
	route.Register()

	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
