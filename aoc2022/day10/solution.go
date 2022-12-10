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
	Output strings.Builder
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
	cpu.Draw()
	cpu.Cycle++
	if (cpu.Cycle-20)%40 == 0 && cpu.Cycle <= 220 {
		//fmt.Printf("cycle %d, register X: %d\n", cpu.Cycle, cpu.X)
		cpu.Result += cpu.Cycle * cpu.X
	}
}

func (cpu *CPU) Draw() {
	symbol := map[bool]rune{
		true:  '#',
		false: '.',
	}
	//fmt.Printf(
	//	"cycle %d, regX value %d, output %c\n",
	//	cpu.Cycle+1,
	//	cpu.X,
	//	symbol[abs(cpu.X - cpu.Cycle%40) <= 1],
	//)
	cpu.Output.WriteRune(
		symbol[abs(cpu.X-cpu.Cycle%40) <= 1],
	)
	if (cpu.Cycle+1)%40 == 0 {
		cpu.Output.WriteString("\n")
	}

	if (cpu.Cycle+1)%(40*6) == 0 {
		cpu.Output.WriteString("\n")
	}
}

func abs(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
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

func Execute(script string) *CPU {
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
	return &cpu
}

func part1(filename string) string {
	cpu := Execute(filename)
	return strconv.Itoa(cpu.Result)
}

func part2(filename string) string {
	cpu := Execute(filename)
	return cpu.Output.String()
}
