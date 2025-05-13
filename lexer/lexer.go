package lexer

import (
	"fmt"
	"os"
	"slices"
	"yozi/token"
)

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlpha(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isIdent(ch byte) bool {
	return isAlnum(ch) || ch == '_'
}

func isPrint(ch byte) bool {
	return ch >= 32 && ch <= 126
}

type Lexer struct {
	pos   token.Pos
	bytes []byte

	ch   byte
	head int
	size int

	peeked bool
	buffer token.Token

	onNewline bool
}

func New(path string) (Lexer, error) {
	l := Lexer{}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return l, err
	}

	l.pos.Path = path
	l.bytes = bytes
	l.size = len(bytes)

	if l.size != 0 {
		l.ch = l.bytes[0]
	}

	return l, nil
}

func (l *Lexer) nextChar() {
	if l.ch == '\n' {
		if l.head+1 < l.size {
			l.pos.Row++
			l.pos.Col = 0
		}
	} else {
		l.pos.Col++
	}

	l.head++
	if l.head < l.size {
		l.ch = l.bytes[l.head]
	} else {
		l.ch = 0
	}
}

func (l *Lexer) peekChar(offset int) byte {
	if l.head+offset < l.size {
		return l.bytes[l.head+offset]
	}

	return 0
}

func (l *Lexer) readChar() byte {
	ch := l.ch
	l.nextChar()
	return ch
}

func (l *Lexer) matchChar(ch byte) bool {
	if l.ch == ch {
		l.nextChar()
		return true
	}
	return false
}

func (l *Lexer) skipWhitespace() {
	for l.head < l.size {
		switch l.ch {
		case ' ', '\t', '\r':
			l.nextChar()

		case '\n':
			l.nextChar()
			l.onNewline = true

		case '/':
			if l.peekChar(1) == '/' {
				for l.head < l.size && l.ch != '\n' {
					l.nextChar()
				}
			} else {
				return
			}

		default:
			return
		}
	}
}

func (l *Lexer) Buffer(tok token.Token) {
	l.peeked = true
	l.buffer = tok
}

func (l *Lexer) Unbuffer() {
	l.peeked = false
}

