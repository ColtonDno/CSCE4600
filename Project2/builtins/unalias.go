package builtins

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidOpt = errors.New("invalid option")
)

func UnsetAlias(aliases map[string]string, args ...string) error {
	var (
		found         bool
		unalias_usage string = "unalias: usage: unalias [-a] name"
	)

	if len(args) == 0 {
		fmt.Println(unalias_usage)
		return nil
	}

	if args[0], found = strings.CutPrefix(args[0], "-"); found {
		if args[0] == "a" {
			clear(aliases)
			return nil
		} else {
			return fmt.Errorf("unalias: -%s: %w\n%s", args[0], ErrInvalidOpt, unalias_usage)
		}
	}

	for _, arg := range args {
		if _, found = CheckForAlias(aliases, arg); found {
			delete(aliases, arg)
		} else {
			fmt.Printf("unalias: %s: not found\n", arg)
		}

	}

	return nil
}
