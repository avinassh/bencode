package main

type DataType int

const (
	StringType DataType = iota
	IntType
	ListType
	MapType
)

type BenStruct struct {
	DataType    DataType
	IntValue    int
	MapValue    map[string]BenStruct
	ListValue   []BenStruct
	StringValue string
	Raw         string
}
