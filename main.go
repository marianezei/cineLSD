package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

func main() {
	file, err := os.Open("actors.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var actorIDs []string
	var idList []string
	var movieIdList []string

	for scanner.Scan() {
		actorIDs = append(actorIDs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//tratamento dos IDs dos atores
	for _, id := range actorIDs {
		id = strings.ReplaceAll(id, `"`, "")
		idList = append(idList, id)
	}

	for _, id := range idList {
		var totalScore float32 = 0
		var totalVotes int = 0
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
		fmt.Println(string(responseData))

		//tratamento dos IDs dos filmes e calculo
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

			fmt.Println(string(responseData))
			totalScore += movie.Score
			totalVotes += movie.Votes

		}

		fmt.Println("Total Score: ", totalScore)
		fmt.Println("Total Votes: ", totalVotes)
	}

}
