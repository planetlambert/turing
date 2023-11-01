package turing

import (
	"slices"
	"strconv"
	"strings"
)

// The m-configurations (or rather, m-functions) below are helper
// functions to be used eventually in Turing's Universal Machine
var (
	// From the m-configuration `f` the machine finds the
	// symbol of form `a` which is farthest to the left (the "first a")
	// and the m-configuration then becomes `C`. If there is no `a`
	// then the m-configuration becomes `B`.
	findLeftMost = []MConfiguration{
		{"f(C, B, a)", []string{"e"}, []string{"L"}, "f1(C, B, a)"},
		{"f(C, B, a)", []string{"!e", " "}, []string{"L"}, "f(C, B, a)"},
		{"f1(C, B, a)", []string{"a"}, []string{}, "C"},
		{"f1(C, B, a)", []string{"!a"}, []string{"R"}, "f1(C, B, a)"},
		{"f1(C, B, a)", []string{" "}, []string{"R"}, "f2(C, B, a)"},
		{"f2(C, B, a)", []string{"a"}, []string{}, "C"},
		{"f2(C, B, a)", []string{"!a"}, []string{"R"}, "f1(C, B, a)"},
		{"f2(C, B, a)", []string{" "}, []string{"R"}, "B"},
	}

	// From `e(C, B, a)` the first `a` is erased and -> `C`.
	// If there is no `a` -> `B`.
	// From `e(B, a)` all letters `a` are erased and -> `B`.
	erase = []MConfiguration{
		{"e(C, B, a)", []string{"*", " "}, []string{}, "f(e1(C, B, a), B, a)"},
		{"e1(C, B, a)", []string{"*", " "}, []string{"E"}, "C"},
		{"e(B, a)", []string{"*", " "}, []string{}, "e(e(B, a), B, a)"},
	}

	// From `pe(C, b)` the machine prints `b` at the end of the sequence
	// of symbols and -> `C`
	printAtTheEnd = []MConfiguration{
		{"pe(C, b)", []string{"*", " "}, []string{}, "f(pe1(C, b), C, e)"},
		{"pe1(C, b)", []string{"*"}, []string{"R", "R"}, "pe1(C, b)"},
		{"pe1(C, b)", []string{" "}, []string{"Pb"}, "C"},
	}

	// From `fl(C, B, a)` it does the same as for `f(C, B, a)`,
	// but moves to the left before -> `C`
	findLeft = []MConfiguration{
		{"l(C)", []string{"*", " "}, []string{"L"}, "C"},
		{"fl(C, B, a)", []string{"*", " "}, []string{}, "f(l(C), B, a)"},
	}

	// From `fr(C, B, a)` it does the same as for `f(C, B, a)`,
	// but moves to the right before -> `C`
	findRight = []MConfiguration{
		{"r(C)", []string{"*", " "}, []string{"R"}, "C"},
		{"fr(C, B, a)", []string{"*", " "}, []string{}, "f(r(C), B, a)"},
	}

	// `c(C, B, a)`. The machine writes at the end the first symbol
	// marked `a` and -> `C`
	copy = []MConfiguration{
		{"c(C, B, a)", []string{"*", " "}, []string{}, "fl(c1(C), B, a)"},
		{"c1(C)", []string{"_b"}, []string{}, "pe(C, _b)"},
	}

	// `ce(B, a)`. The machine copies down in order at the end
	// all symbols marked `a` and erases the letters `a` -> `B`
	copyAndErase = []MConfiguration{
		{"ce(C, B, a)", []string{"*", " "}, []string{}, "c(e(C, B, a), B, a)"},
		{"ce(B, a)", []string{"*", " "}, []string{}, "ce(ce(B, a), B, a)"},
	}

	// `re(C, B, a, b)`. The machine replaces the first `a` by `b` and
	// -> `C` (-> `B` if there is no `a`).
	// `re(B, a, b)`. The machine replaces all letters `a` by `b` -> `B`
	replace = []MConfiguration{
		{"re(C, B, a, b)", []string{"*", " "}, []string{}, "f(re1(C, B, a, b), b, a)"},
		{"re1(C, B, a, b)", []string{"*", " "}, []string{"E", "Pb"}, "C"},
		{"re(B, a, b)", []string{"*", " "}, []string{}, "re(re(B, a, b), B, a, b)"},
	}

	// `cr(B, a)` differs from `ce(B, a)` only in that the letters `a` are not erased.
	// The m-configuration `cr(B, a)` is taken up when no letters `b` are on the tape.
	copyAndReplace = []MConfiguration{
		{"cr(C, B, a, b)", []string{"*", " "}, []string{}, "c(re(C, B, a, b), B, a)"},
		{"cr(B, a, b)", []string{"*", " "}, []string{}, "cr(cr(B, a, b), re(B, a, b), a, b)"},
	}

	// The first symbol marked `a` and the first marked `b` are compared.
	// If there is neither `a` nor `b` -> E. If there are both and the symbols are alike,
	// -> `C`. Otherwise -> `A`.
	compare = []MConfiguration{
		{"cp(C, A, E, a, b)", []string{"*", " "}, []string{}, "fl(cp1(C, A, b), f(A, E, b), a)"},
		{"cp1(C, A, b)", []string{"_y"}, []string{}, "fl(cp2(C, A, _y), A, b)"},
		{"cp2(C, A, y)", []string{"y"}, []string{}, "C"},
		{"cp2(C, A, y)", []string{"!y", " "}, []string{}, "A"},
	}

	// `cpe(C, A, E, a, b)` differs from `cp(C, A, E, a, b)` in that in the case when there is
	// a similarity the first `a` and `b` are erased.
	// `cpe(A, E, a, b)`. The sequence of symbols marked `a` is compared with the sequence
	// marked `b`. -> `C` if they are similar. Otherwise -> `A`. Some of the symbols `a` and `b` are erased.
	compareAndErase = []MConfiguration{
		{"cpe(C, A, E, a, b)", []string{"*", " "}, []string{}, "cp(e(e(C, C, b), C, a), A, E, a, b)"},
		{"cpe(A, E, a, b)", []string{"*", " "}, []string{}, "cpe(cpe(A, E, a, b), A, E, a, b)"},
	}

	// `g(C, a)`. The machine finds the last symbol of form `a` -> `C`.
	findRightMost = []MConfiguration{
		{"g(C)", []string{"*"}, []string{"R"}, "g(C)"},
		{"g(C)", []string{" "}, []string{"R"}, "g1(C)"},
		{"g1(C)", []string{"*"}, []string{"R"}, "g(C)"},
		{"g1(C)", []string{" "}, []string{}, "C"},
		{"g(C, a)", []string{"*", " "}, []string{}, "g(g1(C, a))"},
		{"g1(C, a)", []string{"a"}, []string{}, "C"},
		{"g1(C, a)", []string{"!a", " "}, []string{"L"}, "g1(C, a)"},
	}

	// `pe2(C, a, b)`. The machine prints `a b` at the end.
	printAtTheEnd2 = []MConfiguration{
		{"pe2(C, a, b)", []string{"*", " "}, []string{}, "pe(pe(C, b), a)"},
	}

	// `ce3(B, a, b, y)`. The machine copies down at the end first the symbols
	// marked `a` then those marked `b`, and finally those marked `y`.
	// It erases the symbols `a`, `b`, `y`.
	copyAndErase2 = []MConfiguration{
		{"ce2(B, a, b)", []string{"*", " "}, []string{}, "ce(ce(B, b), a)"},
		{"ce3(B, a, b, y)", []string{"*", " "}, []string{}, "ce(ce2(B, b, y), a)"},
		{"ce4(B, a, b, y, z)", []string{"*", " "}, []string{}, "ce(ce3(B, b, y, z), a)"},
		{"ce5(B, a, b, y, z, w)", []string{"*", " "}, []string{}, "ce(ce4(B, b, y, z, w), a)"},
	}

	// From `e(C)` the marks are erased from all marked symbols -> `C`
	eraseAll = []MConfiguration{
		{"e(C)", []string{"e"}, []string{"R"}, "e1(C)"},
		{"e(C)", []string{"!e", " "}, []string{"L"}, "e(C)"},
		{"e1(C)", []string{"*"}, []string{"R", "E", "R"}, "e1(C)"},
		{"e1(C)", []string{" "}, []string{}, "C"},
	}
)

