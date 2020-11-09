package semantic

const (
	// PrimitiveTypeInt is the "int" primitive type.
	PrimitiveTypeInt PrimitiveType = iota
	// PrimitiveTypeByte is the "byte" primitive type.
	PrimitiveTypeByte
	// PrimitiveTypeBool is the boolean primitive type (i.e. the type of conditions).
	PrimitiveTypeBool
)

// Type is a type in Alan.
type Type interface {
	isType()
}

// DType is a data type (i.e. a primitive or an array).
type DType interface {
	Type
	isDType()
}

// PrimitiveType is a primitive type.
type PrimitiveType int

func (PrimitiveType) isType() {}

func (PrimitiveType) isDType() {}

// ArrayType is an array type.
type ArrayType struct {
	// PrimitiveType is the array's element type.
	PrimitiveType
	// Size is the array's size.
	Size int
}

func (ArrayType) isType() {}

func (ArrayType) isDType() {}

// ParameterType is a function parameter type (i.e a data type with pass-by information).
type ParameterType struct {
	// DType is the parameter's data type (if an array, size is ignored).
	DType
	// IsRef denotes whether the parameter is passed by reference.
	IsRef bool
}

func (ParameterType) isType() {}

// FunctionType is a function type.
type FunctionType struct {
	// Parameters are the parameters' types.
	Parameters []ParameterType
	// Return is the return type (nil if function is a proc).
	Return *PrimitiveType
}

func (FunctionType) isType() {}
