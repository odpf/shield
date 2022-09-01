# Configurations

Shield can be configured with config.yaml file. An example of such is:

```yaml
version: 1

# logging configuration
log:
  # debug, info, warning, error, fatal - default 'info'
  level: debug

app:
  port: 8000
  identity_proxy_header: X-Shield-Email
  # full path prefixed with scheme where resources config yaml files are kept
  # e.g.:
  # local storage file "file:///tmp/resources_config"
  # GCS Bucket "gs://shield/resources_config"
  resources_config_path: file:///tmp/resources_config\
  # secret required to access resources config
  # e.g.:
  # system environment variable "env://TEST_RULESET_SECRET"
  # local file "file:///opt/auth.json"
  # secret string "val://user:password"
  # optional
  resources_config_path_secret: env://TEST_RESOURCE_CONFIG_SECRET

db:
  driver: postgres
  url: postgres://shield:@localhost:5432/shield?sslmode=disable
  max_query_timeout: 500ms

spicedb:
  host: spicedb.localhost
  pre_shared_key: randomkey
  port: 50051

# proxy configuration
proxy:
  services:
    - name: test
      host: 0.0.0.0
      # port where the proxy will be listening on for requests
      port: 5556
      # full path prefixed with scheme where ruleset yaml files are kept
      # e.g.:
      # local storage file "file:///tmp/rules"
      # GCS Bucket "gs://shield/rules"
      ruleset: file:///tmp/rules
      # secret required to access ruleset
      # e.g.:
      # system environment variable "env://TEST_RULESET_SECRET"
      # local file "file:///opt/auth.json"
      # secret string "val://user:password"
      # optional
      ruleset_secret: env://TEST_RULESET_SECRET
```
