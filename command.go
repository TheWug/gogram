package telegram

import (
	"strings"
	"github.com/kballard/go-shellquote"
)

type Command struct {
	Command string
	Target  string
	Line	string
	Argstr  string
	Args  []string
}

func ParseCommand(line string) (Command, error) {
	var c Command
	var err error

	c.Line = line
	if strings.HasPrefix(line, "/") {
		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) == 2 {
			line = tokens[1]
		} else {
			line = ""
		}
		command := tokens[0]
		tokens = strings.SplitN(command, "@", 2)
		if len(tokens) == 2 {
			c.Target = tokens[1]
		} else {
			c.Target = *GetMe().Username
		}
		c.Command = tokens[0]
	}

	c.Argstr = line
	c.Args, err = shellquote.Split(line)
	return c, err
}
