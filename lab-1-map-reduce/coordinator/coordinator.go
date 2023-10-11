package coordinator

import (
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kalverra/lab-1-map-reduce/comms"
	"github.com/rs/zerolog/log"
)

// TODO: I think that having the workers ask for jobs rather than a coordinator assigning them is probably a better way to do this.
var (
	workers       []int
	numReduce     int
	inputDir      string
	done          bool = false
	data               = make(map[string]string)
	workerTimeout      = 5 * time.Second
	pollInterval       = 100 * time.Millisecond

	idleWorkers   = []int{}
	idleWorkersMu = sync.Mutex{}
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

	go monitorWorkers()

	ticker := time.NewTicker(pollInterval)

	// Map phase
	reduceFilesCh := make(chan string, len(data))
	var wg sync.WaitGroup
	for range ticker.C {
		if len(data) == 0 { // No more data to process
			break
		}
		idle, err := getIdleWorker()
		if err != nil {
			continue // No idle workers
		}
		log.Debug().Str("Key", key).Int("Worker", idle).Msg("Got Idle Worker, Assigning Map")
	}

	for key, value := range data {

		go func(key string, value string, idle int) {
			wg.Add(1)
			defer wg.Done()
			call := &comms.KeyValue{
				Key:   key,
				Value: value,
			}
			reply := &comms.WorkerReply{}
			err := comms.Call(idle, "Map", call, reply)
			if err != nil {
				log.Error().Err(err).Msg("Failed to call worker map")
			} else {
				reduceFilesCh <- reply.ResultFile
			}
		}(key, value, idle)
	}

	wg.Wait()
	log.Info().Str("Time Taken", time.Since(startTime).String()).Msg("Map Complete")

	startTime = time.Now()

	log.Info().Msg("Starting Reduce")

	close(reduceFilesCh)
	reduceFiles := []string{}
	for file := range reduceFilesCh {
		reduceFiles = append(reduceFiles, file)
	}

	// Reduce phase
	for _, reduceFile := range reduceFiles {
		idle, err := getIdleWorker()
		if err != nil {
			continue
		}
		log.Debug().Str("File", reduceFile).Int("Worker", idle).Msg("Got Idle Worker, Assigning Reduce")
		go func(reduceFile string, idle int) {
			wg.Add(1)
			defer wg.Done()
			call := &comms.ReduceCall{
				Key:    reduceFile,
				Values: []string{},
			}
			reply := &comms.WorkerReply{}
			err := comms.Call(idle, "Reduce", call, reply)
			if err != nil {
				log.Error().Err(err).Msg("Failed to call worker reduce")
			}
		}(reduceFile, idle)
	}

	wg.Wait()
	log.Info().Str("Time Taken", time.Since(startTime).String()).Msg("Reduce Complete")
}

// monitorWorkers polls the workers for their status and sends any idle ones back on a channel
func monitorWorkers() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for range ticker.C {
		for _, w := range workers {
			worker := w
			go func() {
				reply := &comms.WorkerStatusReply{}
				err := comms.Call(worker, "GetStatus", struct{}{}, reply)
				if err != nil {
					log.Error().Err(err).Int("Worker", worker).Msg("Failed to get status")
					return
				}
				log.Trace().Int("Worker", reply.WorkerID).Str("Status", reply.Status).Msg("Got Status")
				if reply.Status == comms.StatusIdle {
					idleWorkersMu.Lock()
					defer idleWorkersMu.Unlock()
					idleWorkers = append(idleWorkers, reply.WorkerID)
				}
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

// getIdleWorker returns the first idle worker available
func getIdleWorker() (int, error) {
	idleWorkersMu.Lock()
	defer idleWorkersMu.Unlock()
	if len(idleWorkers) == 0 {
		return -1, fmt.Errorf("no idle workers")
	}
	return idleWorkers[0], nil
}
