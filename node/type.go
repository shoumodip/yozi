package node

import "strings"

type TypeKind = byte

const (
	TypeUnit TypeKind = iota
	TypeBool

	TypeI8
	TypeI16
	TypeI32
	TypeI64

	TypeFn
)

type Type struct {
	Kind TypeKind
	Spec Node
	Ref  int
}

// @TypeKind
func (t Type) String() string {
	sb := strings.Builder{}
	for range t.Ref {
		sb.WriteByte('&')
	}

	switch t.Kind {
	case TypeUnit:
		sb.WriteString("()")

	case TypeBool:
		sb.WriteString("bool")

	case TypeI8:
		sb.WriteString("i8")

	case TypeI16:
		sb.WriteString("i16")

	case TypeI32:
		sb.WriteString("i32")

	case TypeI64:
		sb.WriteString("i64")

	case TypeFn:
		fn := t.Spec.(*Fn)

		sb.WriteString("fn (")
		for i, arg := range fn.Args {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(arg.Type.String())
		}
		sb.WriteByte(')')

		if fn.Return != nil {
			sb.WriteByte(' ')
			sb.WriteString(fn.Return.GetType().String())
		}
	}

	return sb.String()
}

func (a Type) Equal(b Type) bool {
	if a.Kind != b.Kind || a.Ref != b.Ref {
		return false
	}

	switch a.Kind {
	case TypeFn:
		aSig := a.Spec.(*Fn)
		bSig := b.Spec.(*Fn)

		if len(aSig.Args) != len(bSig.Args) {
			return false
		}

		for i, aArg := range aSig.Args {
			bArg := bSig.Args[i]
			if !aArg.Type.Equal(bArg.Type) {
				return false
			}
		}

		return aSig.ReturnType().Equal(bSig.ReturnType())

	default:
		return true
	}
}
