package builtins

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrDirStackEmpty = errors.New("directory stack empty")
)

func PopDirectory(dirs *list.List, args ...string) error {
	var (
		found    bool
		new_args []string
	)

	dir := dirs.Front()

	if dir == nil {
		fmt.Println("popd: %w", ErrDirStackEmpty)
		return nil
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-v" {
			new_args = append(new_args, "-v")

		} else if args[i] == "-l" {
			new_args = append(new_args, "-l")

		} else if args[i], found = strings.CutPrefix(args[i], "+"); found {
			index, err := strconv.Atoi(args[i])

			if err != nil {
				return err
			} else if index > dirs.Len() {
				return fmt.Errorf("pushd: +%d: %w", index, ErrInvalidIndex)
			}

			for i := 0; i < index; i++ {
				dir = dir.Next()
			}

		}
	}

	dirs.Remove(dir)
	err := PrintDirectory(dirs, new_args...)

	if err != nil {
		return err
	}

	return nil
}