// @TokenKind
func (l *Lexer) Next() token.Token {
	if l.peeked {
		l.Unbuffer()
		return l.buffer
	}

	l.skipWhitespace()

	head := l.head
	tok := token.Token{
		Pos:       l.pos,
		OnNewline: l.onNewline,
	}
	l.onNewline = false

	if l.head >= l.size {
		tok.Kind = token.Eof
		return tok
	}

	if isDigit(l.ch) {
		for l.head < l.size && isDigit(l.ch) {
			l.nextChar()
		}

		bits := 64
		numStr := string(l.bytes[head:l.head])

		// TODO: unsigned integer literal
		if l.ch == 'i' {
			suffixPos := l.pos
			suffixStart := l.head

			l.nextChar()
			for l.head < l.size && isDigit(l.ch) {
				l.nextChar()
			}

			switch suffix := string(l.bytes[suffixStart:l.head]); suffix {
			case "i8":
				tok.Kind = token.I8
				bits = 8

			case "i16":
				tok.Kind = token.I16
				bits = 16

			case "i32":
				tok.Kind = token.I32
				bits = 32

			case "i64":
				tok.Kind = token.I64
				bits = 64

			default:
				fmt.Fprintf(
					os.Stderr,
					"%s: ERROR: Invalid suffix '%s' to integer literal\n",
					suffixPos,
					suffix,
				)
				os.Exit(1)
			}
		} else {
			tok.Kind = token.Int
		}

		tok.Str = numStr
		tok.ParseInteger(bits)
		return tok
	}

	if isIdent(l.ch) {
		for l.head < l.size && isIdent(l.ch) {
			l.nextChar()
		}

		tok.Str = string(l.bytes[head:l.head])
		switch tok.Str {
		case "true":
			tok.Kind = token.Bool
			tok.I64 = 1

		case "false":
			tok.Kind = token.Bool
			tok.I64 = 0

		case "as":
			tok.Kind = token.As

		case "print":
			tok.Kind = token.Print

		case "if":
			tok.Kind = token.If

		case "else":
			tok.Kind = token.Else

		case "while":
			tok.Kind = token.While

		case "return":
			tok.Kind = token.Return

		case "fn":
			tok.Kind = token.Fn

		case "let":
			tok.Kind = token.Let

		default:
			tok.Kind = token.Ident
		}

		return tok
	}

	switch ch := l.readChar(); ch {
	case '+':
		tok.Kind = token.Add

	case '-':
		tok.Kind = token.Sub

	case '*':
		tok.Kind = token.Mul

	case '/':
		tok.Kind = token.Div

	case '|':
		if l.matchChar('|') {
			tok.Kind = token.LOr
		} else {
			tok.Kind = token.BOr
		}

	case '&':
		if l.matchChar('&') {
			tok.Kind = token.LAnd
		} else {
			tok.Kind = token.BAnd
		}

	case '~':
		tok.Kind = token.BNot

	case '<':
		if l.matchChar('<') {
			tok.Kind = token.Shl
		} else if l.matchChar('=') {
			tok.Kind = token.Le
		} else {
			tok.Kind = token.Lt
		}

	case '>':
		if l.matchChar('>') {
			tok.Kind = token.Shr
		} else if l.matchChar('=') {
			tok.Kind = token.Ge
		} else {
			tok.Kind = token.Gt
		}

	case '=':
		if l.matchChar('=') {
			tok.Kind = token.Eq
		} else {
			tok.Kind = token.Set
		}

	case '!':
		if l.matchChar('=') {
			tok.Kind = token.Ne
		} else {
			tok.Kind = token.LNot
		}

	case '{':
		tok.Kind = token.LBrace

	case '}':
		tok.Kind = token.RBrace

	case '(':
		tok.Kind = token.LParen

	case ')':
		tok.Kind = token.RParen

	case ',':
		tok.Kind = token.Comma

	default:
		message := "%s: ERROR: Invalid character '%c'\n"
		if !isPrint(ch) {
			message = "%s: ERROR: Invalid character %d\n"
		}

		fmt.Fprintf(os.Stderr, message, tok.Pos, ch)
		os.Exit(1)
	}

	tok.Str = string(l.bytes[head:l.head])
	return tok
}

func (l *Lexer) Peek() token.Token {
	if !l.peeked {
		l.Buffer(l.Next())
	}
	return l.buffer
}

func (l *Lexer) Read(kind token.Kind) bool {
	l.Peek()
	l.peeked = l.buffer.Kind != kind
	return !l.peeked
}

func (l *Lexer) Expect(kinds ...token.Kind) token.Token {
	tok := l.Next()
	if slices.Contains(kinds, tok.Kind) {
		return tok
	}

	fmt.Fprintf(os.Stderr, "%s: ERROR: Expected ", tok.Pos)
	for i, kind := range kinds {
		if i > 0 {
			if i == len(kinds)-1 {
				fmt.Fprint(os.Stderr, " or ")
			} else {
				fmt.Fprint(os.Stderr, ", ")
			}
		}

		fmt.Fprint(os.Stderr, token.Names[kind])
	}
	fmt.Fprintln(os.Stderr, ", got", token.Names[tok.Kind])
	os.Exit(1)

	panic("unreachable")
}

// In expressions, && is a single token. In a type, however it is 2 tokens.
// Same goes for >> (for generic types)
//
// Since this lexer is LL(1) a simple workaround is to "split" the token and
// buffer the right part. The left part gets returned
func (l *Lexer) SplitAndBufferRight(tok token.Token) token.Token {
	switch tok.Kind {
	case token.LAnd:
		tok.Kind = token.BAnd
		tok.Str = "&"
		lhs := tok

		tok.Pos.Col++
		l.Buffer(tok)
		return lhs

	case token.Shr:
		tok.Kind = token.Gt
		tok.Str = ">"
		lhs := tok

		tok.Pos.Col++
		l.Buffer(tok)
		return lhs

	default:
		panic("unreachable")
	}
}
