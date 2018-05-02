// +build ignore

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()

	var in io.Reader = os.Stdin
	var out io.Writer = os.Stdout

	if infile := flag.Arg(0); infile != "-" && infile != "" {
		f, err := os.Open(infile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	if outfile := flag.Arg(0); outfile != "-" && outfile != "" {
		f, err := os.Create(outfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		out = f
	}

	fmt.Fprintf(out, "Queue, Size, Cores, Setup, Test, Time\n")

	rx := regexp.MustCompile(`[/\t -]+`)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		tokens := rx.Split(scanner.Text(), -1)
		if !strings.HasPrefix("Benchmark", tokens[0]) {
			continue
		}
		if len(tokens) != 12 {
			log.Fatal("invalid number of results")
			continue
		}

		//  0  Benchmark
		//  1  MPMCcGo
		//  2  b0s256
		//  3  b
		//  4  MPSC
		//  5  ProducerConsumer
		//  6  x100
		//  7  8
		//  8  200
		//  9  34965
		// 10  ns
		// 11  op

		queueName := tokens[1]
		queueSize := tokens[2]
		setup := tokens[4]
		test := tokens[3] + "/" + tokens[5]
		repeatCount, err := strconv.ParseFloat(tokens[6][1:], 64)
		if err != nil {
			log.Fatal(err)
		}

		cores := tokens[7]
		time, err := strconv.ParseFloat(tokens[9], 64)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(out, "%v, %v, %v, %v, %v, %.2f\n", queueName, queueSize, cores, setup, test, time/repeatCount)
	}
}
