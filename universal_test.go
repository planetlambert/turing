package turing

import (
	"testing"
)

func TestUniversalMachineExample1(t *testing.T) {
	input := MachineInput{
		MConfigurations: []MConfiguration{
			{"b", []string{" "}, []string{"P0", "R"}, "c"},
			{"c", []string{" "}, []string{"R"}, "e"},
			{"e", []string{" "}, []string{"P1", "R"}, "k"},
			{"k", []string{" "}, []string{"R"}, "b"},
		},
	}

	m := NewMachine(input)
	st := NewStandardTable(input)

	expected := "0 1 0 1 0 1"

	// Check Machine Tape
	m.MoveN(50)
	checkTape(t, m.TapeString(), expected)

	// Check Universal Machine Tape
	um := NewMachine(NewUniversalMachine(UniversalMachineInput{
		StandardDescription: st.StandardDescription,
		SymbolMap:           st.SymbolMap,
	}))
	um.MoveN(500000)
	checkTape(t, TapeStringFromUniversalMachineTape(um.Tape()), expected)
}
