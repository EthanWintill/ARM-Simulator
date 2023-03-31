package instructions

import (
	"fmt"
	"strconv"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string
	linevalue         uint64
	programCnt        int
	opcode            uint64
	op                string
	rd                uint8
	rn                uint8
	rm                uint8
	im                int16
	address           uint16
	shiftCode         uint16
	field             uint16
	offset            int32
	shamt             uint8
	args              string
}

func (i *Instruction) GetOp() string {
	return i.op
}
func (i *Instruction) GetProgramCount() int {
	return i.programCnt
}
func (i *Instruction) GetArgs() string {
	return i.args
}
func (i *Instruction) GetOffset() int32 {
	return i.offset
}
func (i *Instruction) GetRd() uint8 {
	return i.rd
}
func (i *Instruction) GetRn() uint8 {
	return i.rn
}
func (i *Instruction) GetRm() uint8 {
	return i.rm
}
func (i *Instruction) GetIm() int16 {
	return i.im
}
func (i *Instruction) GetShamt() uint8 {
	return i.shamt
}
func (i *Instruction) GetAddress() uint16 {
	return i.address
}
func (i *Instruction) GetShiftCode() uint16 {
	return i.shiftCode
}
func (i *Instruction) GetField() uint16 {
	return i.field
}

type Data struct {
	rawInstruction string
	linevalue      int32
	programCnt     int
}

func (d *Data) GetLineValue() int {
	return int(d.linevalue)
}
func (d *Data) SetLineValue(value int32) {
	d.linevalue = value
}
func (d *Data) GetProgramCount() int {
	return d.programCnt
}
func (d *Data) SetProgramCount(value int) {
	d.programCnt = value
}

func ConvertLinesToInstructions(lines []string) ([]Instruction, []Data) {
	var instructions []Instruction
	var data []Data
	haveEncounteredBreak := false
	programCount := 96
	for _, line := range lines {
		if haveEncounteredBreak {
			d := BinaryToData(line)
			d.programCnt = programCount
			data = append(data, d)
		} else {
			i := BinaryToInstruction(line)
			i.programCnt = programCount
			instructions = append(instructions, i)
			if instructions[len(instructions)-1].op == "BREAK" {
				haveEncounteredBreak = true
			}
		}
		programCount += 4
	}
	return instructions, data
}

func BinaryToInstruction(binary string) Instruction {
	var result Instruction
	result.rawInstruction = binary
	result.linevalue = BinaryToUnsignedInt(binary)
	result.op, result.typeofInstruction = ValueToOpAndType(result.linevalue)
	ParseInstructionParameters(&result)
	return result
}

func BinaryToData(binary string) Data {
	var result Data
	result.rawInstruction = binary
	value := BinaryToUnsignedInt(binary)
	result.linevalue = TwosComplement(value)
	return result
}

