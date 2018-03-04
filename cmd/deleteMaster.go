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

// deleteMasterCmd represents the deleteMaster command
var deleteMasterCmd = &cobra.Command{
	Use:   "deleteMaster",
	Short: "Tries to delete a K8S master from the specified cluster",
	Long: `This command will list all the masters of a given cluster and it will
destroy one of those instances. If there are more than one master, a random
instance will be selected from the masters list.

For example, the following command will list all the AWS resources that would
be destroyed if it wasn't a dry-run:

~ cluster-test deleteMaster --cluster barcelona --dry-run --all`,
	Run: deleteMaster,
}

func init() {
	rootCmd.AddCommand(deleteMasterCmd)

}

func deleteMaster(cmd *cobra.Command, args []string) {
	fmt.Printf("delete master from cluster %v", cluster)
	initAWSSession()
	masters := getMasters(cluster)
	if all {
		for _, m := range masters {
			deleteInstance(m)
			time.Sleep(time.Duration(wait) * time.Second)
		}
	} else {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		unluckyMaster := masters[r.Intn(len(masters))]
		deleteInstance(unluckyMaster)

	}
}
