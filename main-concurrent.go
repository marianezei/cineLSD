package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Actor struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Movies []string `json:"movies"`
}

type Movie struct {
	ID     string   `json:"movie_id"`
	Title  string   `json:"movie_title"`
	Score  float32  `json:"averageRating"`
	Votes  int      `json:"numberOfVotes"`
	Year   string   `json:"year"`
	Genres []string `json:"genres"`
}

type ActorScore struct {
	Name  string
	Score float32
}

func main() {
	startTime := time.Now()

	actorIDs, err := getActorIDs("actors.txt")
	if err != nil {
		log.Fatalf("Falha ao obter IDs dos atores: %v", err)
	}

	actorScores, err := calculateActorScores(actorIDs)
	if err != nil {
		log.Fatalf("Falha ao calcular scores dos atores: %v", err)
	}

	top10Actors := getTop10(actorScores)

	fmt.Println("Top 10 Atores:")
	for i, actorScore := range top10Actors {
		fmt.Printf("%d. %s - Score: %.1f\n", i+1, actorScore.Name, actorScore.Score)
	}

	executionTime(startTime)
}

func getActorIDs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var actorIDs []string

	for scanner.Scan() {
		actorIDs = append(actorIDs, strings.ReplaceAll(scanner.Text(), `"`, ""))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return actorIDs, nil
}

func getActor(actorID string) (*Actor, error) {
	response, err := http.Get(fmt.Sprintf("http://150.165.15.91:8001/actors/%s", actorID))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var actor Actor
	err = json.Unmarshal(responseData, &actor)
	if err != nil {
		return nil, err
	}

	return &actor, nil
}

func getMovie(movieID string) (*Movie, error) {
	response, err := http.Get(fmt.Sprintf("http://150.165.15.91:8001/movies/%s", movieID))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var movie Movie
	err = json.Unmarshal(responseData, &movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func calculateActorScores(actorIDs []string) (map[string]float32, error) {
	actorScores := make(map[string]float32)
	actorMovieCount := make(map[string]int)

	const numWorkers = 16
	var wg sync.WaitGroup
	var mutex sync.Mutex

	actorCh := make(chan string)

	worker := func() {
		defer wg.Done()

		for actorID := range actorCh {
			actor, err := getActor(actorID)
			if err != nil {
				log.Printf("Falha ao obter informações do ator %s: %v", actorID, err)
				continue
			}

			var totalScore float32
			for _, movieID := range actor.Movies {
				movie, err := getMovie(strings.ReplaceAll(movieID, `"`, ""))
				if err != nil {
					log.Printf("Falha ao obter informações do filme %s: %v", movieID, err)
					continue
				}

				totalScore += movie.Score

				mutex.Lock()
				actorMovieCount[actor.Name]++
				mutex.Unlock()
			}

			mutex.Lock()
			if actorMovieCount[actor.Name] > 0 {
				actorScores[actor.Name] = totalScore / float32(actorMovieCount[actor.Name])
			} else {
				actorScores[actor.Name] = 0
			}
			mutex.Unlock()
		}
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	for _, actorID := range actorIDs {
		actorCh <- actorID
	}
	close(actorCh)

	wg.Wait()

	return actorScores, nil
}

func getTop10(actorScores map[string]float32) []ActorScore {
	sortedActors := make([]ActorScore, 0, len(actorScores))
	for name, score := range actorScores {
		sortedActors = append(sortedActors, ActorScore{Name: name, Score: score})
	}

	sort.Slice(sortedActors, func(i, j int) bool {
		return sortedActors[i].Score > sortedActors[j].Score
	})

	if len(sortedActors) > 10 {
		return sortedActors[:10]
	}
	return sortedActors
}

func executionTime(startTime time.Time) {
	duration := time.Since(startTime)
	fmt.Println("Tempo total de execução:", duration)
}
