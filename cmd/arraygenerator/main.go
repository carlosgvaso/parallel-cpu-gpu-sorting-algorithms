// Package arraygenerator provides a utility to generate an array of n
// pseudo-random integers in range [minInt, maxInt], which it saves to file.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// OutFile is the output file's path
var outFile string = "output.txt"

// N is the length of the output array
var n int = 100

// MaxInt is the maximum integer value allowed in any entry of the output array
var maxInt int = 999

// MinInt is the minimum integer value allowed in any entry of the output array
var minInt int = 0

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
	// Check command-line arguments
	nPtr := flag.Int("n", n, "Length of output array")
	maxIntPtr := flag.Int("max", maxInt, "Maximum value of any array entry")
	minIntPtr := flag.Int("min", minInt, "Minimum value of any array entry")
	outFilePtr := flag.String("output", outFile, "Output file's path")
	flag.Parse()

	n = *nPtr
	maxInt = *maxIntPtr
	minInt = *minIntPtr
	outFile = *outFilePtr

	if maxInt <= 0 {
		fmt.Printf("ERROR: max must be > 0\n")
		os.Exit(exitArg)
	} else if minInt < 0 {
		fmt.Printf("ERROR: min must be >= 0\n")
		os.Exit(exitArg)
	} else if minInt > maxInt {
		fmt.Printf("ERROR: max must be >= min\n")
		os.Exit(exitArg)
	}

	// Initialize arrOut
	arrOut := make([]int, n)

	// Fill array with pseudo-random integers in range [0, maxInt]
	rand.Seed(time.Now().UnixNano()) // Get seed value from the clock
	for i := 0; i < n; i++ {
		arrOut[i] = minInt + rand.Intn(maxInt-minInt+1)
	}

	// Write arrOut to file
	writeOutput(outFile, arrOut)
}
