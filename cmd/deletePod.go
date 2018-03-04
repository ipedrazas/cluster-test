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
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// deletePodCmd represents the deletePod command
var deletePodCmd = &cobra.Command{
	Use:   "deletePod",
	Short: "deletes a pod",
	Long: `This command is used to delete one random pod of a list of pods selected by label. By
	using the --all flag you can delete all the pods that match the label selector.`,
	Run: podsHandler,
}
var filter string
var config *rest.Config
var podNames []string
var namespace string

func init() {
	rootCmd.AddCommand(deletePodCmd)
	deletePodCmd.Flags().StringVarP(&filter, "filter", "f", "", "label selector to fetch pods. Format is: key=value")
	deletePodCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace where the resources are located")
	deletePodCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace where the resources are located")
}

func podsHandler(cmd *cobra.Command, args []string) {
	if incluster {
		config = inCluster()
	}else{
		config = outCluster()
	}
	
	pods := getPods()
	deletePod(pods)
}

func getPods() []string {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: filter})
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func deletePod(podNames []string) error {

	if all {
		for _, n := range podNames {
			deleteResource(n)
			time.Sleep(time.Duration(wait) * time.Second)
		}
	} else {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		unluckyPod := podNames[r.Intn(len(podNames))]
		deleteResource(unluckyPod)
	}

	return nil
}

func deleteResource(podName string) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	err = clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})
	return err
}
