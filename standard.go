package turing

import (
	"bytes"
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type (
	// Our StandardTable is a wrapper for Turing's standard forms
	StandardTable struct {
		// The first is input to a machine that has been standardized.
		MachineInput MachineInput
		// The second is a mapping from our new symbols to our symbols.
		// This is not essential but helps with debugging and testing.
		SymbolMap SymbolMap
		// Turing's Standard Description (S.D.)
		StandardDescription StandardDescription
		// Turing's Description Number (D.N.)
		DescriptionNumber DescriptionNumber
	}

	// Struct to hold shared values when standardizing MachineInput
	standardTableCreator struct {
		input                 MachineInput
		mConfigurationNames   map[string]string
		nameCount             int
		mConfigurationSymbols map[string]string
		symbolCount           int
	}

	// A map of new symbols to old symbols, used to verify Tape output
	SymbolMap map[string]string

	// A string representing the full m-configuration list of a Machine
	StandardDescription string

	// The StandardDescription converted uniquely to a number
	DescriptionNumber string
)

const (
	mConfigurationNamePrefix   string = "q"
	mConfigurationSymbolPrefix string = "S"

	a         byte = 'A'
	c         byte = 'C'
	d         byte = 'D'
	l         byte = 'L'
	r         byte = 'R'
	n         byte = 'N'
	semicolon byte = ';'
)

var (
	sdCharToDNInt = map[byte]int{
		a:         1,
		c:         2,
		d:         3,
		l:         4,
		r:         5,
		n:         6,
		semicolon: 7,
	}

	dnIntToSDChar = map[int]byte{
		1: a,
		2: c,
		3: d,
		4: l,
		5: r,
		6: n,
		7: semicolon,
	}
)

// Standardizes MachineInput so it conforms to Turing's standard form.
func NewStandardTable(input MachineInput) StandardTable {
	s := &standardTableCreator{
		input: input,
	}

	return s.standardize()
}

// Converts a Machine to a Machine that conforms to Turing's standard form.
func (s *standardTableCreator) standardize() StandardTable {
	// The new m-configurations of the machine
	standardMConfigurations := []MConfiguration{}

	// Turing prefers a format where ` ` (None) is S0, `0` is S1, `1` is S2 and so on
	// This ensures ` ` (None) comes first
	s.newMConfigurationSymbol(none)

	// Every m-configuration will be rewritten and potentially introduce further m-configurations
	for _, mConfiguration := range s.input.MConfigurations {
		// Enumerate all symbols for the m-configuration in standard form
		symbols := s.expandStandardSymbols(mConfiguration.Symbols)

		// Split out the operations so they satisfy Turing's acceptable forms:
		// (E), (E, R), (E, L), (Pa), (Pa, R), (Pa, L), (R), (L), (<Nothing>)
		printOperations, moveOperations := s.expandStandardOperations(mConfiguration.Operations)

		// Standardize m-configuration name
		name := s.newMConfigurationName(mConfiguration.Name)

		// Standardize final m-configuration
		finalMConfiguration := s.newMConfigurationName(mConfiguration.FinalMConfiguration)

		// For each symbol, make identical m-configurations
		for _, currentSymbol := range symbols {

			// For each operation for this symbol, calculate the m-configuration
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
					nextName = s.newHiddenMConfigurationName()
					calculatedFinalMConfiguration = nextName
				}

				// If we intend to print a 'Noop', just use the current symbol
				calculatedPrintOperation := s.calculateStandardPrintOperation(printOperations[i], currentSymbol)

				if i == 0 {
					// Only one m-configuration needed
					standardMConfigurations = append(standardMConfigurations, MConfiguration{
						Name:                calculatedName,
						Symbols:             []string{currentSymbol},
						Operations:          []string{calculatedPrintOperation, moveOperations[i]},
						FinalMConfiguration: calculatedFinalMConfiguration,
					})
				} else {
					// When we are in hidden states, we get to the final m-configuration no matter what
					// This means we need to account for all symbols
					for _, calculatedSymbol := range append(s.input.PossibleSymbols, none) {
						// If we intend to print a 'Noop', just use the current symbol
						calculatedPrintOperation := s.calculateStandardPrintOperation(printOperations[i], s.newMConfigurationSymbol(calculatedSymbol))

						standardMConfigurations = append(standardMConfigurations, MConfiguration{
							Name:                calculatedName,
							Symbols:             []string{s.newMConfigurationSymbol(calculatedSymbol)},
							Operations:          []string{calculatedPrintOperation, moveOperations[i]},
							FinalMConfiguration: calculatedFinalMConfiguration,
						})
					}
				}
			}
		}
	}

	machineInput := MachineInput{
		MConfigurations:        standardMConfigurations,
		Tape:                   s.newTape(),
		StartingMConfiguration: s.newStartingMConfiguration(),
		PossibleSymbols:        s.newMConfigurationSymbols(),
		NoneSymbol:             s.newMConfigurationSymbol(none),
	}
	sd := toStandardDescription(machineInput)
	dn := toDescriptionNumber(sd)

	return StandardTable{
		MachineInput:        machineInput,
		SymbolMap:           s.reverseMConfigurationSymbols(),
		StandardDescription: sd,
		DescriptionNumber:   dn,
	}
}

