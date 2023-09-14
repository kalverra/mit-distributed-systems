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

const (
	StatusIdle = iota
	StatusMap
	StatusReduce
)

// Worker represents a worker node
type Worker struct {
	ID         int
	Status     int
	MapFunc    comms.MapFunction
	ReduceFunc comms.ReduceFunction

	coordinatorPort int
}

// New creates a new worker on a given port
func New(workerPort, coordinatorPort int) error {
	server := rpc.NewServer()
	worker := Worker{
		ID:     workerPort,
		Status: StatusIdle,
	}
	server.Register(&worker)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", workerPort))
	if err != nil {
		return fmt.Errorf("worker failed to listen: %w", err)
	}
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
	return nil
}

// RegisterMapReduce registers a map and reduce function to a worker
func (w *Worker) RegisterMapReduce(call *comms.RegisterMapReduce, reply *comms.WorkerReply) error {
	w.MapFunc = call.MapFunction
	w.ReduceFunc = call.ReduceFunction
	reply.WorkerID = w.ID
	return nil
}

func (w *Worker) Map(call *comms.KeyValue, reply *comms.WorkerReply) error {
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

func (w *Worker) Reduce(call *comms.ReduceCall, reply *comms.WorkerReply) error {
	answerFileName := fmt.Sprintf("tmp/%d-%s", w.ID, call.Key)
	answerFile, err := os.OpenFile(answerFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open answer file: %w", err)
	}
	defer answerFile.Close()

	answerFile.WriteString(w.ReduceFunc(call.Key, call.Values))

	reply.WorkerID = w.ID
	reply.ResultFile = answerFileName
	return nil

}
