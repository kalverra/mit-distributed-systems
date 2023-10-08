// Package comms provides the RPC communication structs between the coordinator and workers
package comms

import (
	"fmt"
	"net/rpc"
)

// KeyValue represents a key value pair
type KeyValue struct {
	Key   string
	Value string
}

// MapFunction represents a map function
type MapFunction func(string, string) []KeyValue

// ReduceFunction represents a reduce function
type ReduceFunction func(string, []string) string

type ReduceCall struct {
	Key    string
	Values []string
}

type WorkerReply struct {
	WorkerID   int
	ResultFile string
}

const (
	StatusIdle   = "Idle"
	StatusMap    = "Map"
	StatusReduce = "Reduce"
)

// A NOP to make RPC happy
type WorkerStatusCall struct {
}

// WorkerStatusReply shows the status of a worker
type WorkerStatusReply struct {
	WorkerID int
	Status   string
}

func Call(workerPort int, function string, args any, reply any) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%d", workerPort))
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Call(fmt.Sprintf("Worker.%s", function), args, reply)
}
