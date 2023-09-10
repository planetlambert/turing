package turing

import (
	"testing"
)

func TestUniversalMachineExample1(t *testing.T) {
	m := &Machine{
		MConfigurations: []MConfiguration{
			{"b", []string{" "}, []string{"P0", "R"}, "c"},
			{"c", []string{" "}, []string{"R"}, "e"},
			{"e", []string{" "}, []string{"P1", "R"}, "k"},
			{"k", []string{" "}, []string{"R"}, "b"},
		},
	}
	st := m.ToStandardTable()
	sd := st.ToStandardDescription()

	expected := "0 1 0 1 0 1"

	// Check Machine Tape
	m.MoveN(50)
	checkTape(t, m.TapeString(), expected)

	// Check Universal Machine Tape
	um := NewUniversalMachine(sd, st.SymbolMap)
	um.MoveN(500000)
	checkTape(t, um.CondensedTapeString(), expected)
}
