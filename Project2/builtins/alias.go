package builtins

import (
	"fmt"
	"regexp"
	"strings"
)

func CheckForAlias(aliases map[string]string, name string) (string, bool) {
	// Return the alias command if it exists
	if aliases[name] != "" {
		return aliases[name], true
	}

	return name, false
}

func SetAlias(aliases map[string]string, args ...string) error {
	var (
		found_pre bool
		found_suf bool
		print     bool = false
	)

	if len(args) > 1 {
		if args[0] == "-p" {
			print = true
		}

		args[0] = strings.Join(args, " ")
		args[0], _ = strings.CutSuffix(args[0], "-p")
	}

	// Print aliases
	if print || len(args) == 0 {
		for key, value := range aliases {
			fmt.Printf("%s='%s'\n", key, value)
		}

		return nil
	}

	r := regexp.MustCompile(".+=\".+\"")
	match := r.MatchString(args[0])

	if !match {
		fmt.Println("Failed to match regex")
		return nil //. err
	}

	strs := strings.Split(args[0], "=")

	if strings.ContainsAny(strs[0], " /$'=\"") {
		fmt.Println("Alias contains an invalid character")
		return nil //. err
	}

	strs[1], found_pre = strings.CutPrefix(strs[1], "\"")
	strs[1], found_suf = strings.CutSuffix(strs[1], "\"")

	if !found_pre || !found_suf {
		fmt.Println("Did not start with a \"")
		return nil //. error
	}

	aliases[strs[0]] = strs[1]

	return nil
}

// alias first="cd .."
