package main

import (
	"flag"
	"os"

	"github.com/kalverra/lab-1-map-reduce/coordinator"
	"github.com/kalverra/lab-1-map-reduce/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	numWorkers := flag.Int("numWorkers", 2, "Number of workers to use")
	log.Info().Int("Workers", *numWorkers).Msg("Starting")

	workerPorts := make([]int, *numWorkers)
	for port := 8081; port < 8081+*numWorkers; port++ {
		worker.New(port)
		workerPorts[port-8081] = port
	}

	coordinator.New(8080, workerPorts)
}
