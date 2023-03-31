package main

import (
	"flag"

	f "ARM-Simulator/pkg/files"
	i "ARM-Simulator/pkg/instructions"
	s "ARM-Simulator/pkg/simulator"
)

func main() {
	inputFilePtr := flag.String("i", "", "the input file you want to parse")
	outputFilePtr := flag.String("o", "", "the output file you want to write")
	flag.Parse()

	lines := f.GetLinesFromFile(inputFilePtr)
	instructions, data := i.ConvertLinesToInstructions(lines)
	strings := i.ConvertInstructionsAndDataToStrings(instructions, data)
	f.WriteStringsToFile(outputFilePtr, "_dis.txt", strings)

	simulator := s.InitializeSimulator(data)
	strings = s.RunInstructionsAndReturnStrings(simulator, instructions)
	f.WriteStringsToFile(outputFilePtr, "_sim.txt", strings)
}
