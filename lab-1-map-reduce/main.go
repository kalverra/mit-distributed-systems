package main

import (
	"flag"
	"os"

	"github.com/kalverra/lab-1-map-reduce/coordinator"
	"github.com/kalverra/lab-1-map-reduce/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	numWorkers := flag.Int("numWorkers", 2, "Number of workers to use")
	flag.Parse()

	log.Info().Int("Workers", *numWorkers).Msg("Starting")

	workerPorts := make([]int, *numWorkers)

	eg := errgroup.Group{}
	for p := 8081; p < 8081+*numWorkers; p++ {
		port := p
		workerPorts[port-8081] = port
		eg.Go(func() error {
			return worker.New(port)
		})
	}
	if err := eg.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start workers")
	}

	coordinator.New(8080, workerPorts)
}
