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
	Eof:   "end of file",
	Int:   "integer",
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
