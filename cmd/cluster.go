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
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/wisdommatt/mongodb-data-transfer/internal/cluster"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Transfers data from one cluster to another.",
	Run: func(cmd *cobra.Command, args []string) {
		var wg *sync.WaitGroup = new(sync.WaitGroup)
		log.Println("running")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		if from == "" || to == "" {
			log.Fatalln("`from` and `to` must be provided")
			return
		}
		onlyDBs := []string{}
		only, _ := cmd.Flags().GetString("only")
		if len(only) > 0 {
			onlyDBs = strings.Split(only, ",")
		}
		cluster.CopyDataFromTo(from, to, onlyDBs, wg)
		wg.Wait()
		fmt.Println("Execution Completed !")
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().String("from", "", "Cluster to transfer the data from")
	clusterCmd.Flags().String("to", "", "Cluster to transfer the data to")
	clusterCmd.Flags().StringP("only", "o", "", "Used to strictly specify the only databases to copy separated by comma (,)")
}
