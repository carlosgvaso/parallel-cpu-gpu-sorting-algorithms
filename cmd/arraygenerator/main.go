// Package arraygenerator provides a utility to generate an array of n
// pseudo-random integers in range [0, 999], which it saves to file.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// OutFile is the output file's path
var outFile string = "output.txt"

// MaxInt is the maximum integer value allowed in any entry of the output array
var maxInt int = 999

// Usage contains the usage informations
var usage string = `
Usage:
%s n [-h] [outputFile]

Options:
    -h --help   Print usage information
    n           Length of array as int
    outputFile  Output file's path
`

// Error codes
var exitOk int = 0  // Exit without errors
var exitArg int = 1 // Exit bad arguments
var exitErr int = 2 // Exit unknown error

// WriteOutput writes the arrOut to the output file in CVS format.
//
// It takes the output filename and the array to save to file.
func writeOutput(outFile string, arrOut []int) {
	// Open output file
	fout, err := os.Create(outFile)
	if err != nil {
		log.Fatalln("Could not open the output file", err)
	}

	// Iterate over the array entries, and write them to file
	var n int = len(arrOut)
	for i, v := range arrOut {
		if i == n-1 {
			fmt.Fprintf(fout, "%d", v)
		} else {
			fmt.Fprintf(fout, "%d,", v)
		}
	}

	// Close output file
	fout.Close()
}

// Main generates an array of random numbers of size n, and writes the array to
// the output file.
//
// It saves the results to the specified file, or the default filename.
func main() {
	argv := os.Args
	argc := len(argv)

	// Check command-line arguments
	var n int = 0
	var err error
	switch argc {
	case 2: // If -h flag, print usage; else, save int to n
		if argv[1] == "-h" || argv[1] == "--help" {
			fmt.Printf(usage, argv[0])
			os.Exit(exitOk)
		} else {
			n, err = strconv.Atoi(argv[1])
			if err != nil {
				fmt.Printf("Could not parse the array length, n\n")
				fmt.Printf(usage, argv[0])
				os.Exit(exitArg)
			}
		}
	case 3: // Use 2 args for n and  outFile in this order
		n, err = strconv.Atoi(argv[1])
		if err != nil {
			fmt.Printf("Could not parse the array length, n\n")
			fmt.Printf(usage, argv[0])
			os.Exit(exitArg)
		}
		outFile = argv[2]
	default: // Wrong number of args provided
		fmt.Printf("ERROR: Wrong number of arguments provided\n")
		fmt.Printf(usage, argv[0])
		os.Exit(exitArg)
	}

	// Initialize arrOut
	arrOut := make([]int, n)

	// Fill array with pseudo-random integers in range [0, maxInt]
	rand.Seed(time.Now().UnixNano()) // Get seed value from the clock
	for i := 0; i < n; i++ {
		arrOut[i] = rand.Intn(maxInt)
	}

	// Write arrOut to file
	writeOutput(outFile, arrOut)
}
