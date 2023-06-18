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

type Movies struct {
	MovieId []string `json:"movies"`
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

	for _, id := range actorIDs {
		id = strings.ReplaceAll(id, `"`, "")
		idList = append(idList, id)
	}

	fmt.Println(idList)

	for _, id := range idList {
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

		fmt.Println(string(responseData))

		var actor Actor
		err = json.Unmarshal(responseData, &actor)
		if err != nil {
			log.Fatal(err)
		}

		for _, movieID := range actor.Movies {
			movieID = strings.ReplaceAll(movieID, `"`, "")
			movieIdList = append(movieIdList, movieID)
		}
	}

	for _, movieID := range movieIdList {
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

		fmt.Println(string(responseData))
	}
}
