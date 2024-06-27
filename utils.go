package main

import (
	"strings"
)

func stringToKeyCombination(str string) []string {
	return strings.Split(str, "+")
}
