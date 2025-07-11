package diff

import (
	"os"
	"testing"

	"yspec/pkg/k8s"
)

func TestIntegrationScenarios(t *testing.T) {
	tests := []struct {
		name                string
		beforeFile          string
		afterFile           string
		expectedAdded       int
		expectedRemoved     int
		expectedUpdated     int
		expectedUnchanged   int
	}{
		{
			name:                "identical objects",
			beforeFile:          "examples/test-cases/identical-before.yaml",
			afterFile:           "examples/test-cases/identical-after.yaml",
			expectedAdded:       0,
			expectedRemoved:     0,
			expectedUpdated:     0,
			expectedUnchanged:   1,
		},
		{
			name:                "replica change",
			beforeFile:          "examples/test-cases/replica-change-before.yaml",
			afterFile:           "examples/test-cases/replica-change-after.yaml",
			expectedAdded:       0,
			expectedRemoved:     0,
			expectedUpdated:     1,
			expectedUnchanged:   0,
		},
		{
			name:                "create and delete different objects",
			beforeFile:          "examples/test-cases/create-delete-before.yaml",
			afterFile:           "examples/test-cases/create-delete-after.yaml",
			expectedAdded:       1,
			expectedRemoved:     1,
			expectedUpdated:     0,
			expectedUnchanged:   0,
		},
		{
			name:                "empty file to new object",
			beforeFile:          "examples/test-cases/empty-before.yaml",
			afterFile:           "examples/test-cases/empty-after.yaml",
			expectedAdded:       1,
			expectedRemoved:     0,
			expectedUpdated:     0,
			expectedUnchanged:   0,
		},
		{
			name:                "object to empty file",
			beforeFile:          "examples/test-cases/empty-after.yaml",
			afterFile:           "examples/test-cases/empty-before.yaml",
			expectedAdded:       0,
			expectedRemoved:     1,
			expectedUpdated:     0,
			expectedUnchanged:   0,
		},
		{
			name:                "mixed changes with unchanged resource",
			beforeFile:          "examples/test-cases/mixed-changes-before.yaml",
			afterFile:           "examples/test-cases/mixed-changes-after.yaml",
			expectedAdded:       1,
			expectedRemoved:     1,
			expectedUpdated:     1,
			expectedUnchanged:   1,
		},
		{
			name:                "yaml order independence",
			beforeFile:          "examples/test-cases/yaml-order-before.yaml",
			afterFile:           "examples/test-cases/yaml-order-after.yaml",
			expectedAdded:       0,
			expectedRemoved:     0,
			expectedUpdated:     0,
			expectedUnchanged:   2,
		},
		{
			name:                "no namespace defaults to default",
			beforeFile:          "examples/test-cases/no-namespace-before.yaml",
			afterFile:           "examples/test-cases/no-namespace-after.yaml",
			expectedAdded:       0,
			expectedRemoved:     0,
			expectedUpdated:     1,
			expectedUnchanged:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse before file
			beforeFile, err := os.Open("../../" + tt.beforeFile)
			if err != nil {
				t.Fatalf("failed to open before file: %v", err)
			}
			defer beforeFile.Close()

			beforeObjects, err := k8s.ParseYAMLStream(beforeFile)
			if err != nil {
				t.Fatalf("failed to parse before file: %v", err)
			}

			// Parse after file
			afterFile, err := os.Open("../../" + tt.afterFile)
			if err != nil {
				t.Fatalf("failed to open after file: %v", err)
			}
			defer afterFile.Close()

			afterObjects, err := k8s.ParseYAMLStream(afterFile)
			if err != nil {
				t.Fatalf("failed to parse after file: %v", err)
			}

			// Generate diff
			result, err := Generate(beforeObjects, afterObjects)
			if err != nil {
				t.Fatalf("failed to generate diff: %v", err)
			}

			// Check counts
			if len(result.Added) != tt.expectedAdded {
				t.Errorf("expected %d added resources, got %d", tt.expectedAdded, len(result.Added))
			}
			if len(result.Removed) != tt.expectedRemoved {
				t.Errorf("expected %d removed resources, got %d", tt.expectedRemoved, len(result.Removed))
			}
			if len(result.Updated) != tt.expectedUpdated {
				t.Errorf("expected %d updated resources, got %d", tt.expectedUpdated, len(result.Updated))
			}

			// Calculate unchanged resources
			totalBefore := len(beforeObjects)
			totalAfter := len(afterObjects)
			actualUnchanged := 0

			// Count resources that exist in both and are unchanged
			for key := range beforeObjects {
				if _, exists := afterObjects[key]; exists {
					if _, updated := result.Updated[key]; !updated {
						actualUnchanged++
					}
				}
			}

			if actualUnchanged != tt.expectedUnchanged {
				t.Errorf("expected %d unchanged resources, got %d", tt.expectedUnchanged, actualUnchanged)
			}

			// Verify the math: before + added - removed = after
			expectedAfterTotal := totalBefore + len(result.Added) - len(result.Removed)
			if totalAfter != expectedAfterTotal {
				t.Errorf("resource count mismatch: before=%d + added=%d - removed=%d should equal after=%d, but got %d",
					totalBefore, len(result.Added), len(result.Removed), expectedAfterTotal, totalAfter)
			}
		})
	}
}

