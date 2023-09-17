package turing

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	// The possible symbols for the universal machine
	possibleSymbolsForUniversalMachine = []string{
		"e", "::", ":", "A", "C", "D", "L", "R", "N", ";", "0", "1", "u", "v", "w", "x", "y", "z",
	}

	// Note: By convention any symbol prefixed with an underscore (i.e. `_y`) is a
	// 'symbol parameter'. Similarly, there is one location (`con2`) where MFunction
	// parameters and symbols overlap (in the paper Turing uses capital German letters
	// for MFunction variables, but I chose to keep english.) In this case I also use
	// prefix the variable with an underscore.

	// `con(C, a)`. Starting from an F-square, S say, the sequence C of symbols describing
	// a configuration closest on the right of S is marked out with letters a. -> `C`
	configuration = []MConfiguration{
		{"con(C, a)", []string{"!A", " "}, []string{"R", "R"}, "con(C, a)"},
		{"con(C, a)", []string{"A"}, []string{"L", "Pa", "R"}, "con1(C, a)"},
		{"con1(C, a)", []string{"A"}, []string{"R", "Pa", "R"}, "con1(C, a)"},
		{"con1(C, a)", []string{"D"}, []string{"R", "Pa", "R"}, "con2(C, a)"},
		// Post suggests this final `con1` line is missing from the original paper
		{"con1(C, a)", []string{" "}, []string{"PD", "R", "Pa", "R", "R", "R"}, "C"},
		{"con2(_C, a)", []string{"C"}, []string{"R", "Pa", "R"}, "con2(_C, a)"},
		{"con2(_C, a)", []string{"!C", " "}, []string{"R", "R"}, "_C"},
	}

	// `b`. The machine prints `:`, `D`, `A` on the F-squares after `::` -> `anf`.
	begin = []MConfiguration{
		{"b", []string{"*", " "}, []string{}, "f(b1, b1, ::)"},
		{"b1", []string{"*", " "}, []string{"R", "R", "P:", "R", "R", "PD", "R", "R", "PA"}, "anf"},
	}

	// `anf`. The machine marks the configuration in the last complete configuration with `y`. -> `kom`
	anfang = []MConfiguration{
		{"anf", []string{"*", " "}, []string{}, "g(anf1, :)"},
		{"anf1", []string{"*", " "}, []string{}, "con(kom, y)"},
	}

	// `kom`. The machine finds the last semi-colon not marked with `z`. It marks this semi-colon
	// with `z` and the configuration following it with `x`.
	kom = []MConfiguration{
		{"kom", []string{";"}, []string{"R", "Pz", "L"}, "con(kmp, x)"},
		{"kom", []string{"z"}, []string{"L", "L"}, "kom"},
		{"kom", []string{"!z", "!;", " "}, []string{"L"}, "kom"},
	}

	// `kmp`. The machine compares the sequences marked `x` and `y`. It erases all letters
	// `x` and `y`. -> `sim` if they are alike. Otherwise -> `kom`.
	kmp = []MConfiguration{
		{"kmp", []string{"*", " "}, []string{}, "cpe(e(e(anf, x), y), sim, x, y)"},
	}

	// `sim`. The machine marks out the instructions. That part of the instructions
	// which refers to operations to be carried out is marked with `u`, and the final
	// MConfiguration with `y`. The letters `z` are erased.
	similar = []MConfiguration{
		{"sim", []string{"*", " "}, []string{}, "fl(sim1, sim1, z)"},
		{"sim1", []string{"*", " "}, []string{}, "con(sim2, )"},
		{"sim2", []string{"A"}, []string{}, "sim3"},
		// Post suggests this final `sim2` line should move left before printing `u`
		{"sim2", []string{"!A", " "}, []string{"L", "Pu", "R", "R", "R"}, "sim2"},
		{"sim3", []string{"!A", " "}, []string{"L", "Py"}, "e(mk, z)"},
		{"sim3", []string{"A"}, []string{"L", "Py", "R", "R", "R"}, "sim3"},
	}

	// `mk`. The last complete configuration is marked out into four sections. The
	// configuration is left unmarked. The symbol directly preceding it is marked
	// with `x`. The remainder of the complete configuration is divided into two
	// parts, of which the first is marked with `v` and the last with `w`. A colon
	// is printed after the whole. -> `sh`.
	mark = []MConfiguration{
		{"mk", []string{"*", " "}, []string{}, "g(mk1, :)"},
		{"mk1", []string{"!A", " "}, []string{"R", "R"}, "mk1"},
		{"mk1", []string{"A"}, []string{"L", "L", "L", "L"}, "mk2"},
		{"mk2", []string{"C"}, []string{"R", "Px", "L", "L", "L"}, "mk2"},
		{"mk2", []string{":"}, []string{}, "mk4"},
		{"mk2", []string{"D"}, []string{"R", "Px", "L", "L", "L"}, "mk3"},
		{"mk3", []string{"!:", " "}, []string{"R", "Pv", "L", "L", "L"}, "mk3"},
		{"mk3", []string{":"}, []string{}, "mk4"},
		{"mk4", []string{"*", " "}, []string{}, "con(l(l(mk5)), )"},
		{"mk5", []string{"*"}, []string{"R", "Pw", "R"}, "mk5"},
		{"mk5", []string{" "}, []string{"P:"}, "sh"},
	}

	// `sh`. The instructions (marked `u`) are examined. If it is found that they involve
	// "Print 1", then `0`, `:` or `1`, `:` is printed at the end.
	// Note: See `enhancedShow` for printing symbols beyong binary digits.
	show = []MConfiguration{
		{"sh", []string{"*", " "}, []string{}, "f(sh1, inst, u)"},
		{"sh1", []string{"*", " "}, []string{"L", "L", "L"}, "sh2"},
		{"sh2", []string{"D"}, []string{"R", "R", "R", "R"}, "sh3"},
		{"sh2", []string{"!D", " "}, []string{}, "inst"},
		{"sh3", []string{"C"}, []string{"R", "R"}, "sh4"},
		{"sh3", []string{"!C", " "}, []string{}, "inst"},
		{"sh4", []string{"C"}, []string{"R", "R"}, "sh5"},
		{"sh4", []string{"!C", " "}, []string{}, "pe2(inst, 0, :)"},
		{"sh5", []string{"C"}, []string{}, "inst"},
		{"sh5", []string{"!C", " "}, []string{}, "pe2(inst, 1, :)"},
	}

	// `inst`. The next complete configuration is written down, carrying out the
	// marked instructions. The letters `u`, `v`, `w`, `x`, `y` are erased. -> `anf`.
	// Note that we are using Petzold's recommended inst1 which does not require new machine behavior.
	instruction = []MConfiguration{
		{"inst", []string{"*", " "}, []string{}, "g(l(inst1), u)"},
		{"inst1", []string{"L"}, []string{"R", "E"}, "ce5(ov, v, y, x, u, w)"},
		{"inst1", []string{"R"}, []string{"R", "E"}, "ce5(ov, v, x, u, y, w)"},
		{"inst1", []string{"N"}, []string{"R", "E"}, "ce5(ov, v, x, y, u, w)"},
		{"ov", []string{"*", " "}, []string{}, "e(anf)"},
	}
)

