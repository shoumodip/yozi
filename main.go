package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"yozi/checker"
	"yozi/compiler"
	"yozi/lexer"
	"yozi/parser"
)

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "    yozi [FLAGS] <FILE>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "    -h, -help    Show this help message")
	fmt.Fprintln(w, "    -r, -run     Run the program after compiling it")
}

type Args struct {
	run  bool
	path string
	rest []string
}

func parseArgs() Args {
	args := Args{
		run:  false,
		path: "",
		rest: os.Args[1:],
	}

	for len(args.rest) != 0 {
		arg := args.rest[0]
		args.rest = args.rest[1:]

		switch arg {
		case "-h", "-help", "--help":
			usage(os.Stdout)
			os.Exit(0)

		case "-r", "-run", "--run":
			args.run = true

		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintln(os.Stderr, "ERROR: Invalid flag '"+arg+"'")
				fmt.Fprintln(os.Stderr)
				usage(os.Stderr)
				os.Exit(1)
			}

			args.path = arg
			return args
		}
	}

	fmt.Fprintln(os.Stderr, "ERROR: Input file not provided")
	fmt.Fprintln(os.Stderr)
	usage(os.Stderr)
	os.Exit(1)

	return args
}

func main() {
	args := parseArgs()

	lexer, err := lexer.New(args.path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: Could not open file '"+args.path+"'")
		fmt.Fprintln(os.Stderr)
		usage(os.Stderr)
		os.Exit(1)
	}

	parser := parser.Parser{}
	parser.File(lexer)

	for _, node := range parser.Nodes {
		checker.Check(node)
	}

	exePath := strings.TrimSuffix(args.path, ".yo")
	compiler.Program(parser.Nodes, exePath)
	if args.run {
		cmd := exec.Command("./" + exePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	}
}
