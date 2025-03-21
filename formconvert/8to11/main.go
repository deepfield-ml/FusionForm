// 8 aside to 11 aside

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// RoleType and Role structs 
type RoleType string

const (
	CenterBack        RoleType = "CD"
	FullBack          RoleType = "FB"
	CentralMidfielder RoleType = "CM"
	WideMidfielder    RoleType = "WM"
	Striker           RoleType = "ST"
	Flexible          RoleType = "Flex" // Retaining Flexible role type, might repurpose or ignore in 11-a-side output
)

type Role struct {
	RoleType RoleType
	Name     string
}

// Conversion function from 8-a-side to 11-a-side
func convert8to11Relational(formation8 []Role) []Role {
	roles11 := []Role{}
	remainingRoles8 := formation8

	// Categorize 8-a-side roles for easier processing based on assumed 8-a-side conversion logic output
	cdRoles8 := filterRoles(remainingRoles8, []RoleType{CenterBack})
	cmRoles8 := filterRoles(remainingRoles8, []RoleType{CentralMidfielder})
	stRoles8 := filterRoles(remainingRoles8, []RoleType{Striker})
	flexRoles8 := filterRoles(remainingRoles8, []RoleType{Flexible, WideMidfielder, FullBack, Striker, CentralMidfielder, CenterBack}) // Flex could be anything

	// Rule 1: Expand Defense to 4 (or 3) Defenders - Prioritize adding Full-Backs
	numDefendersToAdd := 4 - len(cdRoles8) // Aim for 4 defenders in 11-a-side as a common base
	if numDefendersToAdd > 0 {
		fbFromFlex := filterRoles(flexRoles8, []RoleType{FullBack})                     // Check if flexible roles can be FB
		selectedFBs := selectRoles(fbFromFlex, numDefendersToAdd, []RoleType{FullBack}) // Prioritize FB from flex if available
		roles11 = append(roles11, selectedFBs...)
		flexRoles8 = removeRoles(flexRoles8, selectedFBs) // Update flex roles
		numDefendersToAdd -= len(selectedFBs)
	}
	// Add Center Backs if still needed (less common to add more CD than FB first in expansion)
	if numDefendersToAdd > 0 {
		cdFromFlex := filterRoles(flexRoles8, []RoleType{CenterBack})
		selectedCDs := selectRoles(cdFromFlex, numDefendersToAdd, []RoleType{CenterBack})
		roles11 = append(roles11, selectedCDs...)
		flexRoles8 = removeRoles(flexRoles8, selectedCDs)
	}
	roles11 = append(roles11, cdRoles8...) // Add the original CD roles from 8-a-side

	// Rule 2: Expand Midfield to 4 (or 5) Midfielders - Prioritize Wide Midfielders for width
	numMidfieldersToAdd := 4 - len(cmRoles8) // Aim for 4 midfielders, can adjust to 5 later if needed
	if numMidfieldersToAdd > 0 {
		wmFromFlex := filterRoles(flexRoles8, []RoleType{WideMidfielder})                       // Check for WM in flex roles
		selectedWMs := selectRoles(wmFromFlex, numMidfieldersToAdd, []RoleType{WideMidfielder}) // Prioritize WM from flex
		roles11 = append(roles11, selectedWMs...)
		flexRoles8 = removeRoles(flexRoles8, selectedWMs)
		numMidfieldersToAdd -= len(selectedWMs)
	}
	// Add Central Midfielders if still needed
	if numMidfieldersToAdd > 0 {
		cmFromFlex := filterRoles(flexRoles8, []RoleType{CentralMidfielder})
		selectedCMs := selectRoles(cmFromFlex, numMidfieldersToAdd, []RoleType{CentralMidfielder})
		roles11 = append(roles11, selectedCMs...)
		flexRoles8 = removeRoles(flexRoles8, selectedCMs)
	}
	roles11 = append(roles11, cmRoles8...) // Add original CM roles

	// Rule 3: Maintain/Adjust Attack - Maybe add one more Striker or Wide Midfielder for attacking width
	numAttackersToAdd := 2 - len(stRoles8) // Maybe aim for 2-3 strikers/attackers total in 11-a-side
	if numAttackersToAdd > 0 {
		stFromFlex := filterRoles(flexRoles8, []RoleType{Striker, WideMidfielder})                     // Consider WM as attacking option
		selectedSTs := selectRoles(stFromFlex, numAttackersToAdd, []RoleType{Striker, WideMidfielder}) // Prioritize ST, then WM
		roles11 = append(roles11, selectedSTs...)
		flexRoles8 = removeRoles(flexRoles8, selectedSTs)
	}
	roles11 = append(roles11, stRoles8...) // Add original ST roles

	// Rule 4: Fill Remaining Slots - Prioritize Full-Backs/Wide Midfielders for balance, then Central Midfielders
	rolesToAddFinal := 10 - len(roles11) // Calculate remaining outfield players needed (10 total outfield in 11-a-side)
	if rolesToAddFinal > 0 {
		remainingFlex := flexRoles8                                                                           // Use remaining flexible roles
		fbWMCMFromFlex := filterRoles(remainingFlex, []RoleType{FullBack, WideMidfielder, CentralMidfielder}) // Prioritize FB, WM, CM in that order
		selectedFinalRoles := selectRoles(fbWMCMFromFlex, rolesToAddFinal, []RoleType{FullBack, WideMidfielder, CentralMidfielder})
		roles11 = append(roles11, selectedFinalRoles...)
	}

	// Goalkeeper
	// roles11 = append(roles11, Role{RoleType: Flexible, Name: "GK_1"})

	return roles11
}

