package db

import (
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
