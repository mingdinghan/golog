package log_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/raft"
	api "github.com/mingdinghan/golog/api/v1"
	"github.com/mingdinghan/golog/internal/log"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
)

// Set up a three-server cluster and shorten the default Raft timeout configs so that Raft elects the leader quickly
func TestMultipleNodes(t *testing.T) {
	var logs []*log.DistributedLog
	nodeCount := 3
	ports := dynaport.Get(nodeCount)

	for i := 0; i < nodeCount; i++ {
		dataDir, err := ioutil.TempDir("", "distributed-log-test")
		require.NoError(t, err)
		defer func(dir string) {
			_ = os.RemoveAll(dir)
		}(dataDir)

		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", ports[i]))
		require.NoError(t, err)

		config := log.Config{}
		config.Raft.StreamLayer = log.NewStreamLayer(ln, nil, nil)
		config.Raft.LocalID = raft.ServerID(fmt.Sprintf("%d", i))
		config.Raft.HeartbeatTimeout = 50 * time.Millisecond
		config.Raft.ElectionTimeout = 50 * time.Millisecond
		config.Raft.LeaderLeaseTimeout = 50 * time.Millisecond
		config.Raft.CommitTimeout = 5 * time.Millisecond

		if i == 0 {
			config.Raft.Bootstrap = true
		}

		l, err := log.NewDistributedLog(dataDir, config)
		require.NoError(t, err)

		// The first server bootstraps the cluster, becomes the leader, and adds the other two servers to the cluster
		// The leader then joins other servers to its cluster
		if i != 0 {
			err = logs[0].Join(fmt.Sprintf("%d", i), ln.Addr().String())
			require.NoError(t, err)
		} else {
			err = l.WaitForLeader(3 * time.Second)
			require.NoError(t, err)
		}
		logs = append(logs, l)
		time.Sleep(500 * time.Millisecond)
	}

	records := []*api.Record{
		{Value: []byte("first")},
		{Value: []byte("second")},
	}

	// test replication by appending records to leader server and check that Raft replicated the records to followers
	for _, record := range records {
		off, err := logs[0].Append(record)
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			for j := 0; j < nodeCount; j++ {
				got, err := logs[j].Read(off)
				if err != nil {
					return false
				}
				record.Offset = off
				if !reflect.DeepEqual(got.Value, record.Value) {
					return false
				}
			}
			return true
		}, 500*time.Millisecond, 50*time.Millisecond)
	}

	// checks that leader stops replicating to a server that's left the cluster,
	// while continuing to replicate to the existing servers.
	err := logs[0].Leave("1")
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	off, err := logs[0].Append(&api.Record{
		Value: []byte("third"),
	})
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	record, err := logs[1].Read(off)
	require.IsType(t, api.ErrOffsetOutOfRange{}, err)
	require.Nil(t, record)

	record, err = logs[2].Read(off)
	require.NoError(t, err)
	require.Equal(t, []byte("third"), record.Value)
	require.Equal(t, off, record.Offset)
}
