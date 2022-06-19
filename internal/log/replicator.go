package log

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	api "github.com/mingdinghan/golog/api/v1"
)

type Replicator struct {
	DialOptions []grpc.DialOption
	LocalServer api.LogClient

	logger *zap.Logger

	mu      sync.Mutex
	servers map[string]chan struct{}
	closed  bool
	close   chan struct{}
}

// Join adds the discovered server's address to the list of servers to replicate
// and kicks off the goroutine to run the actual replication logic
func (r *Replicator) Join(name, addr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.init()

	if r.closed {
		return nil
	}

	if _, ok := r.servers[name]; ok {
		// server exists in the map, so replicator should already be replicating from it. Skip it.
		return nil
	}

	r.servers[name] = make(chan struct{})

	go r.replicate(addr, r.servers[name])

	return nil
}

// replicate connects to the discovered server and replicates messages from it until that server fails or leaves the cluster
func (r *Replicator) replicate(addr string, leave chan struct{}) {
	cc, err := grpc.Dial(addr, r.DialOptions...)
	if err != nil {
		r.logError(err, "failed to dial", addr)
		return
	}
	defer cc.Close()

	// create a gRPC client and open up a stream to consume all logs from the discovered server
	client := api.NewLogClient(cc)

	ctx := context.Background()
	stream, err := client.ConsumeStream(ctx, &api.ConsumeRequest{
		Offset: 0,
	})
	if err != nil {
		r.logError(err, "failed to consume", addr)
		return
	}

	records := make(chan *api.Record)
	go func() {
		for {
			// consume logs from the discovered server in a stream
			recv, err := stream.Recv()
			if err != nil {
				r.logError(err, "failed to receive", addr)
				return
			}
			records <- recv.Record
		}
	}()

	for {
		select {
		case <-r.close:
			return
		case <-leave:
			return
		case record := <-records:
			// produce logs to the local server to save a copy
			_, err = r.LocalServer.Produce(ctx, &api.ProduceRequest{
				Record: record,
			})
			if err != nil {
				r.logError(err, "failed to produce", addr)
				return
			}
		}
	}
}

// Leave removes the leaving server from the list of servers to replicate and closes the server's associated channel
func (r *Replicator) Leave(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.init()

	if _, ok := r.servers[name]; !ok {
		// server is already not in the map, so simply return
		return nil
	}

	close(r.servers[name])
	delete(r.servers, name)
	return nil
}

// init lazily initializes the server map
func (r *Replicator) init() {
	if r.logger == nil {
		r.logger = zap.L().Named("replicator")
	}
	if r.servers == nil {
		r.servers = make(map[string]chan struct{})
	}
	if r.close == nil {
		r.close = make(chan struct{})
	}
}

// Close closes the replicator so that
// it doesn't replicate new servers that join the cluster, and
// it stops replicating existing servers by causing the `replicate()` goroutines to return
func (r *Replicator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.init()

	if r.closed {
		return nil
	}

	r.closed = true
	close(r.close)
	return nil
}

func (r *Replicator) logError(err error, msg, addr string) {
	r.logger.Error(
		msg,
		zap.String("addr", addr),
		zap.Error(err),
	)
}
