package turing

import (
	"slices"
	"strconv"
	"strings"
)

type AbbreviatedTable struct {
	Machine
	mConfigurationCount      int
	newMConfigurationNames   map[string]string
	wasAlreadyInterpretedMap map[string]bool
	newMConfigurations       []MConfiguration
}

const (
	functionOpen           string = "("
	functionClose          string = ")"
	functionParamDelimiter string = ","
)

// Converts an AbbreviatedTable to a valid Machine, which will contain no skeleton tables
func (at *AbbreviatedTable) ToMachine() *Machine {
	// For each MConfiguration that is not an MFunction, begin interpreting
	for _, mConfiguration := range at.MConfigurations {
		if !strings.Contains(mConfiguration.Name, functionOpen) {
			at.interpretMFunction(mConfiguration.Name, []string{})
		}
	}

	var startingMConfiguration string
	if len(at.StartingMConfiguration) != 0 {
		startingMConfiguration = at.newMConfigurationName(at.StartingMConfiguration, []string{})
	}

	return &Machine{
		MConfigurations:        at.sortedNewMConfigurations(),
		Tape:                   at.Tape,
		StartingMConfiguration: startingMConfiguration,
		PossibleSymbols:        at.PossibleSymbols,
		NoneSymbol:             at.NoneSymbol,
		Debug:                  at.Debug,
	}
}

// Given an MFunction call in the form `f(a, b, x(y, z))`, interpret recursively
func (at *AbbreviatedTable) interpretMFunction(name string, params []string) string {
	// Standardize MConfiguration names
	newMConfigurationName := at.newMConfigurationName(name, params)

	// For each MFunction call signature, we only need to interpret once
	if at.wasAlreadyInterpreted(name, params) {
		return newMConfigurationName
	} else {
		at.markAsInterpreted(name, params)
	}

	// For each MFunction that matches our name and param length, recursively interpret
	for _, mFunction := range at.findMFunctions(name, len(params)) {
		// Retrieve the MFunction's parameter names
		_, mFunctionParams := parseMFunction(mFunction.Name)

		// This bit only required to support the scenario Turing outlines in his `c1` (copy) MFunction
		// In this scenario the supplied symbol is read and used as a parameter for operations or the
		// Final MConfiguration.
		symbolValues := []string{}
		symbolParam, isSymbolParam := at.isSymbolParam(mFunction.Symbols, mFunctionParams)
		if isSymbolParam {
			for _, possibleSymbol := range append(at.PossibleSymbols, None) {
				symbolValues = append(symbolValues, possibleSymbol)
			}
		} else {
			symbolValues = append(symbolValues, "")
		}

		for _, symbolValue := range symbolValues {
			// Create a substitution map from parameter names to the provided parameter values
			if isSymbolParam {
				mFunctionParams = append(mFunctionParams, symbolParam)
				params = append(params, symbolValue)
			}
			substitutionMap := createSubstitutionMap(mFunctionParams, params)

			// Parse the FinalMConfiguration (it may be a function)
			finalMFunctionName, finalMFunctionParams := parseMFunction(mFunction.FinalMConfiguration)

			// Perform substitutions on both the FinalMConfiguration name and params
			substitutedFinalMFunctionName := at.substituteFinalMConfigurationName(finalMFunctionName, substitutionMap)
			substitutedFinalMFunctionParams := at.substituteFinalMConfigurationParams(finalMFunctionParams, substitutionMap)

			// This block recursively attempts to interpret whatever the FinalMConfiguration is (potentially an MFunction to follow)
			var newFinalMConfigurationName string
			if len(substitutedFinalMFunctionParams) == 0 {
				// If there were no params, we still might have substituted to an MFunction
				// If this is the case, we want to parse out the name and params
				substitutedFinalMFunctionNameParsedName, substitutedFinalMFunctionNameParsedParams := parseMFunction(substitutedFinalMFunctionName)
				newFinalMConfigurationName = at.interpretMFunction(substitutedFinalMFunctionNameParsedName, substitutedFinalMFunctionNameParsedParams)
			} else {
				// If there were params, go ahead and use those
				newFinalMConfigurationName = at.interpretMFunction(substitutedFinalMFunctionName, substitutedFinalMFunctionParams)
			}

			// Substitute Symbols and Save MConfiguration
			at.saveMConfiguration(MConfiguration{
				Name:                newMConfigurationName,
				Symbols:             at.substituteSymbols(mFunction.Symbols, substitutionMap),
				Operations:          at.substituteOperations(mFunction.Operations, substitutionMap),
				FinalMConfiguration: newFinalMConfigurationName,
			})
		}
	}

	// Bubble up the Standardized MConfiguration name
	return newMConfigurationName
}

