## Chapter 9 - Discover Servers and Load Balance from the Client

Add functionality to allow clients to automatically:
- Discover servers in the cluster
- Direct `append` calls to leaders and `consume` calls to followers
- Balance `consume` calls across followers

Load-Balancing Strategies
1. Server proxying
2. External load balancing
3. Client-side load balancing

### Load Balancing on the Client in gRPC
- gRPC separates server discovery, load balancing, and client request and response handling
- `resolvers` discover servers, and `pickers` load balance by picking what server will handle the current request
- gRPC also has `balancers` that manage sub-connections but defer the load balancing to the pickers
- when `grpc.Dial` is called, gRPC takes the address and passes it to the `resolver` which discovers the servers
  - default `resolver` is the DNS `resolver`
  - for addresses with multiple DNS records, gRPC balances the requests across each of those records' servers
- gRPC uses round-robin load balancing by default - works well when each request requires the same work by the server
- what if load balancing needs to take in additional context on requests, servers and clients, e.g."
  - if server is a distributed service with single writer and multiple readers - read from replicas and write to primary
    - need to know if a request is a read or write, and if a server is a primary or replica
  - globally-distributed service: prioritize local servers needs to know location of clients and servers
  - latency-sensitive: need to track metrics on in-flight or queued requests, or other latency metrics
- write a custom `resolver` which discovers the servers and which server is the leader
- write a custom `picker` which manages directing `produce` calls to the leader, and balancing `consume` calls across followers

### Make Servers Discoverable
- Raft knows the cluster's servers, and which server is the leader
  - this information can be exposed to the `resolver` with an endpoint on the gRPC service

### Resolve the Servers
- `resolver` calls the `GetServers()` gRPC endpoint and passes the response to `picker` to decide where to route requests
- define a new type `Resolver` that implements gRPC `resolver.Builder` and `resolver.Resolver` interfaces
- gRPC's `resolver.Builder` interface comprise two methods:
  - `Build()` sets up a client connection so the resolver can call the `GetServers()` API
    - receives data needed to build a resolver that can discover the servers and the client connection that the resolver will update on
  - `Scheme()` returns the resolver's scheme identifier. gRPC parses it and tries to find a resolver that matches - defaults to DNS resolver
    - register this resolver with gRPC so that gRPC knows about this resolver when it's looking for resolvers that match the target's scheme
- gRPC's `resolver.Resolver` interface comprise two methods:
  - `ResolveNow()` is called by gRPC to resolve the target, discover the servers, and update the client connection with the servers.
    - Update the state with a slice of `resolver.Address` to inform the load balancer what servers it can choose from
  - `Close()` closes the resolver. Close the connection to the server created in `Build()`

### Route and Balance Requests with Pickers
- `pickers` handle the RPC balancing logic - they pick a server from those discovered by the `resolver` to handle each RPC
  - they can route RPCs based on information about the RPC, client, and server, so they can be used to implement any request-routing logic
- gRPC provides a base balancer that takes input from gRPC, manages sub-connections, and collects and aggregates connectivity states
- end-to-end testing:
  - client configures the `resolver` and `picker`
  - `resolver` discovers the servers
  - `picker` picks the sub-connections per RPC request
- documentation on [gRPC Name Resolution](https://github.com/grpc/grpc/blob/master/doc/naming.md)
