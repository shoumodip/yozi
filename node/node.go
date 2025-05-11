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
	IsMemory() bool

	Debug(w io.Writer, depth int) // @Temporary
}

type Atom struct {
	Token token.Token
	Type  Type

	Defined Node
	Memory  bool
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

func (a *Atom) IsMemory() bool {
	return a.Memory
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
	Memory  bool
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

func (u *Unary) IsMemory() bool {
	return u.Memory
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

	Memory bool
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

func (b *Binary) IsMemory() bool {
	return b.Memory
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

func (_ *Print) IsMemory() bool {
	return false
}

// @Temporary
func (p Print) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintln(w, "Print")
	p.Operand.Debug(w, depth+1)
}

type If struct {
	Token token.Token
	Type  Type

	Condition  Node
	Consequent Node
	Antecedent Node
}

func (i *If) Literal() token.Token {
	return i.Token
}

func (i *If) GetType() Type {
	return i.Type
}

func (i *If) SetType(t Type) {
	i.Type = t
}

func (_ *If) IsMemory() bool {
	return false
}

// @Temporary
func (i *If) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintln(w, "If")
	i.Condition.Debug(w, depth+1)
	i.Consequent.Debug(w, depth+1)
	i.Antecedent.Debug(w, depth+1)
}

type While struct {
	Token token.Token
	Type  Type

	Condition Node
	Body      Node
}

func (w *While) Literal() token.Token {
	return w.Token
}

func (w *While) GetType() Type {
	return w.Type
}

func (w *While) SetType(t Type) {
	w.Type = t
}

func (_ *While) IsMemory() bool {
	return false
}

// @Temporary
func (l *While) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintln(w, "While")
	l.Condition.Debug(w, depth+1)
	l.Body.Debug(w, depth+1)
}

type Let struct {
	Token token.Token
	Type  Type

	Assign Node
}

func (l *Let) Literal() token.Token {
	return l.Token
}

func (l *Let) GetType() Type {
	return l.Type
}

func (l *Let) SetType(t Type) {
	l.Type = t
}

func (_ *Let) IsMemory() bool {
	return false
}

// @Temporary
func (l *Let) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintf(w, "Let '%s'\n", l.Token.Str)
	l.Assign.Debug(w, depth+1)
}

type Block struct {
	Token token.Token
	Type  Type

	Body []Node
}

func (b *Block) Literal() token.Token {
	return b.Token
}

func (b *Block) GetType() Type {
	return b.Type
}

func (b *Block) SetType(t Type) {
	b.Type = t
}

func (_ *Block) IsMemory() bool {
	return false
}

// @Temporary
func (b Block) Debug(w io.Writer, depth int) {
	writeIndent(w, depth)
	fmt.Fprintln(w, "Block")
	for _, stmt := range b.Body {
		stmt.Debug(w, depth+1)
	}
}

func writeIndent(w io.Writer, depth int) {
	fmt.Fprintf(w, "%*s", depth*4, "")
}
