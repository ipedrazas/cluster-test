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
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

var port int

// httpServerCmd represents the httpServer command
var httpServerCmd = &cobra.Command{
	Use:   "httpServer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: initHTTPServer,
}

func init() {
	rootCmd.AddCommand(httpServerCmd)
	httpServerCmd.Flags().IntVarP(&port, "port", "p", 8000, "port where the http server will listen to.")

}

func initHTTPServer(cmd *cobra.Command, args []string) {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/pods/{namespace}/{filter}", PodsHandler)
	r.HandleFunc("/nodes/{cluster}/{all}", NodesHandler)
	r.HandleFunc("/master/{cluster}/{all}", MastersHandler)
	r.HandleFunc("/check/master/{cluster}", MastersCheckHandler)
	http.Handle("/", r)
	listenPort := fmt.Sprintf(":%v", strconv.Itoa(port))
	http.ListenAndServe(listenPort, r)
}

// HomeHandler is the default route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "cluster test %v", version)
	lookupAPIServer()
}
