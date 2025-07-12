package main

import future.keywords.contains
import future.keywords.if

# Policy: Deployments should not have more than 1 replica
# This demonstrates how to write policies against skiff flat diff output

# Deny: Deployment replicas > 1
deny contains msg if {
    change := input.resource_changes[_]
    change.type == "Deployment"
    startswith(change.apiVersion, "apps/")
    
    # Check creates and updates
    change.change.actions[_] in ["create", "update"]
    
    replicas := change.change.after.spec.replicas
    replicas > 1
    
    msg := sprintf("Deployment '%s' in namespace '%s' has %d replicas, which exceeds the limit of 1", 
                   [change.name, change.namespace, replicas])
}

# Deny: Deployment replica count increased
deny contains msg if {
    change := input.resource_changes[_]
    change.type == "Deployment"
    startswith(change.apiVersion, "apps/")
    change.change.actions[_] == "update"
    
    # Check if replicas were increased using changes field
    replica_change := change.change.changes["spec.replicas"]
    replica_change.to > replica_change.from
    
    msg := sprintf("Deployment '%s' in namespace '%s' replica count increased from %d to %d", 
                   [change.name, change.namespace, replica_change.from, replica_change.to])
}

# Warn: ConfigMap changes
warn contains msg if {
    change := input.resource_changes[_]
    change.type == "ConfigMap"
    change.change.actions[_] == "update"
    
    # Count how many data fields changed
    changed_fields := [path | 
        some path, _ in change.change.changes
        startswith(path, "data.")
    ]
    
    msg := sprintf("ConfigMap '%s' in namespace '%s' has %d data field(s) modified", 
                   [change.name, change.namespace, count(changed_fields)])
}