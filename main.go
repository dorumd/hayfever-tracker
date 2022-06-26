package main

import (
	"fmt"
	pollen_history_data "hayfever/tracker/pollen-history-data"
	self_assessment "hayfever/tracker/self-assessment"
	"net/http"
)

func main() {
	pollenHistoryDataService := pollen_history_data.New(
		pollen_history_data.NewRepositoryUsingFiles(),
	)
	selfAssessmentService := self_assessment.New(
		self_assessment.NewMySQLRepository(),
	)
	api := NewApi(pollenHistoryDataService, selfAssessmentService)

	http.HandleFunc("/api/pollen-data", api.GetPollenHistoryDataForToday)
	http.HandleFunc("/api/submit-self-assessment", api.SubmitSelfAssessment)
	http.HandleFunc("/api/download-latest-pollen-data", api.DownloadLatestPollenData)

	srv := &http.Server{
		Addr: ":8080",
	}
	fmt.Println("Hay Fever Tracker - Ready")
	srv.ListenAndServe()
}
