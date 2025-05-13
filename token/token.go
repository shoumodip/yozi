package token

import (
	"fmt"
	"os"
	"strconv"
)

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

	I8
	I16
	I32
	I64
	Int // Untyped

	Bool
	Ident

	Add
	Sub
	Mul
	Div

	Shl
	Shr
	BOr
	BAnd
	BNot

	LOr
	LAnd
	LNot

	Set

	Gt
	Ge
	Lt
	Le
	Eq
	Ne

	LBrace
	RBrace
	LParen
	RParen

	Comma

	As

	If
	Else
	While
	Return

	Fn
	Let

	Print

	COUNT
)

// @TokenKind
var Names = [COUNT]string{
	Eof: "end of file",

	I8:  "integer",
	I16: "integer",
	I32: "integer",
	I64: "integer",
	Int: "integer",

	Bool:  "boolean",
	Ident: "identifier",

	Add: "'+'",
	Sub: "'-'",
	Mul: "'*'",
	Div: "'/'",

	Shl:  "'<<'",
	Shr:  "'>>'",
	BOr:  "'|'",
	BAnd: "'&'",
	BNot: "'~'",

	LOr:  "'||'",
	LAnd: "'&&'",
	LNot: "'!'",

	Set: "'='",

	Gt: "'>'",
	Ge: "'>='",
	Lt: "'<'",
	Le: "'<='",
	Eq: "'=='",
	Ne: "'!='",

	LBrace: "'{'",
	RBrace: "'}'",
	LParen: "'('",
	RParen: "')'",

	Comma: "','",

	As: "'as'",

	If:     "'if'",
	Else:   "'else'",
	While:  "'while'",
	Return: "'return'",

	Fn:  "'fn'",
	Let: "'let'",

	Print: "'print'",
}

type Token struct {
	Kind      Kind
	Pos       Pos
	Str       string
	OnNewline bool

	I64 int64
}

// @TokenKind
func (t Token) IsInteger() bool {
	switch t.Kind {
	case I8, I16, I32, I64, Int:
		return true

	default:
		return false
	}
}

func (t *Token) ParseInteger(bits int) {
	value, err := strconv.ParseInt(t.Str, 10, bits)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s: ERROR: Integer literal '%s' is too large for type i%d\n",
			t.Pos,
			t.Str,
			bits,
		)
		os.Exit(1)
	}

	t.I64 = value
}
