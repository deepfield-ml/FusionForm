// 7 to 11

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RoleType string

const (
	CenterBack        RoleType = "CD"
	FullBack          RoleType = "FB"
	CentralMidfielder RoleType = "CM"
	WideMidfielder    RoleType = "WM"
	Striker           RoleType = "ST"
	Flexible          RoleType = "Flex"
)

type Role struct {
	RoleType RoleType
	Name     string
}

// Corrected and Transformative convert7to11RuleBased function:
func convert7to11RuleBased(formation7 []Role) []Role {
	roles11 := []Role{}
	counts7 := make(map[RoleType]int)
	for _, role := range formation7 {
		counts7[role.RoleType]++
	}

	cd7 := counts7[CenterBack]
	cm7 := counts7[CentralMidfielder]
	st7 := counts7[Striker]
	flex7 := counts7[Flexible] // If used

	// Rule 1: Expand Defense to 4 (Transform CD to CDs and add FBs if needed)
	defenders11Needed := 4
	defenders7 := cd7 // Start with existing CDs
	fullbacksToAdd := 0

	if defenders7 <= defenders11Needed {
		fullbacksToAdd = defenders11Needed - defenders7
	} else {
		defenders7 = defenders11Needed // Limit CDs if 7-a-side had too many for a typical 11-a-side back 4
	}

	for i := 0; i < defenders7; i++ {
		roles11 = append(roles11, Role{RoleType: CenterBack, Name: fmt.Sprintf("CD_%d", i+1)})
	}
	for i := 0; i < fullbacksToAdd; i++ {
		roles11 = append(roles11, Role{RoleType: FullBack, Name: fmt.Sprintf("FB_%d", i+1)})
	}

	// Rule 2: Expand Midfield to 4 (or 5 if space allows, Transform CMs and add WMs)
	midfielders11Needed := 4 // Target 4 midfielders initially
	midfielders7 := cm7      // Start with existing CMs
	wideMidfieldersToAdd := 0

	if midfielders7 <= midfielders11Needed {
		wideMidfieldersToAdd = midfielders11Needed - midfielders7
	} else {
		midfielders7 = midfielders11Needed // Limit CMs if 7-a-side had too many
	}

	for i := 0; i < midfielders7; i++ {
		roles11 = append(roles11, Role{RoleType: CentralMidfielder, Name: fmt.Sprintf("CM_%d", i+1)})
	}
	for i := 0; i < wideMidfieldersToAdd; i++ {
		roles11 = append(roles11, Role{RoleType: WideMidfielder, Name: fmt.Sprintf("WM_%d", i+1)})
	}

	// Rule 3: Keep Strikers (Aim for 2, adjust midfield/defense if needed to keep total around 10 outfield)
	strikers11Needed := 2
	strikers7 := st7 // Start with existing strikers
	strikersToAdd := 0

	if strikers7 <= strikers11Needed {
		strikersToAdd = strikers11Needed - strikers7
	} else {
		strikers7 = strikers11Needed // Limit strikers if 7-a-side had too many
	}

	for i := 0; i < strikers7; i++ {
		roles11 = append(roles11, Role{RoleType: Striker, Name: fmt.Sprintf("ST_%d", i+1)})
	}
	for i := 0; i < strikersToAdd; i++ {
		roles11 = append(roles11, Role{RoleType: Striker, Name: fmt.Sprintf("ST_Extra_%d", i+1)})
	}

	// Rule 4:  Handle Flexible Roles (Convert Flex to Wide Midfielder if needed, or adjust midfield count)
	flexibleRolesToAdd := flex7 // Bring over flexible roles
	// Try to convert flexible roles to Wide Midfielders if there's space to reach 10 outfield players
	if flexibleRolesToAdd > 0 && (len(roles11)+flexibleRolesToAdd <= 10) { // Check if adding flex roles keeps total <= 10
		for i := 0; i < flexibleRolesToAdd; i++ {
			roles11 = append(roles11, Role{RoleType: WideMidfielder, Name: fmt.Sprintf("FlexWM_%d", i+1)}) // Convert Flex to WM
		}
	}

	// **Final Adjustment: Ensure Total Outfield Players is Exactly 10**
	currentOutfieldPlayers := len(roles11)
	if currentOutfieldPlayers > 10 {
		roles11 = roles11[:10] // Truncate if somehow we exceeded 10 (safety measure, should not happen with rules above)
	} else if currentOutfieldPlayers < 10 {
		// Add more midfielders (Central Midfielders as default for filling gaps) to reach 10
		midfieldersNeededToFill := 10 - currentOutfieldPlayers
		for i := 0; i < midfieldersNeededToFill; i++ {
			roles11 = append(roles11, Role{RoleType: CentralMidfielder, Name: fmt.Sprintf("CM_Fill_%d", i+1)}) // Fill with CMs
		}
	}

	return roles11
}

