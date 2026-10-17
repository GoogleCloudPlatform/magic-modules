package container

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepersLegacy("gcp_container_cluster", testSweepContainerClusters)
}

func testSweepContainerClusters(region string) error {
	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Fatalf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	// List clusters for all locations by using "-" as the location name
	parent := fmt.Sprintf("projects/%s/locations/-", config.Project)
	listURL := fmt.Sprintf("%s%s/clusters", transport_tpg.BaseUrl(Product, config), parent)
	found, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    listURL,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		log.Printf("error listing container clusters: %s", err)
		return nil
	}

	clusters, _ := found["clusters"].([]interface{})
	if len(clusters) == 0 {
		log.Printf("No container clusters found.")
		return nil
	}

	for _, rawCluster := range clusters {
		cluster, ok := rawCluster.(map[string]interface{})
		if !ok {
			continue
		}
		clusterName, _ := cluster["name"].(string)
		if !sweeper.IsSweepableTestResource(clusterName) {
			continue
		}

		if createTime, _ := cluster["createTime"].(string); createTime != "" {
			createdAt, err := time.Parse(time.RFC3339Nano, createTime)
			if err != nil {
				createdAt, err = time.Parse(time.RFC3339, createTime)
			}
			if err == nil && time.Since(createdAt) < 2*time.Hour {
				log.Printf("Skipping cluster %s, created recently (%s)", clusterName, time.Since(createdAt))
				continue
			}
		}

		status, _ := cluster["status"].(string)
		if status == "PROVISIONING" || status == "RECONCILING" || status == "STOPPING" {
			log.Printf("Skipping cluster %s in non-deletable status %s", clusterName, status)
			continue
		}

		log.Printf("Sweeping Container Cluster: %s", clusterName)
		clusterLocation, _ := cluster["location"].(string)
		clusterURL := fmt.Sprintf("%sprojects/%s/locations/%s/clusters/%s", transport_tpg.BaseUrl(Product, config), config.Project, clusterLocation, clusterName)
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "DELETE",
			Project:   config.Project,
			RawURL:    clusterURL,
			UserAgent: config.UserAgent,
		})

		if err != nil {
			log.Printf("Error, failed to delete cluster %s: %s", clusterName, err)
			continue
		}
	}

	return nil
}
