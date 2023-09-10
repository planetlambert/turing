package turing

import (
	"fmt"
	"slices"
	"strings"
)

type (
	// We may compare a man in the process of computing a real number to a machine...
	Machine struct {

		// ...which is only capable of a finite number of conditions q1, q2, ..., qR which
		// will be called "m-configurations".
		MConfigurations []MConfiguration

		// The machine is supplied with a "tape" (the analogue of paper) running through it,
		// and divided into sections (called "squares") each capable of bearing a "symbol".
		// Our "tape" is a slice of strings because squares can contain multiple characters
		Tape []string

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

		// For use when converting to StandardTable
		mConfigurationNames map[string]string
		// For use when converting to StandardTable
		nameCount int
		// For use when converting to StandardTable
		mConfigurationSymbols map[string]string
		// For use when converting to StandardTable
		symbolCount int
	}

	// An m-configuration contains four components
	MConfiguration struct {

		// The possible behaviour of the machine at any moment is determined by the m-configuration qn ...
		Name string

		// ... and the scanned symbol S(r)
		Symbols []string

		// In some of the configurations in which the scanned square is blank (i.e. bears no symbol)
		// the machine writes down a new symbol on the scanned square: in other configurations it
		// erases the scanned symbol. The machine may also change the square which is being scanned,
		// but only by shifting it one place to right or left.
		Operations []string

		// In addition to any of these operations the m-configuration may be changed.
		FinalMConfiguration string
	}

	// Well-known single-character codes used in an m-configuration's operations.
	operationCode byte
)

const (
	Right operationCode = 'R'
	Left  operationCode = 'L'
	Erase operationCode = 'E'
	Print operationCode = 'P'

	None string = " "
	Not  string = "!"
	Any  string = "*"

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
	mConfiguration, shouldHalt := m.findMConfiguration(m.currentMConfigurationName, symbol)

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
	m.currentMConfigurationName = mConfiguration.FinalMConfiguration
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
	fmt.Println(m.currentMConfigurationName)
}

func (m *Machine) init() {
	if len(m.currentMConfigurationName) == 0 {
		if m.Debug {
			m.printMConfigurations()
		}
		if len(m.StartingMConfiguration) == 0 {
			m.currentMConfigurationName = m.MConfigurations[0].Name
		} else {
			m.currentMConfigurationName = m.StartingMConfiguration
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
	switch operationCode(operation[0]) {
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
