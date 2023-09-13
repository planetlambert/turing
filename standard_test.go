package turing

import (
	"testing"
)

func TestStandardMachineExample1(t *testing.T) {
	st := NewStandardTable(MachineInput{
		MConfigurations: []MConfiguration{
			{"b", []string{" "}, []string{"P0", "R"}, "c"},
			{"c", []string{" "}, []string{"R"}, "e"},
			{"e", []string{" "}, []string{"P1", "R"}, "k"},
			{"k", []string{" "}, []string{"R"}, "b"},
		},
		PossibleSymbols: []string{"0", "1"},
	})
	m := NewMachine(st.MachineInput)
	m.MoveN(100)
	checkTape(t, st.SymbolMap.TranslateTape(m.Tape()), "0 1 0 1 0 1 0 1 0 1 0 1")
	checkStandardDescription(t, st.StandardDescription, ";DADDCRDAA;DAADDRDAAA;DAAADDCCRDAAAA;DAAAADDRDA")
	checkDescriptionNumber(t, st.DescriptionNumber, "73133253117311335311173111332253111173111133531")
}

func TestStandardMachineExample1Short(t *testing.T) {
	st := NewStandardTable(MachineInput{
		MConfigurations: []MConfiguration{
			{"b", []string{" "}, []string{"P0"}, "b"},
			{"b", []string{"0"}, []string{"R", "R", "P1"}, "b"},
			{"b", []string{"1"}, []string{"R", "R", "P0"}, "b"},
		},
		PossibleSymbols: []string{"0", "1"},
	})
	m := NewMachine(st.MachineInput)
	m.MoveN(100)
	checkTape(t, st.SymbolMap.TranslateTape(m.Tape()), "0 1 0 1 0 1 0 1 0 1 0 1")
	// No StandardDescription or DescriptionNumner given
}

func TestStandardMachineExample2(t *testing.T) {
	st := NewStandardTable(MachineInput{
		MConfigurations: []MConfiguration{
			{"b", []string{"*", " "}, []string{"Pe", "R", "Pe", "R", "P0", "R", "R", "P0", "L", "L"}, "o"},
			{"o", []string{"1"}, []string{"R", "Px", "L", "L", "L"}, "o"},
			{"o", []string{"0"}, []string{}, "q"},
			{"q", []string{"0", "1"}, []string{"R", "R"}, "q"},
			{"q", []string{" "}, []string{"P1", "L"}, "p"},
			{"p", []string{"x"}, []string{"E", "R"}, "q"},
			{"p", []string{"e"}, []string{"R"}, "f"},
			{"p", []string{" "}, []string{"L", "L"}, "p"},
			{"f", []string{"*"}, []string{"R", "R"}, "f"},
			{"f", []string{" "}, []string{"P0", "L", "L"}, "o"},
		},
		PossibleSymbols: []string{"0", "1", "e", "x"},
	})
	m := NewMachine(st.MachineInput)
	m.MoveN(1000)
	checkTape(t, st.SymbolMap.TranslateTape(m.Tape()), "ee0 0 1 0 1 1 0 1 1 1 0 1 1 1 1")
	// No StandardDescription or DescriptionNumner given
}

func checkStandardDescription(t *testing.T, actual StandardDescription, expected string) {
	if string(actual) != expected {
		t.Errorf("got %s, want %s", actual, expected)
	}
}

func checkDescriptionNumber(t *testing.T, actual DescriptionNumber, expected string) {
	if string(actual) != expected {
		t.Errorf("got %s, want %s", actual, expected)
	}
}
