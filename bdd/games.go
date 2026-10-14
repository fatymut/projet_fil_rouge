package bdd

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rouge/models"
	"go.mongodb.org/mongo-driver/mongo/options"

)

func GetAllGames(ctx context.Context, limit int64, skip int64, sortField string, sortOrder int, filter bson.M) ([]models.Game, int64, error) {
	totalCount, err := Database.Collection("games").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur de comptage des documents: %v", err)
	}
	

	sort := bson.D{{Key: sortField, Value: sortOrder}}

	if sortField != "_id" {
	sort = append(sort, bson.E{Key: "_id", Value: 1})
	}

	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(sort)// ordre stable

	games := []models.Game{}

	cursor, err := Database.Collection("games").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur recup jeux: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var game models.Game
		if err := cursor.Decode(&game); err != nil {
			return nil, 0, fmt.Errorf("erreur de décodage du jeu: %v", err)
		}
		games = append(games, game)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %v", err)
	}

	return games, totalCount, nil
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