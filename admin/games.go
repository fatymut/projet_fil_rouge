package admin 

import (
	"net/http"
	"rouge/bdd"
	"encoding/json"
	"fmt"
	"strconv"
	"rouge/models"

)

func GetAllGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := bdd.GetAllGames(r.Context())
	if err != nil {
		fmt.Printf("erreur recuperation des jeux: %v", err)
		http.Error(w, "erreur recuperation des jeux", 404)
		return
	}

	response, err := json.Marshal(games)

	if err != nil {
		http.Error(w, "erreur de conversion", 500)
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