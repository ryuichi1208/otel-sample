# otel-sample

## sample

``` yaml
extensions:
  health_check:
  pprof:
    endpoint: 0.0.0.0:1777
  zpages:
    endpoint: 0.0.0.0:55679

receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
  resource:
    attributes:
    - key: auth-token
      action: delete

exporters:
  file:
    path: /tmp/otelcol_data.log
    logLevel: debug

service:
  extensions: [health_check, pprof, zpages]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging]
```