// Finds all MFunctions whose definition matches the name and number of params
func (at *AbbreviatedTable) findMFunctions(name string, numParams int) []MConfiguration {
	mFunctions := []MConfiguration{}
	for _, mFunction := range at.MConfigurations {
		mFunctionName, mFunctionParams := parseMFunction(mFunction.Name)
		if name == mFunctionName && numParams == len(mFunctionParams) {
			mFunctions = append(mFunctions, mFunction)
		}
	}
	return mFunctions
}

func (at *AbbreviatedTable) isSymbolParam(symbols []string, mFunctionParams []string) (string, bool) {
	if len(symbols) != 1 {
		return "", false
	}
	symbol := symbols[0]
	if strings.Contains(symbol, Not) || strings.Contains(symbol, Any) {
		return "", false
	}
	notAPossibleSymbol := !slices.Contains(append(at.PossibleSymbols, None), symbol)
	notAMFunctionParam := !slices.Contains(mFunctionParams, symbol)
	if notAPossibleSymbol && notAMFunctionParam {
		return symbol, true
	}
	return "", false
}

// For the Symbols column of an MFunction, substitute any MFunction params with values
func (at *AbbreviatedTable) substituteSymbols(mFunctionSymbols []string, substitutions map[string]string) []string {
	substitutedSymbols := []string{}
	for _, mFunctionSymbol := range mFunctionSymbols {
		if strings.Contains(mFunctionSymbol, Not) {
			if substitutedSymbol, ok := substitutions[mFunctionSymbol[1:]]; ok {
				substitutedSymbols = append(substitutedSymbols, Not+substitutedSymbol)
			} else {
				substitutedSymbols = append(substitutedSymbols, mFunctionSymbol)
			}
		} else {
			if substitutedSymbol, ok := substitutions[mFunctionSymbol]; ok {
				substitutedSymbols = append(substitutedSymbols, substitutedSymbol)
			} else {
				substitutedSymbols = append(substitutedSymbols, mFunctionSymbol)
			}
		}
	}
	return substitutedSymbols
}

// For the Operations of an MFunction, substitute any MFunction params with values
func (at *AbbreviatedTable) substituteOperations(mFunctionOperations []string, substitutions map[string]string) []string {
	substitutedOperations := []string{}
	for _, mFunctionOperation := range mFunctionOperations {
		switch operationCode(mFunctionOperation[0]) {
		case Print:
			mFunctionOperationSymbol := string(mFunctionOperation[1])
			if substitutedOperation, ok := substitutions[mFunctionOperationSymbol]; ok {
				substitutedOperations = append(substitutedOperations, string(Print)+substitutedOperation)

			} else {
				substitutedOperations = append(substitutedOperations, string(Print)+mFunctionOperationSymbol)
			}
		default:
			substitutedOperations = append(substitutedOperations, mFunctionOperation)
		}
	}
	return substitutedOperations
}

// For a parsed FinalMConfiguration column of an MFunction, attempt to make a parameter substitution if possible for its name
func (at *AbbreviatedTable) substituteFinalMConfigurationName(mFunctionFinalMConfigurationName string, substitutions map[string]string) string {
	if substitutedMFunctionFinalMConfigurationName, ok := substitutions[mFunctionFinalMConfigurationName]; ok {
		return substitutedMFunctionFinalMConfigurationName
	}
	return mFunctionFinalMConfigurationName
}

// For a parsed FinalMConfiguration column of an MFunction, attempt to make a parameter substitution if possible for its values
func (at *AbbreviatedTable) substituteFinalMConfigurationParams(mFunctionFinalMConfigurationParams []string, substitutions map[string]string) []string {
	substitutedMFunctionFinalMConfigurationParams := []string{}
	for _, mFunctionFinalMConfigurationParam := range mFunctionFinalMConfigurationParams {
		potentialInnerName, potentialInnerParams := parseMFunction(mFunctionFinalMConfigurationParam)
		if len(potentialInnerParams) == 0 {
			substitutedMFunctionFinalMConfigurationParams = append(substitutedMFunctionFinalMConfigurationParams, at.substituteFinalMConfigurationName(potentialInnerName, substitutions))
		} else {
			recursiveSubstitution := at.substituteFinalMConfigurationParams(potentialInnerParams, substitutions)
			substitutedMFunctionFinalMConfigurationParams = append(substitutedMFunctionFinalMConfigurationParams, composeMFunction(potentialInnerName, recursiveSubstitution))
		}
	}
	return substitutedMFunctionFinalMConfigurationParams
}

