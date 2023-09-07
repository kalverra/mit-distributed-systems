package worker

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/rs/zerolog/log"
)

func New(port int) error {
	server := rpc.NewServer()
	worker := Worker{ID: port}
	server.Register(&worker)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("worker failed to listen: %w", err)
	}
	go server.Accept(l)
	log.Info().Int("ID", worker.ID).Msg("Worker Running")
	return nil
}

type Worker struct {
	ID int
}

type MapArgs struct {
	JobName  string
	FileName string
}

type MapReply struct {
	Msg string
}

func (w *Worker) Map(args *MapArgs, reply *MapReply) error {
	log.Info().Msgf("Map called with args %v", args)
	reply.Msg = "Map called"
	return nil
}

type ReduceArgs struct {
	JobName string
}

type ReduceReply struct {
	Msg string
}

func (w *Worker) Reduce(args *ReduceArgs, reply *ReduceReply) error {
	log.Info().Msgf("Reduce called with args %v", args)
	reply.Msg = "Reduce called"
	return nil
}
