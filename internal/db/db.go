package db

import (
	"context"
	"fmt"
	"github.com/normegil/evevulcan/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	database *mongo.Database
}

func New(db *mongo.Database) *DB {
	return &DB{
		database: db,
	}
}

func (d *DB) ReplaceAccessToken(token model.StoredAccessToken) error {
	tokenCol := d.database.Collection("access_tokens")
	if _, err := tokenCol.DeleteOne(context.Background(), bson.M{"character_id": token.CharacterID}); nil != err {
		return fmt.Errorf("delete access token for %d: %w", token.CharacterID, err)
	}
	if _, err := tokenCol.InsertOne(context.Background(), token); nil != err {
		return fmt.Errorf("insert access token for %d: %w", token.CharacterID, err)
	}
	return nil
}
