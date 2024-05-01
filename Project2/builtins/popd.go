package builtins

import (
	"container/list"
	"fmt"
	"strings"
)

func PopDirectory(dirs *list.List, args ...string) error {
	var (
		found bool
	)

	if len(args) == 0 {
		if dirs.Len() < 2 {
			return nil
		}

		dirs.Remove(dirs.Front())

		return nil
	} else if len(args) > 1 {
		return nil //.err?
	}

	if args[0], found = strings.CutPrefix(args[0], "+"); found {
		fmt.Println(args[0])
	}

	return nil
}
