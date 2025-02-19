package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deckarep/golang-set/v2"
	"io"
	"log"
	"net/http"
	"strings"
)

var typeChart = [][]float64{
	{1, 1, 1, 1, 1, 0.5, 1, 0, 0.5, 1, 1, 1, 1, 1, 1, 1, 1, 1},           // normal 0
	{2, 1, 0.5, 0.5, 1, 2, 0.5, 0, 2, 1, 1, 1, 1, 0.5, 2, 1, 2, 0.5},     // fighting 1
	{1, 2, 1, 1, 1, 0.5, 2, 1, 0.5, 1, 1, 2, 0.5, 1, 1, 1, 1, 1},         // flying 2
	{1, 1, 1, 0.5, 0.5, 0.5, 1, 0.5, 0, 1, 1, 2, 1, 1, 1, 1, 1, 2},       // poison 3
	{1, 1, 0, 2, 1, 2, 0.5, 1, 2, 2, 1, 0.5, 2, 1, 1, 1, 1, 1},           // ground 4
	{1, 0.5, 2, 1, 0.5, 1, 2, 1, 0.5, 2, 1, 1, 1, 1, 2, 1, 1, 1},         // rock 5
	{1, 0.5, 0.5, 0.5, 1, 1, 1, 0.5, 0.5, 0.5, 1, 2, 1, 2, 1, 1, 2, 0.5}, // bug 6
	{0, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 0.5, 1},             // ghost 7
	{1, 1, 1, 1, 1, 2, 1, 1, 0.5, 0.5, 0.5, 1, 0.5, 1, 2, 1, 1, 2},       // steel  8
	{1, 1, 1, 1, 1, 0.5, 2, 1, 2, 0.5, 0.5, 2, 1, 1, 2, 0.5, 1, 1},       // fire 9
	{1, 1, 1, 1, 2, 2, 1, 1, 1, 2, 0.5, 0.5, 1, 1, 1, 0.5, 1, 1},         // water 10
	{1, 1, 0.5, 0.5, 2, 2, 0.5, 1, 0.5, 0.5, 2, 0.5, 1, 1, 1, 0.5, 1, 1}, // grass 11
	{1, 1, 2, 1, 0, 1, 1, 1, 1, 1, 2, 0.5, 0.5, 1, 1, 0.5, 1, 1},         // electric 12
	{1, 2, 1, 2, 1, 1, 1, 1, 0.5, 1, 1, 1, 1, 0.5, 1, 1, 0, 1},           // psychic 13
	{1, 1, 2, 1, 2, 1, 1, 1, 0.5, 0.5, 0.5, 2, 1, 1, 0.5, 2, 1, 1},       // ice 14
	{1, 1, 1, 1, 1, 1, 1, 1, 0.5, 1, 1, 1, 1, 1, 1, 2, 1, 0},             // dragon 15
	{1, 0.5, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 0.5, 0.5},         // dark 16
	{1, 2, 1, 0.5, 1, 1, 1, 1, 0.5, 0.5, 1, 1, 1, 1, 1, 2, 2, 1},         // fairy 17
}

var typeMapGetString = map[int]string{
	0:  "normal",
	1:  "fighting",
	2:  "flying",
	3:  "poison",
	4:  "ground",
	5:  "rock",
	6:  "bug",
	7:  "ghost",
	8:  "steel",
	9:  "fire",
	10: "water",
	11: "grass",
	12: "electric",
	13: "psychic",
	14: "ice",
	15: "dragon",
	16: "dark",
	17: "fairy",
}

var typeMapGetIndex = map[string]int{
	"normal":   0,
	"fighting": 1,
	"flying":   2,
	"poison":   3,
	"ground":   4,
	"rock":     5,
	"bug":      6,
	"ghost":    7,
	"steel":    8,
	"fire":     9,
	"water":    10,
	"grass":    11,
	"electric": 12,
	"psychic":  13,
	"ice":      14,
	"dragon":   15,
	"dark":     16,
	"fairy":    17,
}

type Move struct {
	Name        string `json:"name"`
	DamageClass struct {
		Name string `json:"name"`
	} `json:"damage_class"`
	MoveType struct {
		Name string `json:"name"`
	} `json:"type"`
}

/*
*Need to translate move to corresponding syntax for PokeAPI
 */
