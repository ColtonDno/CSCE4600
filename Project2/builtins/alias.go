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
	r := regexp.MustCompile("^[^ /$'=\"]+=((\"+[^\"]+\"+$)|([^\" \n]+$))")
	match := r.MatchString(string(str))

	strs := strings.SplitN(string(str), "=", 2)

	if len(strs) == 1 {
		if aliases[strs[0]] == "" {
			fmt.Printf("alias: %s: not found\n", strs[0])
		} else {
			fmt.Printf("%s='%s'\n", strs[0], aliases[strs[0]])
		}
		return
	}

	if !match {
		r = regexp.MustCompile("^[^ /$'=\"]+$")
		match = r.MatchString(strs[0])

		if !match {
			fmt.Printf("alias: '%s': invalid alias name\n", strs[0])
		}

		r = regexp.MustCompile("((\"+[^\"]+\"+$)|([^\" \n]+$))")
		match = r.MatchString(strs[1])

		if !match {
			fmt.Printf("alias: '%s': invalid command\n", strs[1])
		}

		return
	}

	temp := strings.Split(strs[1], "\"")
	strs[1] = strings.Join(temp, "")

	aliases[strs[0]] = strs[1]
}

func SetAlias(aliases map[string]string, args ...string) error {
	var (
		print         bool = false
		commands      [][]rune
		command_count int = 0
		found_quote   int = 0
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

	//Seperate commands
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == ' ' && found_quote == 0 {
			command_count++
			commands = append(commands, make([]rune, 0))
			continue
		}

		commands[command_count] = append(commands[command_count], str[i])

		if str[i] == '"' {
			found_quote++
		} else if str[i] == '=' && found_quote == 2 {
			found_quote = 0
		}
	}

	//Add commands
	slices.Reverse(commands)
	for _, command := range commands {
		slices.Reverse(command)
		AddAlias(aliases, command)
	}

	return nil
}
