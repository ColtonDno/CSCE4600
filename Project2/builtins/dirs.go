package builtins

import (
	"container/list"
	"fmt"
	"strings"
)

func PrintDirectory(dirs *list.List, args ...string) error {

	if dirs.Len() == 0 {
		return nil //.err?
	}

	if len(args) == 0 {
		for dir := dirs.Front(); dir != nil; dir = dir.Next() {
			fmt.Printf("~%s ", dir.Value)
		}
		fmt.Println()
		return nil
	} else if len(args) > 1 {
		return nil //.err?
	}

	if args[0] == "-v" {
		for dir := dirs.Front(); dir != nil; dir = dir.Next() {
			fmt.Printf("~%s\n", dir.Value)
		}
	} else if args[0] == "-c" {
		dirs.Init()
	} else if args[0] == "-l" {
		home_dir := HomeDir
		x := len(HomeDir)
		for i := 0; i < x; i++ {
			if home_dir[i] == '\\' {
				home_dir = replaceAt(home_dir, i, '/')
			}
		}
		home_dir, _ = strings.CutPrefix(home_dir, "C:/")

		for dir := dirs.Front(); dir != nil; dir = dir.Next() {
			fmt.Printf("%s%s ", home_dir, dir.Value)
		}
		fmt.Println()
	}

	return nil
}
