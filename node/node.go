package node

import "yozi/token"

type Node interface {
	Literal() token.Token
	GetType() Type
	SetType(t Type)
	IsMemory() bool
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

type Call struct {
	Token token.Token
	Type  Type

	Fn   Node
	Args []Node
}

func (c *Call) Literal() token.Token {
	return c.Token
}

func (c *Call) GetType() Type {
	return c.Type
}

func (c *Call) SetType(t Type) {
	c.Type = t
}

func (*Call) IsMemory() bool {
	return false
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

type Debug struct {
	Token token.Token
	Type  Type

	Operand Node
}

func (d *Debug) Literal() token.Token {
	return d.Token
}

func (d *Debug) GetType() Type {
	return d.Type
}

func (d *Debug) SetType(t Type) {
	d.Type = t
}

func (*Debug) IsMemory() bool {
	return false
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

func (*If) IsMemory() bool {
	return false
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

func (*While) IsMemory() bool {
	return false
}

type Return struct {
	Token token.Token
	Type  Type

	Operand Node
}

func (r *Return) Literal() token.Token {
	return r.Token
}

func (r *Return) GetType() Type {
	return r.Type
}

func (r *Return) SetType(t Type) {
	r.Type = t
}

func (*Return) IsMemory() bool {
	return false
}

type Fn struct {
	Token token.Token
	Type  Type

	Args   []*Let
	Body   *Block
	Return Node

	Locals []Node
}

func (f *Fn) Literal() token.Token {
	return f.Token
}

func (f *Fn) GetType() Type {
	return f.Type
}

func (f *Fn) SetType(t Type) {
	f.Type = t
}

func (_ *Fn) IsMemory() bool {
	return false
}

func (f *Fn) ReturnType() Type {
	if f.Return == nil {
		return Type{Kind: TypeUnit}
	}

	return f.Return.GetType()
}

type LetKind = byte

const (
	LetGlobal LetKind = iota
	LetLocal
	LetArg
	LetLocalArg
)

type Let struct {
	Token token.Token
	Kind  LetKind
	Type  Type

	// let x = <expr>        // Assign = <expr>, DefType = nil
	// let x <type>          // Assign = nil,    DefType = <type>
	// let x <type> = <expr> // Assign = <expr>, DefType = <type>
	Assign  Node
	DefType Node
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

type Block struct {
	Token token.Token
	Type  Type

	Nodes []Node
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
