package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func main() {
	module := ir.NewModule()
	message := constant.NewCharArrayFromString("Hello, world!\x00")
	messagePtr := module.NewGlobalDef("", message)

	puts := module.NewFunc("puts", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	main := module.NewFunc("main", types.I32)
	zero := constant.NewInt(types.I32, 0)

	entry := main.NewBlock("")
	entry.NewCall(puts, constant.NewGetElementPtr(message.Typ, messagePtr, zero, zero))
	entry.NewRet(zero)

	asmPath := "hello.ll"
	exePath := "hello"

	f, err := os.Create(asmPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	module.WriteTo(f)
	f.Close()

	err = exec.Command("clang", "-o", exePath, asmPath).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	os.Remove(asmPath)
}
