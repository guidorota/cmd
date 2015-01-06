package cmd

import (
	"testing"
)

func TestCommandSetParsed(t *testing.T) {
	cs := &CommandSet{}
	cs.AddCommand("test1", "s", "l")
	cs.AddCommand("test2", "s", "l")

	if cs.Parsed() {
		t.Error("CommandSet has not been parsed yet")
	}

	cs.Parse(nil)
	if !cs.Parsed() {
		t.Error("Parsed flag not set correctly")
	}
}

func TestCommandSetParsing(t *testing.T) {
	cs := &CommandSet{}
	cs.AddCommand("test1", "s", "l")
	cs.AddCommand("test2", "s", "l")
	args := []string{"test1", "remains"}

	c := cs.Parse(args)
	if c == nil {
		t.Error("Command not parsed")
	}
	if c.Name() != "test1" {
		t.Error("Wrong command selected")
	}
	if len(cs.Args()) != 1 {
		t.Error("Wrong remaining args size")
	}
	if cs.Args()[0] != "remains" {
		t.Error("Wrong remaining args element")
	}
}