// Function to parse numerical formation input string for 7-a-side (e.g., "231")
func parseNumericalFormationInput7(input string) ([]Role, error) {
	if len(input) != 3 {
		return nil, fmt.Errorf("input must be a 3-digit number string (e.g., 231)")
	}

	defenders, err := strconv.Atoi(string(input[0]))
	if err != nil || defenders < 0 {
		return nil, fmt.Errorf("invalid number of defenders")
	}
	midfielders, err := strconv.Atoi(string(input[1]))
	if err != nil || midfielders < 0 {
		return nil, fmt.Errorf("invalid number of midfielders")
	}
	attackers, err := strconv.Atoi(string(input[2]))
	if err != nil || attackers < 0 {
		return nil, fmt.Errorf("invalid number of attackers")
	}

	totalOutfield := defenders + midfielders + attackers
	if totalOutfield > 6 { // Allow max 6 outfield players making it seem like it didn't print.
		return nil, fmt.Errorf("total outfield players must be at most 6 for 7-a-side, got %d", totalOutfield)
	}

	formation := make([]Role, 0, 6) // Max 6 outfield for 7-a-side

	// Assume distribution of roles within lines (can be customized)
	// Defenders: Assume all are Center Backs in 7-a-side numerical input
	for i := 0; i < defenders; i++ {
		formation = append(formation, Role{RoleType: CenterBack, Name: fmt.Sprintf("CD_%d", i+1)})
	}

	// Midfielders: Assume all are Central Midfielders in 7-a-side numerical input
	for i := 0; i < midfielders; i++ {
		formation = append(formation, Role{RoleType: CentralMidfielder, Name: fmt.Sprintf("CM_%d", i+1)})
	}

	// Attackers: Assume all are Strikers in 7-a-side numerical input
	for i := 0; i < attackers; i++ {
		formation = append(formation, Role{RoleType: Striker, Name: fmt.Sprintf("ST_%d", i+1)})
	}

	// Add Flexible Role if needed to reach 6 outfield (example: if input was like "230" -> 2-3-0 + 1 Flex = 2-3-1)
	if totalOutfield < 6 {
		formation = append(formation, Role{RoleType: Flexible, Name: "Flex_1"})
	}

	return formation, nil
}

// Uncomment this function to use it if you prefer calling a function for ASCII art.
// func printAsciiArt() {
// 	fmt.Println(`
// 	._______.____     .________.___ ._______  .______  ._______._______  .______  ._____.___
// 	:_ ____/|    |___ |    ___/: __|: .___  \ :      \ :_ ____/: .___  \ : __   \ :         |
// 	|   _/  |    |   ||___    \| : || :   |  ||       ||   _/  | :   |  ||  \____||   \  /  |
// 	|   |   |    :   ||       /|   ||     :  ||   |   ||   |   |     :  ||   :  \ |   |\/   |
// 	|_. |   |        ||__:___/ |   | \_. ___/ |___|   ||_. |    \_. ___/ |   |___\|___| |   |
// 	  :/    |. _____/    :     |___|   :/         |___|  :/       :/     |___|          |___|
// 	  :      :/                        :                 :        :
// 	         :
// 	    FusionForm(7to11) V0.0.1 By Gordon.H
// 	© Deepfield 2025

