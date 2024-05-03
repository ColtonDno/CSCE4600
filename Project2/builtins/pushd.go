package builtins

import (
	"container/list"
	"fmt"
	"os"
	"strconv"
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
		dir      string
	)

	for i := 0; i < len(args); i++ {
		if args[i] == "-v" {
			new_args = append(new_args, "-v")

		} else if args[i] == "-l" {
			new_args = append(new_args, "-l")

		} else if args[i], found = strings.CutPrefix(args[i], "+"); found {
			index, _ := strconv.Atoi(args[i])
			if index > dirs.Len() {
				return nil //. err
			}

		} else if args[i][i] == '/' {
			new_dir, _ := strings.CutPrefix(args[i], "/")
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
		} else {
			fmt.Println("Invalid arguements")
			return nil
		}
	}

	if len(args) == 0 && dirs.Len() > 1 {
		dirs.MoveToFront(dirs.Front().Next())

		dir_err := ChangeDirectory(new_args...)
		if dir_err != nil {
			return dir_err
		}

		dir = dirs.Front().Value.(string)
		dir, _ := strings.CutPrefix(dir, "/")

		if dir == "" {
			fmt.Println("Failed to parse dir")
			return nil
		}
		dir_err = ChangeDirectory(dir)
		if dir_err != nil {
			return dir_err
		}
	}

	PrintDirectory(dirs, new_args...)

	return nil
}
