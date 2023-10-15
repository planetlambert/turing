package turing

// The following m-functions and m-configurations test if a Standard Description
// on the Tape is well-defined. It is assumed that the head of the Tape is at the
// start of the S.D.
var (
	possibleSymbolsForWellDefinedMachine = []string{
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

	wellDefined = []MConfiguration{
		{"b", []string{"*", " "}, []string{}, "checkSemicolon"},
	}

	satisfactory = []MConfiguration{
		{"satisfactory", []string{"*", " "}, []string{}, "decide(s)"},
	}

	unsatisfactory = []MConfiguration{
		{"unsatisfactory", []string{"*", " "}, []string{}, "decide(u)"},
	}

	decide = []MConfiguration{
		{"decide(d)", []string{"*"}, []string{"R", "R"}, "decide(d)"},
		{"decide(d)", []string{" "}, []string{"L", "L"}, "decide1(d)"},
		{"decide1(d)", []string{"*"}, []string{"E", "L", "L"}, "decide(d)"},
		{"decide1(d)", []string{" "}, []string{"Pd"}, "halt"},
	}

	checkSemiColon = []MConfiguration{
		{"checkSemicolon", []string{";"}, []string{"R", "R"}, "checkName"},
		{"checkSemicolon", []string{"!;"}, []string{}, "unsatisfactory"},
		{"checkSemicolon", []string{" "}, []string{"L", "L"}, "checkSemicolon1"},
		{"checkSemicolon1", []string{" "}, []string{}, "unsatisfactory"},
		{"checkSemicolon1", []string{"*"}, []string{"R", "R"}, "satisfactory"},
	}

	checkName = []MConfiguration{
		{"checkName", []string{"D"}, []string{"R", "R"}, "checkName1"},
		{"checkName", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkName1", []string{"A"}, []string{"R", "R"}, "checkName1"},
		{"checkName1", []string{"!A", " "}, []string{}, "checkSymbol"},
	}

	checkSymbol = []MConfiguration{
		{"checkSymbol", []string{"D"}, []string{"R", "R"}, "checkSymbol1"},
		{"checkSymbol", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkSymbol1", []string{"C"}, []string{"R", "R"}, "checkSymbol1"},
		{"checkSymbol1", []string{"!C", " "}, []string{}, "checkPrintOp"},
	}

	checkPrintOp = []MConfiguration{
		{"checkPrintOp", []string{"D"}, []string{"R", "R"}, "checkPrintOp1"},
		{"checkPrintOp", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkPrintOp1", []string{"C"}, []string{"R", "R"}, "checkPrintOp1"},
		{"checkPrintOp1", []string{"!C", " "}, []string{}, "checkMoveOp"},
	}

	checkMoveOp = []MConfiguration{
		{"checkMoveOp", []string{"L", "R", "N"}, []string{"R", "R"}, "checkFinalMConfig"},
		{"checkMoveOp", []string{"!L", "!R", "!N", " "}, []string{}, "unsatisfactory"},
	}

	checkFinalMConfig = []MConfiguration{
		{"checkFinalMConfig", []string{"D"}, []string{"R", "R"}, "checkFinalMConfig1"},
		{"checkFinalMConfig", []string{"!D", " "}, []string{}, "unsatisfactory"},
		{"checkFinalMConfig1", []string{"A"}, []string{"R", "R"}, "checkFinalMConfig1"},
		{"checkFinalMConfig1", []string{"!A", " "}, []string{}, "checkSemicolon"},
	}
)

// TODO: Create non-`D` part of `H`

// TODO: Implement `M1`, `M2`, etc.
