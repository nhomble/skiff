package diff

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		before   map[string]map[string]interface{}
		after    map[string]map[string]interface{}
		expected *Result
	}{
		{
			name:   "identical objects",
			before: map[string]map[string]interface{}{
				"v1/ConfigMap/default/test": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
			},
			after: map[string]map[string]interface{}{
				"v1/ConfigMap/default/test": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
			},
			expected: &Result{
				Added:   map[string]AddedResource{},
				Removed: map[string]RemovedResource{},
				Updated: map[string]UpdatedResource{},
			},
		},
		{
			name: "simple field change",
			before: map[string]map[string]interface{}{
				"apps/v1/Deployment/default/test": {
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "default",
					},
					"spec": map[string]interface{}{
						"replicas": 2,
					},
				},
			},
			after: map[string]map[string]interface{}{
				"apps/v1/Deployment/default/test": {
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "default",
					},
					"spec": map[string]interface{}{
						"replicas": 3,
					},
				},
			},
			expected: &Result{
				Added:   map[string]AddedResource{},
				Removed: map[string]RemovedResource{},
				Updated: map[string]UpdatedResource{
					"apps/v1/Deployment/default/test": {
						Before: map[string]interface{}{
							"apiVersion": "apps/v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name":      "test",
								"namespace": "default",
							},
							"spec": map[string]interface{}{
								"replicas": 2,
							},
						},
						After: map[string]interface{}{
							"apiVersion": "apps/v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name":      "test",
								"namespace": "default",
							},
							"spec": map[string]interface{}{
								"replicas": 3,
							},
						},
						Diff: map[string]interface{}{
							"spec": map[string]interface{}{
								"replicas": map[string]interface{}{
									"__before": 2,
									"__after":  3,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "added resource",
			before: map[string]map[string]interface{}{},
			after: map[string]map[string]interface{}{
				"v1/ConfigMap/default/new": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "new",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
			},
			expected: &Result{
				Added: map[string]AddedResource{
					"v1/ConfigMap/default/new": {
						After: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name":      "new",
								"namespace": "default",
							},
							"data": map[string]interface{}{
								"key": "value",
							},
						},
					},
				},
				Removed: map[string]RemovedResource{},
				Updated: map[string]UpdatedResource{},
			},
		},
		{
			name: "removed resource",
			before: map[string]map[string]interface{}{
				"v1/ConfigMap/default/old": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "old",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
			},
			after: map[string]map[string]interface{}{},
			expected: &Result{
				Added: map[string]AddedResource{},
				Removed: map[string]RemovedResource{
					"v1/ConfigMap/default/old": {
						Before: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name":      "old",
								"namespace": "default",
							},
							"data": map[string]interface{}{
								"key": "value",
							},
						},
					},
				},
				Updated: map[string]UpdatedResource{},
			},
		},
		{
			name: "mixed changes with unchanged resource",
			before: map[string]map[string]interface{}{
				"v1/ConfigMap/default/unchanged": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "unchanged",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
				"v1/ConfigMap/default/old": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "old",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
				"apps/v1/Deployment/default/changed": {
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      "changed",
						"namespace": "default",
					},
					"spec": map[string]interface{}{
						"replicas": 1,
					},
				},
			},
			after: map[string]map[string]interface{}{
				"v1/ConfigMap/default/unchanged": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "unchanged",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
				"v1/ConfigMap/default/new": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "new",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key": "value",
					},
				},
				"apps/v1/Deployment/default/changed": {
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      "changed",
						"namespace": "default",
					},
					"spec": map[string]interface{}{
						"replicas": 2,
					},
				},
			},
			expected: &Result{
				Added: map[string]AddedResource{
					"v1/ConfigMap/default/new": {
						After: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name":      "new",
								"namespace": "default",
							},
							"data": map[string]interface{}{
								"key": "value",
							},
						},
					},
				},
				Removed: map[string]RemovedResource{
					"v1/ConfigMap/default/old": {
						Before: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name":      "old",
								"namespace": "default",
							},
							"data": map[string]interface{}{
								"key": "value",
							},
						},
					},
				},
				Updated: map[string]UpdatedResource{
					"apps/v1/Deployment/default/changed": {
						Before: map[string]interface{}{
							"apiVersion": "apps/v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name":      "changed",
								"namespace": "default",
							},
							"spec": map[string]interface{}{
								"replicas": 1,
							},
						},
						After: map[string]interface{}{
							"apiVersion": "apps/v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name":      "changed",
								"namespace": "default",
							},
							"spec": map[string]interface{}{
								"replicas": 2,
							},
						},
						Diff: map[string]interface{}{
							"spec": map[string]interface{}{
								"replicas": map[string]interface{}{
									"__before": 1,
									"__after":  2,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "nested field changes",
			before: map[string]map[string]interface{}{
				"v1/ConfigMap/default/test": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			after: map[string]map[string]interface{}{
				"v1/ConfigMap/default/test": {
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name":      "test",
						"namespace": "default",
					},
					"data": map[string]interface{}{
						"key1": "value1",
						"key2": "new_value2",
						"key3": "value3",
					},
				},
			},
			expected: &Result{
				Added:   map[string]AddedResource{},
				Removed: map[string]RemovedResource{},
				Updated: map[string]UpdatedResource{
					"v1/ConfigMap/default/test": {
						Before: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name":      "test",
								"namespace": "default",
							},
							"data": map[string]interface{}{
								"key1": "value1",
								"key2": "value2",
							},
						},
						After: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name":      "test",
								"namespace": "default",
							},
							"data": map[string]interface{}{
								"key1": "value1",
								"key2": "new_value2",
								"key3": "value3",
							},
						},
						Diff: map[string]interface{}{
							"data": map[string]interface{}{
								"key2": map[string]interface{}{
									"__before": "value2",
									"__after":  "new_value2",
								},
								"key3": map[string]interface{}{
									"__added": "value3",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Generate(tt.before, tt.after)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check counts
			if len(result.Added) != len(tt.expected.Added) {
				t.Errorf("expected %d added resources, got %d", len(tt.expected.Added), len(result.Added))
			}
			if len(result.Removed) != len(tt.expected.Removed) {
				t.Errorf("expected %d removed resources, got %d", len(tt.expected.Removed), len(result.Removed))
			}
			if len(result.Updated) != len(tt.expected.Updated) {
				t.Errorf("expected %d updated resources, got %d", len(tt.expected.Updated), len(result.Updated))
			}

			// Check specific resources exist
			for key := range tt.expected.Added {
				if _, exists := result.Added[key]; !exists {
					t.Errorf("expected added resource %s not found", key)
				}
			}
			for key := range tt.expected.Removed {
				if _, exists := result.Removed[key]; !exists {
					t.Errorf("expected removed resource %s not found", key)
				}
			}
			for key := range tt.expected.Updated {
				if _, exists := result.Updated[key]; !exists {
					t.Errorf("expected updated resource %s not found", key)
				}
			}
		})
	}
}