// Expands and standardizes the list of symbols (to the form S0, S1, ..., etc.)
func (s *standardTableCreator) expandStandardSymbols(originalSymbols []string) []string {
	// First loop required for multiple Not scenario
	notSymbols := []string{}
	for _, symbol := range originalSymbols {
		if strings.Contains(symbol, not) {
			notSymbols = append(notSymbols, symbol[1:])
		}
	}

	symbols := []string{}
	for _, symbol := range originalSymbols {
		// To support `!` (Not), `*` (Any), etc. we may need multiple m-configurations for this one particular row
		if strings.Contains(symbol, not) {
			for _, possibleSymbol := range s.input.PossibleSymbols {
				if !slices.Contains(notSymbols, possibleSymbol) && !slices.Contains(symbols, possibleSymbol) {
					symbols = append(symbols, s.newMConfigurationSymbol(possibleSymbol))
				}
			}
		} else if symbol == any {
			for _, possibleSymbol := range s.input.PossibleSymbols {
				symbols = append(symbols, s.newMConfigurationSymbol(possibleSymbol))
			}
		} else {
			symbols = append(symbols, s.newMConfigurationSymbol(symbol))
		}
	}
	return symbols
}

// Standardizes the list of options to the form Turing prefers (exactly one Print and one Move operation)
// These are returned in two slices - the Print operation slice and the Move operation slice
func (s *standardTableCreator) expandStandardOperations(originalOperations []string) ([]string, []string) {
	printOperations := []string{}
	moveOperations := []string{}
	if len(originalOperations) == 0 {
		printOperations = append(printOperations, string(printOp))
		moveOperations = append(moveOperations, string(n))
	} else {
		lookingForPrint := true
		for i, operation := range originalOperations {
			operationCode := operationCode(operation[0])
			if lookingForPrint {
				if operationCode == printOp {
					symbol := string(operation[1:])
					var printOperation strings.Builder
					printOperation.WriteByte(byte(printOp))
					printOperation.WriteString(s.newMConfigurationSymbol(symbol))
					printOperations = append(printOperations, printOperation.String())
					lookingForPrint = false
					if i == len(originalOperations)-1 {
						moveOperations = append(moveOperations, string(n))
					}
				} else if operationCode == eraseOp {
					var printOperation strings.Builder
					printOperation.WriteByte(byte(printOp))
					printOperation.WriteString(s.newMConfigurationSymbol(none))
					printOperations = append(printOperations, printOperation.String())
					lookingForPrint = false
					if i == len(originalOperations)-1 {
						moveOperations = append(moveOperations, string(n))
					}
				} else {
					var printOperation strings.Builder
					printOperation.WriteByte(byte(printOp))
					// Printing the current symbol is essentially a Print noop
					// We encode this by just including `P` with no symbol
					printOperations = append(printOperations, printOperation.String())
					moveOperations = append(moveOperations, string(operationCode))
				}
			} else {
				if operationCode == leftOp || operationCode == rightOp {
					moveOperations = append(moveOperations, string(operationCode))
				} else {
					moveOperations = append(moveOperations, string(n))
				}
				lookingForPrint = true
			}
		}
	}
	return printOperations, moveOperations
}

