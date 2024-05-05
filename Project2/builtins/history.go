package builtins

import (
	"fmt"
	"strconv"
)

func PrintHistory(history map[int]string, args ...string) error {
	var (
		print_num  bool = true
		print_from      = 0
		print_to        = len(history)
	)

	for i := 0; i < len(args); i++ {
		if args[i] == "-c" {
			clear(history)
			return nil
		} else if args[i] == "-h" {
			print_num = false
		} else if args[i] == "-r" {
			temp := print_from
			print_from = print_to - 1
			print_to = temp - 1
		} else if _, err := strconv.Atoi(args[i]); err == nil {
			print_from, _ = strconv.Atoi(args[i])
			print_from = len(history) - print_from
		}
	}

	for i := print_from; i != print_to; i++ {
		if !print_num {
			fmt.Printf("\t%s\n", history[i])
		} else {
			fmt.Printf("\t%d  %s\n", i+1, history[i])
		}

		if print_from > print_to {
			i -= 2
		}
	}

	return nil
}
