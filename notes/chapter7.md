## Chapter 7 - Server-to-Server Service Discovery

- Automatically handle when a node is added or removed from the cluster - how to connect to a service
- Keep an up-to-date registry of services, their locations and their health
- Downstream services query this registry to discover the location of upstream services, and connect to them
- Using a load balancer can add cost, increase latency, introduce single points of failure, and require updates as services scale up and down
  - for server-to-server communications, choose service discovery to avoid managing load balancers and DNS records
- Key issues in service discovery
  - how will the servers in the cluster discover each other?
  - how will the clients discover the servers?

### Embedding Service Discovery
- Requirements for a service discovery tool
  - Manage a registry of services containing info such as IPs and ports
  - Help services find other services using the registry
  - Perform health checks on service instances and remove them if they are not healthy
  - Deregister services when they go offline
- Historically, people who built distributed services depended on separate, stand-alone services (clusters) for service discovery
- Serf is a Go library that provides decentralised cluster membership, failure detection and orchestration
  - can be used to embed service discovery functionalities into distributed services
  - created by Hashicorp which uses Serf to power its own service-discovery product, Consul

### Discover Services with Serf
- Serf maintains cluster membership by using an efficient, lightweight gossip protocol to communicate between nodes
  - In contrast, service registry projects like Zookeeper and Consul use a central-registry approach
- To implement service discovery with Serf:
  1. Create a Serf node on each server
  2. Configure each Serf node with an address to listen on and accept connections from other Serf nodes
  3. Configure each Serf node with addresses of other Serf nodes and join their cluster
  4. Handle Serf's cluster discovery events, e.g. when a node joins or fails in the cluster
```bash
$ go get github.com/hashicorp/serf@v0.9.8
```
- Serf Configurable Parameters
  - `NodeName`: a node's unique identifier across the Serf cluster
  - `BindAddr` and `BindPort`: Serf listens on these for gossip protocol
  - `Tags`: shared with other nodes and can be used for cluster management
  - `EventCh`: a channel to receive Serf events when a node joins or leaves the cluster
  - `StartJoinAddrs`: for configuring a new node to connect to one of the existing nodes in the cluster
- Golog service is designed to replicate the data of servers that join the cluster
  - When consensus is added later, Raft will need to know when servers join the cluster, in order to coordinate with them
