// Package comparealgs provides a utility to compare all included parallel
// sorting algorithms' performance using execution time as the metric.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/carlosgvaso/parallel-sort/bitonicsort"
	"github.com/carlosgvaso/parallel-sort/bricksort"
	"github.com/carlosgvaso/parallel-sort/mergesort"
	"github.com/carlosgvaso/parallel-sort/quicksort"
	"github.com/carlosgvaso/parallel-sort/radixsort"
)

// OutFile is the output file's path.
var outFile string = "output.txt"

// InFile is the input file's path.
var inFile string = "input.txt"

// InFileFormat is the input file's format.
//
// Formats are: 0 for CSV, 1 for array entry per line
var inFileFormat int = 0

// Runs is the number of times each algorithm is run to average execution time.
var runs int = 1000

// SleepTime is the time in seconds to sleep between runs to let the CPUs cool down.
var sleepTime time.Duration = 5

// FreeProcs is the number of processor cores to leave free (not use in parallel).
var freeProcs int = 2

// Error codes.
var exitOk int = 0  // Exit without errors
var exitArg int = 1 // Exit bad arguments
var exitErr int = 2 // Exit unknown error

// ArrayFlags defines the type for flags that are an array of entries.
//
// The input of these flags is in the following format:
// $ comparealgs -list1 value1 -list1 value2
type arrayFlags []string

// String returns the string representation of the arrayFlags type.
func (i *arrayFlags) String() string {
	var iStr string

	for j, entry := range *i {
		if j == 0 {
			iStr += "[ " + entry
		} else {
			iStr += ", " + entry
		}
	}
	iStr += " ]"

	return iStr
}

// Set sets the value of the arrayFlag type.
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

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

// ReadInputCsv reads the input file.
//
// It assumes the file is in CSV format with a single line of input integers.
//
// It returns arrays arrIn and arrOut with the comma-separated entries in the
// first line of the file converted to integers, and integer n with the length
// of the arrays.
func readInputCsv(inFile string) ([]int, []int, int) {
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

// ReadInputEntryPerLine reads the input file.
//
// It assumes the file is in a format with a single array integer entry per line
// of the file.
//
// It returns arrays arrIn and arrOut with the comma-separated entries in the
// first line of the file converted to integers, and integer n with the length
// of the arrays.
func readInputEntryPerLine(inFile string) ([]int, []int, int) {
	var n int = 0

	b, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatalln("Could not read the input file", err)
	}

	lines := strings.Split(string(b), "\n")
	// Assign cap to avoid resize on every append.
	n = len(lines)
	arrIn := make([]int, 0, n)
	arrOut := make([]int, 0, n)

	for _, l := range lines {
		// Empty line occurs at the end of the file when we use Split.
		if len(l) == 0 {
			continue
		}
		// Atoi better suits the job when we know exactly what we're dealing
		// with
		num, err := strconv.Atoi(l)
		if err != nil {
			log.Fatalln("Could not parse the input array", err)
		}
		arrIn = append(arrIn, num)
		arrOut = append(arrOut, num)
	}

	n = len(arrIn)

	return arrIn, arrOut, n
}

// MaxNumDigits gets number of digits of largest integer in array of positive
// integers.
//
// arr is the input array of positive integers.
// It return the number of digits of the largest integer in array
func maxNumDigits(arr []int) int {
	var k int = 0
	var max int = 0

	// Find the largest integer in array
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	// Find the number of characters of the largest integer
	for max != 0 {
		max /= 10
		k++
	}

	return k
}

