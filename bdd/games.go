package bdd

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rouge/models"
	"go.mongodb.org/mongo-driver/mongo/options"

)

func GetAllGames(ctx context.Context, limit int, skip int) ([]models.Game, error) {

	limit = r.URL.Query().Get("limit") // recup le sparm de la requete
	skip = r.URL.Query().Get("skip")

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))

	games := []models.Game{}

	cursor, err := Database.Collection("games").Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("erreur recup jeu %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var game models.Game
		if err := cursor.Decode(&game); err != nil {
			return nil, fmt.Errorf("erreur de décodage du jeu: %v", err)
		}
		games = append(games, game)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return games, nil
}

func GetGameByID(id int) (*models.Game, error) {
	var game models.Game
	err := Database.Collection("games").FindOne(context.Background(), bson.M{"_id": id}).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func CreateGame(game models.Game) error {
	_, err := Database.Collection("games").InsertOne(context.Background(), game)
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion du jeu: %w", err)
	}
	return nil
}

func UpdateGame(id int, updatedGame models.Game) error {
	updatedGame.ID = id

	res, err := Database.Collection("games").ReplaceOne(context.Background(), bson.M{"_id": id}, updatedGame)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du jeu: %w", err)
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func DeleteGame(id int) error {
	res, err := Database.Collection("games").DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du jeu: %w", err)
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}