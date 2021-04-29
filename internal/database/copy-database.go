package database

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TransferDataFromTo transfers/copies data from one database to another.
func TransferDataFromTo(fromDB, toDB *mongo.Database, wg *sync.WaitGroup) (err error) {
	collections, err := fromDB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return
	}
	wg.Add(len(collections))
	for _, collectionName := range collections {
		go func(collectionName string) {
			defer wg.Done()

			collection := fromDB.Collection(collectionName)
			toCollection := toDB.Collection(collectionName)

			var records []interface{}
			cursor, err := collection.Find(context.TODO(), bson.M{})
			if err != nil {
				log.Println("An error occured while retrieving "+fromDB.Name()+" - "+collectionName+" records", err.Error())
				return
			}
			defer cursor.Close(context.TODO())
			err = cursor.All(context.TODO(), &records)
			if err != nil {
				log.Println("An error occured while retrieving "+fromDB.Name()+" - "+collectionName+" records", err.Error())
				return
			}
			_, err = toCollection.InsertMany(context.TODO(), records)
			if err != nil {
				log.Println("An error occured while moving data from `"+fromDB.Name()+" - "+collectionName+"` to `"+toDB.Name()+"` - `"+toCollection.Name()+"`", err.Error())
			}
		}(collectionName)
	}
	return
}
