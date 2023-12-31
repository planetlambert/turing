package turing

import (
	"fmt"
	"slices"
	"strings"
)

type (
	// We may compare a man in the process of computing a real number to a machine...
	MachineInput struct {
		// ...which is only capable of a finite number of conditions q1, q2, ..., qR which
		// will be called "m-configurations".
		MConfigurations []MConfiguration

		// The machine is supplied with a "tape" (the analogue of paper) running through it,
		// and divided into sections (called "squares") each capable of bearing a "symbol".
		Tape Tape

		// The m-configuration that the machine should start with. If empty the first m-configuration
		// in the list is chosen.
		StartingMConfiguration string

		// A list of all symbols the machine is capable of reading or printing.
		// This field is used when dealing with special symbols `*` (Any), `!` (Not)
		// Note: The ` ` (None) symbol does not have to be provided (it is assumed).
		PossibleSymbols []string

		// Defaults to ` ` (None), but can optionally be overridden here.
		NoneSymbol string

		// If `true`, the machine's complete configurations are printed at the end of each move.
		Debug bool
	}

	// Turing's Machine
	Machine struct {
		// See corresponding input field
		mConfigurations []MConfiguration

		// See corresponding input field
		tape []string

		// See corresponding input field
		possibleSymbols []string

		// See corresponding input field
		noneSymbol string

		// See corresponding input field
		debug bool

		// At any moment there is just one square, say the r-th, bearing the symbol S(r)
		// which is "in the machine". We may call this square the "scanned square".
		// The symbol on the scanned square may be called the "scanned symbol".
		// The "scanned symbol" is the only one of which the machine is, so to speak, "directly aware".
		scannedSquare int

		// The current m-configuration of the machine.
		currentMConfigurationName string

		// Stores whether the machine has "halted" or not. A machine only halts if it cannot
		// find an m-configuration.
		halted bool
	}

	// An m-configuration contains four components
	MConfiguration struct {
		// The possible behaviour of the machine at any moment is determined by the m-configuration qn...
		Name string

		// ...and the scanned symbol S(r)
		Symbols []string

		// In some of the configurations in which the scanned square is blank (i.e. bears no symbol)
		// the machine writes down a new symbol on the scanned square: in other configurations it
		// erases the scanned symbol. The machine may also change the square which is being scanned,
		// but only by shifting it one place to right or left.
		Operations []string

		// In addition to any of these operations the m-configuration may be changed.
		FinalMConfiguration string
	}

	// Our "tape" is a slice of strings because squares can contain multiple characters
	Tape []string

	// Well-known single-character codes used in an m-configuration's operations.
	operationCode byte
)

const (
	rightOp operationCode = 'R'
	leftOp  operationCode = 'L'
	eraseOp operationCode = 'E'
	printOp operationCode = 'P'

	none string = " "
	not  string = "!"
	any  string = "*"
)

// Returns a new Machine
func NewMachine(input MachineInput) *Machine {
	m := &Machine{
		mConfigurations: input.MConfigurations,
		debug:           input.Debug,
	}

	// Use first m-configuration if starting m-configuration not specified
	if len(input.StartingMConfiguration) == 0 {
		m.currentMConfigurationName = input.MConfigurations[0].Name
	} else {
		m.currentMConfigurationName = input.StartingMConfiguration
	}

	// Use default None character if not specified
	if len(input.NoneSymbol) == 0 {
		m.noneSymbol = none
	} else {
		m.noneSymbol = input.NoneSymbol
	}

	// Initialize tape if nil
	if input.Tape == nil {
		m.tape = []string{}
	} else {
		m.tape = input.Tape
	}

	if m.debug {
		m.printMConfigurationsForDebug()
	}

	return m
}

// Moves the machine n times and stops early if halted. Returns the amount of moves the machine took.
func (m *Machine) MoveN(n int) int {
	for i := 1; i <= n; i++ {
		m.Move()
		if m.halted {
			return i
		}
	}
	return n
}

