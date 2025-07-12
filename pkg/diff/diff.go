package diff

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// TerraformStyleResult represents a flat diff format for easier policy writing
type TerraformStyleResult struct {
	ResourceChanges map[string]ResourceChange `json:"resource_changes"`
}

// ResourceChange represents a single resource change in Terraform style
type ResourceChange struct {
	Type       string `json:"type"`
	APIVersion string `json:"apiVersion"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Change     Change `json:"change"`
}

// FieldChange represents a change to a specific field
type FieldChange struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

// Change represents the before/after state and actions
type Change struct {
	Actions []string               `json:"actions"`
	Before  map[string]interface{} `json:"before,omitempty"`
	After   map[string]interface{} `json:"after,omitempty"`
	Changes map[string]FieldChange `json:"changes,omitempty"`
}

// GenerateTerraformStyle creates a flat diff format for easier policy writing
func GenerateTerraformStyle(before, after map[string]map[string]interface{}) (*TerraformStyleResult, error) {
	result := &TerraformStyleResult{
		ResourceChanges: make(map[string]ResourceChange),
	}

	// Parse all resource keys to extract metadata
	allKeys := make(map[string]bool)
	for key := range before {
		allKeys[key] = true
	}
	for key := range after {
		allKeys[key] = true
	}

	// Process each resource
	for key := range allKeys {
		beforeObj, beforeExists := before[key]
		afterObj, afterExists := after[key]

		// Extract resource metadata from key (apiVersion/kind/namespace/name)
		apiVersion, kind, namespace, name := parseResourceKey(key)

		var change Change
		var actions []string

		if !beforeExists && afterExists {
			// Resource was created
			actions = []string{"create"}
			change = Change{
				Actions: actions,
				After:   afterObj,
			}
		} else if beforeExists && !afterExists {
			// Resource was deleted
			actions = []string{"delete"}
			change = Change{
				Actions: actions,
				Before:  beforeObj,
			}
		} else if beforeExists && afterExists {
			// Resource might be updated
			if !cmp.Equal(beforeObj, afterObj) {
				actions = []string{"update"}
				fieldChanges := generateFieldChanges(beforeObj, afterObj, "")
				change = Change{
					Actions: actions,
					Before:  beforeObj,
					After:   afterObj,
					Changes: fieldChanges,
				}
			} else {
				// No change, skip
				continue
			}
		}

		result.ResourceChanges[key] = ResourceChange{
			Type:       kind,
			APIVersion: apiVersion,
			Namespace:  namespace,
			Name:       name,
			Change:     change,
		}
	}

	return result, nil
}

// parseResourceKey extracts metadata from resource key format: apiVersion/kind/namespace/name
// Note: apiVersion may contain slashes (e.g., "autoscaling/v2")
func parseResourceKey(key string) (apiVersion, kind, namespace, name string) {
	parts := strings.Split(key, "/")
	if len(parts) >= 4 {
		// Handle cases where apiVersion contains slashes
		if len(parts) == 4 {
			// Simple case: v1/Kind/namespace/name
			apiVersion = parts[0]
			kind = parts[1]
			namespace = parts[2]
			name = parts[3]
		} else if len(parts) == 5 {
			// Complex case: group/version/Kind/namespace/name
			apiVersion = parts[0] + "/" + parts[1]
			kind = parts[2]
			namespace = parts[3]
			name = parts[4]
		}
	}
	return
}

// isMapType checks if a value is a map[string]interface{}
func isMapType(val interface{}) bool {
	_, ok := val.(map[string]interface{})
	return ok
}

// generateFieldChanges recursively compares two objects and generates flattened field changes
func generateFieldChanges(before, after map[string]interface{}, prefix string) map[string]FieldChange {
	changes := make(map[string]FieldChange)

	// Get all keys from both objects
	allKeys := make(map[string]bool)
	for key := range before {
		allKeys[key] = true
	}
	for key := range after {
		allKeys[key] = true
	}

	for key := range allKeys {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		beforeVal, beforeExists := before[key]
		afterVal, afterExists := after[key]

		if !beforeExists && afterExists {
			// Field was added
			flattenValue(changes, path, nil, afterVal)
		} else if beforeExists && !afterExists {
			// Field was removed
			flattenValue(changes, path, beforeVal, nil)
		} else if beforeExists && afterExists {
			// Field exists in both, check if changed
			compareValues(changes, path, beforeVal, afterVal)
		}
	}

	return changes
}

// flattenValue adds all paths in a value to changes (for additions/deletions)
func flattenValue(changes map[string]FieldChange, path string, from, to interface{}) {
	if isMapType(to) {
		// If it's a map, recursively add all nested paths
		toMap := to.(map[string]interface{})
		for key, value := range toMap {
			nestedPath := path + "." + key
			flattenValue(changes, nestedPath, nil, value)
		}
	} else if isSliceType(to) {
		// If it's a slice, add paths with indices
		toSlice := to.([]interface{})
		for i, value := range toSlice {
			indexPath := path + "[" + fmt.Sprintf("%d", i) + "]"
			flattenValue(changes, indexPath, nil, value)
		}
	} else {
		// Scalar value
		changes[path] = FieldChange{From: from, To: to}
	}

	// Also handle the case where from is not nil (deletions)
	if isMapType(from) {
		fromMap := from.(map[string]interface{})
		for key, value := range fromMap {
			nestedPath := path + "." + key
			flattenValue(changes, nestedPath, value, nil)
		}
	} else if isSliceType(from) {
		fromSlice := from.([]interface{})
		for i, value := range fromSlice {
			indexPath := path + "[" + fmt.Sprintf("%d", i) + "]"
			flattenValue(changes, indexPath, value, nil)
		}
	} else if from != nil {
		changes[path] = FieldChange{From: from, To: to}
	}
}

// compareValues compares two values and adds changes if they differ
func compareValues(changes map[string]FieldChange, path string, before, after interface{}) {
	if cmp.Equal(before, after) {
		return // No change
	}

	// If both are maps, recursively compare
	if isMapType(before) && isMapType(after) {
		beforeMap := before.(map[string]interface{})
		afterMap := after.(map[string]interface{})
		nestedChanges := generateFieldChanges(beforeMap, afterMap, path)
		for nestedPath, change := range nestedChanges {
			changes[nestedPath] = change
		}
		return
	}

	// If both are slices, compare element by element
	if isSliceType(before) && isSliceType(after) {
		beforeSlice := before.([]interface{})
		afterSlice := after.([]interface{})
		compareSlices(changes, path, beforeSlice, afterSlice)
		return
	}

	// Different types or scalar values - record the change
	changes[path] = FieldChange{From: before, To: after}
}

// compareSlices compares two slices element by element
func compareSlices(changes map[string]FieldChange, path string, before, after []interface{}) {
	maxLen := len(before)
	if len(after) > maxLen {
		maxLen = len(after)
	}

	for i := 0; i < maxLen; i++ {
		indexPath := path + "[" + fmt.Sprintf("%d", i) + "]"

		if i >= len(before) {
			// Element added
			flattenValue(changes, indexPath, nil, after[i])
		} else if i >= len(after) {
			// Element removed
			flattenValue(changes, indexPath, before[i], nil)
		} else {
			// Element exists in both
			compareValues(changes, indexPath, before[i], after[i])
		}
	}
}

// isSliceType checks if a value is a slice
func isSliceType(val interface{}) bool {
	_, ok := val.([]interface{})
	return ok
}
