package pollen_history_data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type PollenHistoryDataService interface {
	DownloadLatestPollenData()
	PollenDataForDay(day time.Time) (*PollenDataForDay, error)
}

type pollenHistoryDataService struct {
	repository PollenDataRepository
}

func New(repository PollenDataRepository) PollenHistoryDataService {
	return &pollenHistoryDataService{
		repository: repository,
	}
}

type PollenDataFromSource struct {
	Message string       `json:"message"`
	Data    []PollenData `json:"data"`
}

type PollenData struct {
	Count struct {
		GrassPollen int64 `json:"grass_pollen"`
		TreePollen  int64 `json:"tree_pollen"`
		WeedPollen  int64 `json:"weed_pollen"`
	} `json:"Count"`
	Species struct {
		Grass struct {
			GrassOrPoaceae int64 `json:"Grass / Poaceae"`
		} `json:"Grass"`
		Others int64
		Tree   struct {
			Alder              int64 `json:"Alder"`
			Birch              int64 `json:"Birch"`
			Cypress            int64 `json:"Cypress"`
			Elm                int64 `json:"Elm"`
			Hazel              int64 `json:"Hazel"`
			Oak                int64 `json:"Oak"`
			Pine               int64 `json:"Pine"`
			Plan               int64 `json:"Plane"`
			PoplarOrCottonWood int64 `json:"Poplar Cottonwood"`
		} `json:"Tree"`
		Weed struct {
			Chenopod int64 `json:"Chenopod"`
			Mugwort  int64 `json:"Mugwort"`
			Nettle   int64 `json:"Nettle"`
			Ragweed  int64 `json:"Ragweed"`
		} `json:"Weed"`
	} `json:"Species"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (service pollenHistoryDataService) DownloadLatestPollenData() {
	url := "https://api.ambeedata.com/latest/pollen/by-lat-lng?lat=52.3676&lng=4.9041"
	req, _ := http.NewRequest("GET", url, nil)

	// TODO: Dependency injection for env files
	apiKey := os.Getenv("API_KEY")

	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("Content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var pollenDataFromSource PollenDataFromSource
	if err := json.Unmarshal(body, &pollenDataFromSource); err != nil {
		fmt.Println(err.Error())
	}

	if pollenDataFromSource.Message != "success" {
		fmt.Println("Expected success")
		return
	}

	data := pollenDataFromSource.Data[0]
	pollenDataForDay := PollenDataForDay{
		Id:  4,
		Day: data.UpdatedAt,
		Data: []PollenDataForSpecies{
			{
				Species: PollenSpecies{
					Name: GrassOrPoaceae,
					Kind: Grass,
				},
				Value: data.Species.Grass.GrassOrPoaceae,
			},
			{
				Species: PollenSpecies{
					Name: OthersType,
					Kind: OthersKind,
				},
				Value: data.Species.Others,
			},
			{
				Species: PollenSpecies{
					Name: Alder,
					Kind: Tree,
				},
				Value: data.Species.Tree.Alder,
			},
			{
				Species: PollenSpecies{
					Name: Birch,
					Kind: Tree,
				},
				Value: data.Species.Tree.Birch,
			},
			{
				Species: PollenSpecies{
					Name: Cypress,
					Kind: Tree,
				},
				Value: data.Species.Tree.Cypress,
			},
			{
				Species: PollenSpecies{
					Name: Elm,
					Kind: Tree,
				},
				Value: data.Species.Tree.Elm,
			},
			{
				Species: PollenSpecies{
					Name: Hazel,
					Kind: Tree,
				},
				Value: data.Species.Tree.Hazel,
			},
			{
				Species: PollenSpecies{
					Name: Oak,
					Kind: Tree,
				},
				Value: data.Species.Tree.Oak,
			},
			{
				Species: PollenSpecies{
					Name: Pine,
					Kind: Tree,
				},
				Value: data.Species.Tree.Pine,
			},
			{
				Species: PollenSpecies{
					Name: Plan,
					Kind: Tree,
				},
				Value: data.Species.Tree.Plan,
			},
			{
				Species: PollenSpecies{
					Name: PoplarOrCottonWood,
					Kind: Tree,
				},
				Value: data.Species.Tree.PoplarOrCottonWood,
			},
			{
				Species: PollenSpecies{
					Name: Chenopod,
					Kind: Weed,
				},
				Value: data.Species.Weed.Chenopod,
			},
			{
				Species: PollenSpecies{
					Name: Mugwort,
					Kind: Weed,
				},
				Value: data.Species.Weed.Mugwort,
			},
			{
				Species: PollenSpecies{
					Name: Nettle,
					Kind: Weed,
				},
				Value: data.Species.Weed.Nettle,
			},
			{
				Species: PollenSpecies{
					Name: Ragweed,
					Kind: Weed,
				},
				Value: data.Species.Weed.Ragweed,
			},
		},
	}

	service.repository.Store(pollenDataForDay)
}

func (service pollenHistoryDataService) PollenDataForDay(day time.Time) (*PollenDataForDay, error) {
	return service.repository.FindForDay(day)
}
