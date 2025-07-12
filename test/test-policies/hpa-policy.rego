package main

import future.keywords.contains
import future.keywords.if

# Deny: HPA scale-up percentage > 5%
deny contains msg if {
    change := input.resource_changes[_]
    change.type == "HorizontalPodAutoscaler"
    change.change.actions[_] == "update"
    
    # Find any scaleUp policy value changes that exceed 5%
    some path, field_change in change.change.changes
    regex.match(`spec\.behavior\.scaleUp\.policies\[\d+\]\.value`, path)
    
    # Get the corresponding type path to check if it's a percentage
    type_path := regex.replace(path, `\.value$`, ".type")
    type_change := change.change.changes[type_path]
    type_change.to == "Percent"
    
    field_change.to > 5
    
    msg := sprintf("HPA '%s' has scale-up percentage %d%% which exceeds limit of 5%% (field: %s)", 
                   [change.name, field_change.to, path])
}

# Deny: HPA max replicas increased too much
deny contains msg if {
    change := input.resource_changes[_]
    change.type == "HorizontalPodAutoscaler"
    change.change.actions[_] == "update"
    
    # Check if maxReplicas changed and increased by more than 50%
    replica_change := change.change.changes["spec.maxReplicas"]
    replica_change.to > replica_change.from * 1.5  # More than 50% increase
    
    msg := sprintf("HPA '%s' maxReplicas increased from %d to %d (more than 50%% increase)", 
                   [change.name, replica_change.from, replica_change.to])
}

# Warn: HPA behavior changes
warn contains msg if {
    change := input.resource_changes[_]
    change.type == "HorizontalPodAutoscaler"
    change.change.actions[_] == "update"
    
    # Check if any behavior fields were added/modified
    behavior_changes := [path |
        some path, _ in change.change.changes
        startswith(path, "spec.behavior.")
    ]
    count(behavior_changes) > 0
    
    msg := sprintf("HPA '%s' has %d behavior field(s) modified - review carefully: %v", 
                   [change.name, count(behavior_changes), behavior_changes])
}