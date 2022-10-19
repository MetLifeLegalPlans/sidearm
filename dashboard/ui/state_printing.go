package ui

import (
	"fmt"
	"io"
	"sort"
)

var (
	bold   = "[::b]"
	reset  = "[-:-:-]"
	green  = "[green]"
	yellow = "[yellow]"
	red    = "[red]"
)

func (b *ResultBucket) Print(out io.Writer) {
	b.PrintSuccessPercent(out)
	b.PrintAverageResponseTime(out)
	b.PrintResponseCodes(out)
}

func (b *ResultBucket) PrintSuccessPercent(out io.Writer) {
	successPercent := b.SuccessPercent
	color := green

	switch {
	case successPercent < 90:
		color = red
	case successPercent < 95:
		color = yellow
	}

	fmt.Fprintf(
		out,
		"%sSuccessful%s: %s%v%s%%\n",
		bold,
		reset,
		color,
		successPercent,
		reset,
	)
}

func (b *ResultBucket) PrintAverageResponseTime(out io.Writer) {
	avgTime := b.AverageResponseTime
	color := green

	switch {
	case avgTime > 650:
		color = red
	case avgTime > 350:
		color = yellow
	}

	fmt.Fprintf(
		out,
		"%sAverage Duration%s: %s%v%sms\n",
		bold,
		reset,
		color,
		avgTime,
		reset,
	)
}

func (b *ResultBucket) PrintResponseCodes(out io.Writer) {
	fmt.Fprint(out, "\n")

	keys := make([]int, 0)
	for key := range b.StatusCodes {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, code := range keys {
		count := b.StatusCodes[code]

		fmt.Fprintf(
			out,
			"%s%v%s: %v\n",
			bold,
			code,
			reset,
			count,
		)
	}
}
