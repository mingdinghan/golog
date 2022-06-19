## Golog

Golog is a distributed commit log built using Go

### Key Concepts and Learning Points
- Serialization and deserialization using Protocol Buffers
- Build a Log of binary Records, stored and indexed as Segments
- Exposing a library as a Service with gRPC
- Securing Services using mutual TLS and access control lists
- Adding Observability to Services: Metrics, Structured Logs and Traces
- Server-to-Server Discovery using [Serf](https://www.serf.io/intro/index.html)
  - a tool for cluster membership, failure detection, and orchestration
  - decentralized, fault-tolerant and highly available
  - uses an efficient `gossip` protocol

### References

- [Distributed Services with Go](https://pragprog.com/titles/tjgo/distributed-services-with-go) by [Travis Jeffery](https://twitter.com/travisjeffery)
