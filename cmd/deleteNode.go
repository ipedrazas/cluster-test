// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// deleteNodeCmd represents the deleteNode command
var deleteNodeCmd = &cobra.Command{
	Use:   "deleteNode",
	Short: "Tries to delete a K8S node from the specified cluster",
	Long: `This command will list all the nodes of a given cluster and it will
destroy one of those instances. If there are more than one node, a random
instance will be selected from the nodes list.

For example, the following command will list all the AWS resources that would
be destroyed if it wasn't a dry-run:

~ cluster-test deleteNode --cluster barcelona --dry-run --all`,
	Run: deleteNode,
}

func init() {
	rootCmd.AddCommand(deleteNodeCmd)
}

func deleteNode(cmd *cobra.Command, args []string) {
	fmt.Printf("delete node from cluster %v", cluster)
	initAWSSession()
	nodes := getNodes(cluster)
	if all {
		for _, n := range nodes {
			deleteInstance(n)
			time.Sleep(time.Duration(wait) * time.Second)
		}
	} else {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		unluckyNode := nodes[r.Intn(len(nodes))]
		deleteInstance(unluckyNode)
	}
}
