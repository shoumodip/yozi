package node

import "strings"

type TypeKind = byte

const (
	TypeNil TypeKind = iota
	TypeBool
	TypeI64
	TypeFn
)

type Type struct {
	Kind TypeKind
	Ref  int
}

// @TypeKind
func (t Type) String() string {
	sb := strings.Builder{}
	for range t.Ref {
		sb.WriteByte('&')
	}

	switch t.Kind {
	case TypeNil:
		sb.WriteString("nil")

	case TypeBool:
		sb.WriteString("bool")

	case TypeI64:
		sb.WriteString("i64")

	case TypeFn:
		sb.WriteString("fn ()")
	}

	return sb.String()
}

func (a Type) Equal(b Type) bool {
	return a.Kind == b.Kind && a.Ref == b.Ref
}
