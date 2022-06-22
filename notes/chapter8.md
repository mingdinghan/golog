## Chapter 8 - Coordinated Services with Consensus

- Consensus algorithms are used to get distributed services to agree on shared state even in the face of failures.
  - Put the servers in leader and follower relationships where the followers replicate the leader's data.
- Use Raft for (1) leader election and (2) replication.

### Leader Election
- A Raft cluster has one leader and the rest of the servers are followers.
- The leader maintains leadership by sending heartbeat requests to its followers.
- If the follower times out while waiting for a heartbeat request from the leader, the follower
  - becomes a candidate and begins an election to decide the next leader
  - votes for itself and requests votes from the other followers.
- If the candidate receives a majority of the votes, it becomes the leader, and it sends heartbeat requests to the followers to establish authority.
- Candidates will hold elections until there's a winner that becomes the new leader.
- Every Raft server has a `term`: a monotonically increasing integer that tells other servers how authoritative and current this server is.
  - Each time a candidate begins an election, it increments its term.
  - If the candidate wins the election and becomes the leader, the followers update their terms to match and don't change until the next election.
- Servers vote once per term for the first candidate that requests votes, as long as the candidate's term is greater than the voters'.
  - These prevent vote splits and ensure the voters elect an up-to-date leader.

### Log Replication
- Replication saves us from losing data when servers fail.
- The leader accepts client requests, each of which represents some command to run across the cluster.
  - For each request, the leader appends the command to its log and then requests its followers to append the command to their logs.
  - After a majority of followers have replicated the command, the leader considers the command committed,
  - executes the command with a finite-state machine and responds to the client with the result.
- The leader tracks the highest committed offset and sends this in the requests to its followers.
- When a follower receives a request, it executes all commands up to the highest committed offset with its finite-state machine.

### Implementing Raft
- Use Raft as a means to replicate a log of commands and then execute those commands with state machines:
  - for a distributed SQL database, replicate and execute the `insert` and `update` SQL commands.
  - for a distributed key-value store, replicate and execute the `set` commands.
  - for Golog, we will replicate the transformation commands, i.e. the `append` commands.
```bash
$ go get github.com/hashicorp/raft@v1.3.9
$ go get github.com/hashicorp/raft-boltdb/v2
```
- A Raft instance comprises:
  - A `finite-state machine` that applies the commands given to Raft.
  - A `log store` where Raft stores those commands.
  - A `stable store` where Raft stores the cluster's configuration — e.g. the servers in the cluster, and their addresses.
  - A `snapshot store` where Raft stores compact snapshots of its data.
  - A `transport` that Raft uses to connect with the server's peers.
- [Bolt](https://github.com/etcd-io/bbolt) is an embedded and persisted key-value database for Go that is used as the stable store.
- Bootstrap a server configured with itself as the only voter, wait until it becomes the leader,
  - then tell the leader to add more servers to the cluster.
- Add public APIs for `DistributedLog` that append records to and read records from the log, and wrap Raft
  - for relaxed consistency, read operations need not go through Raft.
  - if strong consistency is required (i.e. reads must be up-to-date with writes), read operations must go through Raft,
    - so reads will be less efficient and take longer.
- Raft defers the running of your business logic to the Finite-State Machine (FSM)
  - The FSM must access the data it manages. In Golog service, that's a log, and the FSM appends records to the log.
- To implement the `raft.FSM` interface, implement these methods:
  - `Apply(record *raft.Log)`: Raft invokes this method after committing a log entry
  - `Snapshot()`: Raft periodically invokes this method to snapshot its state
  - `Restore(io.ReadCloser)`: Raft invokes this to restore an FSM from a snapshot
- To implement the `raft.FSMSnapshot` interface, implement these methods:
  - `Persist(sink raft.SnapshotSink)`: Raft invokes this method to write its state to a configured sink
  - `Release()`: Raft invokes this method when it is finished with the snapshot
- To implement the `raft.LogStore` interface, implement these methods:
  - `FirstIndex()`
  - `LastIndex()`
  - `GetLog(index uint64, out *raft.Log)`
- Raft uses a stream layer in the `transport` to provide a low-level stream abstraction to connect with Raft servers
  - Implement a stream layer that satisfies Raft's `StreamLayer` interface, i.e. implement these methods:
    - `Accept()`
    - `Close()`
    - `Addr()`
  - `Dial(addr raft.ServerAddress, timeout time.Duration)` makes outgoing connections to other servers in the Raft cluster
    - define a `RaftRPC` byte to identify the connection type, in order to multiplex Raft on the same port as log gRPC requests