// Returns all helper functions
func allhelperFunctions() []MConfiguration {
	helperFunctions := []MConfiguration{}
	helperFunctions = append(helperFunctions, findLeftMost...)
	helperFunctions = append(helperFunctions, erase...)
	helperFunctions = append(helperFunctions, printAtTheEnd...)
	helperFunctions = append(helperFunctions, findLeft...)
	helperFunctions = append(helperFunctions, findRight...)
	helperFunctions = append(helperFunctions, copy...)
	helperFunctions = append(helperFunctions, copyAndErase...)
	helperFunctions = append(helperFunctions, replace...)
	helperFunctions = append(helperFunctions, copyAndReplace...)
	helperFunctions = append(helperFunctions, compare...)
	helperFunctions = append(helperFunctions, compareAndErase...)
	helperFunctions = append(helperFunctions, findRightMost...)
	helperFunctions = append(helperFunctions, printAtTheEnd2...)
	helperFunctions = append(helperFunctions, printAtTheEnd2...)
	helperFunctions = append(helperFunctions, copyAndErase2...)
	helperFunctions = append(helperFunctions, eraseAll...)
	return helperFunctions
}

// Input for an Abbreviated Table
type AbbreviatedTableInput MachineInput

// Gives MachineInput for the abbreviated table. This requires "compiling" the abbreviated table.
func NewAbbreviatedTable(input AbbreviatedTableInput) MachineInput {
	at := &abbreviatedTable{
		input: input,
	}

	return at.toMachineInput()
}