// Returns the standardized print operation, taking into account the "Noop" print situation
func (st *standardTableCreator) calculateStandardPrintOperation(printOperation string, currentSymbol string) string {
	var calculatedPrintOperation string
	// A Print operation with no value (just `P`) means we should perform a "Noop" print,
	// meaning just print whatever symbol is already on the scanned square
	if printOperation == string(printOp) {
		var calculatedPrintOperationBuilder strings.Builder
		calculatedPrintOperationBuilder.WriteByte(byte(printOp))
		calculatedPrintOperationBuilder.WriteString(currentSymbol)
		calculatedPrintOperation = calculatedPrintOperationBuilder.String()
	} else {
		calculatedPrintOperation = printOperation
	}
	return calculatedPrintOperation
}

// Returns the standardized m-configuration name (of the form q1, q2, ..., etc.), and stores it for deduping
func (s *standardTableCreator) newMConfigurationName(name string) string {
	if s.mConfigurationNames == nil {
		s.mConfigurationNames = map[string]string{}
		s.nameCount++
	}
	newName, ok := s.mConfigurationNames[name]
	if !ok {
		newName = mConfigurationNamePrefix + strconv.Itoa(s.nameCount)
		s.nameCount++
		s.mConfigurationNames[name] = newName
	}
	return newName
}

// Returns a new standardized m-configuration name (of the form q1, q2, ..., etc.), without storing it
func (s *standardTableCreator) newHiddenMConfigurationName() string {
	if s.mConfigurationNames == nil {
		s.mConfigurationNames = map[string]string{}
		s.nameCount++
	}
	newName := mConfigurationNamePrefix + strconv.Itoa(s.nameCount)
	s.nameCount++
	return newName
}

// Returns the standardized symbol name (of the form S0, S1, ..., etc.), and stores it for deduping
func (s *standardTableCreator) newMConfigurationSymbol(symbol string) string {
	if s.mConfigurationSymbols == nil {
		s.mConfigurationSymbols = map[string]string{}
	}
	newSymbol, ok := s.mConfigurationSymbols[symbol]
	if !ok {
		newSymbol = mConfigurationSymbolPrefix + strconv.Itoa(s.symbolCount)
		s.symbolCount++
		s.mConfigurationSymbols[symbol] = newSymbol
	}
	return newSymbol
}

// Returns the full set of new symbols in a slice
func (s *standardTableCreator) newMConfigurationSymbols() []string {
	symbols := []string{}
	for _, v := range s.mConfigurationSymbols {
		symbols = append(symbols, v)
	}
	return symbols
}

// Returns a map from new symbols to old symbols
func (s *standardTableCreator) reverseMConfigurationSymbols() SymbolMap {
	mConfigurationSymbols := SymbolMap{}
	for k, v := range s.mConfigurationSymbols {
		mConfigurationSymbols[v] = k
	}
	return mConfigurationSymbols
}

// Returns the starting m-configuration for the standardize machine
func (s *standardTableCreator) newStartingMConfiguration() string {
	if len(s.input.StartingMConfiguration) == 0 {
		return ""
	} else {
		return s.newMConfigurationName(s.input.StartingMConfiguration)
	}
}

// Returns a standardized tape for the machine
func (s *standardTableCreator) newTape() []string {
	newTape := []string{}
	for _, square := range s.input.Tape {
		newTape = append(newTape, s.mConfigurationSymbols[square])
	}
	return newTape
}

