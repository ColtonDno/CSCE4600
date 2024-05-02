package builtins

import (
	"container/list"
	"fmt"
	"os"
	"strings"
)

var (
// ErrInvalidArgCount = errors.New("invalid argument count")
)

func replaceAt(s string, i int, c rune) string {
	r := []rune(s)
	r[i] = c
	return string(r)
}

func PushDirectory(dirs *list.List, args ...string) error {
	var (
		found    bool
		new_args []string
	)

	if len(args) == 0 {
		if dirs.Len() < 2 {
			return nil
		}

		dirs.MoveToFront(dirs.Front().Next())
		PrintDirectory(dirs, new_args...)

		return nil
	} else if len(args) > 1 {
		return nil //.err?
	}

	if args[0] == "-v" {
		new_args = append(new_args, "-v")
		dirs.MoveToFront(dirs.Front().Next())
		PrintDirectory(dirs, new_args...)

	} else if args[0] == "-l" {
		new_args = append(new_args, "-l")
		dirs.MoveToFront(dirs.Front().Next())
		PrintDirectory(dirs, new_args...)

	} else if args[0][0] == '/' {
		new_dir, _ := strings.CutPrefix(args[0], "/")
		dir_err := ChangeDirectory(new_dir)

		if dir_err != nil {
			return dir_err
		}

		cur_dir, err := os.Getwd()

		if err != nil {
			return err
		}

		cur_dir, _ = strings.CutPrefix(cur_dir, HomeDir)

		x := len(cur_dir)
		for i := 0; i < x; i++ {
			if cur_dir[i] == '\\' {
				cur_dir = replaceAt(cur_dir, i, '/')
			}
		}

		dirs.PushBack(cur_dir)
		PrintDirectory(dirs, new_args...)
	} else if args[0], found = strings.CutPrefix(args[0], "+"); found {
		fmt.Println(args[0])
	}

	return nil
}
