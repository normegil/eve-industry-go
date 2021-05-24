package config

import "os"

func MongoDBURL() string {
	return os.Getenv(appPrefix + "MONGO_URL")
}
