package main

import (
	"strings"
)

/**
Strategy:

Look for elements of a variety of modes in VGC and report what we find.
*/

func ModeReport(team []Pokemon) string {

	TRFlag := false
	TailwindFlag := false
	SetupFlag := false
	PsyFlag := false
	PerishFlag := false

	// For balance
	FOFlag := false
	RedirectionFlag := false
	FOCount := 0

	// Weather
	rainFlag := false
	sunFlag := false
	sandFlag := false
	snowFlag := false

	// Check for Trick Room, Tailwind, Psyspam, Perish Song, Weather, and Balance
	for _, pokemon := range team {
		if strings.Contains(pokemon.Ability, "Sand Stream") {
			sandFlag = true
		}
		if strings.Contains(pokemon.Ability, "Drizzle") {
			rainFlag = true
		}
		if strings.Contains(pokemon.Ability, "Drought") {
			sunFlag = true
		}
		if strings.Contains(pokemon.Ability, "Snow Warning") {
			snowFlag = true
		}
		for _, move := range pokemon.Moves {
			if strings.Contains(move, "Trick Room") {
				TRFlag = true
			}
			if strings.Contains(move, "Tailwind") {
				TailwindFlag = true
			}
			if strings.Contains(move, "Perish Song") {
				PerishFlag = true
			}
			if strings.Contains(move, "Expanding Force") {
				for _, poke := range team {
					if strings.Contains(poke.Ability, "Psychic Surge") {
						PsyFlag = true
					}
				}
			}

			// Balance
			if strings.Contains(move, "Follow Me") || strings.Contains(move, "Rage Powder") {
				RedirectionFlag = true
			}
			if strings.Contains(move, "Fake Out") {
				FOFlag = true
				FOCount = FOCount + 1
			}

			// Weather

			if strings.Contains(move, "Sunny Day") {
				sunFlag = true
			}
			if strings.Contains(move, "Rain Dance") {
				rainFlag = true
			}
			if strings.Contains(move, "Sandstorm") {
				sandFlag = true
			}
			if strings.Contains(move, "Snowscape") {
				snowFlag = true
			}

			// Setup
			if strings.Contains(move, "Swords Dance") || strings.Contains(move, "Nasty Plot") || strings.Contains(move, "Calm Mind") ||
				strings.Contains(move, "Bulk Up") || strings.Contains(move, "Belly Drum") || strings.Contains(move, "Decorate") ||
				strings.Contains(move, "Coaching") || strings.Contains(move, "Quiver Dance") || strings.Contains(move, "Victory Dance") ||
				strings.Contains(move, "Dragon Dance") || strings.Contains(move, "Coil") || strings.Contains(move, "Iron Defense") ||
				strings.Contains(move, "Howl") || strings.Contains(move, "Hone Claws") || strings.Contains(move, "Clangorous Soul") ||
				strings.Contains(move, "Work Up") || strings.Contains(move, "Shell Smash") || strings.Contains(move, "Curse") ||
				strings.Contains(move, "Geomancy") || strings.Contains(move, "Shift Gear") || strings.Contains(move, "Minimize") ||
				strings.Contains(move, "Baton Pass") || strings.Contains(move, "Shed Tail") { //TODO add more
				SetupFlag = true
			}
		}
	}

	report := strings.Builder{}
	report.WriteString("\n\nMode report \n -----------------------------")

	if TRFlag {
		report.WriteString("\nTR Mode detected. Make sure you have a slow Pokemon that can take advantage of Trick Room. ")
	} else {
		report.WriteString("\nWe didnt detect a Trick Room mode. If your team is utilizing a lot of slow Pokemon, consider adding the move Trick Room to your team " +
			"so that your slower Pokemon can move first.")
	}

	if TailwindFlag {
		report.WriteString("\nTailwind mode detected.")
	} else {
		report.WriteString("\nWe didn't detect a Tailwind mode. Tailwind doubles the speed of your team for four turns, which gives your Pokemon a better chance to " +
			"move first and attack. A Tailwind setter could fit your team if your Pokemon aren't naturally fast.")
	}

	if PerishFlag {
		report.WriteString("\nPerish Mode detected. This mode relies on using Perish Song to trap your opponents. Good additions to your team would be a Pokemon " +
			"with the Shadow Tag ability, the move Protect on many of your Pokemon, and Pokemon that can beat the Ghost types that Perish Song doesn't trap.")
	} else {
		report.WriteString("\nNo Perish mode detected.")
	}

	if RedirectionFlag {
		if FOFlag {
			report.WriteString("\nWe detected common elements of a Balance mode with Fake Out and Follow Me. This mode is flexible and a variety of Pokemon can " +
				"fit on Balance teams. Consider adding type cores recommended by the Core Report to round out this team.")
		} else {
			report.WriteString("\nWe detected redirection on your team, which is usually an element of a Balance team. While not necessary, Pokemon with Fake Out " +
				"synergize well on teams with redirection. Consider adding a Pokemon with Fake Out.")
		}
	} else if FOCount >= 2 {
		report.WriteString("\nWe detected common elements of a Balance mode with multiple Fake Out users. This mode is flexible and a variety of Pokemon can " +
			"fit on Balance teams. Consider adding type cores recommended by the Core Report to round out this team.")
	} else {
		report.WriteString("\nNo Balance mode detected. Although Balance is an amorphous term in VGC, Balance teams commonly have Pokemon with the moves Fake Out " +
			"or Follow Me/Rage Powder. If you have a strong core of Pokemon with type synergy, consider adding these popular Balance elements.")
	}

	if PsyFlag {
		report.WriteString("\nWe detected a Psyspam mode on your team. Psyspam can excel at spreading damage fast and taking quick knockouts. Make sure your team " +
			"has Pokemon and/or moves that can beat Dark types and types that resist Psychic moves.")
	} else {
		report.WriteString("\nNo Psyspam mode detected. If you are looking for a fast offensive mode that can take knockouts quickly, Psyspam (Psychic Terrain + Expanding Force) " +
			"could be a good mode to add.")
	}

	if SetupFlag {
		report.WriteString("\nWe detected Setup strategies on your team. Consider adding supportive Pokemon, such as Pokemon with Intimidate, Fake Out, or Follow Me/Rage Powder to help " +
			"your setup Pokemon do their job.")
	} else {
		report.WriteString("\nNo setup detected. Using moves to boost your stats can be a powerful mode that will increase your damage output " +
			"or decrease your opponent's damage output. Consider adding a Pokemon that can boost its stats.")
	}

	if sunFlag {
		report.WriteString("\nSun mode detected. Consider adding Pokemon that benefit from the sun, such as Fire types and Pokemon with abilities affected " +
			"by the sun such as Chlorophyll.")
	}

	if rainFlag {
		report.WriteString("\nRain mode detected. Consider adding Pokemon that benefit from the sun, such as Water, Grass, Bug, and Steel types. Pokemon with abilities such as " +
			"Swift Swim also benefit from the sun.")
	}

	if sandFlag {
		report.WriteString("\nSand mode detected. Consider adding Pokemon that benefit from the sand, such as Rock, Ground, and Steel Types. Pokemon with abilities " +
			"such as Sand Veil and Sand Rush also benefit from the sand.")
	}

	if snowFlag {
		report.WriteString("\nSnow mode detected. Consider adding Pokemon that benefit from the snow like Ice types. Pokemon with abilities such as " +
			"Slush Rush and Snow Cloak also benefit from the snow.")
	}

	res := report.String()

	return res
}