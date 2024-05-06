# Project 2: Shell Builtins
A simple shell implemented in GO

## Description

For this project, I added 6 commands to the shell:

- history [-chr] [n]
- pushd [-lv] [+n]
- popd [-lv] [+n]
- dirs [-lvc]
- alias [-p] [name[=value] …]
- unalias [-a] [name … ]

Both the command history and aliases were implemented using a map and the directory stack was implemented with a list. All three data structures are created at the start of the run loop and passed along as needed.

As soon as an input is received, it gets put into the history map. The size of the map is used as the key since no entries can be removed unless the map is cleared.

When a command is executed, it is first checked to determine if it is an alias. If it is not, the program continues as usual. If the command is an alias, it gets split into the command name and arguments. The new command is again checked if its an alias. This process continues until the command is no longer an alias, and it is then executed.

## Steps

1. Clone down the example input/output and skeleton `main.go`:

  `git clone https://github.com/ColtonDno/CSCE4600`
 
2. Copy the `Project2` files to your own git project.

  1. In your `go.mod`, replace "ColtonDno" in the module line with your GitHub name, e.g.:

    - "module github.com/ColtonDno/CSCE4600" changes to "module github.com/CoolProfOrTA123/CSCE4600"

  2. In the `main.go`, replace "ColtonDno" in the imports with your package path, e.g.:

    - "github.com/ColtonDno/CSCE4600/Project2/builtins" changes to "github.com/CoolProfOrTA123/CSCE4600/Project2/builtins"
      
3. To compile the code:

  `go build <file path>`
    
4. To run the code (Windows):

  `./Project2.exe`