package cluster

import (
	"context"
	"log"
	"sync"

	"github.com/wisdommatt/mongodb-data-transfer/internal/database"

	"go.mongodb.org/mongo-driver/bson"
)

// CopyDataFromTo copies data from one cluster to another.
func CopyDataFromTo(fromCluster, toCluster string, wg *sync.WaitGroup) {
	fromDBClient, toDBClient, err := DualConnect(fromCluster, toCluster)
	if err != nil {
		log.Fatalln("Error while connecting to clusters ", err.Error())
		return
	}
	fromDatabases, err := fromDBClient.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalln("An error occured while returning `from` databases", err.Error())
		return
	}
	wg.Add(len(fromDatabases))
	for _, dbName := range fromDatabases {
		db := fromDBClient.Database(dbName)
		toDB := toDBClient.Database(dbName)
		go database.CopyDataFromTo(db, toDB, wg)
		log.Println("Finished copying `" + dbName + "` Database")
	}
}
