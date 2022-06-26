package main

import (
	"encoding/json"
	"errors"
	"fmt"
	pollen_history_data "hayfever/tracker/pollen-history-data"
	self_assessment "hayfever/tracker/self-assessment"
	"net/http"
	"strings"
	"time"
)

const (
	jsonDateFormat = "2006-01-02"
)

type JsonDate time.Time

func (t *JsonDate) UnmarshalJSON(data []byte) (err error) {
	newTime, err := time.ParseInLocation("\""+jsonDateFormat+"\"", string(data), time.Local)
	*t = JsonDate(newTime)
	return
}

func (t JsonDate) MarshalJSON() ([]byte, error) {
	timeStr := fmt.Sprintf("\"%s\"", time.Time(t).Format(jsonDateFormat))
	return []byte(timeStr), nil
}

func (t JsonDate) String() string {
	return time.Time(t).Format(jsonDateFormat)
}

func (t JsonDate) Date() time.Time {
	return time.Time(t)
}

type Api interface {
	GetPollenHistoryDataForToday(w http.ResponseWriter, r *http.Request)
	SubmitSelfAssessment(w http.ResponseWriter, r *http.Request)
	DownloadLatestPollenData(w http.ResponseWriter, r *http.Request)
}

type NetHttpApi struct {
	pollenHistoryDataService pollen_history_data.PollenHistoryDataService
	selfAssessmentService    self_assessment.SelfAssessmentService
}

func NewApi(pollenHistoryDataService pollen_history_data.PollenHistoryDataService, selfAssessmentService self_assessment.SelfAssessmentService) Api {
	return NetHttpApi{
		pollenHistoryDataService: pollenHistoryDataService,
		selfAssessmentService:    selfAssessmentService,
	}
}

type SpeciesData map[string]interface{}

type PollenHistoryDataResponse struct {
	Day  JsonDate      `json:"day" `
	Data []SpeciesData `json:"data"`
}

func (api NetHttpApi) GetPollenHistoryDataForToday(w http.ResponseWriter, r *http.Request) {
	pollenDataForToday, err := api.pollenHistoryDataService.PollenDataForDay(time.Now())

	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}

	data := []SpeciesData{}
	for _, speciesData := range pollenDataForToday.Data {
		data = append(data, SpeciesData{
			"speciesName": string(speciesData.Species.Name),
			"speciesKind": string(speciesData.Species.Kind),
			"value":       speciesData.Value,
		})
	}
	responseData := PollenHistoryDataResponse{
		Day:  JsonDate(pollenDataForToday.Day),
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
	return
}

type SubmitSelfAssessmentRequest struct {
	NoseFeeling   *string `json:"noseFeeling"`
	EyesFeeling   *string `json:"eyesFeeling"`
	ThroatFeeling *string `json:"throatFeeling"`
}

func (r SubmitSelfAssessmentRequest) Valid() (bool, error) {
	if r.NoseFeeling == nil {
		return false, errors.New("noseFeeling is empty")
	}

	if r.EyesFeeling == nil {
		return false, errors.New("eyesFeeling is empty")
	}

	if r.ThroatFeeling == nil {
		return false, errors.New("throatFeeling is empty")
	}

	return true, nil
}

func (api NetHttpApi) SubmitSelfAssessment(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"error\": \"Only POST method supported\"}")
		w.WriteHeader(http.StatusBadRequest)
	}

	// Declare a new Person struct.
	var reportInput SubmitSelfAssessmentRequest

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&reportInput)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		errorMessage := strings.ReplaceAll(err.Error(), `"`, "'")
		fmt.Fprintf(w, fmt.Sprintf(`{"error": "%s"}`, errorMessage))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if valid, err := reportInput.Valid(); valid == false {
		w.Header().Set("Content-Type", "application/json")
		errorMessage := strings.ReplaceAll(err.Error(), `"`, "'")
		fmt.Fprintf(w, fmt.Sprintf(`{"error": "%s"}`, errorMessage))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.selfAssessmentService.SubmitReport(self_assessment.ReportInput{
		ReportedAt:    time.Now(),
		NoseFeeling:   self_assessment.NoseFeeling(*reportInput.NoseFeeling),
		EyesFeeling:   self_assessment.EyesFeeling(*reportInput.EyesFeeling),
		ThroatFeeling: self_assessment.ThroatFeeling(*reportInput.ThroatFeeling),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return
}

func (api NetHttpApi) DownloadLatestPollenData(w http.ResponseWriter, r *http.Request) {
	api.pollenHistoryDataService.DownloadLatestPollenData()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}
