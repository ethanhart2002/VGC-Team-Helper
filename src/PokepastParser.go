package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Pokemon struct {
	Name     string   `json:"pokemon"`
	Type     []string `json:"type"`
	Item     string   `json:"item"`
	Ability  string   `json:"ability"`
	Level    int      `json:"level"`
	TeraType string   `json:"tera_type"`
	EVs      string   `json:"evs"`
	Nature   string   `json:"nature"`
	Moves    []string `json:"moves"`
}

// Type contains a single type information
type Type struct {
	Type InnerType `json:"type"`
}

// InnerType holds the name of the type
type InnerType struct {
	Name string `json:"name"`
}

// PokemonTypeResponse matches only the types field from the PokeAPI response
type PokemonTypeResponse struct {
	Types []Type `json:"types"`
}

// Maps Pokemon Showdown naming conventions to PokeAPI naming conventions
var nameMap = map[string]string{
	"indeedee":              "indeedee-male",
	"indeedee-f":            "indeedee-female",
	"maushold":              "maushold-family-of-three",
	"maushold-four":         "maushold-family-of-four",
	"dudunsparce":           "dudunsparce-two-segment",
	"dudunsparce-three":     "dudunsparce-three-segment",
	"walking wake":          "walking-wake",
	"iron leaves":           "iron-leaves",
	"raging bolt":           "raging-bolt",
	"gouging fire":          "gouging-fire",
	"iron boulder":          "iron-boulder",
	"iron crown":            "iron-crown",
	"roaring moon":          "roaring-moon",
	"iron valiant":          "iron-valiant",
	"iron treads":           "iron-treads",
	"iron bundle":           "iron-bundle",
	"iron hands":            "iron-hands",
	"iron jugulis":          "iron-jugulis",
	"iron moth":             "iron-moth",
	"iron thorns":           "iron-thorns",
	"great tusk":            "great-tusk",
	"scream tail":           "scream-tail",
	"brute bonnet":          "brute-bonnet",
	"flutter mane":          "flutter-mane",
	"slither wing":          "slither-wing",
	"sandy shocks":          "sandy-shocks",
	"tauros-paldea-aqua":    "tauros-paldea-aqua-breed",
	"tauros-paldea-blaze":   "tauros-paldea-blaze-breed",
	"tauros-paldea-combat":  "tauros-paldea-combat-breed",
	"sinistcha-masterpiece": "sinistcha",
	"ogerpon-wellspring":    "ogerpon-wellspring-mask",
	"ogerpon-hearthflame":   "ogerpon-hearthflame-mask",
	"ogerpon-cornerstone":   "ogerpon-cornerstone-mask",
	"tornadus":              "tornadus-incarnate",
	"thundurus":             "thundurus-incarnate",
	"landorus":              "landorus-incarnate",
	"emamorus":              "enamorus-incarnate",
	"basculegion":           "basculegion-male",
	"basculegion-f":         "basculegion-female",
	//casting all the colors of minior doesn't matter, we only need the type
	"minior":    "minior-red",
	"mimikyu":   "mimikyu-disguised",
	"tatsugiri": "tatsugiri-curly",
	"urshifu":   "urshifu-single-strike",
	//TODO
}

// Function to fetch Pokémon data
func fetchPokemonData(pokemonName string) (*http.Response, error) {
	//Attempt 1
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Check if the response was successful. If it wasn't, do a name lookup conversion since PokeAPI has some different naming conventions than Pokemon Showdown & Pokepaste.
	if response.StatusCode != http.StatusOK {
		err := response.Body.Close()
		if err != nil {
			return nil, err
		}
		replacementName, valid := nameMap[pokemonName]
		if valid {
			url = fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", replacementName)
			response2, err2 := http.Get(url)
			if err2 != nil {
				return nil, err2
			}
			if response2.StatusCode != http.StatusOK {
				err := response2.Body.Close()
				if err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("failed to fetch type data for the pokemon named %s; status code =  %d", replacementName, response.StatusCode)
			} else {
				return response2, nil
			}
		} else {
			return nil, fmt.Errorf("Pokemon named %s was not found in the name formatting map; status code =  %d", pokemonName, response.StatusCode)
		}
	} else {
		return response, nil
	}
}

