package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

type Analysis struct {
	Team     string `json:"team"`
	Core     string `json:"core"`
	Mode     string `json:"mode"`
	Coverage string `json:"coverage"`
	Support  string `json:"support"`
}

// CORS middleware function to add CORS headers
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Adjust for your frontend's origin
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

	var core, mode, coverage, support string

	team, teamToText := RunParser(link)
	//fmt.Println(teamToText)

	// Setting up a wait group for goroutines
	var wg sync.WaitGroup
	wg.Add(4)

	// Core analysis
	go func() {
		defer wg.Done()
		core = CoreReport(team)
		//fmt.Print(core)
	}()

	// Mode analysis
	go func() {
		defer wg.Done()
		mode = ModeReport(team)
		//fmt.Print(mode)
	}()

	// Coverage analysis
	go func() {
		defer wg.Done()
		coverage = CoverageReport(team)
		//fmt.Print(coverage)
	}()

	// Support analysis
	go func() {
		defer wg.Done()
		support = supportReport(team)
		//fmt.Print(support)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	res := Analysis{
		Team:     teamToText,
		Core:     core,
		Mode:     mode,
		Coverage: coverage,
		Support:  support,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return
	}

}

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	// react
	http.Handle("/analyze", enableCors(http.HandlerFunc(analyze)))
	//http.HandleFunc("/analyze", analyze)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}