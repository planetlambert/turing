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
		MachineInput        MachineInput
		SymbolMap           SymbolMap
		StandardDescription StandardDescription
		DescriptionNumber   DescriptionNumber

		input MachineInput
		// For use when converting to StandardTable
		mConfigurationNames map[string]string
		// For use when converting to StandardTable
		nameCount int
		// For use when converting to StandardTable
		mConfigurationSymbols map[string]string
		// For use when converting to StandardTable
		symbolCount int
	}

	SymbolMap map[string]string

	StandardDescription string

	DescriptionNumber string
)

const (
	mConfigurationNamePrefix   string = "q"
	mConfigurationSymbolPrefix string = "S"

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

// Standardizes MachineInput so it conforms to Turing's standard form.
func NewStandardTable(input MachineInput) StandardTable {
	st := &StandardTable{
		input: input,
	}

	st.standardize()

	return *st
}

// Converts a Machine to a Machine that conforms to Turing's standard form.
func (st *StandardTable) standardize() {
	// The new MConfigurations of the machine
	standardMConfigurations := []MConfiguration{}

	// Turing prefers a format where ` ` (None) is S0, `0` is S1, `1` is S2 and so on
	// This ensures ` ` (None) comes first
	st.newMConfigurationSymbol(None)

	// Every MConfiguration will be rewritten and potentially introduce further MConfigurations
	for _, mConfiguration := range st.input.MConfigurations {
		// Enumerate all symbols for the MConfiguration in standard form
		symbols := st.expandStandardSymbols(mConfiguration)

		// Split out the operations so they satisfy Turing's acceptable forms:
		// (E), (E, R), (E, L), (Pa), (Pa, R), (Pa, L), (R), (L), (<Nothing>)
		printOperations, moveOperations := st.expandStandardOperations(mConfiguration)

		// Standardize MConfiguration Name
		name := st.newMConfigurationName(mConfiguration.Name)

		// Standardize FinalMConfiguration
		finalMConfiguration := st.newMConfigurationName(mConfiguration.FinalMConfiguration)

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
					nextName = st.newHiddenMConfigurationName()
					calculatedFinalMConfiguration = nextName
				}

				// If we intend to print a 'Noop', just use the current symbol
				calculatedPrintOperation := st.calculateStandardPrintOperation(printOperations[i], currentSymbol)

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
					for _, calculatedSymbol := range append(st.input.PossibleSymbols, None) {
						// If we intend to print a 'Noop', just use the current symbol
						calculatedPrintOperation := st.calculateStandardPrintOperation(printOperations[i], st.newMConfigurationSymbol(calculatedSymbol))

						standardMConfigurations = append(standardMConfigurations, MConfiguration{
							Name:                calculatedName,
							Symbols:             []string{st.newMConfigurationSymbol(calculatedSymbol)},
							Operations:          []string{calculatedPrintOperation, moveOperations[i]},
							FinalMConfiguration: calculatedFinalMConfiguration,
						})
					}
				}
			}
		}
	}

	st.MachineInput = MachineInput{
		MConfigurations:        standardMConfigurations,
		Tape:                   st.newTape(),
		StartingMConfiguration: st.newStartingMConfiguration(),
		PossibleSymbols:        st.newMConfigurationSymbols(),
		NoneSymbol:             st.newMConfigurationSymbol(None),
	}
	st.SymbolMap = st.reverseMConfigurationSymbols()
	st.StandardDescription = toStandardDescription(st.MachineInput)
	st.DescriptionNumber = toDescriptionNumber(st.StandardDescription)
}

func (st *StandardTable) expandStandardSymbols(mConfiguration MConfiguration) []string {
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
			for _, possibleSymbol := range st.input.PossibleSymbols {
				if !slices.Contains(notSymbols, possibleSymbol) && !slices.Contains(symbols, possibleSymbol) {
					symbols = append(symbols, st.newMConfigurationSymbol(possibleSymbol))
				}
			}
		} else if symbol == Any {
			for _, possibleSymbol := range st.input.PossibleSymbols {
				symbols = append(symbols, st.newMConfigurationSymbol(possibleSymbol))
			}
		} else {
			symbols = append(symbols, st.newMConfigurationSymbol(symbol))
		}
	}
	return symbols
}

