package turing

import (
	"bytes"
	"slices"
	"strconv"
	"strings"
)

// Standard Table

type (
	StandardTable struct {
		Machine
		SymbolMap
	}

	SymbolMap map[string]string

	StandardDescription string

	DescriptionNumber string
)

const (
	A         byte = 'A'
	C         byte = 'C'
	D         byte = 'D'
	L         byte = 'L'
	R         byte = 'R'
	N         byte = 'N'
	Semicolon byte = ';'
)

var (
	sDCharToDNInt = map[byte]int{
		A:         1,
		C:         2,
		D:         3,
		L:         4,
		R:         5,
		N:         6,
		Semicolon: 7,
	}
)

// Converts a Machine to a Machine that conforms to Turing's standard form.
func (m *Machine) ToStandardTable() *StandardTable {
	// The new MConfigurations of the machine
	standardMConfigurations := MConfigurations{}

	// Turing prefers a format where ` ` (None) is S0, `0` is S1, `1` is S2 and so on
	// This ensures ` ` (None) comes first
	m.newMConfigurationSymbol(None)

	// Every MConfiguration will be rewritten and potentially introduce further MConfigurations
	for _, mConfiguration := range m.MConfigurations {
		// Enumerate all symbols for the MConfiguration in standard form
		symbols := m.expandStandardSymbols(mConfiguration)

		// Split out the operations so they satisfy Turing's acceptable forms:
		// (E), (E, R), (E, L), (Pa), (Pa, R), (Pa, L), (R), (L), (<Nothing>)
		printOperations, moveOperations := m.expandStandardOperations(mConfiguration)

		// Standardize MConfiguration Name
		name := m.newMConfigurationName(mConfiguration.Name)

		// Standardize FinalMConfiguration
		finalMConfiguration := m.newMConfigurationName(mConfiguration.FinalMConfiguration)

		// For each symbol, make identical MConfigurations
		for _, currentSymbol := range symbols {

			// For each operation for this symbol, calculate the MConfiguration
			var nextName string
			for i := 0; i < len(printOperations); i++ {
				// Only the first operation in a list needs a name that others will recognize
				var calculatedName string
				if i == 0 {
					calculatedName = name
				} else {
					calculatedName = nextName
				}

				// Similarly, only the final operation in a list needs a name that others will recognize
				var calculatedFinalMConfiguration string
				if i == len(printOperations)-1 {
					calculatedFinalMConfiguration = finalMConfiguration
				} else {
					nextName = m.newHiddenMConfigurationName()
					calculatedFinalMConfiguration = nextName
				}

				// If we intend to print a 'Noop', just use the current symbol
				calculatedPrintOperation := m.calculateStandardPrintOperation(printOperations[i], currentSymbol)

				if i == 0 {
					// Only one MConfiguration needed
					standardMConfigurations = append(standardMConfigurations, MConfiguration{
						Name:                calculatedName,
						Symbols:             []string{currentSymbol},
						Operations:          []string{calculatedPrintOperation, moveOperations[i]},
						FinalMConfiguration: calculatedFinalMConfiguration,
					})
				} else {
					// When we are in hidden states, we get to the FinalMConfiguration no matter what
					// This means we need to account for all symbols
					for _, calculatedSymbol := range append(m.PossibleSymbols, None) {
						// If we intend to print a 'Noop', just use the current symbol
						calculatedPrintOperation := m.calculateStandardPrintOperation(printOperations[i], m.newMConfigurationSymbol(calculatedSymbol))

						standardMConfigurations = append(standardMConfigurations, MConfiguration{
							Name:                calculatedName,
							Symbols:             []string{m.newMConfigurationSymbol(calculatedSymbol)},
							Operations:          []string{calculatedPrintOperation, moveOperations[i]},
							FinalMConfiguration: calculatedFinalMConfiguration,
						})
					}
				}
			}
		}
	}

	return &StandardTable{
		Machine: Machine{
			MConfigurations:        standardMConfigurations,
			Tape:                   m.newTape(),
			StartingMConfiguration: m.newStartingMConfiguration(),
			PossibleSymbols:        m.newMConfigurationSymbols(),
			NoneSymbol:             m.newMConfigurationSymbol(None),
		},
		SymbolMap: m.reverseMConfigurationSymbols(),
	}
}

