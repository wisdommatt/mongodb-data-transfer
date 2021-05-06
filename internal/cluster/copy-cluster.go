package cluster

import (
	"context"
	"log"
	"sync"

	"github.com/wisdommatt/mongodb-data-transfer/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

// CopyDataFromTo copies data from one cluster to another.
func CopyDataFromTo(fromCluster, toCluster string, only []string, wg *sync.WaitGroup) {
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
	if len(only) > 0 {
		fromDatabases = only
	}
	for _, dbName := range fromDatabases {
		if dbName == "admin" || dbName == "local" || dbName == "config" {
			continue
		}
		db := fromDBClient.Database(dbName)
		toDB := toDBClient.Database(dbName)
		wg.Add(1)
		go database.CopyDataFromTo(db, toDB, wg)
		log.Println("Finished copying `" + dbName + "` Database")
	}
}

func stringInSlice(str string, slice []string) bool {
	res := false
	for _, value := range slice {
		if value == str {
			res = true
			break
		}
	}
	return res
}
