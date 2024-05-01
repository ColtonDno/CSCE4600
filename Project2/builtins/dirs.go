package builtins

import (
	"container/list"
	"fmt"
)

func PrintDirectory(dirs *list.List, args ...string) error {

	for dir := dirs.Front(); dir != nil; dir = dir.Next() {
		fmt.Println(dir.Value)
	}

	return nil
}
