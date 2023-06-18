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

	file, err := os.Open("actors.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var actorIDs []string
	var idList []string
	var movieIdList []string
	actorScores := make(map[string]float32)
	actorMovieCount := make(map[string]int)

	for scanner.Scan() {
		actorIDs = append(actorIDs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Tratamento dos IDs dos atores
	for _, id := range actorIDs {
		id = strings.ReplaceAll(id, `"`, "")
		idList = append(idList, id)
	}

	for _, id := range idList {
		var totalScore float32 = 0
		response, err := http.Get(fmt.Sprintf("http://150.165.15.91:8001/actors/%s", id))
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		var actor Actor
		err = json.Unmarshal(responseData, &actor)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(string(responseData))

		// Tratamento dos IDs e cálculo do Score
		for _, movieID := range actor.Movies {
			movieID = strings.ReplaceAll(movieID, `"`, "")
			movieIdList = append(movieIdList, movieID)
			response, err := http.Get(fmt.Sprintf("http://150.165.15.91:8001/movies/%s", movieID))
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			var movie Movie
			err = json.Unmarshal(responseData, &movie)
			if err != nil {
				log.Fatal(err)
			}

			//fmt.Println(string(responseData))
			totalScore += movie.Score
			actorMovieCount[actor.Name]++
		}

		if actorMovieCount[actor.Name] > 0 {
			actorScores[actor.Name] = totalScore / float32(actorMovieCount[actor.Name])
		} else {
			actorScores[actor.Name] = 0
		}

		//fmt.Println("Score: ", actorScores[actor.Name])
	}

	// Rank top10
	sortedActors := make([]ActorScore, 0, len(actorScores))
	for name, score := range actorScores {
		sortedActors = append(sortedActors, ActorScore{Name: name, Score: score})
	}

	sort.Slice(sortedActors, func(i, j int) bool {
		return sortedActors[i].Score > sortedActors[j].Score
	})

	fmt.Println("Top 10 Atores:")
	for i, actorScore := range sortedActors[:10] {
		fmt.Printf("%d. %s - Score: %.1f\n", i+1, actorScore.Name, actorScore.Score)
	}

	executionTime(startTime)
}

func executionTime(startTime time.Time) {
	duration := time.Since(startTime)
	fmt.Println("Tempo total de execução:", duration)
}
