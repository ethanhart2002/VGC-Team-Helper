package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Analysis struct {
	//Team     string  `json:"team"`
	Team     []Pokemon `json:"team"`
	Core     string    `json:"core"`
	Mode     string    `json:"mode"`
	Coverage []string  `json:"coverage"`
	Support  string    `json:"support"`
	Score    float64   `json:"score"`
}

// CORS middleware function to add CORS headers
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080") // Adjust for your frontend's origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func analyze(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Link string `json:"link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	link := request.Link

	if len(link) == 0 || !strings.Contains(link, "https://pokepast.es/") {
		http.Error(w, "Invalid link provided.", http.StatusBadRequest)
		return
	}

	if res, err := http.Get(link); err != nil || res.StatusCode != 200 {
		http.Error(w, "Error fetching Pokepaste resource. No resource at this link.", http.StatusBadRequest)
		return
	}

	var core, mode, support string
	var coverage []string
	var coreScore, modeScore, coverageScore, suppScore float64

	team, _ := RunParser(link)

	/**
	Debugging commented out below
	*/

	//_, teamToText := RunParser(link)
	//fmt.Println(teamToText)

	// Setting up a wait group for goroutines
	var wg sync.WaitGroup
	wg.Add(4)

	// Core analysis
	go func() {
		defer wg.Done()
		core, coreScore = CoreReport(team)
	}()

	// Mode analysis
	go func() {
		defer wg.Done()
		mode, modeScore = ModeReport(team)
	}()

	// Coverage analysis
	go func() {
		defer wg.Done()
		coverage, coverageScore = CoverageReport(team)
	}()

	// Support analysis
	go func() {
		defer wg.Done()
		support, suppScore = SupportReport(team)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	//Debugging print for scores commented out below
	//fmt.Printf("core: %d, mode: %d, coverage: %d, supp: %d", coreScore, modeScore, coverageScore, suppScore)

	s := coreScore*.3 + modeScore*.3 + coverageScore*.2 + suppScore*.2
	total := fmt.Sprintf("%.2f", s)
	totalScore, err := strconv.ParseFloat(total, 64)
	if err != nil {
		totalScore = 0
		e := fmt.Errorf("error parsing total analysis score, returning a placeholder 0/10 %s", err.Error())
		log.Panicln(e)
	}

	res := Analysis{
		Team:     team,
		Core:     core,
		Mode:     mode,
		Coverage: coverage,
		Support:  support,
		Score:    totalScore,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return
	}

}

func main() {
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)
	http.HandleFunc("/analyze", analyze)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
