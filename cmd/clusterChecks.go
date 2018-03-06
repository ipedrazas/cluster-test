package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func lookupAPIServer() {
	ips, err := net.LookupIP("api.cerdanyola.k8s.sandbox.nutmeg.co.uk")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("api.cerdanyola.k8s.sandbox.nutmeg.co.uk. IN A %s\n", ip.String())
	}
}
func checkMasters() {
	// Verify the Masters Public IP have been assigned to route53 - api...
	// Verify the Masters private IP have been assigned to route53 - api.internal...
	// Verify EBS volumes are re-attached
	// Verify API_SERVER responds

}

func checkNodes() {
	// Verify that the instance group size is maintained
}

func checkPods() {
	// Verify new pods are scheduled and in running state
}

// MastersCheckHandler returns a json doc with the status of the set of masters
func MastersCheckHandler(w http.ResponseWriter, r *http.Request) {
	initAWSSession()
	vars := mux.Vars(r)
	cluster = vars["cluster"]
	dns, err := listCNAMES()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println(dns)

	instances := GetAllInstances()
	for i, inst := range instances {

		instances[i].DNS = hasDNS(inst, dns)
	}

	result := &CheckResult{
		Instances: instances,
		Route53:   dns,
		Cluster:   cluster,
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func hasDNS(i Instance, entries []RR) []string {
	var result []string
	for _, dns := range entries {
		for _, ips := range dns.Ips {
			fmt.Printf("%v - %v - %v \n", ips, i.PublicIP, i.PrivateIP)
			if ips == i.PublicIP || ips == i.PrivateIP {
				fmt.Println("Found!")
				result = append(result, dns.Name)
			}
		}
	}
	return result
}
