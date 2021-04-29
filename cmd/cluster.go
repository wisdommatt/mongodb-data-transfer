/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/wisdommatt/mongodb-data-transfer/internal/database"

	"github.com/spf13/cobra"
	"github.com/wisdommatt/mongodb-data-transfer/internal/cluster"
	"go.mongodb.org/mongo-driver/bson"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Transfers data from one cluster to another.",
	Run: func(cmd *cobra.Command, args []string) {
		wg := new(sync.WaitGroup)
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		if from == "" || to == "" {
			log.Fatalln("`from` and `to` must be provided")
			return
		}
		fromDBClient, toDBClient, err := cluster.DualConnect(from, to)
		if err != nil {
			log.Fatalln("Error while connecting to clusters ", err.Error())
			return
		}
		fromDatabases, err := fromDBClient.ListDatabaseNames(context.TODO(), bson.M{})
		if err != nil {
			log.Fatalln("An error occured while returning `from` databases", err.Error())
			return
		}
		for _, dbName := range fromDatabases {
			db := fromDBClient.Database(dbName)
			toDB := toDBClient.Database(dbName)
			wg.Add(1)
			go func() {
				err = database.CopyDataFromTo(db, toDB, wg)
				if err != nil {
					log.Println("An error occured while transferring data from database: " + db.Name() + " to " + toDB.Name())
				}
			}()
			log.Println("Finished transferring `" + dbName + "` Database")
			break
		}
		wg.Wait()
		fmt.Println("Execution Completed !")
		// mongo.Connect
		// fmt.Println("cluster called", fromDBClient, toDBClient)
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().String("from", "", "Cluster to transfer the data from")
	clusterCmd.Flags().String("to", "", "Cluster to transfer the data to")
}
