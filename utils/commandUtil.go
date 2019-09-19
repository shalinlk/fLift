package utils

import "strings"

const SizeLength = 60
const NameLength = 100
const PathLength = 60
const Filler = ":"

func FillUpForCommand(command string, length int) string {
	if len(command) < length {
		command = command + strings.Repeat(Filler, length-len(command))
	}
	return command
}