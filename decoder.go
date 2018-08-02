package bencode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func NewBencoder(encoded string) *Bencoder {
	return &Bencoder{
		raw:       []byte(encoded),
		rawString: encoded,
		cursor:    0,
	}
}

func (b *Bencoder) Parse() (*BenStruct, error) {
	return b.encode()
}

func (b *Bencoder) encode() (*BenStruct, error) {

	currentChar := b.currentChar()

	if isDigit(currentChar) {
		return b.extractString()
	}

	switch currentChar {
	case "i":
		return b.extractInt()
	case "l":
		return b.extractList()
	case "d":
		return b.extractMap()
	default:
		return nil, ErrInvalidBenString
	}
	return nil, ErrInvalidBenString
}

func (b *Bencoder) currentChar() string {
	return string(b.raw[b.cursor])
}

func (b *Bencoder) currentByte() byte {
	return b.raw[b.cursor]
}

func (b *Bencoder) increment() {
	b.cursor += 1
}

func (b *Bencoder) incrementBy(offset int) {
	b.cursor += offset
}

// all strings are of the format `size:<string here>`
// where `size` is a base 10 size of string value
// there is no start or end delimeter
// and we don't know how much bytes `size` itself may
// take. it could be `23` or `56565575756`. So we gotta
// read till we encounter `:`. Using a regex would be easy.
func (b *Bencoder) extractString() (*BenStruct, error) {
	// current character is some digit, so we can just start
	// reading the bytes
	var buf bytes.Buffer

	// or may be use regex 🤔
	for {
		currentChar := b.currentChar()
		if currentChar == ":" {
			break
		}
		buf.WriteString(currentChar)
		b.increment()
	}

	sizeString := buf.String()
	logger := log.WithFields(log.Fields{"method": "extractString", "rawSize": sizeString})

	if strings.HasPrefix(sizeString, "-") {
		logger.Error("Size cannot be -ve")
		return nil, ErrSizeString
	}

	size, err := strconv.Atoi(sizeString)
	if err != nil {
		logger.WithError(err).Error("failed to parse the size int")
		return nil, ErrSizeString
	}

	// currently we are at `:`. So lets move the current cursor to next
	b.increment()

	// if size is 0, we just move on
	if size == 0 {
		return &BenStruct{Raw: "0:"}, nil

	}

	// now we need to read the next `size` bytes. But before that
	// we need to check, does it even have those bytes?
	if len(b.raw[b.cursor:]) < size {
		logger.Error("not enough bytes to read")
		return nil, ErrBytesMissing
	}

	// lets read the next `size` bytes
	valueBytes := b.raw[b.cursor : b.cursor+size]
	value := string(valueBytes)

	// and move the cursor by size
	b.incrementBy(size)
	return &BenStruct{StringValue: value, JsonValue: jsonMustMarshal(&value), Raw: fmt.Sprintf("%d:%s", size, value)}, nil
}

// we extract the integer value
// they start with `i`, where the current cursor is now
// at and they end with `e`.
//
// e.g. "i3e" is 3, "i3e" is -3, "i0e" is zero
//
// and "i03e", "i-0e" are invalid
func (b *Bencoder) extractInt() (*BenStruct, error) {
	// current cursor is at `i`, so read till e
	// we move one character and start reading

	var buf bytes.Buffer
	b.increment()

	// or may be use regex 🤔
	for {
		currentChar := b.currentChar()
		if currentChar == "e" {
			break
		}
		buf.WriteString(currentChar)
		b.increment()
	}

	logger := log.WithFields(log.Fields{"method": "extractInt", "raw": buf.String()})
	value := buf.Bytes()
	valueString := buf.String()

	// lets validate
	// case of `0<any number>`
	if len(value) > 1 && strings.HasPrefix(valueString, "0") {
		logger.Error("integer cannot start with 0")
		return nil, ErrInvalidInteger
	}
	// case of `-0` or `-0<anything>`
	if len(value) > 1 && strings.HasPrefix(valueString, "-0") {
		logger.Error("integer cannot have -0")
		return nil, ErrInvalidInteger
	}

	// we have read till `e`, so move to next cursor
	b.increment()

	intValue, err := strconv.Atoi(valueString)

	if err != nil {
		logger.WithError(err).Error("failed to parse the int")
		return nil, ErrInvalidInteger
	}
	return &BenStruct{DataType: IntType, IntValue: intValue, JsonValue: jsonMustMarshal(&intValue), Raw: fmt.Sprintf("i%de", intValue)}, nil
}

func (b *Bencoder) extractList() (*BenStruct, error) {
	// current cursor is at `l`, so read till e
	// we move one character and start reading

	// lets keep track of start and end cursors so that
	// we can build the raw string easily
	startCursor := b.cursor

	// lets move the cursor by 1

	b.increment()
	result := []BenStruct{}

	for {
		currentChar := b.currentChar()
		if currentChar == "e" {
			break
		}

		item, err := b.encode()
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}

	// currently we are at `e`, so lets move ahead
	b.increment()
	endCursor := b.cursor

	// lets construct jsonvalue too
	jsonResults := []json.RawMessage{}
	for _, bs := range result {
		jsonResults = append(jsonResults, bs.JsonValue)
	}
	jsonValue, err := json.Marshal(&jsonResults)
	if err != nil {
		log.WithField("method", "extractList").WithError(err).Error("failed to marshal the jsonlist")
	}

	return &BenStruct{DataType: ListType, ListValue: result, JsonValue: jsonValue, Raw: b.rawString[startCursor:endCursor]}, nil
}

func (b *Bencoder) extractMap() (*BenStruct, error) {
	// current cursor is at `d`, so read till e
	// we move one character and start reading

	// lets keep track of start and end cursors so that
	// we can build the raw string easily
	startCursor := b.cursor

	// lets move the cursor by 1
	b.increment()

	result := map[string]BenStruct{}

	// following is for json
	jsonResult := map[string]json.RawMessage{}

	for {
		currentChar := b.currentChar()
		if currentChar == "e" {
			break
		}

		key, err := b.extractString()
		if err != nil {
			return nil, err
		}
		value, err := b.encode()
		if err != nil {
			return nil, err
		}
		result[key.StringValue] = *value
		jsonResult[key.StringValue] = value.JsonValue
	}

	// currently we are at `e`, so lets move ahead
	b.increment()
	endCursor := b.cursor

	jsonValue, err := json.Marshal(&jsonResult)
	if err != nil {
		log.WithField("method", "extractMap").WithError(err).Error("failed to marshal the jsonmap")
	}

	return &BenStruct{DataType: MapType, MapValue: result, JsonValue: jsonValue, Raw: b.rawString[startCursor:endCursor]}, nil
}
