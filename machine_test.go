package turing

import (
	"strings"
	"testing"
)

func TestMachineExample1(t *testing.T) {
	m := &Machine{
		MConfigurations: []MConfiguration{
			{"b", []string{" "}, []string{"P0", "R"}, "c"},
			{"c", []string{" "}, []string{"R"}, "e"},
			{"e", []string{" "}, []string{"P1", "R"}, "k"},
			{"k", []string{" "}, []string{"R"}, "b"},
		},
	}
	m.MoveN(50)
	checkTape(t, m.TapeString(), "0 1 0 1 0 1 0 1 0 1 0 1")
}

func TestMachineExample1Short(t *testing.T) {
	m := &Machine{
		MConfigurations: []MConfiguration{
			{"b", []string{" "}, []string{"P0"}, "b"},
			{"b", []string{"0"}, []string{"R", "R", "P1"}, "b"},
			{"b", []string{"1"}, []string{"R", "R", "P0"}, "b"},
		},
	}
	m.MoveN(50)
	checkTape(t, m.TapeString(), "0 1 0 1 0 1 0 1 0 1 0 1")
}

func TestMachineExample2(t *testing.T) {
	m := &Machine{
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
	}

	m.Move()
	checkCompleteConfiguration(t, m.CompleteConfiguration(), "eeo0 0")
	m.Move()
	checkCompleteConfiguration(t, m.CompleteConfiguration(), "eeq0 0")
	m.Move()
	checkCompleteConfiguration(t, m.CompleteConfiguration(), "ee0 q0")
	m.Move()
	checkCompleteConfiguration(t, m.CompleteConfiguration(), "ee0 0 q")
	m.Move()
	checkCompleteConfiguration(t, m.CompleteConfiguration(), "ee0 0p 1")
	// ...

	m.MoveN(200)
	checkTape(t, m.TapeString(), "ee0 0 1 0 1 1 0 1 1 1 0 1 1 1 1")
}

func checkTape(t *testing.T, tape string, expectedStart string) {
	if !strings.HasPrefix(tape, expectedStart) {
		var actual string
		if len(expectedStart)+10 <= len(tape) {
			actual = tape[0 : len(expectedStart)+10]
		} else {
			actual = tape
		}
		t.Errorf("got %s, want %s", actual, expectedStart)
	}
}

func checkCompleteConfiguration(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("got %s, want %s", actual, expected)
	}
}
