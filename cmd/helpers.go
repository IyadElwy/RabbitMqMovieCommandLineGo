package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
)

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Genres struct {
	Genres []Genre `json:"genres"`
}

func readGenres(jsonFilePath string) []Genre {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		panic(fmt.Sprintf("%s not found", jsonFilePath))
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic("Error reading json file")
	}
	var genres Genres
	err = json.Unmarshal(byteValue, &genres)
	if err != nil {
		panic("Error unmarshaling json")
	}
	return genres.Genres
}

func createRandomGenreGenerator(jsonFilePath string) func() int {
	genres := readGenres(jsonFilePath)
	lenGenres := len(genres)
	return func() int {
		randomIndex := rand.Intn(lenGenres)
		return genres[randomIndex].Id
	}
}

func returnNRandomGenres(jsonFilePath string, n int) []int {
	genreIds := make([]int, 0)
	genreMap := make(map[int]bool)
	randomGenreGenerator := createRandomGenreGenerator(jsonFilePath)
	for i := 0; i < int(math.Min(float64(n), float64(5))); i++ {
		for {
			randomGenreId := randomGenreGenerator()
			_, ok := genreMap[randomGenreId]
			if !ok {
				genreIds = append(genreIds, randomGenreId)
				break
			}
		}
	}
	return genreIds
}

func getIdByGenreName(jsonFilePath string, genreName string) (int, error) {
	genres := readGenres(jsonFilePath)
	for _, genre := range genres {
		if genre.Name == genreName {
			return genre.Id, nil
		}
	}
	return 0, ErrGenreNotFound
}
