package builtins

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

func CheckForAlias(aliases map[string]string, name string) (string, bool) {
	// Return the alias commands if it exists
	if aliases[name] != "" {
		return aliases[name], true
	}

	return name, false
}

func AddAlias(aliases map[string]string, str []rune) {
	r := regexp.MustCompile("^[^ /$'=\"]+=.+")
	match := r.MatchString(string(str))

	strs := strings.Split(string(str), "=")

	if !match {
		fmt.Printf("alias: '%s': invalid alias name\n", strs[0])
		return
	}

	aliases[strs[0]] = strs[1]
}

func SetAlias(aliases map[string]string, args ...string) error {
	var (
		print         bool = false
		commands      [][]rune
		command_count int  = 0
		found_equal   bool = false
	)

	if len(args) > 1 {
		if args[0] == "-p" {
			print = true
		}

		args[0] = strings.Join(args, " ")
	}

	// Print aliases
	if print || len(args) == 0 {
		for key, value := range aliases {
			fmt.Printf("%s='%s'\n", key, value)
		}

		return nil
	}

	str := []rune(args[0])
	commands = make([][]rune, 1)
	commands[0] = make([]rune, 0)

	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == ' ' && found_equal {
			command_count++
			commands = append(commands, make([]rune, 0))
			found_equal = false
			continue
		}

		if str[i] != '"' {
			commands[command_count] = append(commands[command_count], str[i])
		}

		if str[i] == '=' {
			found_equal = true
		}
	}

	slices.Reverse(commands)
	for _, command := range commands {
		slices.Reverse(command)
		AddAlias(aliases, command)
		// fmt.Printf("%s\n", string(command))
	}

	return nil
}
