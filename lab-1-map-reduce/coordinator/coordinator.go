package coordinator

import (
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/kalverra/lab-1-map-reduce/comms"
	"github.com/rs/zerolog/log"
)

var (
	workers   []int
	numReduce int
	inputDir  string
	done      bool = false
)

func New(workerPorts []int, numReduce int, inputDir string) {
	log.Info().Interface("Worker Ports", workerPorts).Msg("Starting Coordinator")
	numReduce = numReduce

	for _, port := range workerPorts {
		client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to dial worker")
		}
		workers = append(workers, port)
		client.Close()
	}

	log.Info().Msg("Coordinator Started")
}

func RegisterMapReduce(m comms.MapFunction, r comms.ReduceFunction) {
	for _, port := range workers {
		var reply comms.WorkerReply
		err := comms.Call(port, "Worker.RegisterMapReduce", &comms.RegisterMapReduce{MapFunction: m, ReduceFunction: r}, &reply)
		if err != nil {
			log.Fatal().Int("Port", port).Err(err).Msg("Failed to register map reduce functions")
		}
	}
	log.Info().Msg("Map Reduce Functions Registered")
}

func Run() error {
	done = false
	defer func() {
		done = true
	}()

	err := os.Mkdir("tmp", 0755)
	if err != nil {
		return err
	}

	// Run Map
	filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
		return nil
	})

	return nil
}

func IsDone() bool {
	return done
}
