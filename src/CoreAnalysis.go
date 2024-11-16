package main

import (
	"github.com/deckarep/golang-set/v2"
	"strings"
)

/**
Strategy:

Look for common type cores in provided Pokepast.

*/

func CoreReport(team []Pokemon) (string, float64) {
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
	//report.WriteString("\n\nCore report \n -----------------------------\n")
	fwg := false
	dfs := false
	pdf := false

	partFwg := false
	partDFS := false
	partPdf := false

	if fwgSet.Cardinality() == 3 {
		report.WriteString("\nFire-Water-Grass core detected. ")
		fwg = true
	} else if fwgSet.Cardinality() == 2 {
		report.WriteString("\nTwo-thirds of a Fire-Water-Grass core detected.")
		partFwg = true
	}
	if fantasySet.Cardinality() == 3 {
		report.WriteString("\nDragon-Fairy-Steel core detected. ")
		dfs = true
	} else if fantasySet.Cardinality() == 2 {
		partDFS = true
	}
	if PDFSet.Cardinality() == 3 {
		report.WriteString("\nPsychic-Dark-Fighting core detected. ")
		pdf = true
	} else if PDFSet.Cardinality() == 2 {
		partPdf = true
	}

	var score float64

	/**
	Grading strategy: Assign a base score for COMPLETED CORES. Cores are weighted FWG > DFS > PDF based on historical data.
	For each supplementary core that is a partial core (missing 1 member), add a point.
	*/

	if fwg && !dfs && !pdf {
		report.WriteString("\nGood job! Your team has a fire-water-grass core. This is one of the staple offensive " +
			"and defensive cores in Pokemon. Consider adding Pokemon to hit other popular types like Dragon and Steel. ")
		score = 8.0
	} else if dfs && !fwg && !pdf {
		report.WriteString("\nGood job! Your team has a dragon-fairy-steel core. This core can combat many elemental types" +
			" such as fire, water, grass, and electric. Consider adding Pokemon with types/moves that beat the fire, poison, " +
			"and ground types that counter this core.")
		score = 7.8
	} else if pdf && !fwg && !dfs {
		report.WriteString("\nGood job! Your team has a psychic-dark-fighting core. Pokemon with these types have excellent " +
			"offensive synergy. Consider adding Pokemon with defensive tools to support your team, such as redirection, screens, " +
			"intimidate, or more.")
		score = 7.6
	} else if fwg && dfs && !pdf {
		report.WriteString("\nExcellent! Your team has both a fire-water-grass core and a dragon-fairy-steel core. You have consistent " +
			"defensive and offensive types that complement each other well.")
		score = 9.0
	} else if fwg && pdf && !dfs {
		report.WriteString("\nExcellent! Your team has both a fire-water-grass core and a psychic-dark-fighting core. You have consistent " +
			"defensive and offensive types that complement each other well.")
		score = 9.0
	} else if pdf && dfs && !fwg {
		report.WriteString("\nExcellent! your team has both a dragon-fairy-steel core and a psychic-dark-fighting core. You have consistent" +
			" resistances as well as good offensive coverage.")
		score = 9.0
	} else if pdf && dfs && fwg {
		report.WriteString("\nWOW! You've somehow created a team with a fire-water-grass core, a dragon-fairy-steel core, " +
			"and a psychic-dark-fighting core. You have excellent offensive and defensive coverage typing across the board.")
		score = 10.0

	}

	//partial fwg
	if partFwg && !partDFS && !partPdf {
		report.WriteString("\nYou have 2 of the 3 types required for a fire-water-grass core. Consider adding the missing type " +
			"to complete the core.")
		if score != 0 {
			score = score + 1.0
		} else {
			score = 6.7
		}

		//partial dfs
	} else if !partFwg && partDFS && !partPdf {
		report.WriteString("\nYou have 2 of the 3 types required for a dragon-fairy-steel core. Consider adding the missing type " +
			"to complete the core.")
		if score != 0 {
			score = score + 1.0
		} else {
			score = 6.5
		}
		//partial pdf
	} else if !partFwg && !partDFS && partPdf {
		report.WriteString("\nYou have 2 of the 3 types required for a psychic-dark-fighting core. Consider adding the missing type " +
			"to complete the core.")
		if score != 0 {
			score = score + 1.0
		} else {
			score = 6.0
		}

		//partial fwg and dfs
	} else if partFwg && partDFS && !partPdf {
		report.WriteString("\nYou have 2 of the 3 types required for a fire-water-grass core and for a dragon-fairy-steel core. Consider adding the missing types " +
			"to complete the core.")
		if score != 0 {
			score = score + 2.0
		} else {
			score = 7.5
		}

		//partial dfs and pdf
	} else if !partFwg && partDFS && partPdf {
		report.WriteString("\nYou have 2 of the 3 types required for a psychic-dark-fighting core and for a dragon-fairy-steel core. Consider adding the missing types " +
			"to complete the core.")
		if score != 0 {
			score = score + 2.0
		} else {
			score = 7.5
		}

		//partial fwg and pdf
	} else if partFwg && !partDFS && partPdf {
		report.WriteString("\nYou have 2 of the 3 types required for a fire-water-grass core and for a psychic-dark-fighting core. Consider adding the missing types " +
			"to complete the core.")
		if score != 0 {
			score = score + 2.0
		} else {
			score = 7.5
		}

		//all 3 partials
	} else if partFwg && partDFS && partPdf {
		report.WriteString("\nYou have 2 of the 3 types required for a fire-water-grass core, a psychic-dark-fighting core, and a dragon-fairy-steel core. Consider adding the missing types " +
			"to complete the core.")
		if score != 0 {
			score = score + 2.0
		} else {
			score = 8.0
		}
	} else {
		report.WriteString("\nWe didn't detect any common type core. Consider adding popular type cores such as fire-water-grass," +
			" dragon-fairy-steel, or psychic-dark-fighting.")
		if score == 0 {
			score = 4.0
		}
	}

	res := report.String()

	return res, score
}