// GetPokemonType fetches and returns the types of a given Pokémon
func GetPokemonType(pokemonName string) ([]string, error) {

	// Send GET request through a special function
	response, err := fetchPokemonData(pokemonName)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(fmt.Errorf("error closing response body: %v", err))
		}
	}(response.Body)

	// Check if the response was successful
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch type data for the pokemon named %s; status code =  %d", pokemonName, response.StatusCode)
	}

	// Parse JSON into a struct that only captures the types
	var typeResponse PokemonTypeResponse
	if err := json.NewDecoder(response.Body).Decode(&typeResponse); err != nil {
		return nil, err
	}

	// Extract type names into a slice
	var types []string
	for _, t := range typeResponse.Types {
		types = append(types, t.Type.Name)
	}
	return types, nil
}

func parse(pokeInfo string, flag bool) Pokemon {

	//TODO: name regex needs to factor in for a space
	nameItemRegex := regexp.MustCompile(`^(\w+(?:[ -]\w+)*)(\s(\(M\)|\(F\)))* @ (.+)`)
	nameItemRegexFromNickname := regexp.MustCompile(`\((\w+(?:[ -]\w+)*)\)(\s(\(M\)|\(F\)))* @* (.+)*`)
	nameRegex := regexp.MustCompile(`^(\w+(?:[ -]\w+)*)(\s(\(M\)|\(F\)))*`)
	abilityRegex := regexp.MustCompile(`Ability: (.+)`)
	levelRegex := regexp.MustCompile(`Level: (\d+)`)
	teraTypeRegex := regexp.MustCompile(`Tera Type: (\w+)`)
	evsRegex := regexp.MustCompile(`EVs: (.+)`)
	movesRegex := regexp.MustCompile(`-(\s.+)`)
	natureRegex := regexp.MustCompile(`(\w+) Nature`)

	// Initialize a new Pokemon instance
	var p Pokemon

	// Parse name and item, or just name if no item is listed
	if matches := nameItemRegex.FindStringSubmatch(pokeInfo); len(matches) > 2 {
		p.Name = matches[1]
		p.Item = matches[len(matches)-1]
	} else if matches := nameRegex.FindStringSubmatch(pokeInfo); len(matches) > 1 {
		p.Name = matches[1]
	} else {
		// Parse in case for a nickname
		if matches := nameItemRegexFromNickname.FindStringSubmatch(pokeInfo); len(matches) > 2 {
			p.Name = matches[1]
			p.Item = matches[len(matches)-1]
			// No item
		} else if matches := nameRegex.FindStringSubmatch(pokeInfo); len(matches) > 1 {
			p.Name = matches[1]
		} else {
			log.Panicln("Pokemon name and item could not be parsed.")
		}
	}

	// Parse type(s)
	types, err := GetPokemonType(strings.ToLower(p.Name))
	if err != nil {
		p.Type = append(p.Type, "Error: Couldn't lookup the type for this Pokemon.")
	} else {
		p.Type = types
	}
	// Parse ability
	if matches := abilityRegex.FindStringSubmatch(pokeInfo); len(matches) > 1 {
		p.Ability = matches[1]
	}

	// Parse level
	if matches := levelRegex.FindStringSubmatch(pokeInfo); len(matches) > 1 {
		fmt.Sscanf(matches[1], "%d", &p.Level)
	}

	// Parse tera type
	if matches := teraTypeRegex.FindStringSubmatch(pokeInfo); len(matches) > 1 {
		p.TeraType = matches[1]
	}

	// Parse moves
	moves := movesRegex.FindAllStringSubmatch(pokeInfo, -1)
	for _, match := range moves {
		if len(match) >= 1 {
			p.Moves = append(p.Moves, match[1])
		}
	}

	//CTS
	if !flag {
		// Parse EVs
		if matches := evsRegex.FindStringSubmatch(pokeInfo); len(matches) > 1 {
			p.EVs = matches[1]
		}

		// Parse nature
		if matches := natureRegex.FindStringSubmatch(pokeInfo); len(matches) > 1 {
			p.Nature = matches[1]
		}
	}

	return p
}

func RunParser(path string) ([]Pokemon, string) {
	var OtsFlag = false
	var team []Pokemon
	var jsonTxt string

	//debugging line:
	//jsonTxt = jsonTxt + "\nYour team:\n" + "--------------------------------\n"

	// Get the Pokepaste
	req, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(req.Body)

	//initialize goquery
	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Get the OTS/CTS flag
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "(OTS)") {
			OtsFlag = true
		}
	})
	// Parse the team
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		var member Pokemon
		pokeInfo := s.Find("pre").Text()
		member = parse(pokeInfo, OtsFlag)
		team = append(team, member)
		x, err := json.MarshalIndent(member, "", " ")
		if err != nil {
			log.Panicln(err)
		}

		jsonTxt = jsonTxt + string(x) + "\n" + "\n"

	})

	return team, jsonTxt

}
