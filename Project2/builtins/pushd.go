package builtins

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInvalidIndex = errors.New("directory stack index out of range")
	ErrNotDir       = errors.New("no such file or directory")
	ErrMissingDir   = errors.New("directory no longer exists")
)

func replaceAt(s string, i int, c rune) string {
	r := []rune(s)
	r[i] = c
	return string(r)
}

func PushDirectory(dirs *list.List, args ...string) error {
	var (
		found      bool
		empty_args []string
		new_args   []string
		dir        string
	)

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

			target_dir := dirs.Front()
			for i := 0; i < index; i++ {
				target_dir = target_dir.Next()
			}

			dirs.MoveToFront(target_dir)

			cur_dir, err := os.Getwd()

			if err != nil {
				return err
			}

			dir_err := ChangeDirectory(empty_args...)
			if dir_err != nil {
				return dir_err
			}

			dir = dirs.Front().Value.(string)
			dir, _ := strings.CutPrefix(dir, "/")

			dir_err = ChangeDirectory(dir)

			if dir_err != nil {
				dirs.Remove(dirs.Front())
				dir_err = ChangeDirectory(cur_dir)

				if dir_err != nil {
					return fmt.Errorf("pushd: failed to restore original directory")
				}

				return fmt.Errorf("pushd: %w", ErrMissingDir)
			}

		} else {
			new_dir, _ := strings.CutPrefix(args[i], "/")
			dir_err := ChangeDirectory(new_dir)

			if dir_err != nil {
				return fmt.Errorf("pushd: %s: %w", args[i], ErrNotDir)
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
		}
	}

	if len(args) == 0 && dirs.Len() > 1 {
		dirs.MoveToFront(dirs.Front().Next())

		cur_dir, err := os.Getwd()

		if err != nil {
			return err
		}

		//Change to home directory
		dir_err := ChangeDirectory(empty_args...)
		if dir_err != nil {
			return dir_err
		}

		dir = dirs.Front().Value.(string)
		dir, _ := strings.CutPrefix(dir, "/")

		dir_err = ChangeDirectory(dir)

		if dir_err != nil {
			dirs.Remove(dirs.Front())
			dir_err = ChangeDirectory(cur_dir)

			if dir_err != nil {
				return fmt.Errorf("pushd: failed to restore original directory")
			}

			return fmt.Errorf("pushd: %w", ErrMissingDir)
		}
	}

	err := PrintDirectory(dirs, new_args...)

	if err != nil {
		return err
	}

	return nil
}
