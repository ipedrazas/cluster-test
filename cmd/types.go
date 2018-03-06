package cmd

import "time"

// RR is a helper struct to fetch dns entries
type RR struct {
	Name string
	Ips  []string
}

// Instance is a helper struct to fetch EC2 instances
type Instance struct {
	Name       string
	PublicIP   string
	PrivateIP  string
	State      string
	LaunchTime time.Time
	ID         string
	IsMaster   bool
	DNS        []string
}

// CheckResult is a helper struct to return a json payload
type CheckResult struct {
	Instances []Instance
	Route53   []RR
	Cluster   string
}
