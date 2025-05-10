package node

type TypeKind = byte

const (
	TypeBool TypeKind = iota
	TypeI64
)

type Type struct {
	Kind TypeKind
}

// @TypeKind
func (dt Type) String() string {
	switch dt.Kind {
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
