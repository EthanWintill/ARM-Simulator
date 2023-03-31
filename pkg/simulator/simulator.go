package simulator

import (
	"fmt"
	"math"
	"os"

	i "ARM-Simulator/pkg/instructions"
)

type Simulator struct {
	registers [32]int64
	data      []i.Data
}

func InitializeSimulator(data []i.Data) Simulator {
	sim := &Simulator{}
	sim.data = append(sim.data, data...)
	return *sim
}

func GetDataFromAddress(simulator Simulator, address int) int {
	for _, d := range simulator.data {
		if d.GetProgramCount() == address {
			return d.GetLineValue()
		}
	}
	return 0
}

func PutDataIntoAddress(simulator Simulator, address int, value int) {
	for _, d := range simulator.data {
		if d.GetProgramCount() == address {
			d.SetLineValue(int32(value))
			return
		}
	}
	d := i.Data{}
	d.SetProgramCount(address)
	d.SetLineValue(int32(value))
	simulator.data = append(simulator.data, d)
}

func RunInstructionsAndReturnStrings(simulator Simulator, instructions []i.Instruction) []string {
	var results []string
	cycle := 1
	currentInstruction := 0
	instructionCounter := 0

	for {
		if currentInstruction >= len(instructions) || currentInstruction < 0 {
			fmt.Printf("Branched out of the instruction space, attempted to move to %d\n", currentInstruction)
			os.Exit(1)
		}

		offset := ApplyInstructionToSimulator(&simulator, instructions[currentInstruction])
		r := StateToString(cycle, simulator, instructions[currentInstruction])
		results = append(results, r)
		cycle += 1
		if instructions[currentInstruction].GetOp() == "BREAK" {
			break
		}
		currentInstruction += offset
		instructionCounter += 1
		if instructionCounter > 100 {
			fmt.Println("LIKELY LOOP DETECTED OR NO BREAK IN INSTRUCTIONS")
			break
		}
	}

	return results
}

func ApplyInstructionToSimulator(simulator *Simulator, instruction i.Instruction) int {
	result := 1

	switch instruction.GetOp() {
	case "ADD":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] + simulator.registers[instruction.GetRm()]
	case "SUB":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] - simulator.registers[instruction.GetRm()]
	case "ADDI":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] + int64(instruction.GetIm())
	case "SUBI":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] - int64(instruction.GetIm())
	case "EOR":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] ^ simulator.registers[instruction.GetRm()]
	case "ORR":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] | simulator.registers[instruction.GetRm()]
	case "AND":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] & simulator.registers[instruction.GetRm()]
	case "B":
		return int(instruction.GetOffset())
	case "LDUR":
		dataAddress := uint16(instruction.GetRn()) + instruction.GetAddress()*4
		simulator.registers[instruction.GetRd()] = int64(GetDataFromAddress(*simulator, int(dataAddress)))
	case "STUR":
		dataAddress := uint16(instruction.GetRn()) + instruction.GetAddress()*4
		PutDataIntoAddress(*simulator, int(dataAddress), int(simulator.registers[instruction.GetRd()]))
	case "CBZ":
		if simulator.registers[instruction.GetRd()] == 0 {
			return int(instruction.GetOffset())
		}
	case "CBNZ":
		if simulator.registers[instruction.GetRd()] != 0 {
			return int(instruction.GetOffset())
		}
	case "LSL":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] << int(instruction.GetShamt())
	case "ASR":
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRn()] >> int(instruction.GetShamt())
	case "LSR":
		simulator.registers[instruction.GetRd()] = int64(uint64(simulator.registers[instruction.GetRn()]) >> int(instruction.GetShamt()))
	case "MOVZ":
		simulator.registers[instruction.GetRd()] = int64(instruction.GetField()) * int64(math.Pow(2, float64(instruction.GetShiftCode())))
	case "MOVK": // simply subtracts the old target 16 bits from rd then adds the desired field
		oldField := (simulator.registers[instruction.GetRd()] >> int(instruction.GetField())) & 0b1111111111111111
		simulator.registers[instruction.GetRd()] = simulator.registers[instruction.GetRd()] + int64(instruction.GetField())*int64(math.Pow(2, float64(instruction.GetShiftCode()))) - oldField
	}
	return result
}

func StateToString(cycle int, simulator Simulator, instruction i.Instruction) string {
	var result string
	result += "====================\n"
	result += fmt.Sprintf("Cycle:%d\t%d\t%s\n\n",
		cycle,
		instruction.GetProgramCount(),
		instruction.GetArgs(),
	)
	result += fmt.Sprintf("registers:\nr00:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\nr08:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\nr16:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\nr24:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n\n",
		simulator.registers[0], simulator.registers[1], simulator.registers[2], simulator.registers[3], simulator.registers[4], simulator.registers[5], simulator.registers[6], simulator.registers[7],
		simulator.registers[8], simulator.registers[9], simulator.registers[10], simulator.registers[11], simulator.registers[12], simulator.registers[13], simulator.registers[14], simulator.registers[15],
		simulator.registers[16], simulator.registers[17], simulator.registers[18], simulator.registers[19], simulator.registers[20], simulator.registers[21], simulator.registers[22], simulator.registers[23],
		simulator.registers[24], simulator.registers[25], simulator.registers[26], simulator.registers[27], simulator.registers[28], simulator.registers[29], simulator.registers[30], simulator.registers[31],
	)
	result += "data:\n"
	eightCount := 0
	for _, d := range simulator.data {
		if eightCount == 0 {
			result += fmt.Sprintf("%d:\t%d", d.GetProgramCount(), d.GetLineValue())
		} else if eightCount == 7 {
			result += fmt.Sprintf("\t%d\n", d.GetLineValue())
		} else {
			result += fmt.Sprintf("\t%d", d.GetLineValue())
		}
		eightCount = (eightCount + 1) % 8
	}
	if eightCount != 0 {
		result += "\n"
	}
	return result
}
