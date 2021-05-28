package db

import (
	"context"
	"fmt"
	"github.com/normegil/evevulcan/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) LoadIdentity(id int64) (*model.Identity, error) {
	identitiesCol := d.database.Collection("identities")
	var found model.Identity
	err := identitiesCol.FindOne(context.Background(), bson.M{"id": id}).Decode(&found)
	if nil != err && mongo.ErrNoDocuments != err {
		return nil, fmt.Errorf("error when retreiving identity '%d': %w", id, err)
	}
	if mongo.ErrNoDocuments == err {
		return nil, nil
	}
	return &found, nil
}

func (d *DB) InsertOrUpdateIdentity(identity model.Identity) error {
	identitiesCol := d.database.Collection("identities")

	loadedIdentity, err := d.LoadIdentity(identity.ID)
	if err != nil {
		return fmt.Errorf("inserting identity: %w", err)
	}

	if loadedIdentity == nil {
		_, err = identitiesCol.InsertOne(context.Background(), identity)
		if err != nil {
			return fmt.Errorf("inserting new identity: %w", err)
		}
	} else {
		_, err = identitiesCol.UpdateOne(context.Background(), bson.M{"character_id": identity.ID}, bson.D{
			{"$set", bson.D{
				{"refresh_token", identity.RefreshToken}},
			},
		})
		if err != nil {
			return fmt.Errorf("update existing identity: %w", err)
		}
	}
	return nil
}
