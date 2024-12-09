package main

import (
	"fmt"

	"github.com/hellflame/argparse"
)

type config struct {
	seed       int
	iterations int
	interval   int

	style style
}

type style struct {
	width  int
	height int

	boxSize  int
	boxGap   int
	boxRound int
}

func parseArgs() *config {
	p := argparse.NewParser("", "", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,

		WithHint: true,
	})
	seed := p.Int("", "seed", &argparse.Option{Default: "0"})
	iterations := p.Int("i", "iterates", &argparse.Option{Default: "1000", Help: "0 means running forever"})
	interval := p.Int("", "interval", &argparse.Option{Default: "200"})

	width := p.Int("", "width", &argparse.Option{Default: "200"})
	height := p.Int("", "height", &argparse.Option{Default: "100"})
	boxSize := p.Int("s", "box-size", &argparse.Option{Default: "10"})
	boxGap := p.Int("g", "box-gap", &argparse.Option{Default: "1"})
	boxRound := p.Int("r", "box-round", &argparse.Option{Default: "2"})

	if e := p.Parse(nil); e != nil {
		if e != argparse.BreakAfterHelpError {
			fmt.Println(e)
		}
		return nil
	}
	c := config{
		seed:       *seed,
		iterations: *iterations,
		interval:   *interval,

		style: style{
			width:    *width,
			height:   *height,
			boxSize:  *boxSize,
			boxGap:   *boxGap,
			boxRound: *boxRound,
		},
	}

	return &c
}
