package diff

import (
	"os"
	"testing"

	"skiff/pkg/k8s"
)

func TestChangesField(t *testing.T) {
	t.Run("simple field change generates correct changes", func(t *testing.T) {
		beforeFile, err := os.Open("../../test/test-cases/replica-change-before.yaml")
		if err != nil {
			t.Fatalf("failed to open before file: %v", err)
		}
		defer beforeFile.Close()

		afterFile, err := os.Open("../../test/test-cases/replica-change-after.yaml")
		if err != nil {
			t.Fatalf("failed to open after file: %v", err)
		}
		defer afterFile.Close()

		beforeObjects, err := k8s.ParseYAMLStream(beforeFile)
		if err != nil {
			t.Fatalf("failed to parse before YAML: %v", err)
		}

		afterObjects, err := k8s.ParseYAMLStream(afterFile)
		if err != nil {
			t.Fatalf("failed to parse after YAML: %v", err)
		}

		result, err := GenerateTerraformStyle(beforeObjects, afterObjects)
		if err != nil {
			t.Fatalf("failed to generate diff: %v", err)
		}

		// Should have one resource change
		if len(result.ResourceChanges) != 1 {
			t.Errorf("expected 1 resource change, got %d", len(result.ResourceChanges))
		}

		// Find the deployment change
		var deploymentChange *ResourceChange
		for _, change := range result.ResourceChanges {
			if change.Type == "Deployment" {
				deploymentChange = &change
				break
			}
		}

		if deploymentChange == nil {
			t.Fatal("expected deployment change not found")
		}

		// Verify the changes field
		changes := deploymentChange.Change.Changes
		if len(changes) != 1 {
			t.Errorf("expected 1 field change, got %d", len(changes))
		}

		replicaChange, exists := changes["spec.replicas"]
		if !exists {
			t.Error("expected spec.replicas change not found")
		}

		// Check the actual values (YAML can parse as int or float64)
		fromInt, fromOk := replicaChange.From.(int)
		toInt, toOk := replicaChange.To.(int)
		if !fromOk || !toOk || fromInt != 2 || toInt != 5 {
			t.Errorf("expected replica change from 2 to 5, got from %v (%T) to %v (%T)", 
			         replicaChange.From, replicaChange.From, replicaChange.To, replicaChange.To)
		}
	})

	t.Run("complex nested changes generate flattened paths", func(t *testing.T) {
		beforeFile, err := os.Open("../../test/test-cases/hpa-before.yaml")
		if err != nil {
			t.Fatalf("failed to open before file: %v", err)
		}
		defer beforeFile.Close()

		afterFile, err := os.Open("../../test/test-cases/hpa-after.yaml")
		if err != nil {
			t.Fatalf("failed to open after file: %v", err)
		}
		defer afterFile.Close()

		beforeObjects, err := k8s.ParseYAMLStream(beforeFile)
		if err != nil {
			t.Fatalf("failed to parse before YAML: %v", err)
		}

		afterObjects, err := k8s.ParseYAMLStream(afterFile)
		if err != nil {
			t.Fatalf("failed to parse after YAML: %v", err)
		}

		result, err := GenerateTerraformStyle(beforeObjects, afterObjects)
		if err != nil {
			t.Fatalf("failed to generate diff: %v", err)
		}

		// Find the HPA change
		var hpaChange *ResourceChange
		for _, change := range result.ResourceChanges {
			if change.Type == "HorizontalPodAutoscaler" {
				hpaChange = &change
				break
			}
		}

		if hpaChange == nil {
			t.Fatal("expected HPA change not found")
		}

		changes := hpaChange.Change.Changes

		// Should have maxReplicas change
		replicaChange, exists := changes["spec.maxReplicas"]
		if !exists {
			t.Error("expected spec.maxReplicas change")
		}
		fromInt, fromOk := replicaChange.From.(int)
		toInt, toOk := replicaChange.To.(int)
		if !fromOk || !toOk || fromInt != 10 || toInt != 15 {
			t.Errorf("expected maxReplicas change from 10 to 15, got from %v (%T) to %v (%T)", 
			         replicaChange.From, replicaChange.From, replicaChange.To, replicaChange.To)
		}

		// Should have behavior policy changes with array indices
		scaleUpValue, exists := changes["spec.behavior.scaleUp.policies[0].value"]
		if !exists {
			t.Error("expected spec.behavior.scaleUp.policies[0].value change")
		}
		valueInt, valueOk := scaleUpValue.To.(int)
		if scaleUpValue.From != nil || !valueOk || valueInt != 100 {
			t.Errorf("expected new field from nil to 100, got %v to %v (%T)", 
			         scaleUpValue.From, scaleUpValue.To, scaleUpValue.To)
		}

		// Verify flattened paths for nested arrays
		expectedPaths := []string{
			"spec.behavior.scaleUp.policies[0].type",
			"spec.behavior.scaleUp.policies[0].value", 
			"spec.behavior.scaleUp.policies[1].type",
			"spec.behavior.scaleDown.policies[0].type",
		}

		for _, path := range expectedPaths {
			if _, exists := changes[path]; !exists {
				t.Errorf("expected flattened path %s not found", path)
			}
		}
	})

	t.Run("create and delete actions", func(t *testing.T) {
		beforeFile, err := os.Open("../../test/test-cases/create-delete-before.yaml")
		if err != nil {
			t.Fatalf("failed to open before file: %v", err)
		}
		defer beforeFile.Close()

		afterFile, err := os.Open("../../test/test-cases/create-delete-after.yaml")
		if err != nil {
			t.Fatalf("failed to open after file: %v", err)
		}
		defer afterFile.Close()

		beforeObjects, err := k8s.ParseYAMLStream(beforeFile)
		if err != nil {
			t.Fatalf("failed to parse before YAML: %v", err)
		}

		afterObjects, err := k8s.ParseYAMLStream(afterFile)
		if err != nil {
			t.Fatalf("failed to parse after YAML: %v", err)
		}

		result, err := GenerateTerraformStyle(beforeObjects, afterObjects)
		if err != nil {
			t.Fatalf("failed to generate diff: %v", err)
		}

		// Should have multiple resource changes
		if len(result.ResourceChanges) == 0 {
			t.Error("expected resource changes")
		}

		// Check for different action types
		var foundCreate, foundDelete bool
		for _, change := range result.ResourceChanges {
			if len(change.Change.Actions) > 0 {
				action := change.Change.Actions[0]
				if action == "create" {
					foundCreate = true
					// Create should have After but no Before
					if change.Change.Before != nil {
						t.Error("create action should not have Before object")
					}
					if change.Change.After == nil {
						t.Error("create action should have After object")
					}
				} else if action == "delete" {
					foundDelete = true
					// Delete should have Before but no After
					if change.Change.Before == nil {
						t.Error("delete action should have Before object")
					}
					if change.Change.After != nil {
						t.Error("delete action should not have After object")
					}
				}
			}
		}

		if !foundCreate {
			t.Error("expected to find create action")
		}
		if !foundDelete {
			t.Error("expected to find delete action")
		}
	})
}