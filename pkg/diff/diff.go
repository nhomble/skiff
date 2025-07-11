package diff

import (
	"github.com/google/go-cmp/cmp"
)

// Result represents the structured diff output
type Result struct {
	Added   map[string]AddedResource   `json:"added"`
	Removed map[string]RemovedResource `json:"removed"`
	Updated map[string]UpdatedResource `json:"updated"`
}

// AddedResource represents a resource that was added
type AddedResource struct {
	After map[string]interface{} `json:"after"`
}

// RemovedResource represents a resource that was removed
type RemovedResource struct {
	Before map[string]interface{} `json:"before"`
}

// UpdatedResource represents a resource that was modified
type UpdatedResource struct {
	Before map[string]interface{} `json:"before"`
	After  map[string]interface{} `json:"after"`
	Diff   map[string]interface{} `json:"diff"`
}

// Generate compares two sets of K8s objects and produces a structured diff
func Generate(before, after map[string]map[string]interface{}) (*Result, error) {
	result := &Result{
		Added:   make(map[string]AddedResource),
		Removed: make(map[string]RemovedResource),
		Updated: make(map[string]UpdatedResource),
	}

	for key, afterObj := range after {
		if beforeObj, exists := before[key]; exists {
			if !cmp.Equal(beforeObj, afterObj) {
				diff := generateStructuredDiff(beforeObj, afterObj)
				result.Updated[key] = UpdatedResource{
					Before: beforeObj,
					After:  afterObj,
					Diff:   diff,
				}
			}
		} else {
			result.Added[key] = AddedResource{
				After: afterObj,
			}
		}
	}

	for key, beforeObj := range before {
		if _, exists := after[key]; !exists {
			result.Removed[key] = RemovedResource{
				Before: beforeObj,
			}
		}
	}

	return result, nil
}

// generateStructuredDiff creates a structured diff showing what changed
func generateStructuredDiff(before, after map[string]interface{}) map[string]interface{} {
	diff := make(map[string]interface{})
	
	allKeys := make(map[string]bool)
	for key := range before {
		allKeys[key] = true
	}
	for key := range after {
		allKeys[key] = true
	}

	for key := range allKeys {
		beforeVal, beforeExists := before[key]
		afterVal, afterExists := after[key]

		if !beforeExists {
			diff[key] = map[string]interface{}{
				"__added": afterVal,
			}
		} else if !afterExists {
			diff[key] = map[string]interface{}{
				"__removed": beforeVal,
			}
		} else if !cmp.Equal(beforeVal, afterVal) {
			if isMapType(beforeVal) && isMapType(afterVal) {
				nestedDiff := generateStructuredDiff(
					beforeVal.(map[string]interface{}),
					afterVal.(map[string]interface{}),
				)
				if len(nestedDiff) > 0 {
					diff[key] = nestedDiff
				}
			} else {
				diff[key] = map[string]interface{}{
					"__before": beforeVal,
					"__after":  afterVal,
				}
			}
		}
	}

	return diff
}

// isMapType checks if a value is a map[string]interface{}
func isMapType(val interface{}) bool {
	_, ok := val.(map[string]interface{})
	return ok
}