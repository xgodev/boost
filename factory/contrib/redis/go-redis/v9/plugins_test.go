package redis

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeRedis speaks just enough RESP for the factories to connect: PING and
// CLUSTER SLOTS (one shard on itself). Any other command gets an error reply.
func fakeRedis(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	addr := ln.Addr().String()
	host, port, _ := net.SplitHostPort(addr)
	slots := fmt.Sprintf("*1\r\n*3\r\n:0\r\n:16383\r\n*2\r\n$%d\r\n%s\r\n:%s\r\n", len(host), host, port)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go serveRESP(conn, slots)
		}
	}()
	return addr
}

func serveRESP(conn net.Conn, slots string) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		args, err := readCommand(r)
		if err != nil {
			return
		}
		reply := "-ERR unknown command\r\n"
		switch cmd := strings.ToUpper(strings.Join(args, " ")); {
		case cmd == "PING":
			reply = "+PONG\r\n"
		case cmd == "CLUSTER SLOTS":
			reply = slots
		}
		if _, err := io.WriteString(conn, reply); err != nil {
			return
		}
	}
}

func readCommand(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "*")))
	if err != nil {
		return nil, err
	}
	args := make([]string, 0, n)
	for i := 0; i < n; i++ {
		if _, err := r.ReadString('\n'); err != nil { // $len
			return nil, err
		}
		arg, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		args = append(args, strings.TrimRight(arg, "\r\n"))
	}
	return args, nil
}

func TestNewClientWithOptionsRunsPluginsInOrderOnTheBuiltClient(t *testing.T) {
	var calls []string
	var got *redis.Client
	plugin := func(name string) Plugin {
		return func(_ context.Context, c *redis.Client) error {
			calls = append(calls, name)
			got = c
			return nil
		}
	}

	client, err := NewClientWithOptions(context.Background(), &Options{Client: ClientOptions{Addr: fakeRedis(t)}}, plugin("a"), plugin("b"))

	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	assert.Equal(t, []string{"a", "b"}, calls)
	assert.Same(t, client, got)
}

func TestNewClientWithOptionsReturnsPluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := func(context.Context, *redis.Client) error { return boom }
	next := func(context.Context, *redis.Client) error { after = true; return nil }

	client, err := NewClientWithOptions(context.Background(), &Options{Client: ClientOptions{Addr: fakeRedis(t)}}, failing, next)

	assert.ErrorIs(t, err, boom)
	assert.Nil(t, client)
	assert.False(t, after, "plugins after the failing one must not run")
}

func TestNewClusterClientWithOptionsRunsPluginsInOrderOnTheBuiltClient(t *testing.T) {
	var calls []string
	var got *redis.ClusterClient
	plugin := func(name string) ClusterPlugin {
		return func(_ context.Context, c *redis.ClusterClient) error {
			calls = append(calls, name)
			got = c
			return nil
		}
	}

	client, err := NewClusterClientWithOptions(context.Background(), &Options{Cluster: ClusterOptions{Addrs: []string{fakeRedis(t)}}}, plugin("a"), plugin("b"))

	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	assert.Equal(t, []string{"a", "b"}, calls)
	assert.Same(t, client, got)
}

func TestNewClusterClientWithOptionsReturnsPluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := func(context.Context, *redis.ClusterClient) error { return boom }
	next := func(context.Context, *redis.ClusterClient) error { after = true; return nil }

	client, err := NewClusterClientWithOptions(context.Background(), &Options{Cluster: ClusterOptions{Addrs: []string{fakeRedis(t)}}}, failing, next)

	assert.ErrorIs(t, err, boom)
	assert.Nil(t, client)
	assert.False(t, after, "plugins after the failing one must not run")
}
