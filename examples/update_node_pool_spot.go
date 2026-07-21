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

// This example shows how to update spot market options on an existing node pool.

package main

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
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

	// Get the client for a specific node pool:
	clusterID := os.Getenv("OCM_CLUSTER_ID")
	nodePoolID := os.Getenv("OCM_NODE_POOL_ID")
	if clusterID == "" || nodePoolID == "" {
		fmt.Fprintf(os.Stderr, "OCM_CLUSTER_ID and OCM_NODE_POOL_ID environment variables must be set\n")
		os.Exit(1)
	}
	resource := connection.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		NodePools().
		NodePool(nodePoolID)

	// Build the update body with new spot market options:
	body, err := cmv1.NewNodePool().
		AWSNodePool(
			cmv1.NewAWSNodePool().
				SpotMarketOptions(
					cmv1.NewAwsNodePoolSpotMarketOptions().
						MaxPrice("0.08"),
				),
		).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build node pool update: %v\n", err)
		os.Exit(1)
	}

	// Send the update request:
	response, err := resource.Update().
		Body(body).
		SendContext(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't update node pool: %v\n", err)
		os.Exit(1)
	}

	// Print the result:
	np := response.Body()
	fmt.Printf("Node pool '%s' updated\n", np.ID())
	if awsNP, ok := np.GetAWSNodePool(); ok {
		if spot, ok := awsNP.GetSpotMarketOptions(); ok {
			fmt.Printf("  New spot max price: %s\n", spot.MaxPrice())
		}
	}
}