// Main reads the array in the input file, and records the execution times each
// sorting algorithm takes to sort it.
//
// It saves the results to the specified file, or the default filename.
func main() {
	var arrIn, arrOut []int
	var n int
	var procs int = 0
	var err error
	cores := runtime.NumCPU()
	var algs arrayFlags

	// Check command-line arguments
	inFilePtr := flag.String("input", inFile, "Input file's path")
	inFileFormatPtr := flag.Int("format", inFileFormat, "Input file's format")
	outFilePtr := flag.String("output", outFile, "Output file's path")
	procsPtr := flag.Int("procs", (cores - freeProcs), "Maximum number of CPUs to use in parallel")
	runsPtr := flag.Int("runs", runs, "Number of times each algorithm is run to average execution time")
	flag.Var(&algs, "alg", "Specify algorithms to run. This flag should be called multiple times for each algorithm to run. Available algorithms are: bitonicsort, bricksort, mergesort, quicksort and radixsort")
	flag.Parse()

	inFile = *inFilePtr
	inFileFormat = *inFileFormatPtr
	outFile = *outFilePtr
	procs = *procsPtr
	runs = *runsPtr

	if algs == nil {
		algs.Set("bitonicsort")
		algs.Set("bricksort")
		algs.Set("mergesort")
		algs.Set("quicksort")
		algs.Set("radixsort")
	}

	// Read the input file
	switch inFileFormat {
	case 0:
		arrIn, arrOut, n = readInputCsv(inFile)
	case 1:
		arrIn, arrOut, n = readInputEntryPerLine(inFile)
	default:
		log.Fatalln("Unknown input file format")
	}

	// Get number of digits of largest integer in the input array for radix sort
	k := maxNumDigits(arrIn)

	// Open output file
	fout, err := os.Create(outFile)
	if err != nil {
		log.Fatalln("Could not open the output file", err)
	}

	// Print problem parameters
	runtime.GOMAXPROCS(procs)
	fmt.Printf("Input file: %s\nOutput file: %s\nLogical CPUs: %d\nMax procs: %d\nRuns: %d\nProblem size: n = %d\n",
		inFile, outFile, cores, procs, runs, n)
	fmt.Fprintf(fout, "Input=%s\nOutput=%s\nTimes measured in nsec\ncores=%d\nmaxProcs=%d\nruns=%d\nn=%d\narrIn=%v\n\n",
		inFile, outFile, cores, procs, runs, n, arrIn)

	// Print header
	fmt.Fprintf(fout, "Algorithm,")
	for i := 1; i <= runs; i++ {
		fmt.Fprintf(fout, "ExecTime%d,", i)
	}
	fmt.Fprintf(fout, "ExecTimeAvg\n")

	// Initialize variables to calculate the execution time averages
	var execTimeAvg int = 0

	for _, alg := range algs {
		switch alg {
		case "bitonicsort":
			// Run bitonic sort
			fmt.Printf("\tBitonic Sort:\n")
			fmt.Fprintf(fout, "bitonicsort,")

			// Setup variables to calculate the execution time averages
			execTimeAvg = 0

			// Run benchmarks
			for i := 0; i <= runs; i++ {
				// Append 0's to array to make its length exponential of 2
				var diff int
				arrOut, diff = bitonicsort.CheckAndAppendZeros(arrOut)

				// Bitonic sort sorts in place, so pass arrOut to preserve arrIn
				startTime := time.Now()
				arrOut = bitonicsort.Sort(arrOut, diff)
				execTime := time.Since(startTime)

				// Ignore the first run because it is always artificially slower
				if i > 0 {
					fmt.Printf("\t\tExec time %d: %s\n", i, execTime)
					fmt.Fprintf(fout, "%d,", int(execTime))

					// Add all times to average them
					execTimeAvg += int(execTime)
				}

				// Copy arrIn to arrOut for the next iteration
				copy(arrOut, arrIn)

				// Sleep between runs to let the CPUs cool down
				time.Sleep(sleepTime * time.Second)
			}

			// Calculate average
			execTimeAvg = execTimeAvg / runs
			fmt.Printf("\t\tExec time avg: %dns\n", execTimeAvg)
			fmt.Fprintf(fout, "%d\n", execTimeAvg)

		case "bricksort":
			// Run brick sort
			fmt.Printf("\tBrick Sort:\n")
			fmt.Fprintf(fout, "bricksort,")

			// Setup variables to calculate the execution time averages
			execTimeAvg = 0

			// Run benchmarks
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
					execTimeAvg += int(execTime)
				}

				// Copy arrIn to arrOut for the next iteration
				copy(arrOut, arrIn)

				// Sleep between runs to let the CPUs cool down
				time.Sleep(sleepTime * time.Second)
			}

			// Calculate average
			execTimeAvg = execTimeAvg / runs
			fmt.Printf("\t\tExec time avg: %dns\n", execTimeAvg)
			fmt.Fprintf(fout, "%d\n", execTimeAvg)

		case "mergesort":
			// Run mergesort
			fmt.Printf("\tMergesort:\n")
			fmt.Fprintf(fout, "mergesort,")

			// Setup variables to calculate the execution time averages
			execTimeAvg = 0

			// Run benchmarks
			for i := 0; i <= runs; i++ {
				// Mergesort sorts in place, so pass arrOut to preserve arrIn
				startTime := time.Now()
				arrOut = mergesort.Sort(arrOut)
				execTime := time.Since(startTime)

				// Ignore the first run because it is always artificially slower
				if i > 0 {
					fmt.Printf("\t\tExec time %d: %s\n", i, execTime)
					fmt.Fprintf(fout, "%d,", int(execTime))

					// Add all times to average them
					execTimeAvg += int(execTime)
				}

				// Copy arrIn to arrOut for the next iteration
				copy(arrOut, arrIn)

				// Sleep between runs to let the CPUs cool down
				time.Sleep(sleepTime * time.Second)
			}

			// Calculate average
			execTimeAvg = execTimeAvg / runs
			fmt.Printf("\t\tExec time avg: %dns\n", execTimeAvg)
			fmt.Fprintf(fout, "%d\n", execTimeAvg)

		case "quicksort":
			// Run quicksort
			fmt.Printf("\tQuicksort:\n")
			fmt.Fprintf(fout, "quickSort,")

			// Setup variables to calculate the execution time averages
			execTimeAvg = 0

			// Run benchmarks
			for i := 0; i <= runs; i++ {
				// Quicksort sorts in place, so pass arrOut to preserve arrIn
				startTime := time.Now()
				arrOut = quicksort.Sort(arrOut)
				execTime := time.Since(startTime)

				// Ignore the first run because it is always artificially slower
				if i > 0 {
					fmt.Printf("\t\tExec time %d: %s\n", i, execTime)
					fmt.Fprintf(fout, "%d,", int(execTime))

					// Add all times to average them
					execTimeAvg += int(execTime)
				}

				// Copy arrIn to arrOut for the next iteration
				copy(arrOut, arrIn)

				// Sleep between runs to let the CPUs cool down
				time.Sleep(sleepTime * time.Second)
			}

			// Calculate average
			execTimeAvg = execTimeAvg / runs
			fmt.Printf("\t\tExec time avg: %dns\n", execTimeAvg)
			fmt.Fprintf(fout, "%d\n", execTimeAvg)

		case "radixsort":
			// Run radix sort
			fmt.Printf("\tRadix Sort:\n")
			fmt.Fprintf(fout, "radixsort,")

			// Setup variables to calculate the execution time averages
			execTimeAvg = 0

			// Run benchmarks
			for i := 0; i <= runs; i++ {
				// Radix sort overwrites the input array, so pass arrOut to preserve arrIn
				startTime := time.Now()
				arrOut = radixsort.Sort(arrOut, k)
				execTime := time.Since(startTime)

				// Ignore the first run because it is always artificially slower
				if i > 0 {
					fmt.Printf("\t\tExec time %d: %s\n", i, execTime)
					fmt.Fprintf(fout, "%d,", int(execTime))

					// Add all times to average them
					execTimeAvg += int(execTime)
				}

				// Copy arrIn to arrOut for the next iteration
				copy(arrOut, arrIn)

				// Sleep between runs to let the CPUs cool down
				time.Sleep(sleepTime * time.Second)
			}

			// Calculate average
			execTimeAvg = execTimeAvg / runs
			fmt.Printf("\t\tExec time avg: %dns\n", execTimeAvg)
			fmt.Fprintf(fout, "%d\n", execTimeAvg)

		default:
			fmt.Printf("ERROR: %s is not a valid algoritm\nSkipping...\n", alg)
		}
	}

	// Close output file
	fout.Close()
}
