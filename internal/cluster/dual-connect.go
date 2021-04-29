package cluster

import (
	"github.com/wisdommatt/mongodb-data-transfer/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// DualConnect connects to two database clients.
//
// If an error occured while connecting to the first client it doesn't
// proceed to the next client.
func DualConnect(uri1, uri2 string) (client1, client2 *mongo.Client, err error) {
	client1, err = database.Connect(uri1)
	if err != nil {
		return
	}
	client2, err = database.Connect(uri2)
	return
}
