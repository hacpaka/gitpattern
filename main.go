package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	iPattern string
	iDate string
	iNow string
	iDebug bool
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %v\n", r)
			os.Exit(2)
		}
	}()

	flag.StringVar(&iPattern,"pattern", "", "Pattern name")
	flag.StringVar(&iDate, "date", "", "Starting date")
	flag.StringVar(&iNow, "now", "", "Current date")
	flag.BoolVar(&iDebug, "debug", false, "Debug mode")
	flag.Parse()

	data, _ := readPattern(iPattern)
	start, _ := time.Parse("20060102", iDate)

	current := time.Now()
	if len(iNow) > 0 {
		current, _ = time.Parse("20060102", iNow)
	}

	days := int(current.Sub(start).Hours() / 24)
	if days < 0 {
		panic(errors.New("the date interval is invalid"))
	}

	position := days % len(data)

	col := position / 7
	row := (position - col * 7) & 7

	index := row * 10 + col
	if index > len(data) {
		panic(errors.New("the date interval is invalid"))
	}

	if iDebug {
		fmt.Printf("Cur: '%v'\n", current)
		fmt.Printf("Col: '%v'\n", col)
		fmt.Printf("Row: '%v'\n", row)
		fmt.Printf("Ind: '%v'\n", index)
		fmt.Printf("Val: '%v'\n", data[position-1])
	}

	if data[index] != 0 {
		if iDebug {
			fmt.Println("Pattern failed!")
		}

		os.Exit(1)
	}

	if iDebug {
		fmt.Println("Pattern successful!")
	}

	os.Exit(0)
}

func readPattern(name string) ([]int, error) {
	if len(name) < 1 {
		panic(errors.New("invalid pattern"))
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	iPattern = dir + string(os.PathSeparator )+ "patterns" + string(os.PathSeparator) + name
	if iDebug {
		fmt.Printf("Pattern file: %v\n", iPattern)
	}

	if _, err := os.Stat(iPattern); err != nil {
		if os.IsNotExist(err) {
			panic(errors.New("invalid pattern"))
		} else {
			panic(err)
		}
	}

	file, err := os.Open(iPattern)
	if err != nil {
		panic("invalid pattern")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(errors.New("invalid pattern"))
	}

	var output []int
	for  i := 0; i < len(data); i++ {
		if data[i] != 48 && data[i] != 49 {
			continue
		}

		output = append(output, int(data[i])-48)
	}

	if len(output) / 10 != 7 {
		panic("invalid pattern")
	}

	return output, nil
}

