// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sweeper

import (
	"testing"
)

// TestActualSweeperInventoryValidation performs various validations on the actual sweeper inventory
func TestActualSweeperInventoryValidation(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Make sure the sweeper inventory is populated
	if len(sweeperInventory) == 0 {
		t.Skip("sweeper inventory is empty, skipping test")
	}

	t.Logf("Running validations on actual sweeper inventory with %d sweepers", len(sweeperInventory))

	// Test 1: Check that all dependencies exist
	t.Run("dependencies_exist", func(t *testing.T) {
		for name, sweeper := range sweeperInventory {
			for _, dep := range sweeper.Dependencies {
				if _, exists := sweeperInventory[dep]; !exists {
					t.Errorf("Sweeper %s has dependency %s which doesn't exist in the inventory",
						name, dep)
				}
			}
		}
	})

	// Test 2: Check that all parents exist
	t.Run("parents_exist", func(t *testing.T) {
		for name, sweeper := range sweeperInventory {
			for _, parent := range sweeper.Parents {
				if _, exists := sweeperInventory[parent]; !exists {
					t.Errorf("Sweeper %s has parent %s which doesn't exist in the inventory",
						name, parent)
				}
			}
		}
	})

	// Test 3: Check that each resource doesn't list the same resource as both a dependency and a parent
	t.Run("no_self_contradictions", func(t *testing.T) {
		for name, sweeper := range sweeperInventory {
			depMap := make(map[string]bool)
			for _, dep := range sweeper.Dependencies {
				depMap[dep] = true
			}

			for _, parent := range sweeper.Parents {
				if depMap[parent] {
					t.Errorf("Sweeper %s lists %s as both a dependency and a parent, which is contradictory",
						name, parent)
				}
			}
		}
	})

	// Test 4: Check that parent resources have a ListAndAction function
	t.Run("parents_have_list_action", func(t *testing.T) {
		for name, sweeper := range sweeperInventory {
			for _, parentName := range sweeper.Parents {
				parent, exists := sweeperInventory[parentName]
				if !exists {
					continue // Already caught by parents_exist test
				}

				if parent.ListAndAction == nil {
					t.Errorf("Sweeper %s has parent %s, but parent lacks required ListAndAction function",
						name, parentName)
				}
			}
		}
	})

	// Test 5: Verify the topological sort produces a valid ordering
	t.Run("valid_topological_ordering", func(t *testing.T) {
		// Unify relationships
		unified := unifyRelationships(sweeperInventory)

		// Get ordering
		sorted := topologicalSort(unified)

		// Verify ordering is valid (each resource comes after its dependencies)
		resourcePos := make(map[string]int)
		for i, s := range sorted {
			resourcePos[s.Name] = i
		}

		for _, s := range sorted {
			myPos := resourcePos[s.Name]
			for _, dep := range s.Dependencies {
				depPos, exists := resourcePos[dep]
				if !exists {
					t.Errorf("Resource %s depends on %s which isn't in sort result", s.Name, dep)
					continue
				}

				if depPos > myPos {
					t.Errorf("Invalid ordering: %s depends on %s but comes before it in sort order",
						s.Name, dep)
				}
			}
		}
	})

	// Test 6: Ensure no cycles after unification (comprehensive check)
	t.Run("no_cycles_after_unification", func(t *testing.T) {
		unified := unifyRelationships(sweeperInventory)
		err := detectCycles(unified)
		if err != nil {
			t.Fatalf("Cycle detected in unified sweeper inventory: %v", err)
		}
	})
}