func (m *Machine) expandStandardSymbols(mConfiguration MConfiguration) []string {
	// First loop required for multiple Not scenario
	notSymbols := []string{}
	for _, symbol := range mConfiguration.Symbols {
		if strings.Contains(symbol, Not) {
			notSymbols = append(notSymbols, symbol[1:])
		}
	}

	symbols := []string{}
	for _, symbol := range mConfiguration.Symbols {
		// To support `!` (Not), `*` (Any), etc. we may need multiple MConfigurations for this one particular row
		if strings.Contains(symbol, Not) {
			for _, possibleSymbol := range m.PossibleSymbols {
				if !slices.Contains(notSymbols, possibleSymbol) && !slices.Contains(symbols, possibleSymbol) {
					symbols = append(symbols, m.newMConfigurationSymbol(possibleSymbol))
				}
			}
		} else if symbol == Any {
			for _, possibleSymbol := range m.PossibleSymbols {
				symbols = append(symbols, m.newMConfigurationSymbol(possibleSymbol))
			}
		} else {
			symbols = append(symbols, m.newMConfigurationSymbol(symbol))
		}
	}
	return symbols
}

func (m *Machine) expandStandardOperations(mConfiguration MConfiguration) ([]string, []string) {
	printOperations := []string{}
	moveOperations := []string{}
	if len(mConfiguration.Operations) == 0 {
		printOperations = append(printOperations, string(Print))
		moveOperations = append(moveOperations, string(N))
	} else {
		lookingForPrint := true
		for i, operation := range mConfiguration.Operations {
			operationCode := OperationCode(operation[0])
			if lookingForPrint {
				if operationCode == Print {
					symbol := string(operation[1:])
					var printOperation strings.Builder
					printOperation.WriteByte(byte(Print))
					printOperation.WriteString(m.newMConfigurationSymbol(symbol))
					printOperations = append(printOperations, printOperation.String())
					lookingForPrint = false
					if i == len(mConfiguration.Operations)-1 {
						moveOperations = append(moveOperations, string(N))
					}
				} else if operationCode == Erase {
					var printOperation strings.Builder
					printOperation.WriteByte(byte(Print))
					printOperation.WriteString(m.newMConfigurationSymbol(None))
					printOperations = append(printOperations, printOperation.String())
					lookingForPrint = false
					if i == len(mConfiguration.Operations)-1 {
						moveOperations = append(moveOperations, string(N))
					}
				} else {
					var printOperation strings.Builder
					printOperation.WriteByte(byte(Print))
					// Printing the current symbol is essentially a Print noop
					// We encode this by just including `P` with no symbol
					printOperations = append(printOperations, printOperation.String())
					moveOperations = append(moveOperations, string(operationCode))
				}
			} else {
				if operationCode == Left || operationCode == Right {
					moveOperations = append(moveOperations, string(operationCode))
				} else {
					moveOperations = append(moveOperations, string(N))
				}
				lookingForPrint = true
			}
		}
	}
	return printOperations, moveOperations
}

func (m *Machine) calculateStandardPrintOperation(printOperation string, currentSymbol string) string {
	var calculatedPrintOperation string
	if printOperation == string(Print) {
		var calculatedPrintOperationBuilder strings.Builder
		calculatedPrintOperationBuilder.WriteByte(byte(Print))
		calculatedPrintOperationBuilder.WriteString(currentSymbol)
		calculatedPrintOperation = calculatedPrintOperationBuilder.String()
	} else {
		calculatedPrintOperation = printOperation
	}
	return calculatedPrintOperation
}

func (m *Machine) newMConfigurationName(name string) string {
	if m.mConfigurationNames == nil {
		m.mConfigurationNames = map[string]string{}
		m.nameCount++
	}
	newName, ok := m.mConfigurationNames[name]
	if !ok {
		newName = mConfigurationNamePrefix + strconv.Itoa(m.nameCount)
		m.nameCount++
		m.mConfigurationNames[name] = newName
	}
	return newName
}

