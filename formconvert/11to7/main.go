//11 aside to 7 aside

package main

import (
	"bufio"
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
	Flexible          RoleType = "Flex" // For 7-a-side flexible role(maunual decide by coach)
)

type Role struct {
	RoleType RoleType
	Name     string
}

// Conversion function
func convert11to7Relational(formation11 []Role) []Role {
	roles7 := []Role{}
	remainingRoles11 := formation11

	// Rule 1: Preserve Central Defensive Core (2 CD)
	cdRoles := filterRoles(remainingRoles11, []RoleType{CenterBack, FullBack})
	selectedCDs := selectRoles(cdRoles, 2, []RoleType{CenterBack}) // Prioritize CD
	roles7 = append(roles7, selectedCDs...)
	remainingRoles11 = removeRoles(remainingRoles11, selectedCDs)

	// Rule 2: Maintain Central Midfield Control (2 CM)
	cmRoles := filterRoles(remainingRoles11, []RoleType{CentralMidfielder, WideMidfielder, FullBack}) // Include FB again as remaining
	selectedCMs := selectRoles(cmRoles, 2, []RoleType{CentralMidfielder})                             // Prioritize CM
	roles7 = append(roles7, selectedCMs...)
	remainingRoles11 = removeRoles(remainingRoles11, selectedCMs)

	// Rule 3: Create an Attacking Focus (1 ST)
	stRoles := filterRoles(remainingRoles11, []RoleType{Striker, WideMidfielder, FullBack}) // Include FB, WM again
	selectedST := selectRoles(stRoles, 1, []RoleType{Striker})                              // Prioritize ST
	roles7 = append(roles7, selectedST...)
	remainingRoles11 = removeRoles(remainingRoles11, selectedST)

	// Rule 4: Adapt Remaining Roles for Flexibility (1 Flex - if needed to reach 6 outfield)
	if len(roles7) < 6 && len(remainingRoles11) > 0 { // Ensure 6 outfield players in 7-a-side
		flexibleRoles := filterRoles(remainingRoles11, []RoleType{WideMidfielder, FullBack, Striker, CentralMidfielder, CenterBack}) // All remaining
		selectedFlex := selectRoles(flexibleRoles, 1, []RoleType{WideMidfielder, FullBack})                                          // Prioritize WM, FB for width/balance
		if len(selectedFlex) > 0 {
			roles7 = append(roles7, Role{RoleType: Flexible, Name: "Flex_1"}) // Assign a flexible role
		}
	}

	return roles7
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

// Function to parse numerical formation input string to []Role (e.g., "442")
func parseNumericalFormationInput11(input string) ([]Role, error) {
	if len(input) != 3 {
		return nil, fmt.Errorf("input must be a 3-digit number string (e.g., 442)")
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
	if totalOutfield != 10 {
		return nil, fmt.Errorf("total outfield players must be 10 for 11-a-side, got %d", totalOutfield)
	}

	formation := make([]Role, 0, totalOutfield)

	// Assume distribution of roles within lines (can be customized)
	// Defenders: Prioritize CD then FB
	numCD := defenders
	numFB := 0
	if defenders == 4 {
		numCD = 2
		numFB = 2
	} else if defenders == 3 {
		numCD = 3
		numFB = 0
	} else if defenders == 5 { // Example for 5 defenders
		numCD = 3
		numFB = 2
	}

	for i := 0; i < numCD; i++ {
		formation = append(formation, Role{RoleType: CenterBack, Name: fmt.Sprintf("CD_%d", i+1)})
	}
	for i := 0; i < numFB; i++ {
		formation = append(formation, Role{RoleType: FullBack, Name: fmt.Sprintf("FB_%d", i+1)})
	}

	// Midfielders: Prioritize CM then WM
	numCM := midfielders
	numWM := 0
	if midfielders == 4 {
		numCM = 2
		numWM = 2
	} else if midfielders == 5 { // Example for 5 midfielders
		numCM = 3
		numWM = 2
	} else if midfielders == 3 {
		numCM = 3
		numWM = 0
	}

	for i := 0; i < numCM; i++ {
		formation = append(formation, Role{RoleType: CentralMidfielder, Name: fmt.Sprintf("CM_%d", i+1)})
	}
	for i := 0; i < numWM; i++ {
		formation = append(formation, Role{RoleType: WideMidfielder, Name: fmt.Sprintf("WM_%d", i+1)})
	}

	// Attackers: Assume all are Strikers
	for i := 0; i < attackers; i++ {
		formation = append(formation, Role{RoleType: Striker, Name: fmt.Sprintf("ST_%d", i+1)})
	}

	return formation, nil
}

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
																															 
									
	FusionForm(11to7) V0.0.1B By Gordon.H
	© Deepfield 2025
	`)
	fmt.Println("Enter 11-a-side formation in numbers (e.g., 442, 343, 433):")
	inputFormation11NumStr, _ := reader.ReadString('\n')
	inputFormation11NumStr = strings.TrimSpace(inputFormation11NumStr)

	formation11, err := parseNumericalFormationInput11(inputFormation11NumStr)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	formation7 := convert11to7Relational(formation11)

	fmt.Println("\n11-a-side Formation Input:")
	for _, role := range formation11 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println("\n\n7-a-side Formation (Relational Conversion):")
	for _, role := range formation7 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println()

	// Count role types for a summary
	counts := make(map[RoleType]int)
	for _, role := range formation7 {
		counts[role.RoleType]++
	}
	fmt.Println("\n7-a-side Formation Summary:")
	fmt.Printf("%d CD, %d CM, %d ST, %d Flex\n", counts[CenterBack], counts[CentralMidfielder], counts[Striker], counts[Flexible])
	fmt.Printf("Formation Count (approx.): %d-%d-%d (+%d Flex)\n", counts[CenterBack], counts[CentralMidfielder], counts[Striker], counts[Flexible]) // Formation string
}
