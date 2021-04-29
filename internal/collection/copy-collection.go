package collection

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CopyDataFromTo copies data from one collection to another.
func CopyDataFromTo(fromColl, toColl *mongo.Collection, wg *sync.WaitGroup) {
	var records []interface{}
	cursor, err := fromColl.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("An error occured while retrieving "+fromColl.Database().Name()+" - "+fromColl.Name()+" records", err.Error())
		return
	}
	defer cursor.Close(context.TODO())
	err = cursor.All(context.TODO(), &records)
	if err != nil {
		log.Println("An error occured while retrieving "+fromColl.Database().Name()+" - "+fromColl.Name()+" records", err.Error())
		return
	}
	_, err = toColl.InsertMany(context.TODO(), records)
	if err != nil {
		log.Println("An error occured while moving data from `"+fromColl.Database().Name()+" - "+fromColl.Name()+"` to `"+toColl.Database().Name()+"` - `"+toColl.Name()+"`", err.Error())
	}
}