type (
	// Input for the UniversalMachine
	UniversalMachineInput struct {
		StandardDescription
		SymbolMap
	}
)

// If `M` is a Machine that computes a sequence, this function takes the Standard Description of `M` and returns
// MachineInput that will print `M`'s sequence using `U` (the Universal Machine).
func NewUniversalMachine(input UniversalMachineInput) MachineInput {
	// Helper MFunctions
	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, erase...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, findLeft...)
	mConfigurations = append(mConfigurations, findRight...)
	mConfigurations = append(mConfigurations, copy...)
	mConfigurations = append(mConfigurations, copyAndErase...)
	mConfigurations = append(mConfigurations, replace...)
	mConfigurations = append(mConfigurations, copyAndReplace...)
	mConfigurations = append(mConfigurations, compare...)
	mConfigurations = append(mConfigurations, compareAndErase...)
	mConfigurations = append(mConfigurations, findRightMost...)
	mConfigurations = append(mConfigurations, printAtTheEnd2...)
	mConfigurations = append(mConfigurations, printAtTheEnd2...)
	mConfigurations = append(mConfigurations, copyAndErase2...)
	mConfigurations = append(mConfigurations, eraseAll...)

	// Universal Machine MFunctions
	mConfigurations = append(mConfigurations, configuration...)
	mConfigurations = append(mConfigurations, begin...)
	mConfigurations = append(mConfigurations, anfang...)
	mConfigurations = append(mConfigurations, kom...)
	mConfigurations = append(mConfigurations, kmp...)
	mConfigurations = append(mConfigurations, similar...)
	mConfigurations = append(mConfigurations, mark...)
	mConfigurations = append(mConfigurations, getEnhancedShow(input.SymbolMap)...)
	mConfigurations = append(mConfigurations, instruction...)

	// Construct tape
	tapeFromStandardDescription := []string{"e", "e"}
	for _, char := range input.StandardDescription {
		tapeFromStandardDescription = append(tapeFromStandardDescription, string(char))
		tapeFromStandardDescription = append(tapeFromStandardDescription, none)
	}
	tapeFromStandardDescription = append(tapeFromStandardDescription, "::")

	// Return MachineInput of the compiled abbreviated table of `U`
	return NewAbbreviatedTable(AbbreviatedTableInput{
		MConfigurations:        mConfigurations,
		Tape:                   tapeFromStandardDescription,
		StartingMConfiguration: "b",
		PossibleSymbols:        possibleSymbolsForUniversalMachine,
	})
}

