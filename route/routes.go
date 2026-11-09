package route

import (
	"net/http"

	"rouge/admin"
	"rouge/auth"
	"rouge/docs"
)

func Register() {
	// API documentation
	http.HandleFunc("GET /swagger", docs.SwaggerUIHandler)
	http.HandleFunc("GET /swagger/{$}", docs.SwaggerUIHandler)
	http.HandleFunc("GET /swagger/openapi.yaml", docs.OpenAPIHandler)

	// Routes anyone can access
	http.HandleFunc("POST /api/v1/login", auth.LoginHandler)
	http.HandleFunc("GET /api/v1/games", admin.GetAllGamesHandler)
	http.HandleFunc("GET /api/v1/games/{id}", admin.GetGameByIDHandler)

	// routes that require authentication
	http.HandleFunc("POST /api/v1/games", auth.Middleware(admin.CreateGameHandler))
	http.HandleFunc("PUT /api/v1/games/{id}", auth.Middleware(admin.UpdateGameHandler))
	http.HandleFunc("DELETE /api/v1/games/{id}", auth.Middleware(admin.DeleteGameHandler))
}