// Saves a new MConfiguration
func (at *AbbreviatedTable) saveMConfiguration(mConfiguration MConfiguration) {
	if at.newMConfigurations == nil {
		at.newMConfigurations = []MConfiguration{}
	}

	at.newMConfigurations = append(at.newMConfigurations, mConfiguration)
}

// Constructs a new unique MConfiguration name
func (at *AbbreviatedTable) newMConfigurationName(mFunctionName string, mFunctionParams []string) string {
	if at.newMConfigurationNames == nil {
		at.newMConfigurationNames = map[string]string{}
	}

	key := composeMFunction(mFunctionName, mFunctionParams)

	if mConfigurationName, ok := at.newMConfigurationNames[key]; ok {
		return mConfigurationName
	}

	newName := mConfigurationNamePrefix + strconv.Itoa(at.mConfigurationCount)
	at.mConfigurationCount++
	at.newMConfigurationNames[key] = newName
	return newName
}

// Returns true if this MFunction signature was already interpreted
func (at *AbbreviatedTable) wasAlreadyInterpreted(mFunctionName string, mFunctionParams []string) bool {
	if at.wasAlreadyInterpretedMap == nil {
		at.wasAlreadyInterpretedMap = map[string]bool{}
	}

	key := composeMFunction(mFunctionName, mFunctionParams)

	if _, ok := at.wasAlreadyInterpretedMap[key]; ok {
		return true
	}
	return false
}

// Marks an MFunction signature as interpreted
func (at *AbbreviatedTable) markAsInterpreted(mFunctionName string, mFunctionParams []string) {
	if at.wasAlreadyInterpretedMap == nil {
		at.wasAlreadyInterpretedMap = map[string]bool{}
	}

	key := composeMFunction(mFunctionName, mFunctionParams)

	at.wasAlreadyInterpretedMap[key] = true
}

// Returns a sorted slice of the stored interpreted MConfigurations
func (at *AbbreviatedTable) sortedNewMConfigurations() []MConfiguration {
	slices.SortFunc(at.newMConfigurations, func(a, b MConfiguration) int {
		aInt, _ := strconv.Atoi(a.Name[1:])
		bInt, _ := strconv.Atoi(b.Name[1:])
		return aInt - bInt
	})
	return at.newMConfigurations
}

// Parses an MFunction of the form "f(a, b, x(y, z))" into name "f" and params ["a", "b", "x(y, z)"]
func parseMFunction(mFunction string) (string, []string) {
	open := strings.Index(mFunction, functionOpen)
	if open < 0 {
		return mFunction, []string{}
	}

	mFunctionName := mFunction[0:open]
	params := []string{}

	var currentParam strings.Builder
	var recursiveCount int
	for _, char := range mFunction[open+1 : len(mFunction)-1] {
		charAsString := string(char)
		if recursiveCount > 0 || (charAsString != None && charAsString != functionParamDelimiter) {
			currentParam.WriteRune(char)
		}
		if charAsString == functionOpen {
			recursiveCount++
		}
		if charAsString == functionClose {
			recursiveCount--
		}
		if recursiveCount == 0 && charAsString == functionParamDelimiter {
			// Handles the scenario where we want to use ` ` (None) as a parameter
			if currentParam.Len() == 0 {
				currentParam.WriteString(None)
			}
			params = append(params, currentParam.String())
			currentParam.Reset()
		}
	}

	// Handles the scenario where we want to use ` ` (None) as a parameter
	if currentParam.Len() == 0 {
		currentParam.WriteString(None)
	}
	params = append(params, currentParam.String())

	return mFunctionName, params
}

// Composes an MFunction of name "f" and params ["a", "b", "x(y, z)"] into the form "f(a, b, x(y, z))"
func composeMFunction(name string, params []string) string {
	var mFunction strings.Builder
	mFunction.WriteString(name)
	if len(params) > 0 {
		mFunction.WriteString(functionOpen)
		mFunction.WriteString(strings.Join(params, functionParamDelimiter))
		mFunction.WriteString(functionClose)
	}
	return mFunction.String()
}

// Zips up two arrays of strings into a map
func createSubstitutionMap(keys []string, values []string) map[string]string {
	substitutionMap := map[string]string{}
	for i, key := range keys {
		substitutionMap[key] = values[i]
	}
	return substitutionMap
}
