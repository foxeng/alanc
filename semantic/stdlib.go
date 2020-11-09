package semantic

// rInt is an instance of the integer data type, necessary where a *DataType is required.
var rInt = DataTypeInt

// rByte is an instance of the byte data type, necessary where a *DataType is required.
var rByte = DataTypeByte

// stdlib is the collection of standard library functions in Alan.
var stdlib = []FuncDef{
	{
		ID: "writeInteger",
		Parameters: []ParDef{
			{
				VarDef: &PrimVarDef{
					ID:       "n",
					DataType: DataTypeInt,
				},
			},
		},
	},
	{
		ID: "writeByte",
		Parameters: []ParDef{
			{
				VarDef: &PrimVarDef{
					ID:       "b",
					DataType: DataTypeByte,
				},
			},
		},
	},
	{
		ID: "writeChar",
		Parameters: []ParDef{
			{
				VarDef: &PrimVarDef{
					ID:       "b",
					DataType: DataTypeByte,
				},
			},
		},
	},
	{
		ID: "writeString",
		Parameters: []ParDef{
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "s",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
		},
	},
	{
		ID:    "readInteger",
		RType: &rInt,
	},
	{
		ID:    "readByte",
		RType: &rByte,
	},
	{
		ID:    "readChar",
		RType: &rByte,
	},
	{
		ID: "readString",
		Parameters: []ParDef{
			{
				VarDef: &PrimVarDef{
					ID:       "n",
					DataType: DataTypeInt,
				},
			},
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "s",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
		},
	},
	{
		ID: "extend",
		Parameters: []ParDef{
			{
				VarDef: &PrimVarDef{
					ID:       "b",
					DataType: DataTypeByte,
				},
			},
		},
		RType: &rInt,
	},
	{
		ID: "shrink",
		Parameters: []ParDef{
			{
				VarDef: &PrimVarDef{
					ID:       "i",
					DataType: DataTypeInt,
				},
			},
		},
		RType: &rByte,
	},
	{
		ID: "strlen",
		Parameters: []ParDef{
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "s",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
		},
		RType: &rInt,
	},
	{
		ID: "strcmp",
		Parameters: []ParDef{
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "s1",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "s2",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
		},
		RType: &rInt,
	},
	{
		ID: "strcpy",
		Parameters: []ParDef{
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "trg",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "src",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
		},
	},
	{
		ID: "strcat",
		Parameters: []ParDef{
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "trg",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
			{
				VarDef: &ArrayDef{
					PrimVarDef: PrimVarDef{
						ID:       "src",
						DataType: DataTypeByte,
					},
				},
				IsRef: true,
			},
		},
	},
}