func (m *Machine) newHiddenMConfigurationName() string {
	if m.mConfigurationNames == nil {
		m.mConfigurationNames = map[string]string{}
		m.nameCount++
	}
	newName := mConfigurationNamePrefix + strconv.Itoa(m.nameCount)
	m.nameCount++
	return newName
}

func (m *Machine) newMConfigurationSymbol(symbol string) string {
	if m.mConfigurationSymbols == nil {
		m.mConfigurationSymbols = map[string]string{}
	}
	newSymbol, ok := m.mConfigurationSymbols[symbol]
	if !ok {
		newSymbol = mConfigurationSymbolPrefix + strconv.Itoa(m.symbolCount)
		m.symbolCount++
		m.mConfigurationSymbols[symbol] = newSymbol
	}
	return newSymbol
}

func (m *Machine) newMConfigurationSymbols() []string {
	symbols := []string{}
	for _, v := range m.mConfigurationSymbols {
		symbols = append(symbols, v)
	}
	return symbols
}

func (m *Machine) reverseMConfigurationSymbols() SymbolMap {
	mConfigurationSymbols := SymbolMap{}
	for k, v := range m.mConfigurationSymbols {
		mConfigurationSymbols[v] = k
	}
	return mConfigurationSymbols
}

func (m *Machine) newStartingMConfiguration() string {
	if len(m.StartingMConfiguration) == 0 {
		return ""
	} else {
		return m.newMConfigurationName(m.StartingMConfiguration)
	}
}

func (m *Machine) newTape() Tape {
	newTape := Tape{}
	for _, square := range m.Tape {
		newTape = append(newTape, m.mConfigurationSymbols[square])
	}
	return newTape
}

// Translates a tape
func (sm SymbolMap) TranslateTape(tape Tape) string {
	var translatedTape strings.Builder
	for _, square := range tape {
		translatedTape.WriteString(sm[square])
	}
	return translatedTape.String()
}

// Converts a StandardTable to its StandardDescription (S.D.)
func (st *StandardTable) ToStandardDescription() StandardDescription {
	var standardDescription strings.Builder
	for _, standardMConfiguration := range st.Machine.MConfigurations {
		// TODO: Bug in original paper, each Standard Description should begin with
		// a semi-colon.
		standardDescription.WriteByte(Semicolon)

		// Name is `DAAA`
		standardDescription.WriteByte(D)
		nameSuffix := standardMConfiguration.Name[1:]
		nameNum, _ := strconv.Atoi(nameSuffix)
		standardDescription.Write(bytes.Repeat([]byte{A}, nameNum))

		// Symbol is `DCCC`
		standardDescription.WriteByte(D)
		symbolSuffix := standardMConfiguration.Symbols[0][1:]
		symbolNum, _ := strconv.Atoi(symbolSuffix)
		standardDescription.Write(bytes.Repeat([]byte{C}, symbolNum))

		// Print is also is `DCCC`
		standardDescription.WriteByte(D)
		printOperationSuffix := standardMConfiguration.Operations[0][2:]
		printOperationNum, _ := strconv.Atoi(printOperationSuffix)
		standardDescription.Write(bytes.Repeat([]byte{C}, printOperationNum))

		// Move Operations is `L`, `R`, or `N`
		standardDescription.WriteString(standardMConfiguration.Operations[1])

		// Final Configuration is also `DAAA`
		standardDescription.WriteByte(D)
		finalMConfigurationSuffix := standardMConfiguration.FinalMConfiguration[1:]
		finalMConfigurationNum, _ := strconv.Atoi(finalMConfigurationSuffix)
		standardDescription.Write(bytes.Repeat([]byte{A}, finalMConfigurationNum))
	}

	return StandardDescription(standardDescription.String())
}

// Conversts a S.D. to a D.N.
func (sd StandardDescription) ToDescriptionNumber() DescriptionNumber {
	var descriptionNumber strings.Builder
	for _, char := range []byte(sd) {
		descriptionNumber.WriteString(strconv.Itoa(sDCharToDNInt[char]))
	}
	return DescriptionNumber(descriptionNumber.String())
}
