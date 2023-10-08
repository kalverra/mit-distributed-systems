package coordinator

import (
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kalverra/lab-1-map-reduce/comms"
	"github.com/rs/zerolog/log"
)

var (
	workers       []int
	numReduce     int
	inputDir      string
	done          bool = false
	data               = make(map[string]string)
	workerTimeout      = 5 * time.Second
	pollInterval       = 1 * time.Second
)

func New(workerPorts []int, numReduce int, inputDir string) {
	log.Info().Interface("Worker Ports", workerPorts).Msg("Starting Coordinator")

	for _, port := range workerPorts {
		client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to dial worker")
		}
		workers = append(workers, port)
		client.Close()
	}

	log.Info().Msg("Coordinator Started")

	startTime := time.Now()

	log.Info().Msg("Loading Data")

	if err := loadData(inputDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to load data")
	}

	log.Info().Str("Time Taken", time.Since(startTime).String()).Msg("Data Loaded")

	startTime = time.Now()

	log.Info().Msg("Starting Map")

	idleWorkers := make(chan int, 0)
	monitorWorkers(idleWorkers)

	for {
		idle := <-idleWorkers
		log.Debug().Int("ID", idle).Msg("Got Idle Worker")
	}

}

// monitorWorkers polls the workers for their status and sends any idle ones back on a channel
func monitorWorkers(chan int) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for range ticker.C {
		for _, w := range workers {
			worker := w
			go func() {
				reply := comms.WorkerStatusReply{}
				err := comms.Call(worker, "GetStatus", struct{}{}, reply)
				if err != nil {
					log.Error().Err(err).Int("Worker", worker).Msg("Failed to get status")
					return
				}
				log.Debug().Int("Worker", reply.WorkerID).Str("Status", reply.Status).Msg("Got Status")
			}()
		}
	}
}

// loadData loads the data from the input directory into memory
func loadData(inputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		// Check if it's a regular file (not a directory or symlink)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".txt") {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			data[info.Name()] = string(fileContent)
		}

		return nil
	})
}
