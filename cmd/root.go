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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var awsRegion string
var dryrun bool
var all bool
var cluster string
var wait int
var incluster bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cluster-test",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// rootCmd represents the base command when called without any subcommands
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of the app",
	Long:  `Displays the version of the app.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cluster-test v0.1.0")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cluster-test.yaml)")
	rootCmd.PersistentFlags().StringVar(&awsRegion, "region", "eu-west-1", "AWS Region")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dry-run", false, "do a dry-run execution")
	rootCmd.PersistentFlags().BoolVar(&incluster, "in-cluster", false, "flag that indicates if the execution happens inside or outside a kubernetes cluster.")
	rootCmd.PersistentFlags().BoolVar(&all, "all", false, "tries to delete all the resources")
	rootCmd.PersistentFlags().StringVar(&cluster, "cluster", "c", "kubernetes cluster name")
	rootCmd.PersistentFlags().IntVarP(&wait, "wait", "d", 0, "delay between deletion of resources")
	rootCmd.AddCommand(versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cluster-test" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cluster-test")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
