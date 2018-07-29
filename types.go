package main

import "encoding/json"

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
	JsonValue   json.RawMessage
}

type Bencoder struct {
	raw       []byte
	rawString string
	cursor    int // gives the current position
}
