package pollen_history_data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type PollenSpeciesType string

const (
	GrassOrPoaceae     PollenSpeciesType = "grassOrPoaceae"
	OthersType                           = "others"
	Alder                                = "alder"
	Birch                                = "birch"
	Cypress                              = "cypress"
	Elm                                  = "elm"
	Hazel                                = "hazel"
	Oak                                  = "oak"
	Pine                                 = "pine"
	Plan                                 = "plan"
	PoplarOrCottonWood                   = "poplarOrCottonWood"
	Chenopod                             = "chenopod"
	Mugwort                              = "mugwort"
	Nettle                               = "nettle"
	Ragweed                              = "ragweed"
)

type PollenSpeciesKind string

const (
	Grass      PollenSpeciesKind = "grass"
	Tree                         = "tree"
	Weed                         = "weed"
	OthersKind                   = "others"
)

type PollenSpecies struct {
	Name PollenSpeciesType `json:"name"`
	Kind PollenSpeciesKind `json:"kind"`
}

type PollenDataForSpecies struct {
	Species PollenSpecies `json:"species"`
	Value   int64         `json:"value"`
}

type PollenDataForDay struct {
	Id   int64                  `json:"id"`
	Day  time.Time              `json:"day"`
	Data []PollenDataForSpecies `json:"data"`
}

type PollenDataRepository interface {
	Store(pollenDataForDay PollenDataForDay)
	FindForDay(day time.Time) (*PollenDataForDay, error)
}

type pollenDataMySQLRepository struct {
}

func NewMySQLRepository() PollenDataRepository {
	return pollenDataMySQLRepository{}
}

func (repository pollenDataMySQLRepository) Store(day PollenDataForDay) {
	// TODO: Implement MySQL insert query
}

func (repository pollenDataMySQLRepository) FindForDay(day time.Time) (*PollenDataForDay, error) {
	// TODO: Implement MySQL search query

	return &PollenDataForDay{
		1,
		day,
		[]PollenDataForSpecies{
			{
				Species: PollenSpecies{
					Name: GrassOrPoaceae,
					Kind: Grass,
				},
				Value: 50,
			},
		},
	}, nil
}

type pollenDataFileRepository struct {
}

func (p pollenDataFileRepository) Store(pollenDataForDay PollenDataForDay) {
	contentToWrite, err := json.Marshal(pollenDataForDay)
	if err != nil {
		panic(err)
	}

	filePath := fmt.Sprintf("./data/%s.json", pollenDataForDay.Day.Format("2006-01-02"))
	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.Write(contentToWrite)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (p pollenDataFileRepository) FindForDay(day time.Time) (*PollenDataForDay, error) {
	panic("implement me")
}

func NewRepositoryUsingFiles() PollenDataRepository {
	return pollenDataFileRepository{}
}