func BinaryToUnsignedInt(binary string) uint64 {
	result, err := strconv.ParseUint(binary, 2, 32)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func TwosComplement(value uint64) int32 {
	var result int32
	sign := int16(value & 0b1)
	if sign != 0 {
		result = (int32((^value)) + 1) * -1
	} else {
		result = int32(value)
	}
	return result
}

func ValueToOpAndType(opcode uint64) (string, string) {
	opcode = opcode >> 21
	op_code_type, isR := Ropcodes[opcode] //11
	if isR {
		return op_code_type, "R"
	}
	op_code_type, isD := Dopcodes[opcode] //11
	if isD {
		return op_code_type, "D"
	}
	op_code_type, isI := Iopcodes[opcode>>1] //10
	if isI {
		return op_code_type, "I"
	}
	op_code_type, isIM := IMopcodes[opcode>>2] //9
	if isIM {
		return op_code_type, "IM"
	}
	op_code_type, isCB := CBopcodes[opcode>>3] //8
	if isCB {
		return op_code_type, "CB"
	}
	if (opcode >> 5) == 0b000101 {
		return "B", "B"
	}
	if opcode == 0 {
		return "NOP", "NOP"
	}
	if opcode == 0b11111110110 {
		return "BREAK", "BREAK"
	}

	return "", ""
}

func ParseInstructionParameters(instruction *Instruction) {
	switch instruction.typeofInstruction {
	case "R":
		instruction.opcode = instruction.linevalue >> 21
		instruction.rm = uint8(instruction.linevalue >> 16 & 0b11111) //and operation leaves 5 rightmost bits
		instruction.shamt = uint8(instruction.linevalue >> 10 & 0b111111)
		instruction.rn = uint8(instruction.linevalue >> 5 & 0b11111)
		instruction.rd = uint8(instruction.linevalue & 0b11111)
	case "D":
		instruction.opcode = instruction.linevalue >> 21
		instruction.address = uint16(instruction.linevalue >> 12 & 0b111111111)
		instruction.rn = uint8(instruction.linevalue >> 5 & 0b11111)
		instruction.rd = uint8(instruction.linevalue & 0b11111)
		instruction.args = fmt.Sprintf("%s	R%d, [R%d, #%d]",
			instruction.op,
			instruction.rd,
			instruction.rn,
			instruction.address,
		)
	case "I":
		instruction.opcode = instruction.linevalue >> 22
		sign := int16(instruction.linevalue >> 21 & 0b1)
		if sign == 0 {
			instruction.im = int16(instruction.linevalue >> 10 & 0b11111111111)
		} else {
			instruction.im = (int16((^instruction.linevalue>>10)&0b11111111111) + 1) * -1
		}
		instruction.rn = uint8(instruction.linevalue >> 5 & 0b11111)
		instruction.rd = uint8(instruction.linevalue & 0b11111)
		instruction.args = fmt.Sprintf("%s	R%d, R%d, #%d",
			instruction.op,
			instruction.rd,
			instruction.rn,
			instruction.im,
		)
	case "IM":
		instruction.opcode = instruction.linevalue >> 23
		instruction.shiftCode = 16 * (uint16(instruction.linevalue >> 21 & 0b11))
		instruction.field = uint16(instruction.linevalue >> 5 & 0b1111111111111111)
		instruction.rd = uint8(instruction.linevalue & 0b11111)
		instruction.args = fmt.Sprintf("%s	R%d, %d, %s %d",
			instruction.op,
			instruction.rd,
			instruction.field,
			"LSL",
			instruction.shiftCode,
		)
	case "CB":
		instruction.opcode = instruction.linevalue >> 24
		sign := int16(instruction.linevalue >> 23 & 0b1)
		if sign == 0 {
			instruction.offset = int32(instruction.linevalue >> 5 & 0b111111111111111111)
		} else {
			instruction.offset = (int32((^instruction.linevalue>>5)&0b111111111111111111) + 1) * -1
		}
		instruction.rd = uint8(instruction.linevalue & 0b11111)
		instruction.args = fmt.Sprintf("%s	R%d, #%d",
			instruction.op,
			instruction.rd,
			instruction.offset,
		)
	case "B":
		instruction.opcode = instruction.linevalue >> 26
		sign := int16(instruction.linevalue >> 25 & 0b1)
		if sign == 0 {
			instruction.offset = int32(instruction.linevalue & 0b1111111111111111111111111)
		} else {
			instruction.offset = (int32((^instruction.linevalue)&0b1111111111111111111111111) + 1) * -1
		}
		instruction.args = fmt.Sprintf("%s	#%d",
			instruction.op,
			instruction.offset,
		)
	case "NOP":
		instruction.opcode = 0
		instruction.args = instruction.op
	case "BREAK":
		instruction.args = instruction.op
	default:
	}
	if instruction.op == "AND" || instruction.op == "ADD" || instruction.op == "EOR" || instruction.op == "ORR" || instruction.op == "SUB" {
		instruction.args = fmt.Sprintf("%s	R%d, R%d, R%d",
			instruction.op,
			instruction.rd,
			instruction.rn,
			instruction.rm)
	}
	if instruction.op == "LSL" || instruction.op == "ASR" || instruction.op == "LSR" {
		instruction.args = fmt.Sprintf("%s	R%d, R%d, #%d",
			instruction.op,
			instruction.rd,
			instruction.rn,
			instruction.shamt,
		)
	}
}

func InstructionToString(input Instruction) string {
	var result string
	if input.op == "AND" || input.op == "ADD" || input.op == "EOR" || input.op == "ORR" || input.op == "SUB" {
		result = fmt.Sprintf("%s %s %s %s %s	%d	%s\n",
			input.rawInstruction[0:11],
			input.rawInstruction[11:16],
			input.rawInstruction[16:22],
			input.rawInstruction[22:27],
			input.rawInstruction[27:32],
			input.programCnt,
			input.args,
		)
	} else if input.op == "LSL" || input.op == "ASR" || input.op == "LSR" {
		result = fmt.Sprintf("%s %s %s %s %s	%d	%s\n",
			input.rawInstruction[0:11],
			input.rawInstruction[11:16],
			input.rawInstruction[16:22],
			input.rawInstruction[22:27],
			input.rawInstruction[27:32],
			input.programCnt,
			input.args,
		)
	} else if input.typeofInstruction == "D" {
		result = fmt.Sprintf("%s %s %s %s %s	%d	%s\n",
			input.rawInstruction[0:11],
			input.rawInstruction[11:20],
			input.rawInstruction[20:22],
			input.rawInstruction[22:27],
			input.rawInstruction[27:32],
			input.programCnt,
			input.args,
		)
	} else if input.typeofInstruction == "I" {
		result = fmt.Sprintf("%s %s %s %s	%d	%s\n",
			input.rawInstruction[0:10],
			input.rawInstruction[10:22],
			input.rawInstruction[22:27],
			input.rawInstruction[27:32],
			input.programCnt,
			input.args,
		)
	} else if input.typeofInstruction == "B" {
		result = fmt.Sprintf("%s %s	%d	%s\n",
			input.rawInstruction[0:6],
			input.rawInstruction[6:32],
			input.programCnt,
			input.args,
		)
	} else if input.typeofInstruction == "CB" {
		result = fmt.Sprintf("%s %s %s	%d	%s\n",
			input.rawInstruction[0:8],
			input.rawInstruction[8:27],
			input.rawInstruction[27:32],
			input.programCnt,
			input.args,
		)
	} else if input.typeofInstruction == "IM" {
		result = fmt.Sprintf("%s %s %s %s	%d	%s\n",
			input.rawInstruction[0:9],
			input.rawInstruction[9:11],
			input.rawInstruction[11:27],
			input.rawInstruction[27:32],
			input.programCnt,
			input.args,
		)
	} else if input.typeofInstruction == "NOP" {
		result = fmt.Sprintf("%s	%d	%s\n",
			input.rawInstruction,
			input.programCnt,
			input.op,
		)
	} else if input.typeofInstruction == "BREAK" {
		result = fmt.Sprintf("%s %s %s %s %s %s %s	%d	%s\n",
			input.rawInstruction[0:1],
			input.rawInstruction[1:6],
			input.rawInstruction[6:11],
			input.rawInstruction[11:16],
			input.rawInstruction[16:21],
			input.rawInstruction[21:26],
			input.rawInstruction[26:32],
			input.programCnt,
			input.op,
		)
	}
	return result
}

func DataToString(d Data) string {
	return fmt.Sprintf("%s	%d	%d\n",
		d.rawInstruction,
		d.programCnt,
		d.linevalue,
	)
}

func ConvertInstructionsAndDataToStrings(iSlice []Instruction, dSlice []Data) []string {
	var results []string
	for _, i := range iSlice {
		results = append(results, InstructionToString(i))
	}
	for _, d := range dSlice {
		results = append(results, DataToString(d))
	}
	return results
}

var Ropcodes = map[uint64]string{
	0b10001010000: "AND",
	0b10001011000: "ADD",
	0b10101010000: "ORR",
	0b11001011000: "SUB",
	0b11010011011: "LSL",
	0b11010011010: "LSR",
	0b11010011100: "ASR",
	0b11101010000: "EOR"}
var Dopcodes = map[uint64]string{
	0b11111000000: "STUR",
	0b11111000010: "LDUR"}
var Iopcodes = map[uint64]string{
	0b1001000100: "ADDI",
	0b1101000100: "SUBI"}
var IMopcodes = map[uint64]string{
	0b110100101: "MOVZ",
	0b111100101: "MOVK"}
var CBopcodes = map[uint64]string{
	0b10110100: "CBZ",
	0b10110101: "CBNZ"}
