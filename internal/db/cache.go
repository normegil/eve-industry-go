package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const cacheCollectionName = "api-cache"

type APICacheObject struct {
	QueryID    string    `bson:"query_id"`
	Expiration time.Time `bson:"expiration"`
	Object     []byte    `bson:"object"`
}

func (d *DB) FromCache(queryId string) ([]byte, error) {
	cacheCol := d.database.Collection(cacheCollectionName)
	var obj APICacheObject
	err := cacheCol.FindOne(context.Background(), bson.M{"query_id": queryId}).Decode(&obj)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			return nil, nil
		}
		return nil, fmt.Errorf("query cache for '%s': %w", queryId, err)
	}
	return obj.Object, err
}

func (d *DB) ToCache(obj APICacheObject) error {
	cacheCol := d.database.Collection(cacheCollectionName)
	if _, err := cacheCol.InsertOne(context.Background(), obj); nil != err {
		return fmt.Errorf("inserting into cache '%s': %w", obj.QueryID, err)
	}
	return nil
}
