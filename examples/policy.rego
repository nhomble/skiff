package main

import future.keywords.contains
import future.keywords.if

# Policy: Deployments should not have more than 1 replica
# This demonstrates how to write policies against yspec output

# Get all proposed resources (added and updated)
proposed_resources contains resource if {
    some key
    resource := input.added[key].after
}

proposed_resources contains resource if {
    some key
    resource := input.updated[key].after
}

# Deny: Deployment replicas > 1
deny contains msg if {
    resource := proposed_resources[_]
    resource.kind == "Deployment"
    startswith(resource.apiVersion, "apps/")
    
    replicas := resource.spec.replicas
    replicas > 1
    
    msg := sprintf("Deployment '%s' in namespace '%s' has %d replicas, which exceeds the limit of 1", 
                   [resource.metadata.name, resource.metadata.namespace, replicas])
}

# Deny: Deployment replica count increased
deny contains msg if {
    some key
    updated := input.updated[key]
    resource := updated.after
    
    resource.kind == "Deployment"
    startswith(resource.apiVersion, "apps/")
    
    # Check if replicas were increased
    before_replicas := updated.before.spec.replicas
    after_replicas := updated.after.spec.replicas
    after_replicas > before_replicas
    
    msg := sprintf("Deployment '%s' in namespace '%s' replica count increased from %d to %d", 
                   [resource.metadata.name, resource.metadata.namespace, before_replicas, after_replicas])
}

# Warn: ConfigMap changes
warn contains msg if {
    some key
    updated := input.updated[key]
    resource := updated.after
    
    resource.kind == "ConfigMap"
    
    msg := sprintf("ConfigMap '%s' in namespace '%s' has been modified", 
                   [resource.metadata.name, resource.metadata.namespace])
}