func moveTranslate(move string) string {
	s := strings.TrimSpace(move)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

/**
Get the move class and type from Pokeapi
*/

func getMoveData(moveName string, p Pokemon) (DamageClass string, MoveType string, err error) {

	//Handling of special case moves that are changed by tera type for example

	if moveName == "tera-blast" {
		return "special", strings.ToLower(p.TeraType), nil
	} else if moveName == "ivy-cudgel" {
		if p.Name == "Ogerpon-Wellspring" {
			return "physical", "water", nil
		} else if p.Name == "Ogerpon-Cornerstone" {
			return "physical", "rock", nil
		} else if p.Name == "Ogerpon-Hearthflame" {
			return "physical", "fire", nil
		} else {
			return "physical", "grass", nil
		}
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/move/%s/", moveName)

	response, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(fmt.Errorf("error closing response body: %v", err))
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {

		return "", "", errors.New(response.Status)
	}

	// Decode the JSON response
	var move Move
	if err := json.NewDecoder(response.Body).Decode(&move); err != nil {
		return "", "", errors.New("couldn't decode json")
	}

	return move.DamageClass.Name, move.MoveType.Name, nil

}

/**
Strategy:

	-for each pokemon, get each move.

	-for each move, check to see if it is an attacking move. If it is, get the type of the move.

	-for each typed move, check to see what types it hits super effectively, and store the
	types it hits in a set global to this function

*/

func CoverageReport(team []Pokemon) ([]string, mapset.Set[string], map[string]int, float64) {

	typesHitSuperEffectively := mapset.NewSet[string]()
	allTypes := mapset.NewSet[string](
		"normal",
		"fighting",
		"flying",
		"poison",
		"ground",
		"rock",
		"bug",
		"ghost",
		"steel",
		"fire",
		"water",
		"grass",
		"electric",
		"psychic",
		"ice",
		"dragon",
		"dark",
		"fairy",
	)

	foundCoverageTypes := mapset.NewSet[string]()
	var foundTypeFrequency = map[string]int{
		"normal":   0,
		"fighting": 0,
		"flying":   0,
		"poison":   0,
		"ground":   0,
		"rock":     0,
		"bug":      0,
		"ghost":    0,
		"steel":    0,
		"fire":     0,
		"water":    0,
		"grass":    0,
		"electric": 0,
		"psychic":  0,
		"ice":      0,
		"dragon":   0,
		"dark":     0,
		"fairy":    0,
	}

	for _, pokemon := range team {
		for _, moveString := range pokemon.Moves {
			move := moveTranslate(moveString)
			dmgClass, moveType, err := getMoveData(strings.ToLower(move), pokemon)
			if err != nil {
				fmt.Errorf("unable to get move data for move called %s, Error: %s", move, err)
			}
			// do type checking
			if dmgClass == "physical" || dmgClass == "special" {
				typeNum, valid := typeMapGetIndex[strings.ToLower(moveType)]
				if !valid {
					fmt.Errorf("Couldn't find move: %s  in map", moveType)
				} else {
					foundCoverageTypes.Add(strings.ToLower(moveType))
					foundTypeFrequency[strings.ToLower(moveType)]++
					typeEffectiveness := typeChart[typeNum]
					for i, num := range typeEffectiveness {
						if num == 2 {
							s := typeMapGetString[i]
							typesHitSuperEffectively.Add(s)
						}
					}
				}
			}
		}
	}

	//create a set that will contain types the team can't hit super effectively.
	difference := allTypes.SymmetricDifference(typesHitSuperEffectively)

	s := strings.Builder{}

	var score float64

	var missingTypes []string

	if difference.Cardinality() == 0 {
		s.WriteString("\nYour team has coverage options to hit all 18 types!")
		score = 10
		return missingTypes, foundCoverageTypes, foundTypeFrequency, score
	} else {
		for missingType := range difference.Iter() {
			missingTypes = append(missingTypes, missingType)
		}

		/**
		Grading strategy: score by taking off .55 points for every type that the team can't hit. .55 originates from this calc's
		grading scale being out of 10, which is divided by 18 types in the game.
		*/

		score = 10 - (float64(difference.Cardinality()) * .55)

		return missingTypes, foundCoverageTypes, foundTypeFrequency, score

	}

}
