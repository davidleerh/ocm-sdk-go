/*
Copyright (c) 2026 Red Hat, Inc.

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

// This example shows how to list node pools and display their spot market options.

package main

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/logging"
)

func main() {
	// Create a context:
	ctx := context.Background()

	// Create a logger that has the debug level enabled:
	logger, err := logging.NewGoLoggerBuilder().
		Debug(true).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build logger: %v\n", err)
		os.Exit(1)
	}

	// Create the connection, and remember to close it:
	token := os.Getenv("OCM_TOKEN")
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(token).
		BuildContext(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build connection: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	// Get the client for the node pools collection of a specific cluster:
	clusterID := os.Getenv("OCM_CLUSTER_ID")
	if clusterID == "" {
		fmt.Fprintf(os.Stderr, "OCM_CLUSTER_ID environment variable must be set\n")
		os.Exit(1)
	}
	collection := connection.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		NodePools()

	// List all node pools:
	response, err := collection.List().
		SendContext(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't list node pools: %v\n", err)
		os.Exit(1)
	}

	// Print each node pool and its spot configuration:
	for _, np := range response.Items().Slice() {
		fmt.Printf("Node pool: %s\n", np.ID())
		awsNP, ok := np.GetAWSNodePool()
		if !ok {
			continue
		}
		fmt.Printf("  Instance type: %s\n", awsNP.InstanceType())
		spot, ok := awsNP.GetSpotMarketOptions()
		if !ok {
			fmt.Printf("  Spot: not configured (on-demand)\n")
			continue
		}
		maxPrice, ok := spot.GetMaxPrice()
		if ok {
			fmt.Printf("  Spot max price: %s\n", maxPrice)
		} else {
			fmt.Printf("  Spot: enabled (using on-demand price as ceiling)\n")
		}
	}
}
