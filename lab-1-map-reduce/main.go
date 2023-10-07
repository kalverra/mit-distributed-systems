package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

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
	numReduce := flag.Int("numReduce", 2, "Number of reduce tasks to use")
	job := flag.String("job", "word-count", "Job to run")
	dumb := flag.Bool("dumb", false, "Use dumb word count to generate answers")
	flag.Parse()

	if *dumb {
		dumbWordCount()
		return
	}

	var mapReduceJob MapReduceJob
	switch *job {
	case "word-count":
		mapReduceJob = &WordCount{}
	default:
		log.Warn().Str("Job", *job).Msg("Unknown job, defaulting to word-count")
		mapReduceJob = &WordCount{}
	}

	log.Info().Int("Workers", *numWorkers).Msg("Starting")

	workerPorts := make([]int, *numWorkers)

	eg := errgroup.Group{}
	for p := 8081; p < 8081+*numWorkers; p++ {
		port := p
		workerPorts[port-8081] = port
		eg.Go(func() error {
			return worker.New(port, 8080, mapReduceJob.Map, mapReduceJob.Reduce)
		})
	}
	if err := eg.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start workers")
	}

	coordinator.New(workerPorts, *numReduce, "./books")
}

// dumbWordCount is a dumb word count implementation that doesn't use map reduce.
// It's used to generate the answers for the lab.
func dumbWordCount() {
	startTime := time.Now()
	booksDir := "./books"
	files, err := os.ReadDir(booksDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read books directory")
	}

	wordCount := map[string]int{}
	wordRegex := regexp.MustCompile("[^a-zA-Z]+")
	for _, file := range files {
		fileContents, err := os.ReadFile(filepath.Join(booksDir, file.Name()))
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read file %s", file.Name())
		}
		fileString := string(fileContents)
		for _, word := range strings.Split(fileString, " ") {
			word = wordRegex.ReplaceAllString(word, "")
			word = strings.ToLower(word)
			word = strings.Trim(word, ".,!?[]\n")

			if word == "\n" || word == "" {
				continue
			}

			if _, ok := wordCount[word]; !ok {
				wordCount[word] = 1
			} else {
				wordCount[word]++
			}
		}
	}

	sortedWords := []string{}

	for word := range wordCount {
		sortedWords = append(sortedWords, word)
	}
	sort.Slice(sortedWords, func(i, j int) bool {
		return wordCount[sortedWords[i]] > wordCount[sortedWords[j]]
	})

	answerFile, err := os.OpenFile("answers.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open answers.txt")
	}

	for _, word := range sortedWords {
		_, err := answerFile.WriteString(fmt.Sprintf("%s: %d\n", word, wordCount[word]))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to write to answers.txt")
		}
	}

	log.Info().Str("Time", time.Since(startTime).String()).Msg("Dumb word count complete")
}
