# Opentelemetry Agent

This project aims to run opentelemtry collector using one script and provide an api to update the config of the collector and spin up the collector instance. This project leverages the opentelmetry collector builder project to create a custom instace of the opentelmetry collector which can have modules from core or contrib collector repo. Visit opentelmetry collector builder to know how to add a custom exporter, reciecer or processor.

## Steps to run the agent

1) Run below command in bash terminal
```bash
sh install.sh
```

2) Create a file with name `otelcol.yaml` and paste below content in it

```yaml
extensions:
  zpages:
    endpoint: localhost:55679

receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
      http:
        endpoint: localhost:4318

processors:
  batch:
  memory_limiter:
    # 75% of maximum memory up to 2G
    limit_mib: 1536
    # 25% of limit up to 2G
    spike_limit_mib: 512
    check_interval: 5s

exporters:
  debug:
    verbosity: detailed

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [debug]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [debug]

  extensions: [zpages]
```

3) Now call below url and add the above file as binary in the payload like below with appropriate file location:

```bash
curl --location 'localhost:4343/collector/config' \
--header 'Content-Type: text/yaml' \
--data-binary '<replace-path>/otelcol.yaml'
```

The opentelemetry collector is setup and is ready to receive telemetry data.

The above config can be configured as per convience considering core collector in mind as the builder config in this repo consists of only core collector components but they can be extended to use contri collector components as well.