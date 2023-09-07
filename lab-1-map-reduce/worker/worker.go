package worker

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/rs/zerolog/log"
)

func New(port int) {
	server := rpc.NewServer()
	worker := Worker{ID: port}
	server.Register(&worker)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}
	go server.Accept(l)
	log.Info().Int("ID", worker.ID).Msg("Worker Running")
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
