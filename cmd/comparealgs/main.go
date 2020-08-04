// Package comparealgs provides a utility to compare all included parallel
// sorting algorithms' performance using execution time as the metric.
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/carlosgvaso/parallel-sort/bricksort"
)

// OutFile is the output file's path
var outFile string = "output.txt"

// InFile is the input file's path
var inFile string = "input.txt"

// Runs is the number of times each algorithm is run to average execution time
var runs int = 10

// Usage contains the usage informations
var usage string = `
Usage:
%s [inputFile outputFile]

Options:
    inputFile   Input file's path
    outputFile  Output file's path
`

// Error codes
var exitOk int = 0  // Exit without errors
var exitArg int = 1 // Exit bad arguments
var exitErr int = 2 // Exit unknown error

// LoadArrays loads the arrIn and arrOut arrays in parallel with the provided
// value val converted to an int.
func loadArrays(val string, arrIn []int, arrOut []int, i int,
	waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	var err error

	// Convert val from string to int, and save it to arrIn
	arrIn[i], err = strconv.Atoi(val)
	if err != nil {
		log.Fatalln("Could not parse the input array", err)
	}

	// Copy arrIn to arrOut
	arrOut[i] = arrIn[i]
}

// ReadInput reads the input file.
//
// It assumes the file is in CSV format with a single line of input integers.
//
// It returns arrays arrIn and arrOutwith the comma-separated entries in the
// first line of the file converted to integers, and integer n with the length
// of the arrays.
func readInput(inFile string) ([]int, []int, int) {
	// Open the file
	csvfile, err := os.Open(inFile)
	if err != nil {
		log.Fatalln("Could not open the input file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	// Read first line only
	record, err := r.Read()
	if err != nil {
		log.Fatalln("Could not read the input file", err)
	}

	// Close file
	csvfile.Close()

	// Setup in/out arrays
	var n int = len(record)
	arrIn := make([]int, n)
	arrOut := make([]int, n)

	// Load the arrIn and arrOut arrays in parallel with the values in record,
	// and convert those values from string to int. Note, arrIn = arrOut, since
	// some methods sort the array in place. In those cases, we pass the arrOut
	// as the input.
	var waitGroup sync.WaitGroup // Wait group to synchronize parallel goroutines
	for i, v := range record {
		waitGroup.Add(1)
		go loadArrays(v, arrIn, arrOut, i, &waitGroup)
	}
	waitGroup.Wait()

	return arrIn, arrOut, n
}

// Main reads the array in the input file, and records the execution times each
// sorting algorithm takes to sort it.
//
// It saves the results to the specified file, or the default filename.
func main() {
	argv := os.Args
	argc := len(argv)

	// Check command-line arguments
	if argc != 3 && argc != 1 {
		fmt.Printf("ERROR: Wrong number of arguments provided\n")
		fmt.Printf(usage, argv[0])
		os.Exit(exitArg)
	} else if argc == 3 {
		inFile = argv[1]
		outFile = argv[2]
	}

	// Read the input file
	arrIn, arrOut, n := readInput(inFile)

	// Open output file
	fout, err := os.Create(outFile)
	if err != nil {
		log.Fatalln("Could not open the output file", err)
	}

	// Setup variables to calculate the execution time averages
	var execTimeAvgBrickSort int = 0

	// Print problem parameters
	fmt.Printf("Problem size: n = %d\n", n)
	fmt.Fprintf(fout, "n = %d\narrIn = %v\nTimes measured in nsec\n\n", n, arrIn)

	// Print header
	fmt.Fprintf(fout, "Algorithm,")
	for i := 1; i <= runs; i++ {
		fmt.Fprintf(fout, "ExecTime%d,", i)
	}
	fmt.Fprintf(fout, "ExecTimeAvg\n")

	// Run Brick Sort
	fmt.Printf("\tBrick Sort:\n")
	fmt.Fprintf(fout, "BrickSort,")

	for i := 0; i <= runs; i++ {
		// Brick sort sorts in place, so pass arrOut to preserve arrIn
		startTime := time.Now()
		arrOut = bricksort.Sort(arrOut)
		execTime := time.Since(startTime)

		// Ignore the first run because it is always artificially slower
		if i > 0 {
			fmt.Printf("\t\tExec time %d: %s\n", i, execTime)
			fmt.Fprintf(fout, "%d,", int(execTime))

			// Add all times to average them
			execTimeAvgBrickSort += int(execTime)
		}
	}

	// Calculate average
	execTimeAvgBrickSort = execTimeAvgBrickSort / runs
	fmt.Fprintf(fout, "%d\n", execTimeAvgBrickSort)

	// Close output file
	fout.Close()
}
