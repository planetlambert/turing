package turing

import (
	"strings"
	"testing"
)

func TestFirstCircularDN(t *testing.T) {
	_, err := NewMachineFromDescriptionNumber(DescriptionNumber("1"))
	if err == nil {
		t.Error("expecting D.N. 1 to be circular")
	}
}

func TestFirstCircleFreeDN(t *testing.T) {
	machineInput, err := NewMachineFromDescriptionNumber(DescriptionNumber("731332531"))
	if err != nil {
		t.Error(err)
	}

	m := NewMachine(machineInput)
	m.MoveN(100)
	checkTape(t, m.TapeString(), "S1S1S1S1S1S1S1")
}

func TestWellDefinedness(t *testing.T) {
	mConfigurations := []MConfiguration{}
	mConfigurations = append(mConfigurations, wellDefined...)
	mConfigurations = append(mConfigurations, satisfactory...)
	mConfigurations = append(mConfigurations, unsatisfactory...)
	mConfigurations = append(mConfigurations, decide...)
	mConfigurations = append(mConfigurations, checkSemiColon...)
	mConfigurations = append(mConfigurations, checkName...)
	mConfigurations = append(mConfigurations, checkSymbol...)
	mConfigurations = append(mConfigurations, checkPrintOp...)
	mConfigurations = append(mConfigurations, checkMoveOp...)
	mConfigurations = append(mConfigurations, checkFinalMConfig...)

	m := NewMachine(NewAbbreviatedTable(AbbreviatedTableInput{
		MConfigurations:        mConfigurations,
		Tape:                   strings.Split("; D A D A D A D", ""),
		StartingMConfiguration: "b",
		PossibleSymbols:        possibleSymbolsForWellDefinedMachine,
	}))

	m.MoveN(100000)
	if m.TapeString()[0] != 'u' {
		t.Errorf("expecting unsatisfactory but got %s", m.TapeString())
	}

	m = NewMachine(NewAbbreviatedTable(AbbreviatedTableInput{
		MConfigurations:        mConfigurations,
		Tape:                   strings.Split("; D A D D C R D A", ""),
		StartingMConfiguration: "b",
		PossibleSymbols:        possibleSymbolsForWellDefinedMachine,
		Debug:                  true,
	}))

	m.MoveN(100000)
	if m.TapeString()[0] != 's' {
		t.Errorf("expecting satisfactory but got %s", m.TapeString())
	}
}

// TODO: Test non-`D` part of `H`

// TODO: Test `M1`, `M2`, etc.
