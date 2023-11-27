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
	m := NewMachine(NewAbbreviatedTable(AbbreviatedTableInput{
		MConfigurations:        wellDefinedMachineMConfigurations,
		Tape:                   strings.Split("; D A D A D A D", ""),
		StartingMConfiguration: "b",
		PossibleSymbols:        wellDefinedMachinePossibleSymbols,
	}))

	m.MoveN(100000)
	if m.TapeString()[0] != 'u' {
		t.Errorf("expecting unsatisfactory but got %s", m.TapeString())
	}

	m = NewMachine(NewAbbreviatedTable(AbbreviatedTableInput{
		MConfigurations:        wellDefinedMachineMConfigurations,
		Tape:                   strings.Split("; D A D D C R D A", ""),
		StartingMConfiguration: "b",
		PossibleSymbols:        wellDefinedMachinePossibleSymbols,
	}))

	m.MoveN(100000)
	if m.TapeString()[0] != 's' {
		t.Errorf("expecting satisfactory but got %s", m.TapeString())
	}
}

// TODO: Test non-`D` part of `H`

// TODO: Test `M1`, `M2`, etc.