// Translates a tape to the original symbol set.
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
		// There is a bug in original paper, each m-configuration should begin with a semi-colon.
		standardDescription.WriteByte(semicolon)

		// Name is `DAAA`
		standardDescription.WriteByte(d)
		nameSuffix := standardMConfiguration.Name[1:]
		nameNum, _ := strconv.Atoi(nameSuffix)
		standardDescription.Write(bytes.Repeat([]byte{a}, nameNum))

		// Symbol is `DCCC`
		standardDescription.WriteByte(d)
		symbolSuffix := standardMConfiguration.Symbols[0][1:]
		symbolNum, _ := strconv.Atoi(symbolSuffix)
		standardDescription.Write(bytes.Repeat([]byte{c}, symbolNum))

		// Print is also is `DCCC`
		standardDescription.WriteByte(d)
		printOperationSuffix := standardMConfiguration.Operations[0][2:]
		printOperationNum, _ := strconv.Atoi(printOperationSuffix)
		standardDescription.Write(bytes.Repeat([]byte{c}, printOperationNum))

		// Move Operations is `L`, `R`, or `N`
		standardDescription.WriteString(standardMConfiguration.Operations[1])

		// Final Configuration is also `DAAA`
		standardDescription.WriteByte(d)
		finalMConfigurationSuffix := standardMConfiguration.FinalMConfiguration[1:]
		finalMConfigurationNum, _ := strconv.Atoi(finalMConfigurationSuffix)
		standardDescription.Write(bytes.Repeat([]byte{a}, finalMConfigurationNum))
	}

	return StandardDescription(standardDescription.String())
}

// Conversts a S.D. to a D.N.
func toDescriptionNumber(sd StandardDescription) DescriptionNumber {
	var descriptionNumber strings.Builder
	for _, char := range []byte(sd) {
		descriptionNumber.WriteString(strconv.Itoa(sdCharToDNInt[char]))
	}
	return DescriptionNumber(descriptionNumber.String())
}

// Converts a D.N. to a Machine. Returns an error if the D.N. is not well-defined.
func NewMachineFromDescriptionNumber(dn DescriptionNumber) (MachineInput, error) {
	matched, _ := regexp.MatchString("^(?:731+32*32*[456]31+)+$", string(dn))
	if !matched {
		return MachineInput{}, errors.New("not a well defined Description Number")
	}

	var standardDescription strings.Builder
	for _, char := range []byte(dn) {
		i, err := strconv.Atoi(string(char))
		if err != nil {
			return MachineInput{}, err
		}
		standardDescription.WriteString(string(dnIntToSDChar[i]))
	}

	mConfigurations := []MConfiguration{}
	for _, section := range strings.Split(standardDescription.String()[1:], string(semicolon)) {
		subsections := strings.Split(section[1:], string(d))
		name := mConfigurationNamePrefix + strconv.Itoa(len(subsections[0]))
		symbol := mConfigurationSymbolPrefix + strconv.Itoa(len(subsections[1]))
		printOperation := string(printOp) + mConfigurationSymbolPrefix + strconv.Itoa(len(subsections[2])-1)
		moveOperation := string(subsections[2][len(subsections[2])-1])
		finalMConfiguration := mConfigurationNamePrefix + strconv.Itoa(len(subsections[len(subsections)-1]))

		mConfigurations = append(mConfigurations, MConfiguration{
			Name:                name,
			Symbols:             []string{symbol},
			Operations:          []string{printOperation, moveOperation},
			FinalMConfiguration: finalMConfiguration,
		})
	}

	possibleSymbols := []string{}
	for i := 0; i <= maxCharsRepeated([]byte(standardDescription.String()), c); i++ {
		possibleSymbols = append(possibleSymbols, mConfigurationSymbolPrefix+strconv.Itoa(i))
	}

	return MachineInput{
		MConfigurations: mConfigurations,
		PossibleSymbols: possibleSymbols,
		NoneSymbol:      mConfigurationSymbolPrefix + strconv.Itoa(0),
	}, nil
}

func maxCharsRepeated(s []byte, ch byte) int {
	var maxCount int
	var runningCount int
	for _, b := range s {
		if b == ch {
			runningCount += 1
			if runningCount > maxCount {
				maxCount = runningCount
			}
		} else {
			runningCount = 0
		}
	}
	return maxCount
}
