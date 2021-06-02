package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionNameAccessToken = "access_tokens"
const collectionNameCache = "api-cache"
const collectionNameIdentities = "identities"

type DB struct {
	database *mongo.Database
}

func New(db *mongo.Database) *DB {
	return &DB{
		database: db,
	}
}

func (d *DB) CreateIndexes() error {
	if err := createExpirationIndex(d.database.Collection(collectionNameAccessToken), "accessTokenExpiration", "created", 900); nil != err {
		return fmt.Errorf("create expiration index for '%s': %w", collectionNameAccessToken, err)
	}
	if err := createExpirationIndex(d.database.Collection(collectionNameCache), "cacheExpiration", "expiration", 0); nil != err {
		return fmt.Errorf("create expiration index for '%s': %w", collectionNameAccessToken, err)
	}
	return nil
}

func createExpirationIndex(collection *mongo.Collection, name string, key string, expirationSeconds int) error {
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{key, 1}},
		Options: options.Index().SetName(name).SetExpireAfterSeconds(int32(expirationSeconds)),
	})
	if err != nil {
		return fmt.Errorf("create index '%s': %w", name, err)
	}
	return nil
}
