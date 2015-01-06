package cmd

import (
	"testing"
)

func Test(t *testing.T) {
	cs := &CommandSet{}

	cs.AddCommand("test1", "s", "l")
	cs.AddCommand("test2", "s", "l")

	if cs.Parsed() {
		t.Error("CommandSet has not been parsed yet")
	}
}
