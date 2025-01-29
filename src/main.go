package main

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Analysis struct {
	Team     []Pokemon `json:"team"`
	Core     string    `json:"core"`
	Mode     string    `json:"mode"`
	Coverage []string  `json:"coverage"`
	Support  string    `json:"support"`
	Meta     Meta      `json:"meta_matchups"`
	Score    float64   `json:"score"`
}

type Meta struct {
	GoodMU []string `json:"goodMU"`
	OkMU   []string `json:"okMU"`
	BadMU  []string `json:"badMU"`
}

// CORS middleware function to add CORS headers
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
		w.Header().Set("Access-Control-Allow-Origin", "https://vgcteamhelper.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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
	var missingCoverageTypes []string
	var foundCoverageTypes mapset.Set[string]
	var foundTypesFrequency map[string]int
	var coreScore, modeScore, coverageScore, suppScore, metaScore float64
	var goodMU, okMU, badMU []string

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

	// Coverage analysis and Metagame Matchup analysis
	go func() {
		defer wg.Done()
		missingCoverageTypes, foundCoverageTypes, foundTypesFrequency, coverageScore = CoverageReport(team)
		goodMU, okMU, badMU, metaScore = MetagameMatchups(foundCoverageTypes, foundTypesFrequency)
	}()

	// Support analysis
	go func() {
		defer wg.Done()
		support, suppScore = SupportReport(team)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	//Debugging print for scores commented out below
	fmt.Printf("\n core: %.2f, mode: %.2f, coverage: %.2f, supp: %.2f, meta: %.2f \n", coreScore, modeScore, coverageScore, suppScore, metaScore)

	//Score calculation; Each section has 20% weight
	s := coreScore*.2 + modeScore*.2 + coverageScore*.2 + suppScore*.2 + metaScore*.2
	total := fmt.Sprintf("%.2f", s)
	totalScore, err := strconv.ParseFloat(total, 64)
	if err != nil {
		totalScore = 0
		e := fmt.Errorf("error parsing total analysis score, returning a placeholder 0/10 %s", err.Error())
		log.Panicln(e)
	}

	metaReport := Meta{
		GoodMU: goodMU,
		OkMU:   okMU,
		BadMU:  badMU,
	}

	res := Analysis{
		Team:     team,
		Core:     core,
		Mode:     mode,
		Coverage: missingCoverageTypes,
		Support:  support,
		Meta:     metaReport,
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

	http.Handle("/", enableCors(fs))

	http.HandleFunc("/analyze", analyze)

	//For running locally, go to script.js to uncomment the local debug line that reads 'hostPath = "http://localhost:443/analyze";'
	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
		return
	}

}
