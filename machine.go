package turing

import (
	"fmt"
	"slices"
	"strings"
)

type (
	Machine struct {
		MConfigurations
		Tape
		StartingMConfiguration   string
		PossibleSymbols          []string // ` ` (None) is assumed
		NoneSymbol               string   // Required as standard machines use `S0`
		Debug                    bool
		currentConfigurationName string
		scannedSquare            int
		halted                   bool

		// For use when converting to StandardTable
		mConfigurationNames   map[string]string
		nameCount             int
		mConfigurationSymbols map[string]string
		symbolCount           int
	}

	MConfigurations []MConfiguration

	MConfiguration struct {
		Name                string
		Symbols             []string
		Operations          []string
		FinalMConfiguration string
	}

	Tape []string

	OperationCode byte
)

const (
	None string = " "
	Not  string = "!"
	Any  string = "*"

	Right OperationCode = 'R'
	Left  OperationCode = 'L'
	Erase OperationCode = 'E'
	Print OperationCode = 'P'

	mConfigurationNamePrefix   string = "q"
	mConfigurationSymbolPrefix string = "S"
)

// Moves the Machine n times
func (m *Machine) MoveN(n int) {
	for i := 1; i <= n; i++ {
		m.Move()
		if m.halted {
			return
		}
	}
}

// Moves the Machine once
func (m *Machine) Move() {
	if m.halted {
		return
	}

	// Initialize
	m.init()

	// Scan Symbol from the Tape
	symbol := m.scan()

	// Find MConfiguration
	mConfiguration, shouldHalt := m.findMConfiguration(m.currentConfigurationName, symbol)

	// An MConfiguration could not be found
	if shouldHalt {
		m.halted = true
		return
	}

	// Perform operations
	for _, operation := range mConfiguration.Operations {
		m.performOperation(operation)
	}

	if m.Debug {
		m.printCompleteConfiguration()
	}

	// Move to specified final MConfiguration
	m.currentConfigurationName = mConfiguration.FinalMConfiguration
}

// Return the Tape represented as a string
func (m *Machine) TapeString() string {
	return strings.Join([]string(m.Tape), "")
}

func (m *Machine) printMConfigurations() {
	for _, mConfiguration := range m.MConfigurations {
		fmt.Printf("%s %v %v %s\n", mConfiguration.Name, mConfiguration.Symbols, mConfiguration.Operations, mConfiguration.FinalMConfiguration)
	}
}

// Prints the complete configuration for the Machine
func (m *Machine) printCompleteConfiguration() {
	for _, square := range m.Tape {
		fmt.Print(strings.Repeat("-", len(square)))
	}
	fmt.Println("-")
	fmt.Println(m.TapeString())
	for i, square := range m.Tape {
		if i >= m.scannedSquare {
			break
		}
		fmt.Print(strings.Repeat(" ", len(square)))
	}
	fmt.Println(m.currentConfigurationName)
}

func (m *Machine) init() {
	if len(m.currentConfigurationName) == 0 {
		if m.Debug {
			m.printMConfigurations()
		}
		if len(m.StartingMConfiguration) == 0 {
			m.currentConfigurationName = m.MConfigurations[0].Name
		} else {
			m.currentConfigurationName = m.StartingMConfiguration
		}
	}
	if len(m.NoneSymbol) == 0 {
		m.NoneSymbol = None
	}
	if m.Tape == nil {
		m.Tape = []string{}
	}
}

// Scan the Tape
func (m *Machine) scan() string {
	m.extendTape()
	return m.Tape[m.scannedSquare]
}

// The Machine's Tape is infinite, so we extend it as-needed
func (m *Machine) extendTape() {
	if m.scannedSquare >= len(m.Tape) {
		m.Tape = append(m.Tape, m.NoneSymbol)
	}
	if m.scannedSquare < 0 {
		m.Tape = append([]string{m.NoneSymbol}, m.Tape...)
		m.scannedSquare++
	}
}

// Find the appropriate full MConfiguration given the current MConfiguration and the scanned symbol
func (m *Machine) findMConfiguration(mConfigurationName string, symbol string) (MConfiguration, bool) {
	for _, mConfiguration := range m.MConfigurations {
		if mConfiguration.Name == mConfigurationName {
			// Scenario 1: The provided symbol is contained exactly in the MConfiguration
			if slices.Contains(mConfiguration.Symbols, symbol) {
				return mConfiguration, false
			}
			if symbol != m.NoneSymbol {
				// Scenario 3: The MConfiguration contains `*`
				// Note that `*` does not include ` ` (None), which must be specified manually
				if slices.Contains(mConfiguration.Symbols, Any) {
					return mConfiguration, false
				}

				// Scenario 3: The MConfiguration contains `!x` where `x` is not the provided symbol
				// Note that `!` does not include ` ` (None), which must be specified manually
				notSymbols := []string{}
				for _, mConfigurationSymbol := range mConfiguration.Symbols {
					if strings.Contains(mConfigurationSymbol, Not) {
						notSymbols = append(notSymbols, mConfigurationSymbol[1:])
					}
				}
				if len(notSymbols) > 0 && !slices.Contains(notSymbols, symbol) {
					return mConfiguration, false
				}
			}
		}
	}
	return MConfiguration{}, true
}

// Successively perform each operation
func (m *Machine) performOperation(operation string) {
	m.extendTape()
	switch OperationCode(operation[0]) {
	case Right:
		m.scannedSquare++
	case Left:
		m.scannedSquare--
	case Erase:
		m.Tape[m.scannedSquare] = m.NoneSymbol
	case Print:
		m.Tape[m.scannedSquare] = string(operation[1:])
	}
}
