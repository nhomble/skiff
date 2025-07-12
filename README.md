# skiff

structured kubernetes diff

## Usage

```sh
skiff test/test-cases/before.yaml test/test-cases/after.yaml
```

or with image

```sh
docker run --rm -v $(pwd):/app -w /app skiff:latest test/test-cases/before.yaml test/test-cases/after.yaml
```

## Download

From [docker hub](https://hub.docker.com/r/hombro/skiff)

## Example

```
{
  "resource_changes": {
    "apps/v1/Deployment/default/web-app": {
      "type": "Deployment",
      "apiVersion": "apps/v1",
      "namespace": "default",
      "name": "web-app",
      "change": {
        "actions": [
          "update"
        ],
        "before": {
          "apiVersion": "apps/v1",
          "kind": "Deployment",
          "metadata": {
            "name": "web-app",
            "namespace": "default",
            "resourceVersion": "67890",
            "uid": "def-456"
          },
          "spec": {
            "replicas": 2,
            "selector": {
              "matchLabels": {
                "app": "web-app"
              }
            },
            "template": {
              "metadata": {
                "labels": {
                  "app": "web-app"
                }
              },
              "spec": {
                "containers": [
                  {
                    "env": [
                      {
                        "name": "ENV",
                        "value": "production"
                      }
                    ],
                    "image": "nginx:1.20",
                    "name": "web",
                    "ports": [
                      {
                        "containerPort": 80
                      }
                    ]
                  }
                ]
              }
            }
          },
          "status": {
            "availableReplicas": 2,
            "readyReplicas": 2
          }
        },
        "after": {
          "apiVersion": "apps/v1",
          "kind": "Deployment",
          "metadata": {
            "name": "web-app",
            "namespace": "default",
            "resourceVersion": "67891",
            "uid": "def-456"
          },
          "spec": {
            "replicas": 3,
            "selector": {
              "matchLabels": {
                "app": "web-app"
              }
            },
            "template": {
              "metadata": {
                "labels": {
                  "app": "web-app"
                }
              },
              "spec": {
                "containers": [
                  {
                    "env": [
                      {
                        "name": "ENV",
                        "value": "production"
                      },
                      {
                        "name": "DEBUG",
                        "value": "true"
                      }
                    ],
                    "image": "nginx:1.21",
                    "name": "web",
                    "ports": [
                      {
                        "containerPort": 80
                      }
                    ]
                  }
                ]
              }
            }
          },
          "status": {
            "availableReplicas": 3,
            "readyReplicas": 3
          }
        },
        "changes": {
          "metadata.resourceVersion": {
            "from": "67890",
            "to": "67891"
          },
          "spec.replicas": {
            "from": 2,
            "to": 3
          },
          "spec.template.spec.containers[0].env[1].name": {
            "from": null,
            "to": "DEBUG"
          },
          "spec.template.spec.containers[0].env[1].value": {
            "from": null,
            "to": "true"
          },
          "spec.template.spec.containers[0].image": {
            "from": "nginx:1.20",
            "to": "nginx:1.21"
          },
          "status.availableReplicas": {
            "from": 2,
            "to": 3
          },
          "status.readyReplicas": {
            "from": 2,
            "to": 3
          }
        }
      }
    },
    "v1/ConfigMap/default/app-config": {
      "type": "ConfigMap",
      "apiVersion": "v1",
      "namespace": "default",
      "name": "app-config",
      "change": {
        "actions": [
          "update"
        ],
        "before": {
          "apiVersion": "v1",
          "data": {
            "database_url": "postgres://old-db:5432/app",
            "log_level": "info"
          },
          "kind": "ConfigMap",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"ConfigMap\"...}\n"
            },
            "name": "app-config",
            "namespace": "default",
            "resourceVersion": "12345",
            "uid": "abc-123"
          }
        },
        "after": {
          "apiVersion": "v1",
          "data": {
            "cache_enabled": "true",
            "database_url": "postgres://new-db:5432/app",
            "log_level": "debug"
          },
          "kind": "ConfigMap",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"ConfigMap\"...}\n"
            },
            "name": "app-config",
            "namespace": "default",
            "resourceVersion": "12346",
            "uid": "abc-123"
          }
        },
        "changes": {
          "data.cache_enabled": {
            "from": null,
            "to": "true"
          },
          "data.database_url": {
            "from": "postgres://old-db:5432/app",
            "to": "postgres://new-db:5432/app"
          },
          "data.log_level": {
            "from": "info",
            "to": "debug"
          },
          "metadata.resourceVersion": {
            "from": "12345",
            "to": "12346"
          }
        }
      }
    },
    "v1/Service/default/api-service": {
      "type": "Service",
      "apiVersion": "v1",
      "namespace": "default",
      "name": "api-service",
      "change": {
        "actions": [
          "create"
        ],
        "after": {
          "apiVersion": "v1",
          "kind": "Service",
          "metadata": {
            "name": "api-service",
            "namespace": "default"
          },
          "spec": {
            "ports": [
              {
                "port": 8080,
                "targetPort": 8080
              }
            ],
            "selector": {
              "app": "api"
            },
            "type": "ClusterIP"
          }
        }
      }
    },
    "v1/Service/default/web-service": {
      "type": "Service",
      "apiVersion": "v1",
      "namespace": "default",
      "name": "web-service",
      "change": {
        "actions": [
          "delete"
        ],
        "before": {
          "apiVersion": "v1",
          "kind": "Service",
          "metadata": {
            "name": "web-service",
            "namespace": "default"
          },
          "spec": {
            "ports": [
              {
                "port": 80,
                "targetPort": 80
              }
            ],
            "selector": {
              "app": "web-app"
            },
            "type": "ClusterIP"
          }
        }
      }
    }
  }
}
```