func (st *StandardTable) expandStandardOperations(mConfiguration MConfiguration) ([]string, []string) {
	printOperations := []string{}
	moveOperations := []string{}
	if len(mConfiguration.Operations) == 0 {
		printOperations = append(printOperations, string(Print))
		moveOperations = append(moveOperations, string(N))
	} else {
		lookingForPrint := true
		for i, operation := range mConfiguration.Operations {
			operationCode := operationCode(operation[0])
			if lookingForPrint {
				if operationCode == Print {
					symbol := string(operation[1:])
					var printOperation strings.Builder
					printOperation.WriteByte(byte(Print))
					printOperation.WriteString(st.newMConfigurationSymbol(symbol))
					printOperations = append(printOperations, printOperation.String())
					lookingForPrint = false
					if i == len(mConfiguration.Operations)-1 {
						moveOperations = append(moveOperations, string(N))
					}
				} else if operationCode == Erase {
					var printOperation strings.Builder
					printOperation.WriteByte(byte(Print))
					printOperation.WriteString(st.newMConfigurationSymbol(None))
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

func (st *StandardTable) calculateStandardPrintOperation(printOperation string, currentSymbol string) string {
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

func (st *StandardTable) newMConfigurationName(name string) string {
	if st.mConfigurationNames == nil {
		st.mConfigurationNames = map[string]string{}
		st.nameCount++
	}
	newName, ok := st.mConfigurationNames[name]
	if !ok {
		newName = mConfigurationNamePrefix + strconv.Itoa(st.nameCount)
		st.nameCount++
		st.mConfigurationNames[name] = newName
	}
	return newName
}

func (st *StandardTable) newHiddenMConfigurationName() string {
	if st.mConfigurationNames == nil {
		st.mConfigurationNames = map[string]string{}
		st.nameCount++
	}
	newName := mConfigurationNamePrefix + strconv.Itoa(st.nameCount)
	st.nameCount++
	return newName
}

func (st *StandardTable) newMConfigurationSymbol(symbol string) string {
	if st.mConfigurationSymbols == nil {
		st.mConfigurationSymbols = map[string]string{}
	}
	newSymbol, ok := st.mConfigurationSymbols[symbol]
	if !ok {
		newSymbol = mConfigurationSymbolPrefix + strconv.Itoa(st.symbolCount)
		st.symbolCount++
		st.mConfigurationSymbols[symbol] = newSymbol
	}
	return newSymbol
}

func (st *StandardTable) newMConfigurationSymbols() []string {
	symbols := []string{}
	for _, v := range st.mConfigurationSymbols {
		symbols = append(symbols, v)
	}
	return symbols
}

func (st *StandardTable) reverseMConfigurationSymbols() SymbolMap {
	mConfigurationSymbols := SymbolMap{}
	for k, v := range st.mConfigurationSymbols {
		mConfigurationSymbols[v] = k
	}
	return mConfigurationSymbols
}

func (st *StandardTable) newStartingMConfiguration() string {
	if len(st.input.StartingMConfiguration) == 0 {
		return ""
	} else {
		return st.newMConfigurationName(st.input.StartingMConfiguration)
	}
}

func (st *StandardTable) newTape() []string {
	newTape := []string{}
	for _, square := range st.input.Tape {
		newTape = append(newTape, st.mConfigurationSymbols[square])
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
func toStandardDescription(input MachineInput) StandardDescription {
	var standardDescription strings.Builder
	for _, standardMConfiguration := range input.MConfigurations {
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
func toDescriptionNumber(sd StandardDescription) DescriptionNumber {
	var descriptionNumber strings.Builder
	for _, char := range []byte(sd) {
		descriptionNumber.WriteString(strconv.Itoa(sDCharToDNInt[char]))
	}
	return DescriptionNumber(descriptionNumber.String())
}
