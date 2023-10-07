package coordinator

import (
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

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

func Run() error {
	done = false
	defer func() {
		done = true
	}()

	err := os.Mkdir("tmp", 0755)
	if err != nil {
		return err
	}

	toReduce := mapF()
	reduceF(toReduce)

	return nil
}

func IsDone() bool {
	return done
}

// mapF walks the input directory and gives each worker a file to map
// it returns a list of all the answer files for reduce tasks to take on
func mapF() []string {
	start := time.Now()
	filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
		return nil
	})

	log.Info().Str("Time Taken", time.Since(start).String()).Msg("Map Finished")
	return nil
}

// reduceF
func reduceF(filesToReduce []string) {

}
