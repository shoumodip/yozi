package node

type TypeKind = byte

const (
	TypeNil TypeKind = iota
	TypeBool
	TypeI64
)

type Type struct {
	Kind TypeKind
}

// @TypeKind
func (dt Type) String() string {
	switch dt.Kind {
	case TypeNil:
		return "nil"

	case TypeBool:
		return "bool"

	case TypeI64:
		return "i64"
	}

	panic("unreachable")
}

func (a Type) Equal(b Type) bool {
	return a.Kind == b.Kind
}
