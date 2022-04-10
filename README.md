# Otel demo with golang application

This is a distributed tracing demonstration. The demonstration consists of a golang application
being instrumented, and able to send the send the trace to the otel collector, then the collector will
export the trace to multiple third party vendors.

In this demo the otel collector is configured to send the trace to both **datadog** and **jaeger**

- First `cd` into the `otlp/demo/` folder and run `docker-compose up` to run the
  nessesary image for this demo.
- After running otel collector with `docker-compose up`, then go back to the root directory
  and run `go run main.go`
