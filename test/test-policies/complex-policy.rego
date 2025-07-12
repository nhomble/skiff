package main

import future.keywords.contains
import future.keywords.if

# Complex policies demonstrating flat format capabilities

# Deny: HPA maxReplicas > 12
deny contains msg if {
    change := input.resource_changes[_]
    change.type == "HorizontalPodAutoscaler"
    change.change.actions[_] in ["create", "update"]
    
    maxReplicas := change.change.after.spec.maxReplicas
    maxReplicas > 12
    
    msg := sprintf("HPA '%s' in namespace '%s' has maxReplicas %d, which exceeds the limit of 12", 
                   [change.name, change.namespace, maxReplicas])
}

# Warn: HPA behavior changes
warn contains msg if {
    change := input.resource_changes[_]
    change.type == "HorizontalPodAutoscaler"
    change.change.actions[_] == "update"
    
    # Check if any behavior fields were added using changes
    behavior_changes := [path |
        some path, field_change in change.change.changes
        startswith(path, "spec.behavior.")
        field_change.from == null  # Field was added
    ]
    count(behavior_changes) > 0
    
    msg := sprintf("HPA '%s' in namespace '%s' has %d new scaling behavior field(s)", 
                   [change.name, change.namespace, count(behavior_changes)])
}

# Warn: Service LoadBalancer source ranges added
warn contains msg if {
    change := input.resource_changes[_]
    change.type == "Service" 
    change.change.actions[_] == "update"
    change.change.after.spec.type == "LoadBalancer"
    
    # Check if loadBalancerSourceRanges was added using changes field
    ranges_change := change.change.changes["spec.loadBalancerSourceRanges"]
    ranges_change.from == null  # Field was added
    
    msg := sprintf("Service '%s' in namespace '%s' now has loadBalancerSourceRanges restrictions", 
                   [change.name, change.namespace])
}

# Deny: Service with more than 5 ports
deny contains msg if {
    change := input.resource_changes[_]
    change.type == "Service"
    change.change.actions[_] in ["create", "update"]
    
    port_count := count(change.change.after.spec.ports)
    port_count > 5
    
    msg := sprintf("Service '%s' in namespace '%s' has %d ports, which exceeds the limit of 5", 
                   [change.name, change.namespace, port_count])
}