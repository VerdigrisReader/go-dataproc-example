// Package main calculates the daily AOV from an incoming stream of order records
package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

const datefmt = "2006-01-02 15:04:05"

// type Order represents the monetary value of an order on a date
type Order struct {
	date  string
	value float64
}

// type Thread Manages the input & output thread of a running goroutine
type Thread struct {
	input  chan float64
	output chan float64
}

// NewOrder constructs a new order object by parsing the raw date/monetary values from strings
func NewOrder(date, value string) (Order, error) {
	parsedDate, err := time.Parse(datefmt, date)
	if err != nil {
		return *new(Order), errors.New(fmt.Sprintf("Error parsing date: %v %v\n", date, err))
	}
	parsedFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return *new(Order), errors.New(fmt.Sprintf("Error parsing value: %v %v\n", value, err))
	}
	return Order{parsedDate.Truncate(24 * time.Hour).Format("2006-01-02"), parsedFloat}, nil
}

// NewThread is a constructor for 'Thread' objects, which contains two fresh channels.
// @todo - Change RunningAvg to recieve a Thread object rather than two channels
func NewThread() *Thread {
	return &Thread{
		input:  make(chan float64),
		output: make(chan float64),
	}
}

// RunningAvg receives values from a channel and passes a 
// running average to the output channel once
// all values are exhausted
func RunningAvg(ch chan float64, output chan float64) {
	var count float64
	var total float64
	for x := range ch {
		count++
		total += x
	}
	output <- total / count
	close(output)
}

// RowProcessor recieves order objects and ensures that the revenue from each
// is passed to a RunningAvg goroutine for its date.
func RowProcessor(ch chan Order, output chan map[string]float64) {
	var threads = make(map[string]*Thread)
	for o := range ch {
		if thread, ok := threads[o.date]; ok {
			thread.input <- o.value
		} else {
			thread := NewThread()
			go RunningAvg(thread.input, thread.output)
			threads[o.date] = thread
		}
	}
	var result = make(map[string]float64)
	for key, thread := range threads {
		close(thread.input)
		result[key] = <-thread.output
	}
	output <- result
	close(output)
}

// Prints the results returned from RowProcessor, a list of dates & values
func PrintResults(result map[string]float64) {
	keys := make([]string, 0, len(result))
	for k, _ := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		avg, _ := result[k]
		fmt.Printf("%v  : %v\n", k, avg)
	}
}

func main() {
	stat, err := os.Stdin.Stat()
	if (stat.Mode() & os.ModeNamedPipe) == 0 || err != nil {
		fmt.Println("Nothing in stdin")
		os.Exit(0)
	}

	r := csv.NewReader(bufio.NewReader(os.Stdin))

	input := make(chan Order)
	output := make(chan map[string]float64)
	go RowProcessor(input, output)

	for i := 0; ; i++ {
		if i == 0 {
			continue
		}
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		// Skip any parse errors
		order, err := NewOrder(record[6], record[5])
		if err == nil {
			input <- order
		}

	}
	close(input)
	result := <-output

	PrintResults(result)
}
