package worker

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"

	"github.com/kalverra/lab-1-map-reduce/comms"
)

var worker *Worker

// Worker represents a worker node
type Worker struct {
	ID         int
	Status     string
	MapFunc    comms.MapFunction
	ReduceFunc comms.ReduceFunction
}

// New creates a new worker on a random free port and returns the port
func New(mapFunc comms.MapFunction, reduceFunc comms.ReduceFunction) (int, error) {
	server := rpc.NewServer()
	worker = &Worker{
		Status:     comms.StatusIdle,
		MapFunc:    mapFunc,
		ReduceFunc: reduceFunc,
	}
	server.Register(worker)
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("worker failed to listen: %w", err)
	}
	worker.ID = l.Addr().(*net.TCPAddr).Port
	go server.Accept(l)

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt, os.Kill)
	go func() {
		// We'd want to save our place, or at least alert the coordinator normally, but we'll leave that for later
		<-shutDown
		log.Warn().Int("ID", worker.ID).Msg("Worker Shutting Down")
		os.Exit(0)
	}()

	log.Info().Int("ID", worker.ID).Msg("Worker Running")
	return worker.ID, nil
}

// Map runs the map function on the worker
func (w *Worker) Map(call *comms.KeyValue, reply *comms.WorkerReply) error {
	w.Status = comms.StatusMap
	answerFileName := fmt.Sprintf("tmp/%d-%s", w.ID, call.Key)
	answerFile, err := os.OpenFile(answerFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open answer file: %w", err)
	}
	defer answerFile.Close()

	for _, kv := range w.MapFunc(call.Key, call.Value) {
		answerFile.WriteString(fmt.Sprintf("%s: %s\n", kv.Key, kv.Value))
	}
	reply.WorkerID = w.ID
	reply.ResultFile = answerFileName
	return nil
}

// Reduce runs the reduce function on the worker
func (w *Worker) Reduce(call *comms.ReduceCall, reply *comms.WorkerReply) error {
	w.Status = comms.StatusReduce
	answerFileName := fmt.Sprintf("tmp/%d-%s", w.ID, call.Key)
	answerFile, err := os.OpenFile(answerFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open answer file: %w", err)
	}
	defer answerFile.Close()

	_, err = answerFile.WriteString(w.ReduceFunc(call.Key, call.Values))
	reply.WorkerID = w.ID
	reply.ResultFile = answerFileName
	return err
}

// GetStatus returns the status of the worker
func (w *Worker) GetStatus(_ *struct{}, reply *comms.WorkerStatusReply) error {
	reply.WorkerID = w.ID
	reply.Status = w.Status
	return nil
}
