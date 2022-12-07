package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

const CommandPrompt = "$ "

func (c *Command) Parse(line string) error {
	words := strings.Fields(line)
	if len(words) < 2 {
		return fmt.Errorf("command line too short: %d", line)
	}
	if words[0] != strings.TrimSpace(CommandPrompt) {
		return fmt.Errorf("command must start with prompt %q: %s", CommandPrompt, line)
	}
	c.Name = words[1]
	if len(words) > 2 {
		c.Args = words[2:]
	} else {
		c.Args = []string{}
	}

	nargs := map[string]int{
		"cd": 1,
		"ls": 0,
	}
	expected, ok := nargs[c.Name]
	if !ok {
		return fmt.Errorf("cannot validate args for unknown command: %s", c.Name)
	}
	if expected != len(c.Args) {
		return fmt.Errorf(
			"%q: invalid number of arguments for %s: expected %d, got %d %q",
			line,
			c.Name,
			expected,
			len(c.Args),
			c.Args,
		)
	}
	return nil
}

type FSItem struct {
	Name     string
	Type     FSItemType
	Children map[string]*FSItem
	Parent   *FSItem
	fileSize int
}

type FSItemType uint8

const (
	Directory FSItemType = iota
	File
)

func (fi *FSItem) IsDir() bool {
	return fi.Type == Directory
}

func (fi *FSItem) Size() (size int) {
	switch fi.Type {
	case File:
		return fi.fileSize
	case Directory:
		for _, child := range fi.Children {
			size += child.Size()
		}
		return size
	default:
		panic(fmt.Sprintf("Size() not implemented for file type %d", fi.Type))
	}
	return size
}

type Shell struct {
	Running    Command
	CurrentDir *FSItem
	root       *FSItem
}

func (s *Shell) Execute(cmd *Command) error {
	switch cmd.Name {
	default:
		return fmt.Errorf("Execute() not implemented for command: %s", cmd.Name)
	case "cd":
		var dest *FSItem
		var arg string
		arg = cmd.Args[0]
		switch arg {
		case "/":
			dest = s.root
		case "..":
			dest = s.CurrentDir.Parent
			if dest == nil {
				return fmt.Errorf("can not go up from %d", s.CurrentDir.Name)
			}
		default:
			var ok bool
			dest, ok = s.CurrentDir.Children[cmd.Args[0]]
			if !ok {
				return fmt.Errorf("%s: destination does not exist: %s", cmd.Name, cmd.Args[0])
			}
		}
		s.CurrentDir = dest
		s.Running = Command{}
	case "ls":
		s.Running = *cmd
	}
	return nil
}

func (s *Shell) ParseLs(line string) (err error) {
	meta, name, ok := strings.Cut(line, " ")
	if !ok {
		return fmt.Errorf("invalid ls output: %s", line)
	}
	item := &FSItem{Name: name}
	switch meta {
	case "dir":
		item.Type = Directory
		item.Parent = s.CurrentDir
	default:
		item.Type = File
		item.fileSize, err = strconv.Atoi(meta)
		if err != nil {
			return fmt.Errorf("unable to parse file size for %s: %w", item.Name, err)
		}
	}
	children := &s.CurrentDir.Children
	if *children == nil {
		*children = make(map[string]*FSItem)
	}
	(*children)[item.Name] = item
	return nil
}

func ParseShellOutput(filename string) (root FSItem) {
	root = FSItem{Name: "/", Type: Directory}
	var shell Shell
	shell = Shell{
		root: &root,
	}
	var command *Command
	command = &Command{}

	var err error
	var lineNo uint
	for line := range ReadLines(filename) {
		lineNo++
		if strings.HasPrefix(line, CommandPrompt) {
			err = command.Parse(line)
			if err != nil {
				log.Fatal(err)
			}
			err = shell.Execute(command)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		if shell.Running.Name == "ls" {
			shell.ParseLs(line)
			continue
		}
		log.Fatalf("shell broke at line %d: %s", lineNo, line)
	}
	return root
}

func (fs *FSItem) SpecialSize1() (sum int) {
	var size int
	for _, child := range fs.Children {
		if child.Type != Directory {
			continue
		}
		sum += child.SpecialSize1()
	}
	size = fs.Size()
	if size <= 100000 && fs.Name != "/" {
		sum += size
	}
	return sum
}

func part1(filename string) (result string) {
	fs := ParseShellOutput(filename)
	return strconv.Itoa(fs.SpecialSize1())
}

func part2(filename string) (result string) {
	return ""
}
