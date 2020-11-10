package semantic

// rInt is an instance of the int primitive type, necessary where a *PrimitiveType is required.
var rInt = PrimitiveTypeInt

// rByte is an instance of the byte primitive type, necessary where a *PrimitiveType is required.
var rByte = PrimitiveTypeByte

// stdlib is the collection of standard library functions in Alan.
var stdlib = []struct {
	ID
	FunctionType
}{
	{
		ID: "writeInteger",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: PrimitiveTypeInt,
				},
			},
		},
	},
	{
		ID: "writeByte",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: PrimitiveTypeByte,
				},
			},
		},
	},
	{
		ID: "writeChar",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: PrimitiveTypeByte,
				},
			},
		},
	},
	{
		ID: "writeString",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
			},
		},
	},
	{
		ID: "readInteger",
		FunctionType: FunctionType{
			Parameters: []ParameterType{},
			Return:     &rInt,
		},
	},
	{
		ID: "readByte",
		FunctionType: FunctionType{
			Parameters: []ParameterType{},
			Return:     &rByte,
		},
	},
	{
		ID: "readChar",
		FunctionType: FunctionType{
			Parameters: []ParameterType{},
			Return:     &rByte,
		},
	},
	{
		ID: "readString",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: PrimitiveTypeInt,
				},
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
			},
		},
	},
	{
		ID: "extend",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: PrimitiveTypeByte,
				},
			},
			Return: &rInt,
		},
	},
	{
		ID: "shrink",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: PrimitiveTypeInt,
				},
			},
			Return: &rByte,
		},
	},
	{
		ID: "strlen",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
			},
			Return: &rInt,
		},
	},
	{
		ID: "strcmp",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
			},
			Return: &rInt,
		},
	},
	{
		ID: "strcpy",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
			},
		},
	},
	{
		ID: "strcat",
		FunctionType: FunctionType{
			Parameters: []ParameterType{
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
				{
					DType: ArrayType{
						PrimitiveType: PrimitiveTypeByte,
					},
					IsRef: true,
				},
			},
		},
	},
}