// Helper struct to compile the abbreviated table
type abbreviatedTable struct {
	input                    AbbreviatedTableInput
	mConfigurationCount      int
	newMConfigurationNames   map[string]string
	wasAlreadyInterpretedMap map[string]bool
	newMConfigurations       []MConfiguration
}

// Used when parsing m-functions
const (
	functionOpen           string = "("
	functionClose          string = ")"
	functionParamDelimiter string = ","
)

// Converts an AbbreviatedTable to a valid Machine, which will contain no skeleton tables
func (at *abbreviatedTable) toMachineInput() MachineInput {
	// For each m-configuration that is not an m-function, begin interpreting
	for _, mConfiguration := range at.input.MConfigurations {
		if !strings.Contains(mConfiguration.Name, functionOpen) {
			at.interpretMFunction(mConfiguration.Name, []string{})
		}
	}

	var startingMConfiguration string
	if len(at.input.StartingMConfiguration) != 0 {
		startingMConfiguration = at.newMConfigurationName(at.input.StartingMConfiguration, []string{})
	}

	return MachineInput{
		MConfigurations:        at.sortedNewMConfigurations(),
		Tape:                   at.input.Tape,
		StartingMConfiguration: startingMConfiguration,
		PossibleSymbols:        at.input.PossibleSymbols,
		NoneSymbol:             at.input.NoneSymbol,
		Debug:                  at.input.Debug,
	}
}

