package main

import (
	"testing"
)

var runningAvgTestData = []struct {
	input    []float64
	expected float64
}{
	{[]float64{14, 53.2, 6, 0, -10, 12.6, 6, 12}, 11.725},
	{[]float64{1, 2, 3, 5, 8, -10}, 1.5},
}

func TestRunningAvg(t *testing.T) {
	for _, test := range runningAvgTestData {
		input := make(chan float64)
		result := make(chan float64)
		go RunningAvg(input, result)
		for _, x := range test.input {
			input <- x
		}
		close(input)
		if actual := <-result; actual != test.expected {
			t.Errorf("RunningAvg(%d) expected %d, Actual %d", test.input, test.expected, actual)
		}
	}
}

var orderTestData = []struct {
	input    [2]string
	expected Order
}{
	{[2]string{"2015-06-01 00:27:24", "121.2"}, Order{"2015-06-01", float64(121.2)}},
	{[2]string{"2017-03-01 10:12:34", "127.2"}, Order{"2017-03-01", float64(127.2)}},
}

func TestOrderConstructor(t *testing.T) {
	for _, test := range orderTestData {
		if actual, _ := NewOrder(test.input[0], test.input[1]); actual != test.expected {
			t.Errorf("NewOrder(%d) expected %d, Actual %d", test.input, test.expected, actual)
		}
	}
}
