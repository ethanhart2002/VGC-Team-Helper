package main

import (
	"github.com/deckarep/golang-set/v2"
	"strings"
)

/**
Strategy:

Look for common type cores in provided Pokepast.

*/

func CoreReport(team []Pokemon) string {
	fwgSet := mapset.NewSet[string]()
	fantasySet := mapset.NewSet[string]()
	PDFSet := mapset.NewSet[string]()
	for _, pokemon := range team {
		for _, typeString := range pokemon.Type {
			if typeString == "fire" || typeString == "water" || typeString == "grass" {
				fwgSet.Add(typeString)
			} else if typeString == "dragon" || typeString == "fairy" || typeString == "steel" {
				fantasySet.Add(typeString)
			} else if typeString == "psychic" || typeString == "dark" || typeString == "fighting" {
				PDFSet.Add(typeString)
			}
		}
	}
	report := strings.Builder{}
	report.WriteString("\n\nCore report \n -----------------------------\n")
	fwg := false
	dfs := false
	pdf := false
	if fwgSet.Cardinality() == 3 {
		report.WriteString("Fire-Water-Grass core detected. \n")
		fwg = true
	}
	if fantasySet.Cardinality() == 3 {
		report.WriteString("Dragon-Fairy-Steel core detected. \n")
		dfs = true
	}
	if PDFSet.Cardinality() == 3 {
		report.WriteString("Psychic-Dark-Fighting core detected. \n")
		pdf = true
	}

	if !fwg && !dfs && !pdf {
		report.WriteString("We didn't detect any common type core. Consider adding popular type cores such as fire-water-grass," +
			" dragon-fairy-steel, or psychic-dark-fighting.")
	} else if fwg && !dfs && !pdf {
		report.WriteString("Good job! Your team has a fire-water-grass core. This is one of the staple offensive " +
			"and defensive cores in Pokemon. Consider adding Pokemon to hit other popular types like Dragon and Steel. \n")
	} else if dfs && !fwg && !pdf {
		report.WriteString("Good job! Your team has a dragon-fairy-steel core. This core can combat many elemental types" +
			" such as fire, water, grass, and electric. Consider adding Pokemon with types/moves that beat the fire, poison, " +
			"and ground types that counter this core.")
	} else if pdf && !fwg && !dfs {
		report.WriteString("Good job! Your team has a psychic-dark-fighting core. Pokemon with these types have excellent " +
			"offensive synergy. Consider adding Pokemon with defensive tools to support your team, such as redirection, screens, " +
			"intimidate, or more.")
	} else if fwg && dfs && !pdf {
		report.WriteString("Excellent! Your team has both a fire-water-grass core and a dragon-fairy-steel core. You have consistent " +
			"defensive and offensive types that complement each other well.")
	} else if fwg && pdf && !dfs {
		report.WriteString("Excellent! Your team has both a fire-water-grass core and a psychic-dark-fighting core. You have consistent " +
			"defensive and offensive types that complement each other well.")
	} else if pdf && dfs && !fwg {
		report.WriteString("Excellent! your team has both a dragon-fairy-steel core and a psychic-dark-fighting core. You have consistent" +
			" resistances as well as good offensive coverage.")
	} else {
		report.WriteString("WOW! You've somehow created a team with a fire-water-grass core, a dragon-fairy-steel core, " +
			"and a psychic-dark-fighting core. You have excellent offensive and defensive coverage typing across the board.")
	}

	res := report.String()

	return res
}
