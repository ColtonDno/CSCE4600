package builtins

import "fmt"

func UnsetAlias(aliases map[string]string, args ...string) error {

	if len(args) != 1 {
		return fmt.Errorf("invalid argument count: expected zero or one arguments (directory)")
	}

	if args[0] == "-a" {
		clear(aliases)
	}

	delete(aliases, args[0])

	return nil
}