// Given an m-function call in the form `f(a, b, x(y, z))`, interpret recursively
func (at *abbreviatedTable) interpretMFunction(name string, params []string) string {
	// Standardize m-configuration names
	newMConfigurationName := at.newMConfigurationName(name, params)

	// For each m-function call signature, we only need to interpret once
	if at.wasAlreadyInterpreted(name, params) {
		return newMConfigurationName
	} else {
		at.markAsInterpreted(name, params)
	}

	// For each m-function that matches our name and param length, recursively interpret
	for _, mFunction := range at.findMFunctions(name, len(params)) {
		// Retrieve the m-function's parameter names
		_, mFunctionParams := parseMFunction(mFunction.Name)

		// This bit only required to support the scenario Turing outlines in his `c1` (copy) m-function
		// In this scenario the supplied symbol is read and used as a parameter for operations or the
		// final m-configuration.
		symbolValues := []string{}
		symbolParam, isSymbolParam := at.isSymbolParam(mFunction.Symbols, mFunctionParams)
		if isSymbolParam {
			for _, possibleSymbol := range append(at.input.PossibleSymbols, none) {
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

			// Parse the final m-configuration (it may be a function)
			finalMFunctionName, finalMFunctionParams := parseMFunction(mFunction.FinalMConfiguration)

			// Perform substitutions on both the final m-configuration name and params
			substitutedFinalMFunctionName := at.substituteFinalMConfigurationName(finalMFunctionName, substitutionMap)
			substitutedFinalMFunctionParams := at.substituteFinalMConfigurationParams(finalMFunctionParams, substitutionMap)

			// This block recursively attempts to interpret whatever the final m-configuration is (potentially an m-function to follow)
			var newFinalMConfigurationName string
			if len(substitutedFinalMFunctionParams) == 0 {
				// If there were no params, we still might have substituted to an m-function
				// If this is the case, we want to parse out the name and params
				substitutedFinalMFunctionNameParsedName, substitutedFinalMFunctionNameParsedParams := parseMFunction(substitutedFinalMFunctionName)
				newFinalMConfigurationName = at.interpretMFunction(substitutedFinalMFunctionNameParsedName, substitutedFinalMFunctionNameParsedParams)
			} else {
				// If there were params, go ahead and use those
				newFinalMConfigurationName = at.interpretMFunction(substitutedFinalMFunctionName, substitutedFinalMFunctionParams)
			}

			// Substitute Symbols and Save m-configuration
			at.saveMConfiguration(MConfiguration{
				Name:                newMConfigurationName,
				Symbols:             at.substituteSymbols(mFunction.Symbols, substitutionMap),
				Operations:          at.substituteOperations(mFunction.Operations, substitutionMap),
				FinalMConfiguration: newFinalMConfigurationName,
			})
		}
	}

	// Bubble up the Standardized m-configuration name
	return newMConfigurationName
}

// Finds all m-functions whose definition matches the name and number of params
func (at *abbreviatedTable) findMFunctions(name string, numParams int) []MConfiguration {
	mFunctions := []MConfiguration{}
	for _, mFunction := range at.input.MConfigurations {
		mFunctionName, mFunctionParams := parseMFunction(mFunction.Name)
		if name == mFunctionName && numParams == len(mFunctionParams) {
			mFunctions = append(mFunctions, mFunction)
		}
	}
	return mFunctions
}

func (at *abbreviatedTable) isSymbolParam(symbols []string, mFunctionParams []string) (string, bool) {
	if len(symbols) != 1 {
		return "", false
	}
	symbol := symbols[0]
	if strings.Contains(symbol, not) || strings.Contains(symbol, any) {
		return "", false
	}
	notAPossibleSymbol := !slices.Contains(append(at.input.PossibleSymbols, none), symbol)
	notAMFunctionParam := !slices.Contains(mFunctionParams, symbol)
	if notAPossibleSymbol && notAMFunctionParam {
		return symbol, true
	}
	return "", false
}

// For the Symbols column of an m-function, substitute any m-function params with values
func (at *abbreviatedTable) substituteSymbols(mFunctionSymbols []string, substitutions map[string]string) []string {
	substitutedSymbols := []string{}
	for _, mFunctionSymbol := range mFunctionSymbols {
		if strings.Contains(mFunctionSymbol, not) {
			if substitutedSymbol, ok := substitutions[mFunctionSymbol[1:]]; ok {
				substitutedSymbols = append(substitutedSymbols, not+substitutedSymbol)
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

// For the Operations of an m-function, substitute any m-function params with values
func (at *abbreviatedTable) substituteOperations(mFunctionOperations []string, substitutions map[string]string) []string {
	substitutedOperations := []string{}
	for _, mFunctionOperation := range mFunctionOperations {
		switch operationCode(mFunctionOperation[0]) {
		case printOp:
			mFunctionOperationSymbol := string(mFunctionOperation[1])
			if substitutedOperation, ok := substitutions[mFunctionOperationSymbol]; ok {
				substitutedOperations = append(substitutedOperations, string(printOp)+substitutedOperation)

			} else {
				substitutedOperations = append(substitutedOperations, string(printOp)+mFunctionOperationSymbol)
			}
		default:
			substitutedOperations = append(substitutedOperations, mFunctionOperation)
		}
	}
	return substitutedOperations
}

// For a parsed final m-configuration column of an m-function, attempt to make a parameter substitution if possible for its name
func (at *abbreviatedTable) substituteFinalMConfigurationName(mFunctionFinalMConfigurationName string, substitutions map[string]string) string {
	if substitutedMFunctionFinalMConfigurationName, ok := substitutions[mFunctionFinalMConfigurationName]; ok {
		return substitutedMFunctionFinalMConfigurationName
	}
	return mFunctionFinalMConfigurationName
}

// For a parsed final m-configuration column of an m-function, attempt to make a parameter substitution if possible for its values
func (at *abbreviatedTable) substituteFinalMConfigurationParams(mFunctionFinalMConfigurationParams []string, substitutions map[string]string) []string {
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

// Saves a new m-configuration
func (at *abbreviatedTable) saveMConfiguration(mConfiguration MConfiguration) {
	if at.newMConfigurations == nil {
		at.newMConfigurations = []MConfiguration{}
	}

	at.newMConfigurations = append(at.newMConfigurations, mConfiguration)
}

// Constructs a new unique m-configuration name
func (at *abbreviatedTable) newMConfigurationName(mFunctionName string, mFunctionParams []string) string {
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

// Returns true if this m-function signature was already interpreted
func (at *abbreviatedTable) wasAlreadyInterpreted(mFunctionName string, mFunctionParams []string) bool {
	if at.wasAlreadyInterpretedMap == nil {
		at.wasAlreadyInterpretedMap = map[string]bool{}
	}

	key := composeMFunction(mFunctionName, mFunctionParams)

	if _, ok := at.wasAlreadyInterpretedMap[key]; ok {
		return true
	}
	return false
}

// Marks an m-function signature as interpreted
func (at *abbreviatedTable) markAsInterpreted(mFunctionName string, mFunctionParams []string) {
	if at.wasAlreadyInterpretedMap == nil {
		at.wasAlreadyInterpretedMap = map[string]bool{}
	}

	key := composeMFunction(mFunctionName, mFunctionParams)

	at.wasAlreadyInterpretedMap[key] = true
}

// Returns a sorted slice of the stored interpreted m-configurations
func (at *abbreviatedTable) sortedNewMConfigurations() []MConfiguration {
	slices.SortFunc(at.newMConfigurations, func(a, b MConfiguration) int {
		aInt, _ := strconv.Atoi(a.Name[1:])
		bInt, _ := strconv.Atoi(b.Name[1:])
		return aInt - bInt
	})
	return at.newMConfigurations
}

// Parses an m-function of the form "f(a, b, x(y, z))" into name "f" and params ["a", "b", "x(y, z)"]
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
		if recursiveCount > 0 || (charAsString != none && charAsString != functionParamDelimiter) {
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
				currentParam.WriteString(none)
			}
			params = append(params, currentParam.String())
			currentParam.Reset()
		}
	}

	// Handles the scenario where we want to use ` ` (None) as a parameter
	if currentParam.Len() == 0 {
		currentParam.WriteString(none)
	}
	params = append(params, currentParam.String())

	return mFunctionName, params
}

// Composes an m-function of name "f" and params ["a", "b", "x(y, z)"] into the form "f(a, b, x(y, z))"
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
