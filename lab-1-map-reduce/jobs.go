package main

import "github.com/kalverra/lab-1-map-reduce/comms"

var registeredJobs = map[string]MapReduceJob{
	"word-count": &WordCount{},
}

type MapReduceJob interface {
	Map(filename string, contents string) []comms.KeyValue
	Reduce(key string, values []string) string
}

type WordCount struct {
}

func (wc *WordCount) Map(filename string, contents string) []comms.KeyValue {
	return []comms.KeyValue{
		{Key: filename, Value: "map"},
	}
}

func (wc *WordCount) Reduce(key string, values []string) string {
	return "reduce"
}
