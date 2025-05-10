package node

import (
	"fmt"
	"io"
	"yozi/token"
)

type Node interface {
	Literal() token.Token
	GetType() Type
	SetType(t Type)

	Debug(w io.Writer, depth int) // @Temporary
}

type Atom struct {
	Token token.Token
	Type  Type
}

func (a *Atom) Literal() token.Token {
	return a.Token
}

func (a *Atom) GetType() Type {
	return a.Type
}

func (a *Atom) SetType(t Type) {
	a.Type = t
}

// @Temporary
func (a Atom) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintf(w, "Atom '%s'\n", a.Token.Str)
}

type Unary struct {
	Token token.Token
	Type  Type

	Operand Node
}

func (u *Unary) Literal() token.Token {
	return u.Token
}

func (u *Unary) GetType() Type {
	return u.Type
}

func (u *Unary) SetType(t Type) {
	u.Type = t
}

// @Temporary
func (u Unary) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintln(w, "Unary")
	u.Operand.Debug(w, depth+1)
}

type Binary struct {
	Token token.Token
	Type  Type

	Lhs Node
	Rhs Node
}

func (b *Binary) Literal() token.Token {
	return b.Token
}

func (b *Binary) GetType() Type {
	return b.Type
}

func (b *Binary) SetType(t Type) {
	b.Type = t
}

// @Temporary
func (b Binary) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintln(w, "Binary")
	b.Lhs.Debug(w, depth+1)
	b.Rhs.Debug(w, depth+1)
}

type Print struct {
	Token token.Token
	Type  Type

	Operand Node
}

func (p *Print) Literal() token.Token {
	return p.Token
}

func (p *Print) GetType() Type {
	return p.Type
}

func (p *Print) SetType(t Type) {
	p.Type = t
}

// @Temporary
func (p Print) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintln(w, "Print")
	p.Operand.Debug(w, depth+1)
}

func writeIndent(w io.Writer, depth int) {
	fmt.Fprintf(w, "%*s", depth*4, "")
}
