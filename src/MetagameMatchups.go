package main

import (
	"bufio"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
	"math"
	"os"
	"strings"
)

/*
*
Strategy for this analysis file:
 1. Pull usage stats generated by smogon and pulled by separate python file
    - 1 Pokemon per line in usage stats text file, line# matches rank of usage
 2. For each Pokemon on the inputted team, check to see what pokemon in the top __ of usage they struggle to hit
 3. For each Pokemon on the inputted team, check to see what pokemon in the top __ of usage they struggle against defensively
 4. Return
*/
func calculateScore(goodMUlen int, okMUlen int, badMUlen int) float64 {
	goodWeight := 1.0
	okWeight := 0.5
	badWeight := -0.3

	score := 10 * (((float64(goodMUlen) * goodWeight) + (float64(okMUlen) * okWeight) + (float64(badMUlen) * badWeight)) / 50)
	return math.Max(0, math.Min(10, score))
}

func MetagameMatchups(teamCoverage mapset.Set[string], foundTypesFrequency map[string]int) ([]string, []string, []string, float64) {

	topPokemonUsage := [50]string{}

	//TODO
	file, err := os.Open("./usage/gen9vgc2025regi-0-2025-05-Usage.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		topPokemonUsage[lineCount] = line
		lineCount++
		if lineCount > 50 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// for each meta pokemon, get their type(s), convert it similarly to coverage file from string to number
	// for each move of each Pokemon in the inputted team, check against each meta Pokemon to see if you can hit it super effectively, neutral, or not effectively
	// store the bad matchups and return them

	//array for pokemon the team has a good MU into
	var goodMatchups []string
	//array for pokemon the team has a okay MU into
	var okMatchups []string
	//array for pokemon the team has a bad MU into
	var badMatchups []string

	for i := range topPokemonUsage {
		//get their type from API call
		//TODO may have to do some conversion checking for names
		//print("\n\nCurrent Pokemon: " + topPokemonUsage[i] + "\n")
		metaPokemonType, err := GetPokemonType(strings.ToLower(topPokemonUsage[i]))
		var metaPokemonTypeIndices []int
		for _, typeString := range metaPokemonType {
			metaPokemonTypeIndices = append(metaPokemonTypeIndices, typeMapGetIndex[typeString])
		}
		if err != nil {
			log.Fatal(err)
		}
		effectiveCounter := 0.0

		/**
		-loop through typesfound
		-if there is a type in types found that hits the current pokemon super effectively, add to a counter
		*/

		for teamCoverageType := range teamCoverage.Iterator().C {
			//fmt.Printf("Inspecting the %s type against %s \n", teamCoverageType, topPokemonUsage[i])
			coverageTypeIndex := typeMapGetIndex[teamCoverageType]
			// case that opposing Pokemon has 1 type
			if len(metaPokemonTypeIndices) == 1 {
				if typeChart[coverageTypeIndex][metaPokemonTypeIndices[0]] == 2 {
					//fmt.Printf("The %s type appears %d times on the team, and is super effective against the meta Pokemon %s \n", teamCoverageType, foundTypesFrequency[teamCoverageType], topPokemonUsage[i])
					effectiveCounter += float64(foundTypesFrequency[teamCoverageType])
				}
				// case that opposing Pokemon has 2 types
			} else {
				metaMonFirstType := metaPokemonTypeIndices[0]
				metaMonSecondType := metaPokemonTypeIndices[1]
				combinedEffectiveness := typeChart[coverageTypeIndex][metaMonFirstType] * typeChart[coverageTypeIndex][metaMonSecondType]
				if combinedEffectiveness > 1 {
					//fmt.Printf("The %s type appears %d times on the team, and is super effective against the meta Pokemon %s \n", teamCoverageType, foundTypesFrequency[teamCoverageType], topPokemonUsage[i])
					effectiveCounter += float64(foundTypesFrequency[teamCoverageType])
				}
			}
		}

		//fmt.Printf("Final effective counter for %s: %f \n", topPokemonUsage[i], effectiveCounter)
		//if counter is less than 1, the team has a bad matchup into current pokemon. if counter is 1, team has okay matchup. if counter is 2+, team has good matchup.
		if effectiveCounter < 1 {
			badMatchups = append(badMatchups, topPokemonUsage[i])
		} else if effectiveCounter == 1 {
			okMatchups = append(okMatchups, topPokemonUsage[i])
		} else if effectiveCounter > 1 {
			goodMatchups = append(goodMatchups, topPokemonUsage[i])
		}
	}

	// TODO Defensive MU

	// TODO score
	var score float64
	score = calculateScore(len(goodMatchups), len(okMatchups), len(badMatchups))

	//print("\n\nGood Matchups: \n")
	//for i := range goodMatchups {
	//	print(goodMatchups[i] + ", ")
	//}
	//print("\n\nOK Matchups: \n")
	//for i := range okMatchups {
	//	print(okMatchups[i] + ", ")
	//}
	//print("\n\nBad Matchups: \n")
	//for i := range badMatchups {
	//	print(badMatchups[i] + ", ")
	//}

	return goodMatchups, okMatchups, badMatchups, score
}
