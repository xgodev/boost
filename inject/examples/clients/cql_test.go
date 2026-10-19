package main

import (
	"encoding/binary"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// cqlReleaseVersion is the release_version the fake node reports.
const cqlReleaseVersion = "4.1.0"

// CQL native protocol opcodes the fake node reads and writes.
const (
	opError     = 0x00
	opStartup   = 0x01
	opReady     = 0x02
	opOptions   = 0x05
	opSupported = 0x06
	opQuery     = 0x07
	opResult    = 0x08
	opPrepare   = 0x09
	opExecute   = 0x0A
	opRegister  = 0x0B
)

// startCQLNode starts a fake Cassandra node on a free port and returns the
// port. It answers what a gocql session sends while it connects (OPTIONS,
// STARTUP, REGISTER, the system.local and schema queries) and selects on
// system.local, plain or prepared; any other query gets an empty result.
func startCQLNode(t *testing.T) int {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go serveCQL(conn)
		}
	}()

	return ln.Addr().(*net.TCPAddr).Port
}

func serveCQL(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 9) // version, flags, stream (2), opcode, length (4)
	for {
		if _, err := io.ReadFull(conn, header); err != nil {
			return
		}
		body := make([]byte, binary.BigEndian.Uint32(header[5:]))
		if _, err := io.ReadFull(conn, body); err != nil {
			return
		}

		var op byte
		var resp []byte
		switch header[4] {
		case opOptions:
			op, resp = opSupported, cqlShort(nil, 0) // no options
		case opStartup, opRegister:
			op = opReady
		case opQuery:
			n := binary.BigEndian.Uint32(body)
			op, resp = opResult, cqlQueryResult(string(body[4:4+n]))
		case opPrepare:
			n := binary.BigEndian.Uint32(body)
			op, resp = opResult, cqlPrepared(string(body[4:4+n]))
		case opExecute: // the id of a prepared statement is its query
			n := binary.BigEndian.Uint16(body)
			op, resp = opResult, cqlQueryResult(string(body[2:2+n]))
		default:
			op = opError
			resp = binary.BigEndian.AppendUint32(nil, 0x000A) // protocol error
			resp = cqlString(resp, "unsupported opcode")
		}

		out := []byte{header[0] | 0x80, 0, header[2], header[3], op}
		out = binary.BigEndian.AppendUint32(out, uint32(len(resp)))
		if _, err := conn.Write(append(out, resp...)); err != nil {
			return
		}
	}
}

// localColumns returns the system.local columns query selects, or nil
// when it does not read system.local.
func localColumns(query string) [][2]string {
	if !strings.Contains(query, "FROM system.local") {
		return nil
	}

	local := [][2]string{{"data_center", "dc1"}, {"rack", "rack1"}, {"release_version", cqlReleaseVersion}}
	selected := strings.TrimSpace(query[len("SELECT "):strings.Index(query, " FROM")])
	var cols [][2]string
	for _, col := range local {
		if selected == "*" || strings.Contains(selected, col[0]) {
			cols = append(cols, col)
		}
	}
	return cols
}

// cqlQueryResult returns the system.local row with the selected columns,
// or a void result for any other query.
func cqlQueryResult(query string) []byte {
	cols := localColumns(query)
	if cols == nil {
		return binary.BigEndian.AppendUint32(nil, resultVoid)
	}

	b := binary.BigEndian.AppendUint32(nil, resultRows)
	b = append(b, cqlRowsMetadata(cols)...)
	b = binary.BigEndian.AppendUint32(b, 1) // one row
	for _, col := range cols {
		b = binary.BigEndian.AppendUint32(b, uint32(len(col[1])))
		b = append(b, col[1]...)
	}
	return b
}

// cqlPrepared prepares query without bind markers; its id is the query.
func cqlPrepared(query string) []byte {
	b := binary.BigEndian.AppendUint32(nil, resultPrepared)
	b = append(cqlShort(b, uint16(len(query))), query...)
	b = binary.BigEndian.AppendUint32(b, 0) // bind metadata flags
	b = binary.BigEndian.AppendUint32(b, 0) // bind markers
	b = binary.BigEndian.AppendUint32(b, 0) // partition key indexes
	return append(b, cqlRowsMetadata(localColumns(query))...)
}

// cqlRowsMetadata describes cols, all varchar columns of system.local.
func cqlRowsMetadata(cols [][2]string) []byte {
	const globalTablesSpec, varchar = 0x0001, 0x000D

	if len(cols) == 0 {
		return binary.BigEndian.AppendUint64(nil, 0) // no flags, no columns
	}
	b := binary.BigEndian.AppendUint32(nil, globalTablesSpec)
	b = binary.BigEndian.AppendUint32(b, uint32(len(cols)))
	b = cqlString(cqlString(b, "system"), "local")
	for _, col := range cols {
		b = cqlShort(cqlString(b, col[0]), varchar)
	}
	return b
}

// CQL result kinds.
const (
	resultVoid     = 0x0001
	resultRows     = 0x0002
	resultPrepared = 0x0004
)

func cqlShort(b []byte, v uint16) []byte {
	return binary.BigEndian.AppendUint16(b, v)
}

func cqlString(b []byte, s string) []byte {
	return append(cqlShort(b, uint16(len(s))), s...)
}
