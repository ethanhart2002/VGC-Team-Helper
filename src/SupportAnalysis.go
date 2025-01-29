package main

import (
	"strings"
)

/**
Strategy:

Let's check for common moves in VGC and report what we find.
*/

func SupportReport(team []Pokemon) (string, float64) {
	fakeOutFlag := false
	redirectionFlag := false
	screensFlag := false
	protectCount := 0
	speedControlFlag := false
	statusFlag := false

	totalMoveCount := 0
	supportMoveCount := 0
	for _, pokemon := range team {
		for _, move := range pokemon.Moves {
			totalMoveCount++
			if strings.Contains(move, "Fake Out") {
				fakeOutFlag = true
				supportMoveCount++
			}
			if strings.Contains(move, "Rage Powder") || strings.Contains(move, "Follow Me") {
				redirectionFlag = true
				supportMoveCount++
			}
			if strings.Contains(move, "Light Screen") || strings.Contains(move, "Reflect") ||
				strings.Contains(move, "Aurora Veil") {
				screensFlag = true
				supportMoveCount++
			}
			if strings.Contains(move, "Protect") || strings.Contains(move, "Spiky Shield") ||
				strings.Contains(move, "Burning Bulwark") || strings.Contains(move, "Baneful Bunker") ||
				strings.Contains(move, "Detect") || strings.Contains(move, "King's Shield") ||
				strings.Contains(move, "Obstruct") || strings.Contains(move, "Silk Trap") {
				protectCount++
				supportMoveCount++
			}
			if strings.Contains(move, "Tailwind") || strings.Contains(move, "Trick Room") ||
				strings.Contains(move, "Icy Wind") || strings.Contains(move, "Electroweb") ||
				strings.Contains(move, "Glaciate") || strings.Contains(move, "Sticky Web") ||
				strings.Contains(move, "Thunder Wave") {
				speedControlFlag = true
				supportMoveCount++
			}

			//Todo
			if strings.Contains(move, "Taunt") || strings.Contains(move, "Helping Hand") ||
				strings.Contains(move, "Encore") || strings.Contains(move, "Disable") ||
				strings.Contains(move, "Spore") || strings.Contains(move, "Will-O-Wisp") ||
				strings.Contains(move, "Thunder Wave") || strings.Contains(move, "Sleep Powder") ||
				strings.Contains(move, "Hypnosis") || strings.Contains(move, "Decorate") ||
				strings.Contains(move, "Haze") || strings.Contains(move, "Parting Shot") ||
				strings.Contains(move, "Coaching") || strings.Contains(move, "Charm") ||
				strings.Contains(move, "Eerie Impulse") || strings.Contains(move, "Wide Guard") ||
				strings.Contains(move, "Life Dew") || strings.Contains(move, "Yawn") ||
				strings.Contains(move, "Quick Guard") || strings.Contains(move, "Grassy Terrain") ||
				strings.Contains(move, "Misty Terrain") || strings.Contains(move, "Psychic Terrain") ||
				strings.Contains(move, "Electric Terrain") || strings.Contains(move, "Stealth Rock") ||
				strings.Contains(move, "Spikes") || strings.Contains(move, "Scary Face") ||
				strings.Contains(move, "Howl") {
				statusFlag = true
				supportMoveCount++

			}
		}
	}

	report := strings.Builder{}
	var score float64
	//report.WriteString("\n\nSupport report \n -----------------------------\n")

	if totalMoveCount < 24 {
		report.WriteString("Not all of your Pokemon have 4 moves, so your score will be reduced. Make sure each of your Pokemon have 4 moves.")
		score = 0.0
	} else {

		if fakeOutFlag {
			report.WriteString("\nFake out detected.")
		} else {
			report.WriteString("\nThere isn't any fake out on your team. Consider adding the move to a Pokemon that " +
				"learns it and adding that Pokemon to your team, as it is a valuable tool in VGC that can prevent an opponent from attacking.")
		}

		if redirectionFlag {
			report.WriteString("\nRedirection detected.")
		} else {
			report.WriteString("\nYour team has no redirection moves. Consider adding Rage Powder or Follow Me to a Pokemon " +
				"that learns it and add it to your team. Redirection moves can protect frailer offensive Pokemon so they can " +
				"survive longer.")
		}

		if screensFlag {
			report.WriteString("\nScreens detected.")
		} else {
			report.WriteString("\nYour team has no screens to reduce damage. Screens aren't necessary for VGC teams, but if your " +
				"team doesn't have a lot of defense, consider adding Light Screen, Reflect, or Aurora Veil.")
		}

		if protectCount <= 2 {
			report.WriteString("\nYour team doesn't have a lot of Pokemon that are carrying the move Protect. Protect is the most popular " +
				"move in VGC, and is a staple to well-rounded teams. It may be worth adding Protect to more Pokemon.")
		}

		if speedControlFlag {
			report.WriteString("\nSpeed Control detected.")
		} else {
			report.WriteString("\nNone of your Pokemon have a way to affect the speed control of your team. Speed control " +
				"is vital in VGC so that you can make your Pokemon move as early in the turn as possible. It would be good to add " +
				"methods of speed control such as Tailwind or Trick Room, or spread moves that affect speed such as Icy Wind or " +
				"Electroweb.")
		}

		if statusFlag {
			report.WriteString("\nStatus moves detected.")
		} else {
			report.WriteString("\nYour team doesn't have any popular status moves, such as Taunt, status-infliction moves, Helping " +
				"Hand, or stat raising/decreasing moves. Which moves fit your team may vary, but status moves can be disruptive towards " +
				"your opponent's Pokemon.")
		}

		/**
		Grading strategy: Based off of statistics of top tournament teams, 8 support moves (including protects) will be the benchmark.
		Lower amounts of supporting moves will be decremented accordingly.
		*/

		if supportMoveCount >= 8 {
			score = 10.0
		} else if supportMoveCount == 7 {
			score = 8.5
		} else if supportMoveCount == 6 {
			score = 7.0
		} else if supportMoveCount == 5 {
			score = 6.0
		} else {
			score = 1.0
		}
	}

	res := report.String()
	return res, score
}
