package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {

	path := os.Args[1]
	if len(os.Args) <= 0 || len(os.Args) > 2 {
		fmt.Println("Program takes a single parameter which should be a Pokepaste link.")
		return
	} else if len(path) == 0 {
		fmt.Println("No valid Pokepaste provided.")
		return
	} else if !strings.Contains(path, "https://pokepast.es/") {
		fmt.Println("Link provided is not a Pokepaste resource. Pokepast domain: https://pokepast.es/")
		return
	} else if res, err := http.Get(path); err != nil || res.StatusCode != 200 {
		fmt.Println("Error fetching Pokepaste. No resource at this link. Error:", err)
		return
	}

	team, teamToText := RunParser(path)
	fmt.Println(teamToText)

	// Setting up a wait group for goroutines
	var wg sync.WaitGroup
	wg.Add(4)

	// Core analysis
	go func() {
		defer wg.Done()
		fmt.Print(CoreReport(team))
	}()

	// Mode analysis
	go func() {
		defer wg.Done()
		fmt.Print(ModeReport(team))
	}()

	// Coverage analysis
	go func() {
		defer wg.Done()
		fmt.Print(CoverageReport(team))
	}()

	// Support analysis
	go func() {
		defer wg.Done()
		fmt.Print(supportReport(team))
	}()

	// Wait for all goroutines to finish
	wg.Wait()

}
