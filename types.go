package bencode

import "encoding/json"

type DataType int

const (
	StringType DataType = iota
	IntType
	ListType
	MapType
)

func (d DataType) String() string {
	switch d {
	case StringType:
		return "string"
	case IntType:
		return "int"
	case ListType:
		return "list"
	case MapType:
		return "map"
	default:
		return ""
	}
}

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
