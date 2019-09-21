package utils

import (
	"errors"
	"strings"
)

const SizeLength = 60
const NameLength = 100
const PathLength = 60
const Filler = ":"

func FillUpForCommand(command string, length int) (string, error) {
	if len(command) < length {
		command = command + strings.Repeat(Filler, length-len(command))
	}else if len(command) > length{
		return "", errors.New("value exceeds size limits")
	}
	return command, nil
}