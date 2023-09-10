package turing

import (
	"reflect"
	"testing"
)

var (
	// `ph`. Prints the provided character, and halts.
	printAndHalt = MConfiguration{"ph(b)", []string{"*", " "}, []string{"Pb"}, "halt"}
	// `prh`. Prints the provided character to the right, and halts.
	printToRightAndHalt = MConfiguration{"prh(b)", []string{"*", " "}, []string{"R", "Pb"}, "halt"}
)

func checkParseMFunction(t *testing.T, mFunction string, expectedName string, expectedParams []string) {
	actualName, actualParams := parseMFunction(mFunction)
	if actualName != expectedName {
		t.Errorf("got %s, want %s", actualName, expectedName)
	}
	if !reflect.DeepEqual(actualParams, expectedParams) {
		t.Errorf("got %s, want %s", actualParams, expectedParams)
	}
}

func TestParseMFunctionRecursiveFirst(t *testing.T) {
	checkParseMFunction(t, "f(x(y, z), a, b)", "f", []string{"x(y, z)", "a", "b"})
}

func TestParseMFunctionRecursiveMiddle(t *testing.T) {
	checkParseMFunction(t, "f(a, x(y, z), b)", "f", []string{"a", "x(y, z)", "b"})
}

func TestParseMFunctionRecursiveLast(t *testing.T) {
	checkParseMFunction(t, "f(a, b, x(y, z))", "f", []string{"a", "b", "x(y, z)"})
}

func TestParseMFunctionRecursiveTwice(t *testing.T) {
	checkParseMFunction(t, "f(x(y, z), x(y, z), b)", "f", []string{"x(y, z)", "x(y, z)", "b"})
}

func TestParseMFunctionBlank(t *testing.T) {
	checkParseMFunction(t, "f(x, )", "f", []string{"x", " "})
}

func TestParseMFunctionNoFunction(t *testing.T) {
	checkParseMFunction(t, "f", "f", []string{})
}

func TestFindLeftMost(t *testing.T) {
	mConfigurations := []MConfiguration{
		// Invokes the findLeftMost MFunction (`f`), printing `x` or `y` depending on if `0` is found
		{"b", []string{"*", " "}, []string{"R", "R", "R"}, "f(ph(x), ph(y), 0)"},
	}

	mConfigurations = append(mConfigurations, printAndHalt)
	mConfigurations = append(mConfigurations, findLeftMost...)
	possibleSymbols := []string{"e", "x", "y", "0", "1"}

	t.Run("FindFirstZero", func(t *testing.T) {
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "1", " ", "1", " ", "0", " ", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee1 1 x 0")
	})

	t.Run("NoZero", func(t *testing.T) {
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "1", " ", "1"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee1 1  y")
	})
}

