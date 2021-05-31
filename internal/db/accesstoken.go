package db

import (
	"context"
	"fmt"
	"github.com/normegil/evevulcan/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const accessTokenCollection = "access_tokens"

func (d *DB) AccessToken(characterID int64) (*model.StoredAccessToken, error) {
	tokenCol := d.database.Collection(accessTokenCollection)
	var token model.StoredAccessToken
	if err := tokenCol.FindOne(context.Background(), bson.M{"character_id": characterID}).Decode(&token); nil != err {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("retrieving token for character id '%d': %w", characterID, err)
	}
	return &token, nil
}

func (d *DB) ReplaceAccessToken(token model.StoredAccessToken) error {
	tokenCol := d.database.Collection(accessTokenCollection)
	if _, err := tokenCol.DeleteOne(context.Background(), bson.M{"character_id": token.CharacterID}); nil != err {
		return fmt.Errorf("delete access token for %d: %w", token.CharacterID, err)
	}
	if _, err := tokenCol.InsertOne(context.Background(), token); nil != err {
		return fmt.Errorf("insert access token for %d: %w", token.CharacterID, err)
	}
	return nil
}
