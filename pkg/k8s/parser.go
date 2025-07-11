package k8s

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// ParseYAMLStream parses a multi-document YAML stream and returns a map of K8s objects
// keyed by their unique identifier (apiVersion/kind/namespace/name)
func ParseYAMLStream(reader io.Reader) (map[string]map[string]interface{}, error) {
	objects := make(map[string]map[string]interface{})
	decoder := yaml.NewDecoder(reader)

	for {
		var obj map[string]interface{}
		err := decoder.Decode(&obj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to decode YAML document: %w", err)
		}

		if len(obj) == 0 {
			continue
		}

		key, err := GenerateObjectKey(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to generate object key: %w", err)
		}

		objects[key] = obj
	}

	return objects, nil
}


// GenerateObjectKey creates a unique identifier for a K8s object
// Format: apiVersion/kind/namespace/name (uses "default" when namespace not specified)
func GenerateObjectKey(obj map[string]interface{}) (string, error) {
	apiVersion, ok := obj["apiVersion"].(string)
	if !ok || apiVersion == "" {
		return "", fmt.Errorf("missing or invalid apiVersion")
	}

	kind, ok := obj["kind"].(string)
	if !ok || kind == "" {
		return "", fmt.Errorf("missing or invalid kind")
	}

	metadata, ok := obj["metadata"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}

	name, ok := metadata["name"].(string)
	if !ok || name == "" {
		return "", fmt.Errorf("missing or invalid metadata.name")
	}

	namespace := "default"
	if ns, ok := metadata["namespace"].(string); ok && ns != "" {
		namespace = ns
	}

	return fmt.Sprintf("%s/%s/%s/%s", apiVersion, kind, namespace, name), nil
}