// Moves the machine once
func (m *Machine) Move() {
	if m.halted {
		return
	}

	// Scan symbol from the tape
	symbol := m.scan()

	// Find the the correct m-configuration depending on the scanned synbol
	mConfiguration, shouldHalt := m.findMConfiguration(m.currentMConfigurationName, symbol)

	// If an m-configuration could not be found, halt the machine
	if shouldHalt {
		m.halted = true
		return
	}

	// Perform operations
	for _, operation := range mConfiguration.Operations {
		m.performOperation(operation)
	}

	if m.debug {
		m.printCompleteConfigurationForDebug()
	}

	// Move to specified final-m-configuration
	m.currentMConfigurationName = mConfiguration.FinalMConfiguration
}

// Returns the Machine's Tape
func (m *Machine) Tape() Tape {
	return m.tape
}

// Return the Tape represented as a string
func (m *Machine) TapeString() string {
	return strings.Join([]string(m.tape), "")
}

// Returns the machine's Complete Configuration of the single-line form
func (m *Machine) CompleteConfiguration() string {
	var completeConfiguration strings.Builder
	for i, square := range m.tape {
		if i == m.scannedSquare {
			completeConfiguration.WriteString(m.currentMConfigurationName)
		}
		completeConfiguration.WriteString(square)
	}
	if m.scannedSquare == len(m.tape) {
		completeConfiguration.WriteString(m.currentMConfigurationName)
	}
	return completeConfiguration.String()
}

// Scans the tape for the scanned symbol
func (m *Machine) scan() string {
	m.extendTapeIfNeeded()
	return m.tape[m.scannedSquare]
}

// The Machine's Tape is infinite, so we extend it as-needed
func (m *Machine) extendTapeIfNeeded() {
	if m.scannedSquare >= len(m.tape) {
		m.tape = append(m.tape, m.noneSymbol)
	}
	if m.scannedSquare < 0 {
		m.tape = append([]string{m.noneSymbol}, m.tape...)
		m.scannedSquare++
	}
}

// Find the appropriate full m-configuration given the current m-configuration name and the scanned symbol
func (m *Machine) findMConfiguration(mConfigurationName string, symbol string) (MConfiguration, bool) {
	for _, mConfiguration := range m.mConfigurations {
		if mConfiguration.Name == mConfigurationName {
			// Scenario 1: The provided symbol is contained exactly in the m-configuration
			if slices.Contains(mConfiguration.Symbols, symbol) {
				return mConfiguration, false
			}

			if symbol != m.noneSymbol {
				// Scenario 2: The m-configuration contains `*`
				// Note that `*` does not include ` ` (None), which must be specified manually
				if slices.Contains(mConfiguration.Symbols, any) {
					return mConfiguration, false
				}

				// Scenario 3: The MConfiguration contains `!x` where `x` is not the provided symbol
				// Note that `!` does not include ` ` (None), which must be specified manually
				notSymbols := []string{}
				// First loop is required in the scenario we have multiple (`!x` and `!y`)
				for _, mConfigurationSymbol := range mConfiguration.Symbols {
					if strings.Contains(mConfigurationSymbol, not) {
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

// Perform an operation
func (m *Machine) performOperation(operation string) {
	m.extendTapeIfNeeded()
	switch operationCode(operation[0]) {
	case rightOp:
		m.scannedSquare++
	case leftOp:
		m.scannedSquare--
	case eraseOp:
		m.tape[m.scannedSquare] = m.noneSymbol
	case printOp:
		m.tape[m.scannedSquare] = string(operation[1:])
	}
}

// Prints the m-configurations of the machine nicely for debugging
func (m *Machine) printMConfigurationsForDebug() {
	for _, mConfiguration := range m.mConfigurations {
		fmt.Printf("%s %v %v %s\n", mConfiguration.Name, mConfiguration.Symbols, mConfiguration.Operations, mConfiguration.FinalMConfiguration)
	}
}

// Prints the complete configuration for the machine nicely for debugging
func (m *Machine) printCompleteConfigurationForDebug() {
	for _, square := range m.tape {
		fmt.Print(strings.Repeat("-", len(square)))
	}
	fmt.Println("-")
	fmt.Println(m.TapeString())
	for i, square := range m.tape {
		if i >= m.scannedSquare {
			break
		}
		fmt.Print(strings.Repeat(" ", len(square)))
	}
	fmt.Println(m.currentMConfigurationName)
}
