package database

import (
	"context"
	"log"
	"sync"

	"github.com/wisdommatt/mongodb-data-transfer/internal/collection"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CopyDataFromTo copies data from one database to another.
func CopyDataFromTo(fromDB, toDB *mongo.Database, wg *sync.WaitGroup) {
	defer wg.Done()
	collections, err := fromDB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Println("An error occured while transferring data from database: " + fromDB.Name() + " to " + toDB.Name())
		return
	}
	wg.Add(len(collections))
	for _, collectionName := range collections {
		fromColl := fromDB.Collection(collectionName)
		toColl := toDB.Collection(collectionName)
		go collection.CopyDataFromTo(fromColl, toColl, wg)
	}
}
