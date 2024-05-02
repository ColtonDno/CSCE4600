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

	if len(args) == 0 {
		if dirs.Len() == 0 {
			fmt.Println("popd: directory stack empty")
			return nil
		}

		dirs.Remove(dirs.Front())
		PrintDirectory(dirs, new_args...)
		return nil
	} else if len(args) > 1 {
		return nil //.err?
	}

	if args[0] == "-v" {
		new_args = append(new_args, "-v")
		dirs.Remove(dirs.Front())
		PrintDirectory(dirs, new_args...)

	} else if args[0] == "-l" {
		new_args = append(new_args, "-l")
		dirs.Remove(dirs.Front())
		PrintDirectory(dirs, new_args...)

	} else if args[0], found = strings.CutPrefix(args[0], "+"); found {
		entry, err := strconv.Atoi(args[0])

		if err != nil {
			return err
		} else if entry > dirs.Len() {
			return nil //.err?
		}

		dir := dirs.Front()
		for i := 1; i < entry; i++ {
			dir = dir.Next()
		}

		dirs.Remove(dir)
		PrintDirectory(dirs, new_args...)
	}

	return nil
}
