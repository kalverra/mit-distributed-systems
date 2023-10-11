package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/kalverra/lab-1-map-reduce/coordinator"
	"github.com/kalverra/lab-1-map-reduce/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
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

	mapReduceJob, ok := registeredJobs[*job]
	if !ok {
		log.Fatal().Str("Job", *job).Msg("Job not found")
	}

	log.Info().Int("Workers", *numWorkers).Msg("Starting")

	workerPorts := []int{}
	for i := 0; i < *numWorkers; i++ {
		port, err := worker.New(mapReduceJob.Map, mapReduceJob.Reduce)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start worker")
		}
		workerPorts = append(workerPorts, port)
	}

	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		if err := os.Mkdir("./tmp", 0755); err != nil {
			log.Fatal().Err(err).Msg("Failed to create tmp directory")
		}
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
	ff := func(r rune) bool { return !unicode.IsLetter(r) }
	for _, file := range files {
		fileContents, err := os.ReadFile(filepath.Join(booksDir, file.Name()))
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read file %s", file.Name())
		}
		fileString := string(fileContents)
		words := strings.FieldsFunc(fileString, ff)
		for _, word := range words {
			word = strings.ToLower(word)

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
