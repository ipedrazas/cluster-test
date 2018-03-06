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
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	Run: runNode,
}

func init() {
	rootCmd.AddCommand(deleteNodeCmd)
}

func runNode(cmd *cobra.Command, args []string) {
	deleteNode()
}

func deleteNode() error {
	deletedResources = nil
	fmt.Printf("delete node from cluster %v", cluster)
	nodes := GetNodes(cluster)
	if all {
		for _, n := range nodes {
			d, err := deleteInstance(n)
			if err != nil {
				return err
			}
			if !stringInSlice(deletedResources, d) {
				deletedResources = append(deletedResources, d)
			}
			time.Sleep(time.Duration(wait) * time.Second)
		}
	} else {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		unluckyNode := nodes[r.Intn(len(nodes))]
		d, err := deleteInstance(unluckyNode)
		if err != nil {
			return err
		}
		if !stringInSlice(deletedResources, d) {
			deletedResources = append(deletedResources, d)
		}
	}
	return nil
}

// NodesHandler is the default route
func NodesHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cluster = vars["cluster"]
	if _, ok := vars["all"]; ok {
		all = true
	}
	if debug {
		fmt.Println("cmd.deleteNode.NodesHandler:")
		fmt.Printf("cluster: %v/n", cluster)
		fmt.Printf("all: %v/n", all)
		fmt.Printf("dry-run: %v/n", dryrun)
	}

	err := deleteNode()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting AWS Resources: %v", deletedResources)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "AWS Resources deleted: %v", deletedResources)
	}

}
