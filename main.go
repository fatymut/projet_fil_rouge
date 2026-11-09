package main

import (
	"log"
	"net/http"
	"os"

	"rouge/admin"
	"rouge/auth"
	"rouge/bdd"
)

func main() {

	if len(os.Args) > 1 {
		if auth.HandleCli(os.Args) {
			return
		}
	}
	bdd.NewDB()
	// Routes anyone can access
	http.HandleFunc("POST /api/v1/login", auth.LoginHandler)
	http.HandleFunc("GET /api/v1/games", admin.GetAllGamesHandler)
	http.HandleFunc("GET /api/v1/games/{id}", admin.GetGameByIDHandler)

	// routes that require authentication
	http.HandleFunc("POST /api/v1/games", auth.Middleware(admin.CreateGameHandler))
	http.HandleFunc("PUT /api/v1/games/{id}", auth.Middleware(admin.UpdateGameHandler))
	http.HandleFunc("DELETE /api/v1/games/{id}", auth.Middleware(admin.DeleteGameHandler))

	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