// 	`)
// }

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(`


	 ▄▄▄▄▄▄▄▄▄▄▄ ▄         ▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄        ▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄       ▄▄ 
	▐░░░░░░░░░░░▐░▌       ▐░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░▌      ▐░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░▌     ▐░░▌
	▐░█▀▀▀▀▀▀▀▀▀▐░▌       ▐░▐░█▀▀▀▀▀▀▀▀▀ ▀▀▀▀█░█▀▀▀▀▐░█▀▀▀▀▀▀▀█░▐░▌░▌     ▐░▐░█▀▀▀▀▀▀▀▀▀▐░█▀▀▀▀▀▀▀█░▐░█▀▀▀▀▀▀▀█░▐░▌░▌   ▐░▐░▌
	▐░▌         ▐░▌       ▐░▐░▌              ▐░▌    ▐░▌       ▐░▐░▌▐░▌    ▐░▐░▌         ▐░▌       ▐░▐░▌       ▐░▐░▌▐░▌ ▐░▌▐░▌
	▐░█▄▄▄▄▄▄▄▄▄▐░▌       ▐░▐░█▄▄▄▄▄▄▄▄▄     ▐░▌    ▐░▌       ▐░▐░▌ ▐░▌   ▐░▐░█▄▄▄▄▄▄▄▄▄▐░▌       ▐░▐░█▄▄▄▄▄▄▄█░▐░▌ ▐░▐░▌ ▐░▌
	▐░░░░░░░░░░░▐░▌       ▐░▐░░░░░░░░░░░▌    ▐░▌    ▐░▌       ▐░▐░▌  ▐░▌  ▐░▐░░░░░░░░░░░▐░▌       ▐░▐░░░░░░░░░░░▐░▌  ▐░▌  ▐░▌
	▐░█▀▀▀▀▀▀▀▀▀▐░▌       ▐░▌▀▀▀▀▀▀▀▀▀█░▌    ▐░▌    ▐░▌       ▐░▐░▌   ▐░▌ ▐░▐░█▀▀▀▀▀▀▀▀▀▐░▌       ▐░▐░█▀▀▀▀█░█▀▀▐░▌   ▀   ▐░▌
	▐░▌         ▐░▌       ▐░▌         ▐░▌    ▐░▌    ▐░▌       ▐░▐░▌    ▐░▌▐░▐░▌         ▐░▌       ▐░▐░▌     ▐░▌ ▐░▌       ▐░▌
	▐░▌         ▐░█▄▄▄▄▄▄▄█░▌▄▄▄▄▄▄▄▄▄█░▌▄▄▄▄█░█▄▄▄▄▐░█▄▄▄▄▄▄▄█░▐░▌     ▐░▐░▐░▌         ▐░█▄▄▄▄▄▄▄█░▐░▌      ▐░▌▐░▌       ▐░▌
	▐░▌         ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░▌      ▐░░▐░▌         ▐░░░░░░░░░░░▐░▌       ▐░▐░▌       ▐░▌
	 ▀           ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀        ▀▀ ▀           ▀▀▀▀▀▀▀▀▀▀▀ ▀         ▀ ▀         ▀ 
																															 
									
	FusionForm(7to11) V0.0.1 By Gordon.H
	© Deepfield 2025
	`)
	fmt.Println("Enter 7-a-side formation in numbers (e.g., 231, 222):")
	inputFormation7NumStr, _ := reader.ReadString('\n')
	inputFormation7NumStr = strings.TrimSpace(inputFormation7NumStr)

	formation7, err := parseNumericalFormationInput7(inputFormation7NumStr)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	formation11 := convert7to11RuleBased(formation7)

	fmt.Println("\n7-a-side Formation Input:")
	for _, role := range formation7 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println("\n\n11-a-side Formation (Rule-Based Expansion):")
	for _, role := range formation11 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println()

	// Count role types for summary
	counts := make(map[RoleType]int)
	for _, role := range formation11 {
		counts[role.RoleType]++
	}
	fmt.Println("\n11-a-side Formation Summary:")
	fmt.Printf("%d CD, %d FB, %d CM, %d WM, %d ST\n", counts[CenterBack], counts[FullBack], counts[CentralMidfielder], counts[WideMidfielder], counts[Striker])
	fmt.Printf("Formation Count (approx.): %d-%d-%d-%d\n", counts[CenterBack]+counts[FullBack], counts[CentralMidfielder]+counts[WideMidfielder], counts[Striker], 0) // Formation string (simplified)
}
