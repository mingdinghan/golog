package loadbalance_test

import (
	"net"
	"net/url"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	api "github.com/mingdinghan/golog/api/v1"
	"github.com/mingdinghan/golog/internal/config"
	"github.com/mingdinghan/golog/internal/loadbalance"
	"github.com/mingdinghan/golog/internal/server"
	"github.com/stretchr/testify/require"
)

// Set up a server for test resolver to try and discover some servers from
func TestResolver(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	tlsConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	serverCreds := credentials.NewTLS(tlsConfig)

	srv, err := server.NewGRPCServer(&server.Config{
		ServerGetter: &getServers{},
	}, grpc.Creds(serverCreds))
	require.NoError(t, err)

	go srv.Serve(l)

	conn := &clientConn{}
	tlsConfig, err = config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	clientCreds := credentials.NewTLS(tlsConfig)
	opts := resolver.BuildOptions{DialCreds: clientCreds}

	// create and build the test resolver, and configure its target endpoint to point to the server set up earlier
	r := &loadbalance.Resolver{}
	_, err = r.Build(
		resolver.Target{
			URL: url.URL{
				Path: l.Addr().String(),
			},
		},
		conn,
		opts,
	)
	require.NoError(t, err)

	// check that the resolver updated the client connection with the servers and data expected
	// i.e. two servers, with the port:9001 server marked as the leader
	wantState := resolver.State{
		Addresses: []resolver.Address{{
			Addr:       "localhost:9001",
			Attributes: attributes.New("is_leader", true),
		}, {
			Addr:       "localhost:9002",
			Attributes: attributes.New("is_leader", false),
		}},
	}
	require.Equal(t, wantState, conn.state)

	conn.state.Addresses = nil
	r.ResolveNow(resolver.ResolveNowOptions{})
	require.Equal(t, wantState, conn.state)
}

// getServers (mock) implements `ServerGetter` interface and returns a known server set for the resolver to find
type getServers struct{}

func (s *getServers) GetServers() ([]*api.Server, error) {
	return []*api.Server{{
		Id:       "leader",
		RpcAddr:  "localhost:9001",
		IsLeader: true,
	}, {
		Id:      "follower",
		RpcAddr: "localhost:9002",
	}}, nil
}

// clientConn (mock) implements `resolver.ClientConn` interface and keeps a reference to the state that the resolver updated with,
// to verify the resolver updated the client connection with the correct data
type clientConn struct {
	resolver.ClientConn
	state resolver.State
}

func (c *clientConn) UpdateState(state resolver.State) error {
	c.state = state
	return nil
}

func (c *clientConn) ReportError(err error) {}

func (c *clientConn) NewAddress(addrs []resolver.Address) {}

func (c *clientConn) NewServiceConfig(config string) {}

func (c *clientConn) ParseServiceConfig(config string) *serviceconfig.ParseResult {
	return nil
}
