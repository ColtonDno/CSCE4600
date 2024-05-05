package builtins

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidArg = errors.New("invalid number")
)

func PrintDirectory(dirs *list.List, args ...string) error {
	var (
		str        string = "~%s "
		home_dir          = HomeDir
		new_line   bool   = false
		print_home bool   = false
		dirs_usage string = "dirs: usage: dirs [-clv]"
	)

	for i := 0; i < len(args); i++ {

		if args[i] == "-v" {
			if strings.Contains(str, " ") {
				str = replaceAt(str, strings.Index(str, " "), '\n')
			}
			new_line = true

		} else if args[i] == "-c" {
			dirs.Init()
			return nil

		} else if args[i] == "-l" {
			x := len(home_dir)
			for i := 0; i < x; i++ {
				if home_dir[i] == '\\' {
					home_dir = replaceAt(home_dir, i, '/')
				}
			}
			home_dir, _ = strings.CutPrefix(home_dir, "C:/")
			str = strings.Replace(str, "~", "%s", -1)

			print_home = true

		} else {
			return fmt.Errorf("dirs: %s: %w\n%s", args[i], ErrInvalidArg, dirs_usage)
		}
	}

	if dirs.Len() == 0 {
		return nil
	}

	for dir := dirs.Front(); dir != nil; dir = dir.Next() {
		if print_home {
			fmt.Printf(str, home_dir, dir.Value)
		} else {
			fmt.Printf(str, dir.Value)
		}
	}

	if !new_line {
		fmt.Println()
	}

	return nil
}
