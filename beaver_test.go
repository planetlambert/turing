package turing

import "testing"

func TestFirstBusyBeaver(t *testing.T) {
	testBusyBeaver(t, 1, 1, false)
}

func TestSecondBusyBeaver(t *testing.T) {
	testBusyBeaver(t, 2, 4, false)
}

// BB-3 and beyond take quite a bit of time...

// func TestThirdBusyBeaver(t *testing.T) {
// 	testBusyBeaver(t, 3, 6, false)
// }

// func TestFourthBusyBeaver(t *testing.T) {
// 	testBusyBeaver(t, 4, 13, false)
// }

func testBusyBeaver(t *testing.T, n int, expected int, debug bool) {
	actual, _ := busyBeaver(n, debug)
	if actual != expected {
		t.Errorf("Incorrect BB-%d number %d, expected %d", n, actual, expected)
	}
}
