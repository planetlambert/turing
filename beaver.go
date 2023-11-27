package turing

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	haltMConfigurationName = "halt"
	maxMoves               = 1000
)

// Finds the m-configuration and number of `1`'s of the `n`'th busy beaver.
func busyBeaver(n int, debug bool) (int, MachineInput) {
	// Initialize sets of m-configurations
	var mConfigurations []MConfiguration
	for i := 0; i < n; i++ {
		mConfigurations = append(mConfigurations, MConfiguration{
			Name:                strconv.Itoa(i),
			Symbols:             []string{"0"},
			Operations:          []string{"P0", "L"}, // Print, then Move
			FinalMConfiguration: "0",
		})
		mConfigurations = append(mConfigurations, MConfiguration{
			Name:                strconv.Itoa(i),
			Symbols:             []string{"1"},
			Operations:          []string{"P0", "L"}, // Print, then Move
			FinalMConfiguration: "0",
		})
	}

	// Keep track of the best so far
	var best int
	var bestMConfigurations []MConfiguration

	// The main bit
	for {
		// Run the current set of m-configurations
		if atLeastOneHaltState(mConfigurations) {
			result := simulateBusyBeaver(mConfigurations)
			if debug {
				mConfigurationsString := getMConfigurationsString(mConfigurations)
				fmt.Printf("best %d | result %d | %s\n", best, result, mConfigurationsString)
			}
			if result > best {
				best = result
				bestMConfigurations = mConfigurations
			}
		}

		var over bool
		for i := 0; i < n*2; i++ {
			// Iterate to the next the m-configuration
			nextMConfiguration, reset := nextMConfiguration(n, mConfigurations[i])
			mConfigurations[i] = nextMConfiguration

			// If we don't need to reset the next m-configuration, move on
			if !reset {
				break
			}

			// We just reset all of our m-configurations (we've looped through everything)
			if i == (n*2)-1 {
				over = true
			}
		}

		// Check to see if we are done
		if over {
			break
		}
	}

	// Return the best we have
	return best, getBusyBeaverMachineInput(bestMConfigurations)
}

// Iterate through all of the variables of an m-configuration, return true if we did a full loop
func nextMConfiguration(n int, mConfiguration MConfiguration) (MConfiguration, bool) {
	// Print Operations: P0, P1
	if mConfiguration.Operations[0] == "P0" {
		return MConfiguration{
			Name:                mConfiguration.Name,
			Symbols:             mConfiguration.Symbols,
			Operations:          []string{"P1", mConfiguration.Operations[1]},
			FinalMConfiguration: mConfiguration.FinalMConfiguration,
		}, false
	}

	// Move Operations: L, R
	if mConfiguration.Operations[1] == "L" {
		return MConfiguration{
			Name:                mConfiguration.Name,
			Symbols:             mConfiguration.Symbols,
			Operations:          []string{"P0", "R"},
			FinalMConfiguration: mConfiguration.FinalMConfiguration,
		}, false
	}

	// Final m-configurations: 0...n, halt
	if mConfiguration.FinalMConfiguration != haltMConfigurationName {
		finalMConfigurationInt, _ := strconv.Atoi(mConfiguration.FinalMConfiguration)
		var finalMConfiguration string
		if finalMConfigurationInt == n-1 {
			finalMConfiguration = haltMConfigurationName
		} else {
			finalMConfiguration = strconv.Itoa(finalMConfigurationInt + 1)
		}
		return MConfiguration{
			Name:                mConfiguration.Name,
			Symbols:             mConfiguration.Symbols,
			Operations:          []string{"P0", "L"},
			FinalMConfiguration: finalMConfiguration,
		}, false
	}

	// If we reset the full m-configuration, signal upwards that the next m-configuration in the list should iterate
	return MConfiguration{
		Name:                mConfiguration.Name,
		Symbols:             mConfiguration.Symbols,
		Operations:          []string{"P0", "L"},
		FinalMConfiguration: "0",
	}, true
}

// No need to simulate if we know the MConfiguration will never halt
func atLeastOneHaltState(mConfigurations []MConfiguration) bool {
	for _, mConfiguration := range mConfigurations {
		if mConfiguration.FinalMConfiguration == haltMConfigurationName {
			return true
		}
	}
	return false
}

// Return the amount of `1`'s the machine prints up to `maxMoves`
func simulateBusyBeaver(mConfigurations []MConfiguration) int {
	m := NewMachine(getBusyBeaverMachineInput(mConfigurations))
	moves := m.MoveN(maxMoves)
	if moves == maxMoves {
		return 0
	}

	var count int
	for _, square := range m.Tape() {
		if square == "1" {
			count++
		}
	}
	return count
}

// For a set of our m-configurations, give a runnable MachineInput
func getBusyBeaverMachineInput(mConfigurations []MConfiguration) MachineInput {
	return MachineInput{
		MConfigurations: mConfigurations,
		PossibleSymbols: []string{"1"},
		NoneSymbol:      "0",
	}
}

// Shortens m-configurations for printing/debugging
func getMConfigurationsString(mConfigurations []MConfiguration) string {
	var s strings.Builder

	for i, mConfiguration := range mConfigurations {
		if i%2 == 0 {
			s.WriteString(fmt.Sprintf("%s[", mConfiguration.Name))
		}

		s.WriteString(fmt.Sprintf(" %s:%s;%s;%s", mConfiguration.Symbols[0], mConfiguration.Operations[0], mConfiguration.Operations[1], mConfiguration.FinalMConfiguration))

		if i%2 == 0 {
			s.WriteString(fmt.Sprintf(","))
		} else {
			s.WriteString(fmt.Sprintf(" ] "))
		}
	}

	return s.String()
}
