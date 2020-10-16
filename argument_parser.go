package main

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

const usage = "[options] <from> <to> [<path>]"

type arguments struct {
	Bare         bool   `short:"b" long:"bare" description:"Use patterns as they are"`
	RawCaseNames string `short:"c" long:"cases" description:"Comma-separated names of enabled cases (options: camel, upper-camel, kebab, upper-kebab, snake, upper-snake, space, upper-space)"`
	Verbose      bool   `short:"v" long:"verbose" description:"Be verbose"`
	Help         bool   `short:"h" long:"help" description:"Show this help"`
	Version      bool   `long:"version" description:"Show version"`
	From         string
	To           string
	Path         string
	CaseNames    map[caseName]struct{}
}

type argumentParser struct {
	workingDirectory string
}

func newArgumentParser(d string) *argumentParser {
	return &argumentParser{d}
}

func (p *argumentParser) Parse(ss []string) (*arguments, error) {
	args := arguments{}
	parser := flags.NewParser(&args, flags.PassDoubleDash)
	parser.Usage = usage

	ss, err := parser.ParseArgs(ss)
	if err != nil {
		return nil, err
	} else if args.Help || args.Version {
		return &args, nil
	} else if len(ss) < 2 || len(ss) > 3 {
		return nil, errors.New("invalid number of arguments")
	}

	args.From, args.To = ss[0], ss[1]

	if len(ss) == 3 {
		args.Path = p.resolvePath(ss[2])
	} else {
		args.Path = p.workingDirectory
	}

	if args.RawCaseNames != "" {
		args.CaseNames = map[caseName]struct{}{}

		for _, n := range strings.Split(args.RawCaseNames, ",") {
			n := caseName(n)

			if _, ok := allCaseNames[n]; !ok {
				return nil, fmt.Errorf("invalid case name: %v", n)
			}

			args.CaseNames[n] = struct{}{}
		}
	}

	return &args, nil
}

func (*argumentParser) Help() string {
	p := flags.NewParser(&arguments{}, flags.PassDoubleDash)
	p.Usage = usage

	// Parse() is run here to show default values in help.
	// This seems to be a bug in go-flags.
	p.Parse() // nolint:errcheck

	b := &bytes.Buffer{}
	p.WriteHelp(b)
	return b.String()
}

func (p *argumentParser) resolvePath(s string) string {
	if filepath.IsAbs(s) {
		return s
	}

	return filepath.Join(p.workingDirectory, s)
}