func TestErase(t *testing.T) {
	// Invokes the erase MFunction (`e`), printing `x` or `y` depending on if `z` is found and erased
	eraseOnceTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "e(ph(x), ph(y), z)"}
	// Invokes the erase MFunction (`e`), printing `x` or `y` depending on if all `z` symbols are found and erased
	eraseAllTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "e(ph(x), z)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, printAndHalt)
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, erase...)
	possibleSymbols := []string{"e", "0", "z", "x", "y"}

	t.Run("EraseX", func(t *testing.T) {
		mConfigurations := append(mConfigurations, eraseOnceTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "z", "0", "z"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee0x0z")
	})

	t.Run("EraseXDoesNotExist", func(t *testing.T) {
		mConfigurations := append(mConfigurations, eraseOnceTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee  y")
	})

	t.Run("EraseAll", func(t *testing.T) {
		mConfigurations := append(mConfigurations, eraseAllTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "", "z", " ", "z"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(30)
		checkTape(t, m.TapeString(), "ee  x")
	})
}

func TestPrintAtTheEnd(t *testing.T) {
	// Invokes the printAtTheEnd MFunction (`pe`), printing `x` at the end of the sequence
	printAtTheEndTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "pe(halt, x)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, printAtTheEndTest)
	possibleSymbols := []string{"e", "0", "x"}

	t.Run("PrintAtTheEnd", func(t *testing.T) {
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee0 0 x")
	})
}

func TestFindLeft(t *testing.T) {
	// Invokes the left MFunction (`l`), printing an `x` to the left of the current tape location
	leftTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R", "R"}, "l(ph(0))"}
	// Invokes the findLeft MFunction (`fl`) printing an `x` to the left of the first ocurrence of the symb ol
	findLeftTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "fl(ph(x), ph(y), 0)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, printAndHalt)
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findLeft...)
	possibleSymbols := []string{"e", "0", "1", "x", "y"}

	t.Run("Left", func(t *testing.T) {
		mConfigurations := append(mConfigurations, leftTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee0")
	})

	t.Run("FindLeft", func(t *testing.T) {
		mConfigurations := append(mConfigurations, findLeftTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "1", " ", "1", " ", "0", " ", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee1 1x0 0")
	})
}

func TestFindRight(t *testing.T) {
	// Invokes the right MFunction (`r`), printing an `x` to the right of the current tape location
	rightTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "r(ph(x))"}
	// Invokes the findRight MFunction (`fr`) printing an `x` to the right of the first ocurrence of the symb ol
	findRightTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "fr(ph(x), ph(y), 0)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, printAndHalt)
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findRight...)
	possibleSymbols := []string{"e", "0", "1", "x", "y"}

	t.Run("Right", func(t *testing.T) {
		mConfigurations := append(mConfigurations, rightTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee0x")
	})

	t.Run("FindRight", func(t *testing.T) {
		mConfigurations := append(mConfigurations, findRightTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "1", " ", "1", " ", "0", " ", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee1 1 0x0")
	})
}

func TestCopy(t *testing.T) {
	// Invokes the copy MFunction (`c`), copying the `0` to the left of `x` to the end of the sequence.
	copyTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "c(halt, halt, x)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findLeft...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, copy...)
	possibleSymbols := []string{"e", "0", "x"}

	t.Run("Copy", func(t *testing.T) {
		mConfigurations := append(mConfigurations, copyTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "0", "x"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(30)
		checkTape(t, m.TapeString(), "ee0 0x0")
	})
}

func TestCopyAndErase(t *testing.T) {
	// Invokes the copyAndErase MFunction (`ce`), copying the marked figure and erasing the marker.
	copyAndEraseOnceTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "ce(halt, halt, x)"}
	// Invokes the copyAndErase MFunction (`ce`), copying all of the marked figures and erasing all of the markers.
	copyAndEraseAllTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "ce(halt, x)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findLeft...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, erase...)
	mConfigurations = append(mConfigurations, copy...)
	mConfigurations = append(mConfigurations, copyAndErase...)
	possibleSymbols := []string{"e", "0", "1", "x"}

	t.Run("CopyAndEraseOnce", func(t *testing.T) {
		mConfigurations := append(mConfigurations, copyAndEraseOnceTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "0", "x"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(50)
		checkTape(t, m.TapeString(), "ee0 0 0")
	})

	t.Run("CopyAndEraseAll", func(t *testing.T) {
		mConfigurations := append(mConfigurations, copyAndEraseAllTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "1", "x", "0", "x"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0 1 0 1 0")
	})
}

func TestReplace(t *testing.T) {
	// Invokes the replace MFunction (`re`), replacing a marker with another.
	replaceOnceTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "re(halt, halt, x, y)"}
	// Invokes the replace MFunction (`re`), replace all markers with another.
	replaceAllTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "re(halt, x, y)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, replace...)
	possibleSymbols := []string{"e", "0", "x", "y"}

	t.Run("ReplaceOnce", func(t *testing.T) {
		mConfigurations := append(mConfigurations, replaceOnceTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "0", "x"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0y0x")
	})

	t.Run("ReplaceAll", func(t *testing.T) {
		mConfigurations := append(mConfigurations, replaceAllTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "0", "x"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0y0y")
	})
}

func TestCopyAndReplace(t *testing.T) {
	// Invokes the copyAndReplace MFunction (`cr`), copying the marked figure and replacing the mark.
	copyAndReplaceOnceTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "cr(halt, halt, x, y)"}
	// Invokes the copyAndReplace MFunction (`cr`), copying the marked figures and replacing the marks.
	copyAndReplaceAllTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "cr(halt, x, y)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findLeft...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, erase...)
	mConfigurations = append(mConfigurations, copy...)
	mConfigurations = append(mConfigurations, copyAndErase...)
	mConfigurations = append(mConfigurations, replace...)
	mConfigurations = append(mConfigurations, copyAndReplace...)
	possibleSymbols := []string{"e", "0", "1", "x", "y"}

	t.Run("CopyAndReplaceOnce", func(t *testing.T) {
		mConfigurations := append(mConfigurations, copyAndReplaceOnceTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "0", "x"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0 0y0")
	})

	t.Run("CopyAndReplaceAll", func(t *testing.T) {
		mConfigurations := append(mConfigurations, copyAndReplaceAllTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "1", "x", "0", "x"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0 1y0y1 0")
	})
}

func TestCompare(t *testing.T) {
	// Invokes the compare MFunction (`cp`), printing `z` at the end if neither markers exist.
	compareNeitherExistTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "cp(halt, halt, pe(halt, z), x, y)"}
	// Invokes the compare MFunction (`cp`), printing `z` at the end if the figures are not equal.
	compareNotEqualTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "cp(halt, pe(halt, z), halt, x, y)"}
	// Invokes the compare MFunction (`cp`), printing `z` at the end if the figures are equal.
	compareEqualTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "cp(pe(halt, z), halt, halt, x, y)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findLeft...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, compare...)
	possibleSymbols := []string{"e", "0", "1", "x", "y", "z"}

	t.Run("CompareNeitherExist", func(t *testing.T) {
		mConfigurations := append(mConfigurations, compareNeitherExistTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0 0 z")
	})

	t.Run("CompareNotEqual", func(t *testing.T) {
		mConfigurations := append(mConfigurations, compareNotEqualTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "1", "y"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0x1yz")
	})

	t.Run("CompareEqual", func(t *testing.T) {
		mConfigurations := append(mConfigurations, compareEqualTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "0", "y"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0x0yz")
	})
}

func TestCompareAndErase(t *testing.T) {
	// Invokes the compareAndErase MFunction (`cpe`), erasing the markers if the figures are equal.
	compareAndEraseOnceTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "cpe(halt, halt, halt, x, y)"}
	// Invokes the compareAndErase MFunction (`cpe`), erasing the markers if the sequence of figures are equal.
	compareAndEraseAllTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "cpe(halt, halt, x, y)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findLeft...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, compare...)
	mConfigurations = append(mConfigurations, erase...)
	mConfigurations = append(mConfigurations, compareAndErase...)
	possibleSymbols := []string{"e", "0", "1", "x", "y"}

	t.Run("CompareAndEraseOnce", func(t *testing.T) {
		mConfigurations := append(mConfigurations, compareAndEraseOnceTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "0", "y"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(200)
		checkTape(t, m.TapeString(), "ee0 0 ")
	})

	t.Run("CompareAndEraseAll", func(t *testing.T) {
		mConfigurations := append(mConfigurations, compareAndEraseAllTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "1", "x", "0", "y", "1", "y"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(200)
		checkTape(t, m.TapeString(), "ee0 1 0 1 ")
	})
}

func TestFindRightMost(t *testing.T) {
	// Invokes the findRightMost MFunction (`g`) (which finds the end of the tape), and prints.
	findEndOfTapeTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "g(ph(x))"}
	// Invokes the findRightMost MFunction (`g`), and prints to the right of the found character.
	findRightMostTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "g(prh(x), 0)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, printAndHalt)
	mConfigurations = append(mConfigurations, printToRightAndHalt)
	mConfigurations = append(mConfigurations, findRightMost...)
	possibleSymbols := []string{"e", "x", "0", "1"}

	t.Run("FindEndOfTape", func(t *testing.T) {
		mConfigurations := append(mConfigurations, findEndOfTapeTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "1", " ", "0", " ", "1"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee0 1 0 1 x")
	})

	t.Run("FindRightMost", func(t *testing.T) {
		mConfigurations := append(mConfigurations, findRightMostTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "1", " ", "0", " ", "1"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(20)
		checkTape(t, m.TapeString(), "ee0 1 0x1")
	})
}

func TestPrintAtTheEnd2(t *testing.T) {
	// Invokes the printAtTheEnd2 MFunction (`pe2`), printing `x` and `y` at the end of the sequence
	printAtTheEndTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "pe2(halt, x, y)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, printAtTheEnd2...)
	mConfigurations = append(mConfigurations, printAtTheEndTest)
	possibleSymbols := []string{"e", "0", "x", "y"}

	t.Run("PrintAtTheEnd2", func(t *testing.T) {
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", " ", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(30)
		checkTape(t, m.TapeString(), "ee0 0 x y")
	})
}

func TestCopyAndErase2(t *testing.T) {
	// Invokes the copyAndErase5 MFunction (`ce5`), copying all of the marked figures and erasing all of the markers.
	copyAndEraseAll2Test := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "ce5(halt, x, s, t, u, v)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, findLeftMost...)
	mConfigurations = append(mConfigurations, findLeft...)
	mConfigurations = append(mConfigurations, printAtTheEnd...)
	mConfigurations = append(mConfigurations, erase...)
	mConfigurations = append(mConfigurations, copy...)
	mConfigurations = append(mConfigurations, copyAndErase...)
	mConfigurations = append(mConfigurations, copyAndErase2...)
	possibleSymbols := []string{"e", "0", "1", "x", "s", "t", "u", "v"}

	t.Run("CopyAndEraseAll2", func(t *testing.T) {
		mConfigurations := append(mConfigurations, copyAndEraseAll2Test)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "0", "x", "0", "s", "1", "s", "1", "t", "0", "t", "1", "u", "1", "u", "0", "v", "0", "v", "0"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(1500)
		checkTape(t, m.TapeString(), "ee0 0 0 1 1 0 1 1 0 0")
	})
}

func TestEraseAll(t *testing.T) {
	// Invokes the eraseAll MFunction (`e`), erasing all markers.
	eraseAllTest := MConfiguration{"b", []string{"*", " "}, []string{"R", "R"}, "e(halt)"}

	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, eraseAll...)
	possibleSymbols := []string{"e", "0", "x", "y"}

	t.Run("EraseAll", func(t *testing.T) {
		mConfigurations := append(mConfigurations, eraseAllTest)
		at := &AbbreviatedTable{
			Machine: Machine{
				MConfigurations:        mConfigurations,
				Tape:                   []string{"e", "e", "0", "x", "0", " ", "0", "y"},
				PossibleSymbols:        possibleSymbols,
				StartingMConfiguration: "b",
			},
		}
		m := at.ToMachine()
		m.MoveN(100)
		checkTape(t, m.TapeString(), "ee0 0 0")
	})
}
