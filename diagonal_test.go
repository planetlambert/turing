package turing

import (
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

// TODO: Test `W` (test it on the above). If it doesn't take too long, show `313,325,317` is the first well-defined.
// TODO: Test non-`D` part of `H`

// TODO: Test `M1`, `M2`, etc.
