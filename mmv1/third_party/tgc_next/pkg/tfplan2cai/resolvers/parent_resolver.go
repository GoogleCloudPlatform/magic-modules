package resolvers

import (
	"fmt"
	"sort"
	"strings"

	"go.uber.org/zap"

	provider "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/provider"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/tfplan"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ParentResourceResolver struct {
	schema *schema.Provider

	// For logging error / status information that doesn't warrant an outright failure
	errorLogger *zap.Logger
}

func NewParentResourceResolver(errorLogger *zap.Logger) *ParentResourceResolver {
	return &ParentResourceResolver{
		schema:      provider.Provider(),
		errorLogger: errorLogger,
	}
}

func (r *ParentResourceResolver) Resolve(jsonPlan []byte) map[string][]string {
	parentToChildMap := make(map[string][]string)

	// Read elements from the resouce config
	resourceConfig, err := tfplan.ReadResourceConfigurations(jsonPlan)
	if err != nil {
		return parentToChildMap
	}

	for _, resource := range resourceConfig.RootModule.Resources {
		for _, expression := range resource.Expressions {
			if expression.ExpressionData.NestedBlocks != nil {
				for _, innerExexpression := range expression.ExpressionData.NestedBlocks {
					for _, v := range innerExexpression {
						reference := v.References
						if reference != nil {
							if strings.HasSuffix(reference[0], ".id") {
								parentToChildMap[reference[1]] = append(parentToChildMap[reference[1]], resource.Address)
							}
						}
					}
				}
			}
			reference := expression.ExpressionData.References
			if reference != nil {
				if strings.HasSuffix(reference[0], ".id") {
					parentToChildMap[reference[1]] = append(parentToChildMap[reference[1]], resource.Address)
				}
			}
		}
	}

	return parentToChildMap
}

func sortTraversalOrder(graph map[string][]string) (map[int][]string, error) {
	inDegree := make(map[string]int)
	allNodes := make(map[string]bool)

	// Step 1: Identify all unique nodes and calculate in-degrees.
	for u, neighbors := range graph {
		allNodes[u] = true
		if _, exists := inDegree[u]; !exists {
			inDegree[u] = 0
		}
		for _, v := range neighbors {
			allNodes[v] = true
			if _, exists := inDegree[v]; !exists {
				inDegree[v] = 0
			}
			inDegree[v]++
		}
	}

	// Step 2: Initialize queue with nodes having an in-degree of 0.
	nodeLevels := make(map[string]int)
	queue := []string{}

	for node := range allNodes {
		if inDegree[node] == 0 {
			queue = append(queue, node)
			nodeLevels[node] = 0 // Source nodes are at level 0.
		}
	}
	sort.Strings(queue) // Start with sorted source nodes for deterministic behavior.

	// Step 3: Process nodes to determine levels (longest path from any source).
	head := 0
	processedCount := 0
	for head < len(queue) {
		u := queue[head]
		head++
		processedCount++

		currentLevel := nodeLevels[u]

		// Sort neighbors to ensure deterministic queue order if multiple nodes
		// simultaneously reach an in-degree of 0.
		neighbors := graph[u]
		if neighbors == nil {
			neighbors = []string{}
		}
		sort.Strings(neighbors)

		for _, v := range neighbors {
			// Update level of v to be the maximum of its current vs. new path length.
			if newLevel := currentLevel + 1; newLevel > nodeLevels[v] {
				nodeLevels[v] = newLevel
			}

			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	// Step 4: Check for cycles.
	if processedCount < len(allNodes) {
		cycleRelatedNodes := []string{}
		for node := range allNodes {
			if _, visited := nodeLevels[node]; !visited {
				cycleRelatedNodes = append(cycleRelatedNodes, node)
			}
		}
		sort.Strings(cycleRelatedNodes)
		return nil, fmt.Errorf("cycle detected in graph, nodes not leveled: %v", cycleRelatedNodes)
	}

	// Step 5: Transform nodeLevels (map[string]int) to levelsToNodes (map[int][]string).
	levelsToNodes := make(map[int][]string)
	for node, level := range nodeLevels {
		levelsToNodes[level] = append(levelsToNodes[level], node)
	}

	// Step 6: Sort the node lists within each level for deterministic output.
	for level := range levelsToNodes {
		sort.Strings(levelsToNodes[level])
	}

	return levelsToNodes, nil
}
