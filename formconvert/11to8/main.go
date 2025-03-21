// 11 aside to 8 aside
package main

import (
	"bufio"
	"flag"
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
)

type Role struct {
	RoleType RoleType
	Name     string
}

func convert11to8Relational(formation11 []Role) []Role {
	roles8 := []Role{}
	remainingRoles11 := formation11

	cdRoles := filterRoles(remainingRoles11, []RoleType{CenterBack, FullBack})
	selectedCDs := selectRoles(cdRoles, 2, []RoleType{CenterBack})
	roles8 = append(roles8, selectedCDs...)
	remainingRoles11 = removeRoles(remainingRoles11, selectedCDs)

	cmRoles := filterRoles(remainingRoles11, []RoleType{CentralMidfielder, WideMidfielder, FullBack})
	selectedCMs := selectRoles(cmRoles, 3, []RoleType{CentralMidfielder, WideMidfielder})
	roles8 = append(roles8, selectedCMs...)
	remainingRoles11 = removeRoles(remainingRoles11, selectedCMs)

	stRoles := filterRoles(remainingRoles11, []RoleType{Striker, WideMidfielder, FullBack})
	selectedSTs := selectRoles(stRoles, 2, []RoleType{Striker, WideMidfielder})
	roles8 = append(roles8, selectedSTs...)
	remainingRoles11 = removeRoles(remainingRoles11, selectedSTs)

	if len(roles8) < 7 && len(remainingRoles11) > 0 {
		flexibleRoles := filterRoles(remainingRoles11, []RoleType{WideMidfielder, FullBack, Striker, CentralMidfielder, CenterBack})
		selectedFlex := selectRoles(flexibleRoles, 1, []RoleType{WideMidfielder, FullBack, CentralMidfielder, Striker, CenterBack})
		if len(selectedFlex) > 0 {
			// roles8 = append(roles8, Role{RoleType: Flexible, Name: "Flex_Outfield_1"})
		}
	}

	// roles8 = append(roles8, Role{RoleType: Goalkeeper, Name: "GK_1"})
	return roles8
}

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
	for _, rt := range preferredRoleTypes {
		for _, role := range roles {
			if role.RoleType == rt && len(selected) < n {
				selected = append(selected, role)
			}
		}
	}
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

	numCD := defenders
	numFB := 0
	if defenders == 4 {
		numCD = 2
		numFB = 2
	} else if defenders == 3 {
		numCD = 3
		numFB = 0
	} else if defenders == 5 {
		numCD = 3
		numFB = 2
	}

	for i := 0; i < numCD; i++ {
		formation = append(formation, Role{RoleType: CenterBack, Name: fmt.Sprintf("CD_%d", i+1)})
	}
	for i := 0; i < numFB; i++ {
		formation = append(formation, Role{RoleType: FullBack, Name: fmt.Sprintf("FB_%d", i+1)})
	}

	numCM := midfielders
	numWM := 0
	if midfielders == 4 {
		numCM = 2
		numWM = 2
	} else if midfielders == 5 {
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

	for i := 0; i < attackers; i++ {
		formation = append(formation, Role{RoleType: Striker, Name: fmt.Sprintf("ST_%d", i+1)})
	}

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

					FusionForm(11to8) V0.0.1D By Gordon.H
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

	formation8 := convert11to8Relational(formation11)

	fmt.Println("\n11-a-side Formation Input:")
	for _, role := range formation11 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println("\n\n8-a-side Formation (Relational Conversion):")
	for _, role := range formation8 {
		fmt.Printf("%s (%s) ", role.Name, role.RoleType)
	}
	fmt.Println()

	counts := make(map[RoleType]int)
	for _, role := range formation8 {
		counts[role.RoleType]++
	}
	fmt.Println("\n8-a-side Formation Summary:")
	fmt.Printf("%d CD, %d CM, %d WM, %d ST\n", counts[CenterBack], counts[CentralMidfielder], counts[WideMidfielder], counts[Striker])
	fmt.Printf("Formation Count (approx. Outfield): %d-%d-%d\n", counts[CenterBack], counts[CentralMidfielder]+counts[WideMidfielder], counts[Striker])
}
