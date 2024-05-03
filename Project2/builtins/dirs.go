package builtins

import (
	"container/list"
	"fmt"
	"strings"
)

func PrintDirectory(dirs *list.List, args ...string) error {
	var (
		str        string = "~%s "
		home_dir          = HomeDir
		new_line   bool   = false
		print_home bool   = false
	)

	for i := 0; i < len(args); i++ {

		if args[i] == "-v" {
			if new_line {
				return nil //. Dup argument err?
			}

			str = replaceAt(str, strings.Index(str, " "), '\n')
			new_line = true

		} else if args[i] == "-c" {
			dirs.Init()
			return nil

		} else if args[i] == "-l" {
			if print_home {
				return nil //. Dup argument err?
			}

			x := len(HomeDir)
			for i := 0; i < x; i++ {
				if home_dir[i] == '\\' {
					home_dir = replaceAt(home_dir, i, '/')
				}
			}
			home_dir, _ = strings.CutPrefix(home_dir, "C:/")

			str = strings.Replace(str, "~", "%s", -1)

			print_home = true

		} else {
			fmt.Println("Invalid arguements")
			return nil
		}
	}

	if dirs.Len() == 0 {
		return nil //.err?
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
