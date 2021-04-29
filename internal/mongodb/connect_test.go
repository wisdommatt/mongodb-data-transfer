package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestConnect(t *testing.T) {
	client, err := Connect("mongodb://localhost:27017")
	require.Nil(t, err, err)
	require.NotNil(t, client)
	err = client.Ping(context.TODO(), readpref.Primary())
	require.Nil(t, err, err)
}
