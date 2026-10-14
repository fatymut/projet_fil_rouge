package admin 

import (
	"net/http"
	"rouge/bdd"
	"encoding/json"
	"fmt"
	"strconv"
	"rouge/models"
	"go.mongodb.org/mongo-driver/bson"

)

func GetAllGamesHandler(w http.ResponseWriter, r *http.Request) {

	var page = 1
	var err error
	var sortField = "_id"
	var sortOrder = 1
	var filter = bson.M{}
	
	if r.URL.Query().Get("release_year") != "" {
		releaseYear, err := strconv.Atoi(r.URL.Query().Get("release_year"))
		if err != nil {
			http.Error(w, "release_year doit être un nombre", http.StatusBadRequest)
			return
		}
		filter["release_year"] = releaseYear
	}

	if r.URL.Query().Get("title") != "" {
		title := r.URL.Query().Get("title")
		filter ["title"] = bson.M{"$regex": title, "$options": "i"}
		if err != nil {
			fmt.Println("erreur recuperation des jeux:", err)
			http.Error(w, "erreur recuperation des jeux", http.StatusInternalServerError)
			return
		}
	}

	if r.URL.Query().Get("ordering") != "" {
		ordering := r.URL.Query().Get("ordering")
		if ordering != "title" && ordering != "-title" && ordering != "release_year" && ordering != "-release_year" {
			http.Error(w, "ordering doit être 'title', '-title', 'release_year' ou '-release_year'", http.StatusBadRequest)
			return
		}

		if ordering[0] == '-' {
			sortOrder = -1
			sortField = ordering[1:]
		} else {
			sortOrder = 1
			sortField = ordering
		}
	}

	if r.URL.Query().Get("page") != "" {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			http.Error(w, "page doit être un nombre", http.StatusBadRequest)
			return
		}

		if page < 1 {
			http.Error(w, "page doit être supérieur à 0", http.StatusBadRequest)
			return
		}
	}

	var limit = 5
	if r.URL.Query().Get("limit") != "" {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "limit doit être un nombre", http.StatusBadRequest)
			return
		}

		if limit < 1 || limit > 100 {
			http.Error(w, "limit doit être compris entre 1 et 100", http.StatusBadRequest)
			return
		}
	}

	var skip = int64((page - 1) * limit)

	games, total, err := bdd.GetAllGames(r.Context(), int64(limit), skip, sortField, sortOrder, filter)
	if err != nil {
		fmt.Println("erreur recuperation des jeux:", err)
		http.Error(w, "erreur recuperation des jeux", http.StatusInternalServerError)
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)
	fmt.Println("total:", total, "pages:", totalPages)

	data := map[string]interface{}{
	"page":        page,
	"limit":       limit,
	"total":       total,
	"total_pages": totalPages,
	"results":     games,
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "erreur de conversion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", response)
}

func GetGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("Id du jeu invalide: %v", err)
		http.Error(w, "id du jeu invalide", http.StatusBadRequest)
		return


	}

	game, err := bdd.GetGameByID(id)
	if err != nil {
		fmt.Printf("Erreur lors de la récupération du jeu: %v", err)
		http.Error(w, "Jeu non trouvé", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(game)
	if err != nil {
		fmt.Printf("Erreur lors de la conversion du jeu en JSON: %v", err)
		http.Error(w, "Erreur de conversion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		fmt.Printf("Erreur lors du décodage du jeu: %v", err)
		http.Error(w, "erreur decodage du jeu", http.StatusBadRequest)
		return
	}

	err = bdd.CreateGame(game)
	if err != nil {
		fmt.Printf("erreur creation du jeu %v", err)
		http.Error(w, "Erreur lors de la création du jeu", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateGameHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("Id du jeu invalide: %v", err)
		http.Error(w, "Id du jeu invalide", http.StatusBadRequest)
		return
	}

	var updatedGame models.Game
	err = json.NewDecoder(r.Body).Decode(&updatedGame)
	if err != nil {
		http.Error(w, "erreur decodage du jeu", http.StatusBadRequest)
		return
	}

	err = bdd.UpdateGame(id, updatedGame)
	if err != nil {
		fmt.Printf("erreur mise a jour du jeu %v", err)
		http.Error(w, "Erreur lors de la mise à jour du jeu", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("id du jeu invalide: %v", err)
		http.Error(w, "id du jeu invalide", http.StatusBadRequest)
		return
	}

	err = bdd.DeleteGame(id)
	if err != nil {
		fmt.Printf("erreur suppression du jeu %v", err)
		http.Error(w, "Erreur lors de la suppression du jeu", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}