// Helper functions
func filterRoles(roles []Role, roleTypes []RoleType) []Role {
	filtered := []Role{}
	for _, role := range roles {
		for _, rt := range roleTypes {
			if role.RoleType == rt {
				filtered = append(filtered, role)
				break
			}
		}
	}
	return filtered
}

func selectRoles(roles []Role, n int, preferredRoleTypes []RoleType) []Role {
	selected := []Role{}
	// First, try to select preferred roles
	for _, rt := range preferredRoleTypes {
		for _, role := range roles {
			if role.RoleType == rt && len(selected) < n {
				selected = append(selected, role)
			}
		}
	}
	// If still need more roles, select from the rest
	for _, role := range roles {
		alreadySelected := false
		for _, selRole := range selected {
			if role == selRole {
				alreadySelected = true
				break
			}
		}
		if !alreadySelected && len(selected) < n {
			selected = append(selected, role)
		}
	}
	return selected
}

func removeRoles(allRoles []Role, rolesToRemove []Role) []Role {
	remaining := []Role{}
	for _, role := range allRoles {
		isRemoved := false
		for _, roleToRemove := range rolesToRemove {
			if role == roleToRemove {
				isRemoved = true
				break
			}
		}
		if !isRemoved {
			remaining = append(remaining, role)
		}
	}
	return remaining
}

