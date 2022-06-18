## Chapter 6 - Observe Your Systems

- Three types of telemetry data
  1. metrics
  2. structured logs
  3. traces

### Metrics
- measure numerical data over time
- can be used to define service-level indicators (SLI), objectives (SLO), and agreements (SLA)
- Three kinds of metrics:
  1. Counters - track the number of times an event happened
    - use counters to get rates - the number of times an event happened in an interval
  2. Histograms - distribution of data, percentiles
  3. Gauges - measure the current value of something
    - useful for saturation-type metrics (e.g. resources)
- Google 4 golden signals to measure
  1. Latency - time it takes to process requests -> informs decisions to scale vertically of horizontally
  2. Traffic - amount of demand on a service, e.g. requests processed per second, number of concurrent users
  3. Errors - request failure rate, esp internal server errors
  4. Saturation - a measure of a service's capacity, e.g. cpu, memory, disk

### Structured Logs
- Logs describe events in a system - log any event that gives useful insight into each service
- Logs help in troubleshooting, auditing, and profiling, in order to learn what went wrong and why
- In distributed systems, the request ID is helpful for piecing together a complete picture of a request that is handled by multiple services
- Structured logs
  - a set of name-value ordered-pairs encoded in a consistent schema and format that is easily read by programs
  - enable separation of log capturing, transport, persistence and querying
    - e.g. capture and transport logs as protocol buffers, re-encode into Parquet format, then persist in a columnar datastore

### Traces
- Traces capture request lifecycles and allow tracking requests as they flow through the system
  - some tracing tools provide a visual representation of where requests spend time in a system
  - especially useful in distributed systems as requests execute over multiple services
- Tag traces with details to know more about each request, e.g. user ID
- Traces comprise one or more spans - can have parent/child relationships or be linked as siblings
  - Each span represents a part of the request's execution
  - First, go wide to trace requests across all services end-to-end, with spans that begin and end at the entry and exit points of a system
  - Then, go deep in each service and trace important methods calls

### Making the Service Observable
- use OpenCensus for gRPC metrics and traces
- use Uber's Zap library for logging
- most Go networking APIs support middleware (interceptors)
  - wrap request handling with custom logic for metrics, logs and traces
```bash
$ go get -u go.opencensus.io
# go: downloading go.opencensus.io v0.23.0

$ go get -u go.uber.org/zap
# go: downloading go.uber.org/zap v1.21.0
```
- When tracing in production, consider writing a custom sampler that always traces important requests, and samples a percentage of the remaining requests
- configure gRPC service to
  - apply the Zap interceptors that log the gRPC calls
  - attach OpenCensus as the server's stat handler so that OpenCensus can record stats on the server's request handling

```bash
$ cd internal/server
$ go test -v -debug=true
# server_test.go:52: metrics log file: /var/folders/51/2g9hblfs0rv3x31nhzkp6wxr0000gn/T/metrics-730717193.log
# server_test.go:52: traces log file: /var/folders/51/2g9hblfs0rv3x31nhzkp6wxr0000gn/T/traces-174595540.log
# ...
```
