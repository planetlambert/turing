package turing

// The following m-functions and m-configurations test if a Standard Description
// on the Tape is well-defined. It is assumed that the head of the Tape is at the
// start of the S.D.
var (
	wellDefinedMachinePossibleSymbols = []string{
		string(a),
		string(c),
		string(d),
		string(l),
		string(r),
		string(n),
		string(semicolon),
		"s",
		"u",
	}

	wellDefinedMachineMConfigurations = []MConfiguration{
		// The start of the machine
		{"b", []string{"*", " "}, []string{}, "checkSemicolon"},

		// The decision of well-definedness is satisfactory
		{"satisfactory", []string{"*", " "}, []string{}, "decide(s)"},

		// The decision of well-definedness is unsatisfactory
		{"unsatisfactory", []string{"*", " "}, []string{}, "decide(u)"},

		// Erase everything and print the decision
		{"decide(d)", []string{"*"}, []string{"R", "R"}, "decide(d)"},
		{"decide(d)", []string{" "}, []string{"L", "L"}, "decide1(d)"},
		{"decide1(d)", []string{"*"}, []string{"E", "L", "L"}, "decide(d)"},
		{"decide1(d)", []string{" "}, []string{"Pd"}, "halt"},

		// Check the semicolon that deliminates the S.D.
		{"checkSemicolon", []string{";"}, []string{"R", "R"}, "checkName"},
		{"checkSemicolon", []string{"!;"}, []string{}, "unsatisfactory"},
		{"checkSemicolon", []string{" "}, []string{"L", "L"}, "checkSemicolon1"},
		{"checkSemicolon1", []string{" "}, []string{}, "unsatisfactory"},
		{"checkSemicolon1", []string{"*"}, []string{"R", "R"}, "satisfactory"},

		// Check the name portion of the S.D. subsegment
		{"checkName", []string{"D"}, []string{"R", "R"}, "checkName1"},
		{"checkName", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkName1", []string{"A"}, []string{"R", "R"}, "checkName1"},
		{"checkName1", []string{"!A", " "}, []string{}, "checkSymbol"},

		// Check the symbol portion of the S.D. subsegment
		{"checkSymbol", []string{"D"}, []string{"R", "R"}, "checkSymbol1"},
		{"checkSymbol", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkSymbol1", []string{"C"}, []string{"R", "R"}, "checkSymbol1"},
		{"checkSymbol1", []string{"!C", " "}, []string{}, "checkPrintOp"},

		// Check the print operation portion of the S.D. subsegment
		{"checkPrintOp", []string{"D"}, []string{"R", "R"}, "checkPrintOp1"},
		{"checkPrintOp", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkPrintOp1", []string{"C"}, []string{"R", "R"}, "checkPrintOp1"},
		{"checkPrintOp1", []string{"!C", " "}, []string{}, "checkMoveOp"},

		// Check the move operation portion of the S.D. subsegment
		{"checkMoveOp", []string{"L", "R", "N"}, []string{"R", "R"}, "checkFinalMConfig"},
		{"checkMoveOp", []string{"!L", "!R", "!N", " "}, []string{}, "unsatisfactory"},

		// Check the final m-config portion of the S.D. subsegment
		{"checkFinalMConfig", []string{"D"}, []string{"R", "R"}, "checkFinalMConfig1"},
		{"checkFinalMConfig", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkFinalMConfig1", []string{"A"}, []string{"R", "R"}, "checkFinalMConfig1"},
		{"checkFinalMConfig1", []string{"!A", " "}, []string{}, "checkSemicolon"},
	}
)

// The following defines Turing's `H` machine. The entire machine is implemented
// with the exception of the `D` machine (which is not possible).
var (
	hMachinePossibleSymbols = []string{}

	hMachineMConfigurations = []MConfiguration{
		// TODO: The start of the machine
		{"beginH", []string{"*", " "}, []string{"Pe", "R", "Pe", "R", "P:::"}, "TODO"},

		// TODO: Find the 2nd to last `:` and copy what comes after to the end of the tape,
		// but increment this number by 1. Afterwards, print `::` and move to `convert`.
		{"iter", []string{"*", " "}, []string{}, "TODO"},

		// TODO: Find the latest `::` and convert the D.N. to the left into a S.D. on the right.
		// Afterwards, print `:::` invoke `D`.
		{"convert", []string{"*", " "}, []string{}, "TODO"},

		// TODO: Fake `D`
		{"D", []string{"*", " "}, []string{}, "TODO"},

		// TODO: Check for `s` or `u`. If `u`, print `:`, and move back to `iter`.
		// If `s`, print `::::`, and move to `R`.
		{"check", []string{"*", " "}, []string{}, "TODO"},

		// TODO: Find the most recent `R` after the 2nd to last `::::` and copy it to the end
		// of the tape. Add one more symbol to increment. Print `:::::` and move to `simulate`.
		{"R", []string{"*", " "}, []string{}, "TODO"},

		// TODO: Copy the S.D. after the most recent `::` to the end of the tape.
		// Use `U` to simulate the machine, with the modification of only printing
		// `R` characters. After this has happened, print `::::::` move to `print`.
		{"simulate", []string{"*", " "}, []string{}, "TODO"},

		// TODO: Pluck the `R`'th character from the complete configuration after `:::::`
		// and print it after `::::::`. Now print `:` and move back to `iter`.
		{"print", []string{"*", " "}, []string{}, "TODO"},
	}
)

// The following defines Turing's `G` machine. The entire machine is implemented
// with the exception of the `E` machine (which is not possible).
var (
	gMachinePossibleSymbols = []string{}

	gMachineMConfigurations = []MConfiguration{
		// TODO: The start of the machine
		{"b", []string{"*", " "}, []string{"Pe", "R", "Pe", "R", "P:"}, "TODO"},

		// TODO: Enumerate m-functions
	}
)