// Function to parse numerical formation input string to []Role for 8-a-side (e.g., "232") -  Defenders-Midfielders-Attackers (assuming flexible is handled separately or part of input)
func parseNumericalFormationInput8(input string) ([]Role, error) {
	if len(input) != 3 {
		return nil, fmt.Errorf("input must be a 3-digit number string for 8-a-side (e.g., 232)")
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
	if totalOutfield > 7 { // Max 7 outfield in 8-a-side formation input (excluding GK and flexible for now in num input)
		return nil, fmt.Errorf("total outfield players should not exceed 7 for 8-a-side input, got %d", totalOutfield)
	}

	formation := make([]Role, 0, totalOutfield)

	// Defenders: Assume all are Center Backs initially for simplicity in 8-a-side input
	for i := 0; i < defenders; i++ {
		formation = append(formation, Role{RoleType: CenterBack, Name: fmt.Sprintf("CD_8_%d", i+1)}) // CD_8 to differentiate from 11-a-side roles if needed
	}

	// Midfielders: Assume all are Central Midfielders initially
	for i := 0; i < midfielders; i++ {
		formation = append(formation, Role{RoleType: CentralMidfielder, Name: fmt.Sprintf("CM_8_%d", i+1)})
	}

	// Attackers: Assume all are Strikers
	for i := 0; i < attackers; i++ {
		formation = append(formation, Role{RoleType: Striker, Name: fmt.Sprintf("ST_8_%d", i+1)})
	}

	// Add 1 Flexible Role -  8-a-side often has 7 outfield + 1 flexible or specific 8th player (GK)
	formation = append(formation, Role{RoleType: Flexible, Name: "Flex_8_1"}) // Adding a default flexible role for 8-a-side input

	return formation, nil
}

func main() {

	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(`

	
     ___     	 ▄▄▄▄▄▄▄▄▄▄▄ ▄         ▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄        ▄ 
    /\  \    	▐░░░░░░░░░░░▐░▌       ▐░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░▌      ▐░▐
   |::\  \   	▐░█▀▀▀▀▀▀▀▀▀▐░▌       ▐░▐░█▀▀▀▀▀▀▀▀▀ ▀▀▀▀█░█▀▀▀▀▐░█▀▀▀▀▀▀▀█░▐░▌░▌     ▐░▐
   |:|:\  \  	▐░▌         ▐░▌       ▐░▐░▌              ▐░▌    ▐░▌       ▐░▐░▌▐░▌    ▐░▐
 __|:|\:\  \ 	▐░█▄▄▄▄▄▄▄▄▄▐░▌       ▐░▐░█▄▄▄▄▄▄▄▄▄     ▐░▌    ▐░▌       ▐░▐░▌ ▐░▌   ▐░▐
/DFML|_\:\__\	▐░░░░░░░░░░░▐░▌       ▐░▐░░░░░░░░░░░▌    ▐░▌    ▐░▌       ▐░▐░▌  ▐░▌  ▐░▐
\:\~~\  \/__/	▐░█▀▀▀▀▀▀▀▀▀▐░▌       ▐░▌▀▀▀▀▀▀▀▀▀█░▌    ▐░▌    ▐░▌       ▐░▐░▌   ▐░▌ ▐░▐
 \:\  \      	▐░▌         ▐░▌       ▐░▌         ▐░▌    ▐░▌    ▐░▌       ▐░▐░▌    ▐░▌▐░▐
  \:\  \     	▐░▌         ▐░█▄▄▄▄▄▄▄█░▌▄▄▄▄▄▄▄▄▄█░▌▄▄▄▄█░█▄▄▄▄▐░█▄▄▄▄▄▄▄█░▐░▌     ▐░▐░▐
   \:\__\    	▐░▌         ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░▌      ▐░░▐
    \/__/     	 ▀           ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀        ▀▀ 

		▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄       ▄▄
		░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░▌     ▐░░▌
		░█▀▀▀▀▀▀▀▀▀▐░█▀▀▀▀▀▀▀█░▐░█▀▀▀▀▀▀▀█░▐░▌░▌   ▐░▐░▌
		░▌         ▐░▌       ▐░▐░▌       ▐░▐░▌▐░▌ ▐░▌▐░▌
		░█▄▄▄▄▄▄▄▄▄▐░▌       ▐░▐░█▄▄▄▄▄▄▄█░▐░▌ ▐░▐░▌ ▐░▌
		░░░░░░░░░░░▐░▌       ▐░▐░░░░░░░░░░░▐░▌  ▐░▌  ▐░▌
		░█▀▀▀▀▀▀▀▀▀▐░▌       ▐░▐░█▀▀▀▀█░█▀▀▐░▌   ▀   ▐░▌
		░▌         ▐░▌       ▐░▐░▌     ▐░▌ ▐░▌       ▐░▌
		░▌         ▐░█▄▄▄▄▄▄▄█░▐░▌      ▐░▌▐░▌       ▐░▌
		░▌         ▐░░░░░░░░░░░▐░▌       ▐░▐░▌       ▐░▌
		▀           ▀▀▀▀▀▀▀▀▀▀▀ ▀         ▀ ▀         ▀

					FusionForm(8to11) V0.0.1B By Gordon.H
					© Deepfield 2025
	`)
	fmt.Println("Enter 8-a-side formation in numbers (e.g., 232):")
	inputFormation8NumStr, _ := reader.ReadString('\n')
	inputFormation8NumStr = strings.TrimSpace(inputFormation8NumStr)

	formation8, err := parseNumericalFormationInput8(inputFormation8NumStr)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	formation11 := convert8to11Relational(formation8)

	fmt.Println("\n8-a-side Formation Input:")
	for _, role := range formation8 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println("\n\n11-a-side Formation (Relational Conversion):")
	for _, role := range formation11 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println()

	// Count role types for summary in 11-a-side
	counts := make(map[RoleType]int)
	for _, role := range formation11 {
		counts[role.RoleType]++
	}
	fmt.Println("\n11-a-side Formation Summary:")
	fmt.Printf("%d CD, %d FB, %d CM, %d WM, %d ST, %d Flex\n", counts[CenterBack], counts[FullBack], counts[CentralMidfielder], counts[WideMidfielder], counts[Striker], counts[Flexible])
	fmt.Printf("Formation Count (approx.): %d-%d-%d-%d\n", counts[CenterBack]+counts[FullBack], counts[CentralMidfielder]+counts[WideMidfielder], counts[Striker], counts[Flexible]) // Formation string approximation
}

func displaySettings() {
	fmt.Println("Settings:")
	// Add settings display logic here
	// For example, display current configuration or options
	fmt.Println("Option 1: ...")
	fmt.Println("Option 2: ...")
	// ...add more settings as needed...
}
