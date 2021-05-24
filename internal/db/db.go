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

func (d *DB) InsertOrUpdateIdentity(identity model.Identity) error {
	identitiesCol := d.database.Collection("identities")
	var found model.Identity
	err := identitiesCol.FindOne(context.Background(), bson.M{"character_id": identity.ID}).Decode(&found)
	if nil != err && mongo.ErrNoDocuments != err {
		return fmt.Errorf("error when inserting identity: %w", err)
	}

	if mongo.ErrNoDocuments == err {
		_, err = identitiesCol.InsertOne(context.Background(), identity)
		if err != nil {
			return fmt.Errorf("inserting new identity: %w", err)
		}
	}

	_, err = identitiesCol.UpdateOne(context.Background(), bson.M{"character_id": identity.ID}, bson.D{
		{"$set", bson.D{
			{"refresh_token", identity.RefreshToken}},
		},
	})
	if err != nil {
		return fmt.Errorf("update existing identity: %w", err)
	}
	return nil
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
