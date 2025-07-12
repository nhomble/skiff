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

From [releases](https://github.com/nhomble/skiff/releases)

## Example

```
{
  "resource_changes": {
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
            "name": "app-config",
            "namespace": "default"
          }
        },
        "after": {
          "apiVersion": "v1",
          "data": {
            "database_url": "postgres://new-db:5432/app",
            "log_level": "debug"
          },
          "kind": "ConfigMap",
          "metadata": {
            "name": "app-config",
            "namespace": "default"
          }
        },
        "changes": {
          "data.database_url": {
            "from": "postgres://old-db:5432/app",
            "to": "postgres://new-db:5432/app"
          },
          "data.log_level": {
            "from": "info",
            "to": "debug"
          }
        }
      }
    }
  }
}
```