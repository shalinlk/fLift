package utils

import "strings"

const CommandLength = 60
const Filler = ":"

func FillUpForCommad(command string) string {
	if len(command) < CommandLength {
		command = command + strings.Repeat(Filler, CommandLength-len(command))
	}
	return command
}
