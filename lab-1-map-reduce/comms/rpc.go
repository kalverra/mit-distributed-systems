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

// RegisterMapReduce represents a way to register a map and reduce function to a worker
type RegisterMapReduce struct {
	MapFunction    MapFunction
	ReduceFunction ReduceFunction
}

type ReduceCall struct {
	Key    string
	Values []string
}

type WorkerReply struct {
	WorkerID   int
	ResultFile string
}

func Call(workerPort int, function string, args interface{}, reply interface{}) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%d", workerPort))
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Call(function, args, reply)
}
