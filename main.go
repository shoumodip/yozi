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
	fmt.Fprintln(w, "    -h           Show this help message")
	fmt.Fprintln(w, "    -r           Run the program after compiling it")
	fmt.Fprintln(w, "    -o <name>    Set the name of the output executable")
}

type Args struct {
	run  bool
	rest []string

	inputPath  string
	outputPath string
}

func parseArgs() Args {
	args := Args{
		run:  false,
		rest: os.Args[1:],

		inputPath:  "",
		outputPath: "",
	}

	for len(args.rest) != 0 {
		arg := args.rest[0]
		args.rest = args.rest[1:]

		switch arg {
		case "-h", "-help", "--help":
			usage(os.Stdout)
			os.Exit(0)

		case "-r":
			args.run = true

		case "-o":
			if len(args.rest) == 0 {
				fmt.Fprintln(os.Stderr, "ERROR: Output file not provided")
				fmt.Fprintln(os.Stderr)
				usage(os.Stderr)
				os.Exit(1)
			}

			args.outputPath = args.rest[0]
			args.rest = args.rest[1:]

		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintln(os.Stderr, "ERROR: Invalid flag '"+arg+"'")
				fmt.Fprintln(os.Stderr)
				usage(os.Stderr)
				os.Exit(1)
			}

			args.inputPath = arg
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

	lexer, err := lexer.New(args.inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: Could not open file '"+args.inputPath+"'")
		fmt.Fprintln(os.Stderr)
		usage(os.Stderr)
		os.Exit(1)
	}

	parser := parser.Parser{}
	parser.File(lexer)

	context := checker.NewContext()
	for _, node := range parser.Nodes {
		context.Check(node)
	}

	if args.outputPath == "" {
		args.outputPath = strings.TrimSuffix(args.inputPath, ".yo")
	}

	compiler.Program(&context, parser.Nodes, args.outputPath)
	if args.run {
		if !strings.HasPrefix(args.outputPath, "/") {
			args.outputPath = "./" + args.outputPath
		}

		cmd := exec.Command(args.outputPath)
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
