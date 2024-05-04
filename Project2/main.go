package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/ColtonDno/CSCE4600/Project2/builtins"
)

func main() {
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}) {
	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	dirs := list.New()
	aliases := make(map[string]string)
	history := make(map[int]string)

	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, history, dirs, aliases, exit); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string, history map[int]string, dirs *list.List, aliases map[string]string, exit chan<- struct{}) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)
	// fmt.Printf("Adding %s. Size of history: %d\n", input, len(history))

	history[len(history)] = input

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check if the incoming command is an alias
	new_command, isAlias := builtins.CheckForAlias(aliases, name)

	// name was an alias so split the new command
	if isAlias {
		args = strings.Split(new_command, " ")
		name, args = args[0], args[1:]
	}

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	switch name {
	case "alias":
		return builtins.SetAlias(aliases, args...)
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "dirs":
		return builtins.PrintDirectory(dirs, args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "history":
		return builtins.PrintHistory(history, args...)
	case "pushd":
		return builtins.PushDirectory(dirs, args...)
	case "popd":
		return builtins.PopDirectory(dirs, args...)
	case "quit":
		fallthrough
	case "exit":
		exit <- struct{}{}
		return nil
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
