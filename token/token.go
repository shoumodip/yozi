package token

import "fmt"

type Pos struct {
	Path string
	Row  int
	Col  int
}

func (p Pos) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Path, p.Row+1, p.Col+1)
}

type Kind = byte

const (
	Eof Kind = iota
	Int
	Bool
	Ident

	Add
	Sub
	Mul
	Div

	LBrace
	RBrace

	Print
	COUNT
)

// @TokenKind
var Names = [COUNT]string{
	Eof:   "end of file",
	Int:   "integer",
	Bool:  "boolean",
	Ident: "identifier",

	Add: "'+'",
	Sub: "'-'",
	Mul: "'*'",
	Div: "'/'",

	LBrace: "'{'",
	RBrace: "'}'",

	Print: "'print'",
}

type Token struct {
	Kind      Kind
	Pos       Pos
	Str       string
	OnNewline bool

	I64 int64
}
