package turing

import (
	"fmt"
	"strconv"
	"strings"
)

type UniversalMachine struct {
	Machine
}

const (
	Colon       string = ":"
	DoubleColon string = "::"
	Underscore  string = "_"
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

	// From the MConfiguration `f` the machine finds the
	// symbol of form `a` which is farthest to the left (the "first a")
	// and the MConfiguration then becomes `C`. If there is no `a`
	// then the MConfiguration becomes `B`.
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
	// The MConfiguration `cr(B, a)` is taken up when no letters `b` are on the tape.
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
	instruction = []MConfiguration{
		{"inst", []string{"*", " "}, []string{}, "g(l(inst1), u)"},
		{"inst1", []string{"L"}, []string{"R", "E"}, "ce5(ov, v, y, x, u, w)"},
		{"inst1", []string{"R"}, []string{"R", "E"}, "ce5(ov, v, x, u, y, w)"},
		{"inst1", []string{"N"}, []string{"R", "E"}, "ce5(ov, v, x, y, u, w)"},
		{"ov", []string{"*", " "}, []string{}, "e(anf)"},
	}
)

func NewUniversalMachine(sd StandardDescription, symbolMap SymbolMap) *UniversalMachine {
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
	mConfigurations = append(mConfigurations, getEnhancedShow(symbolMap)...)
	mConfigurations = append(mConfigurations, instruction...)

	// Construct tape
	tapeFromStandardDescription := []string{"e", "e"}
	for _, char := range sd {
		tapeFromStandardDescription = append(tapeFromStandardDescription, string(char))
		tapeFromStandardDescription = append(tapeFromStandardDescription, None)
	}
	tapeFromStandardDescription = append(tapeFromStandardDescription, "::")

	at := AbbreviatedTable{
		Machine: Machine{
			MConfigurations:        mConfigurations,
			Tape:                   tapeFromStandardDescription,
			StartingMConfiguration: "b",
			PossibleSymbols:        possibleSymbolsForUniversalMachine,
		},
	}

	return &UniversalMachine{
		Machine: *at.ToMachine(),
	}
}

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

func (um *UniversalMachine) CondensedTapeString() string {
	var tapeString strings.Builder

	// We essentially need to find only the symbols between two colons
	var started bool
	var skip bool
	var squareMinusTwo string
	var squareMinusOne string
	for _, square := range um.Tape {
		if !started {
			if square == DoubleColon {
				started = true
				skip = true
			}
			continue
		}
		if skip {
			skip = !skip
			continue
		}
		if squareMinusTwo == Colon && square == Colon {
			if squareMinusOne == Underscore {
				tapeString.WriteString(None)
			} else {
				tapeString.WriteString(strings.TrimPrefix(squareMinusOne, Underscore))
			}
		}
		squareMinusTwo = squareMinusOne
		squareMinusOne = square
		skip = !skip
	}
	return tapeString.String()
}
