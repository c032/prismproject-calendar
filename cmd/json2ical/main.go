package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	ppc "github.com/c032/prismproject-calendar"
)

var (
	inputFile  string
	outputFile string
)

func main() {
	flag.StringVar(&inputFile, "i", "-", "JSON file to read.")
	flag.StringVar(&outputFile, "o", "-", "iCalendar file to write.")

	flag.Parse()

	var (
		r io.Reader
		w io.Writer
	)

	if inputFile == "-" {
		r = os.Stdin
	} else {
		var (
			err error
			f   *os.File
		)

		f, err = os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %#v\n", inputFile)
		}
		defer f.Close()

		r = f
	}

	if outputFile == "-" {
		w = os.Stdout
	} else {
		var (
			err error
			f   *os.File
		)

		f, err = os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file to write: %#v\n", outputFile)
		}
		defer f.Close()

		w = f
	}

	var (
		err error

		feed *ppc.Feed
	)

	feed, err = ppc.ParseFeed(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing feed.\n")

		os.Exit(1)
	}

	calendarStr, err := feed.CalendarString(ppc.DefaultCalendarID, "PRISM Project")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, calendarStr)
}
