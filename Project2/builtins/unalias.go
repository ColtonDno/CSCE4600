package builtins

func UnsetAlias(aliases map[string]string, args ...string) error {

	if len(args) != 1 {
		return nil //.
	}

	if args[0] == "-a" {
		clear(aliases)
	}

	delete(aliases, args[0])

	return nil
}

// alias first="cd .."
