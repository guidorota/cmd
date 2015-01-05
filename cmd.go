package cmd

import (
	"errors"
	"fmt"
	"os"
)

type FlagSet interface {
	Parse(args []string) error
	Arg(i int) string
	Args() []string
}

var (
	ErrNoCommand     = errors.New("No command specified")
	ErrAlreadyParsed = errors.New("Parse can be invoked at most once")
	ErrNoRunFunction = errors.New("No run function specified")
)

type CommandSet struct {
	cmds     map[string]*Command
	args     []string
	parsed   bool
	selected *Command
}

func (c *CommandSet) AddCommand(name string, short, long string) (*Command, error) {
	if c.cmds == nil {
		c.cmds = make(map[string]*Command)
	}

	if name == "" {
		return nil, fmt.Errorf("Empty command name")
	}
	if _, ok := c.cmds[name]; ok {
		return nil, fmt.Errorf("Duplicate command '%v'", name)
	}

	cmd := NewCommand(name, short, long)
	c.cmds[name] = cmd
	return cmd, nil
}

func (c *CommandSet) Run(args []string) error {
	cmd := c.Parse(args)
	if cmd == nil {
		return ErrNoCommand
	}

	if cmd.Run == nil {
		return ErrNoRunFunction
	}

	cmd.Run(cmd, args)
	return nil
}

func (c *CommandSet) Parse(args []string) *Command {
	if c.parsed {
		return c.selected
	}
	c.parsed = true

	if len(args) == 0 {
		return nil
	}

	cmd := args[0]
	for n, sub := range c.cmds {
		if n == cmd {
			c.selected = sub
			c.args = args[1:]
			return sub
		}
	}

	return nil
}

func (c *CommandSet) Parsed() bool {
	return c.parsed
}

func (c *CommandSet) Selected() *Command {
	return c.selected
}

type Command struct {
	name   string
	short  string
	long   string
	sub    *CommandSet
	parsed bool
	Run    func(cmd *Command, args []string)
	Flags  FlagSet
}

func NewCommand(name string, short, long string) *Command {
	return &Command{
		name:  name,
		short: short,
		long:  long,
	}
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) SubCommand(name string, short, long string) (*Command, error) {
	return c.sub.AddCommand(name, short, long)
}

func (c *Command) Parse(args []string) error {
	if c.parsed {
		return ErrAlreadyParsed
	}
	c.parsed = true

	if c.Flags != nil {
		if err := c.Flags.Parse(args); err != nil {
			return err
		}
	}

	c.sub.Parse(c.Flags.Args())
	return nil
}

func (c *Command) Parsed() bool {
	return c.parsed
}

var CommandLine = &Command{}

func SubCommand(name string, short, long string) (*Command, error) {
	return CommandLine.SubCommand(name, short, long)
}

func Parse() error {
	return CommandLine.Parse(os.Args[1:])
}
