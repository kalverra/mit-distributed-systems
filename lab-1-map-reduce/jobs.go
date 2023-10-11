package main

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/kalverra/lab-1-map-reduce/comms"
)

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
	// function to detect word separators.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	// split contents into an array of words.
	words := strings.FieldsFunc(contents, ff)

	kva := []comms.KeyValue{}
	for _, w := range words {
		w = strings.ToLower(w)
		kv := comms.KeyValue{Key: w, Value: "1"}
		kva = append(kva, kv)
	}
	return kva
}

func (wc *WordCount) Reduce(key string, values []string) string {
	return strconv.Itoa(len(values))
}
