package files

import (
	"bufio"
	"fmt"
	"os"
)

func GetLinesFromFile(filename *string) []string {
	readFile, err := os.Open(*filename)
	if err != nil {
		fmt.Println("There was an error opening the input file.")
		fmt.Println(err)
		os.Exit(1)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var results []string
	for fileScanner.Scan() {
		results = append(results, fileScanner.Text())
	}
	return results
}

func WriteStringsToFile(filename *string, extension string, strings []string) {
	writeFile, err := os.Create(*filename + extension)
	if err != nil {
		fmt.Println("There was an error opening the output file.")
		fmt.Println(err)
		os.Exit(1)
	}
	defer writeFile.Close()

	for _, s := range strings {
		_, err := writeFile.WriteString(s)
		if err != nil {
			fmt.Println("There was an error writing the output file.")
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
