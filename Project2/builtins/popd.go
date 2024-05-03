package builtins

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

func PopDirectory(dirs *list.List, args ...string) error {
	var (
		found    bool
		new_args []string
	)

	dir := dirs.Front()

	if dir == nil {
		fmt.Println("popd: directory stack empty")
		return nil
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-v" {
			new_args = append(new_args, "-v")

		} else if args[i] == "-l" {
			new_args = append(new_args, "-l")

		} else if args[i], found = strings.CutPrefix(args[i], "+"); found {
			entry, err := strconv.Atoi(args[i])

			if err != nil {
				return err
			} else if entry > dirs.Len() {
				return nil //.err?
			}

			for i := 0; i < entry; i++ {
				dir = dir.Next()
			}

		}
	}

	dirs.Remove(dir)
	PrintDirectory(dirs, new_args...)

	return nil
}
