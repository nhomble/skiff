package k8s

import (
	"strings"
	"testing"
)

func TestParseYAMLStream(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected int
		wantErr  bool
	}{
		{
			name: "single document",
			yaml: `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value`,
			expected: 1,
			wantErr:  false,
		},
		{
			name: "multi document",
			yaml: `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
  namespace: default
spec:
  replicas: 1`,
			expected: 2,
			wantErr:  false,
		},
		{
			name: "empty document",
			yaml: `---
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value`,
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "completely empty",
			yaml:     "",
			expected: 0,
			wantErr:  false,
		},
		{
			name: "invalid yaml",
			yaml: `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
  invalid: [unclosed`,
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.yaml)
			objects, err := ParseYAMLStream(reader)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(objects) != tt.expected {
				t.Errorf("expected %d objects, got %d", tt.expected, len(objects))
			}
		})
	}
}

func TestGenerateObjectKey(t *testing.T) {
	tests := []struct {
		name     string
		object   map[string]interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "namespaced resource",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name":      "test-config",
					"namespace": "default",
				},
			},
			expected: "v1/ConfigMap/default/test-config",
			wantErr:  false,
		},
		{
			name: "resource without namespace",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "test-config",
				},
			},
			expected: "v1/ConfigMap/default/test-config",
			wantErr:  false,
		},
		{
			name: "missing apiVersion",
			object: map[string]interface{}{
				"kind": "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "test-config",
				},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "missing kind",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "test-config",
				},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "missing metadata",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "missing name",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"namespace": "default",
				},
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateObjectKey(tt.object)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if key != tt.expected {
				t.Errorf("expected key %q, got %q", tt.expected, key)
			}
		})
	}
}
