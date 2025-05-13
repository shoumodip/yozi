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
	U8
	U16
	U32
	U64
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

	DebugAlloc
	DebugPrint

	COUNT
)

// @TokenKind
var Names = [COUNT]string{
	Eof: "end of file",

	I8:  "integer",
	I16: "integer",
	I32: "integer",
	I64: "integer",
	U8:  "integer",
	U16: "integer",
	U32: "integer",
	U64: "integer",
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

	DebugAlloc: "'#alloc'",
	DebugPrint: "'#print'",
}

type Token struct {
	Kind      Kind
	Pos       Pos
	Str       string
	OnNewline bool

	Int uint64
}

// @TokenKind
func (t Token) IsInteger() bool {
	switch t.Kind {
	case I8, I16, I32, I64, U8, U16, U32, U64, Int:
		return true

	default:
		return false
	}
}

func (t *Token) ParseInteger(bits int) {
	var err error
	var typeName string

	switch t.Kind {
	case I8, I16, I32, I64, Int:
		typeName = fmt.Sprintf("i%d", bits)

		var temp int64
		temp, err = strconv.ParseInt(t.Str, 10, bits)

		// Integer literals are always positive hence int64 fits within uint64
		t.Int = uint64(temp)

	case U8, U16, U32, U64:
		typeName = fmt.Sprintf("u%d", bits)
		t.Int, err = strconv.ParseUint(t.Str, 10, bits)

	default:
		panic("unreachable")
	}

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s: ERROR: Integer literal '%s' is too large for type %s\n",
			t.Pos,
			t.Str,
			typeName,
		)
		os.Exit(1)
	}
}
