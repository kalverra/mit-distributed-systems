package coordinator

import (
	"fmt"
	"net/rpc"

	"github.com/kalverra/lab-1-map-reduce/worker"
	"github.com/rs/zerolog/log"
)

var workers []*rpc.Client

func New(port int, workerPorts []int) {
	log.Info().Int("Coordinator Port", port).Interface("Worker Ports", workerPorts).Msg("Starting Coordinator")

	for _, port := range workerPorts {
		client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to dial worker")
		}
		workers = append(workers, client)
	}

	log.Info().Msg("Coordinator Started")

	for _, client := range workers {
		reply := worker.MapReply{}
		err := client.Call("Worker.Map", &worker.MapArgs{JobName: "test", FileName: "test"}, &reply)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to call worker")
		}
		log.Info().Msgf("Worker replied with %s", reply)
	}
}
