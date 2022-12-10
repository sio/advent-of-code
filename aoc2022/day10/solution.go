package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Instruction int8

const (
	Noop Instruction = 1 << iota
	Addx
)

type CPU struct {
	X      int
	Cycle  int
	Result int
}

func (cpu *CPU) Execute(op Instruction, arg ...int) {
	switch op {
	default:
		panic(fmt.Sprintf("unsupported instruction: %b", op))
	case Noop:
		cpu.noop(arg)
	case Addx:
		cpu.addx(arg)
	}
}

func (cpu *CPU) Tick() {
	cpu.Cycle++
	if (cpu.Cycle-20)%40 == 0 && cpu.Cycle <= 220 {
		fmt.Printf("cycle %d, register X: %d\n", cpu.Cycle, cpu.X)
		cpu.Result += cpu.Cycle * cpu.X
	}
}

func (cpu *CPU) noop(arg []int) {
	if len(arg) != 0 {
		panic("noop instruction takes no arguments")
	}
	cpu.Tick()
}

func (cpu *CPU) addx(arg []int) {
	if len(arg) != 1 {
		panic("addx instruction takes no exactly one argument")
	}
	cpu.Tick()
	cpu.Tick()
	cpu.X += arg[0]
}

func Execute(script string) int {
	cpu := CPU{X: 1}
	var command, value string
	var parsed int
	var found bool
	var err error
	var arg []int
	var op Instruction

	opcode := map[string]Instruction{
		"addx": Addx,
		"noop": Noop,
	}
	for line := range ReadLines(script) {
		command, value, found = strings.Cut(line, " ")
		if !found {
			command = line
			arg = []int{}
		} else {
			parsed, err = strconv.Atoi(value)
			if err != nil {
				log.Fatalf("could not parse command argument: %s", line)
			}
			arg = []int{parsed}
		}
		op, found = opcode[command]
		if !found {
			log.Fatalf("invalid CPU instruction: %s", command)
		}
		cpu.Execute(op, arg...)
	}
	return cpu.Result
}

func part1(filename string) string {
	return strconv.Itoa(Execute(filename))
}

func part2(filename string) string {
	return ""
}