// Rather than using Turing's original `show` m-function, we create our own version
// that is capable of printing all characters the Machine requires (not just `0` and `1`).
func getEnhancedShow(symbolMap SymbolMap) []MConfiguration {
	enhancedShow := []MConfiguration{}

	// First four `show` MConfigurations are valid
	enhancedShow = append(enhancedShow, show[0:4]...)

	// Pick up where `show` left off...
	for symbolKey, symbolValue := range symbolMap {
		symbolNumber, _ := strconv.Atoi(symbolKey[1:])
		// The blank symbol (S0) is `sh3`, and so on.
		showNumber := symbolNumber + 3
		showName := fmt.Sprintf("sh%d", showNumber)

		// To continue our `enhancedShow` MConfigurations move to the next afterwards,
		// unless it is the last symbol.
		var nextShowName string
		if symbolNumber >= len(symbolMap)-1 {
			nextShowName = "inst"
		} else {
			nextShowName = fmt.Sprintf("sh%d", showNumber+1)
		}

		// Turing's convention is that F-squares do not contain blanks
		// unless they are the end of the tape. This poses a problem for our
		// `enhancedShow` MFunction, which would like to show blanks, etc.
		// In addition, sometimes we might be simulating a Machine that prints
		// `e`, etc. or other characters that the Universal Machine itself uses.
		// To combat this, we prepend "shown" values with `_` (underscore).
		// The `CondensedTapeString` will remove this prepended underscore.
		printSymbol := fmt.Sprintf("pe2(inst, _%s, :)", symbolValue)

		enhancedShow = append(enhancedShow, []MConfiguration{
			{showName, []string{"C"}, []string{"R", "R"}, nextShowName},
			{showName, []string{"!C", " "}, []string{}, printSymbol},
		}...)
	}

	return enhancedShow
}

// Helper function to isolate the computed sequence between the colons
func (m *Machine) TapeStringFromUniversalMachine() string {
	var tapeString strings.Builder

	// We essentially need to find only the symbols between two colons
	var started bool
	var skip bool
	var squareMinusTwo string
	var squareMinusOne string
	for _, square := range m.tape {
		if !started {
			if square == "::" {
				started = true
				skip = true
			}
			continue
		}
		if skip {
			skip = !skip
			continue
		}
		if squareMinusTwo == ":" && square == ":" {
			if squareMinusOne == "_" {
				tapeString.WriteString(none)
			} else {
				tapeString.WriteString(strings.TrimPrefix(squareMinusOne, "_"))
			}
		}
		squareMinusTwo = squareMinusOne
		squareMinusOne = square
		skip = !skip
	}
	return tapeString.String()
}
