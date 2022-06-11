## Chapter 4 - Serve Requests with gRPC

- gRPC is an open source high-performance RPC (remote procedure call) framework
- it is built with protobuf and HTTP/2
  - protobuf performs very well at serialization
  - HTTP/2 provides a means for long-lasting connections
- supports various kinds of load-balancing:
  - thick client-side load balancing
  - proxy load balancing
  - look-aside load balancing
  - service mesh
- A gRPC service consists of a group of related RPC endpoints
  - e.g. enable clients to write to and read from their log
- Creating a gRPC service involves
  - protobuf definitions: service endpoints, request and response types
  - compiling the protocol buffers into code comprising the client and server stubs to be implemented
    - tell the protobuf compiler to use the gRPC plugin
      ```bash
      $ go get google.golang.org/grpc@v1.47.0
      $ go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
      ```
- gRPC `requests` and `responses` are messages that the compiler turns into Go structs
- To implement a gRPC service, build a struct whose methods match the service definition in the protobuf
  - implement the `Produce` and `Consume` handlers and streaming APIs
- Add a custom error `ErrOffsetOutOfRange` that the server will send back to the client when the client tries to consume an offset that's out of the log's range
- Dependency Inversion with Interfaces
  - design a service that depends on a log interface, rather than a specific log implementation
    - in production env, the service may need to persist user data
    - in testing env, an in-memory log implementation may suffice
  - the service can use any given log implementation that satisfies the `CommitLog` interface