func TestSpecificScenarios(t *testing.T) {
	t.Run("replica change shows correct diff", func(t *testing.T) {
		beforeFile, err := os.Open("../../examples/test-cases/replica-change-before.yaml")
		if err != nil {
			t.Fatalf("failed to open before file: %v", err)
		}
		defer beforeFile.Close()

		beforeObjects, err := k8s.ParseYAMLStream(beforeFile)
		if err != nil {
			t.Fatalf("failed to parse before file: %v", err)
		}

		afterFile, err := os.Open("../../examples/test-cases/replica-change-after.yaml")
		if err != nil {
			t.Fatalf("failed to open after file: %v", err)
		}
		defer afterFile.Close()

		afterObjects, err := k8s.ParseYAMLStream(afterFile)
		if err != nil {
			t.Fatalf("failed to parse after file: %v", err)
		}

		result, err := Generate(beforeObjects, afterObjects)
		if err != nil {
			t.Fatalf("failed to generate diff: %v", err)
		}

		// Should have exactly one updated resource
		if len(result.Updated) != 1 {
			t.Fatalf("expected 1 updated resource, got %d", len(result.Updated))
		}

		// Get the updated resource
		var updatedResource UpdatedResource
		for _, resource := range result.Updated {
			updatedResource = resource
			break
		}

		// Check that the diff shows replica change
		if diffMap, ok := updatedResource.Diff["spec"].(map[string]interface{}); ok {
			if replicaDiff, ok := diffMap["replicas"].(map[string]interface{}); ok {
				if beforeVal, ok := replicaDiff["__before"].(int); ok {
					if beforeVal != 2 {
						t.Errorf("expected before replicas to be 2, got %d", beforeVal)
					}
				} else {
					t.Error("expected __before in replica diff")
				}
				if afterVal, ok := replicaDiff["__after"].(int); ok {
					if afterVal != 5 {
						t.Errorf("expected after replicas to be 5, got %d", afterVal)
					}
				} else {
					t.Error("expected __after in replica diff")
				}
			} else {
				t.Error("expected replicas in spec diff")
			}
		} else {
			t.Error("expected spec in diff")
		}
	})

	t.Run("mixed changes identifies correct resources", func(t *testing.T) {
		beforeFile, err := os.Open("../../examples/test-cases/mixed-changes-before.yaml")
		if err != nil {
			t.Fatalf("failed to open before file: %v", err)
		}
		defer beforeFile.Close()

		beforeObjects, err := k8s.ParseYAMLStream(beforeFile)
		if err != nil {
			t.Fatalf("failed to parse before file: %v", err)
		}

		afterFile, err := os.Open("../../examples/test-cases/mixed-changes-after.yaml")
		if err != nil {
			t.Fatalf("failed to open after file: %v", err)
		}
		defer afterFile.Close()

		afterObjects, err := k8s.ParseYAMLStream(afterFile)
		if err != nil {
			t.Fatalf("failed to parse after file: %v", err)
		}

		result, err := Generate(beforeObjects, afterObjects)
		if err != nil {
			t.Fatalf("failed to generate diff: %v", err)
		}

		// Check specific resources
		if _, exists := result.Added["v1/ConfigMap/default/new-config"]; !exists {
			t.Error("expected new-config to be added")
		}
		if _, exists := result.Removed["v1/ConfigMap/default/old-config"]; !exists {
			t.Error("expected old-config to be removed")
		}
		if _, exists := result.Updated["apps/v1/Deployment/default/changed-app"]; !exists {
			t.Error("expected changed-app to be updated")
		}

		// Check that unchanged-config is not in any diff section
		unchangedKey := "v1/ConfigMap/default/unchanged-config"
		if _, exists := result.Added[unchangedKey]; exists {
			t.Error("unchanged-config should not be in added")
		}
		if _, exists := result.Removed[unchangedKey]; exists {
			t.Error("unchanged-config should not be in removed")
		}
		if _, exists := result.Updated[unchangedKey]; exists {
			t.Error("unchanged-config should not be in updated")
		}
	})
}