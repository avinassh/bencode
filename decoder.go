package main

import (
	"bytes"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func NewBencoder(encoded string) *Bencoder {
	return &Bencoder{
		raw:       []byte(encoded),
		rawString: encoded,
		cursor:    0,
	}
}

func (b *Bencoder) Parse() *BenStruct {
	return b.encode()
}

func (b *Bencoder) encode() *BenStruct {

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
		return b.extractDict()
	default:
		return nil
	}
	return nil
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

func (b *Bencoder) extractString() *BenStruct {
	return nil
}

// we extract the integer value
// they start with `i`, where the current cursor is now
// at and they end with `e`.
//
// e.g. "i3e" is 3, "i3e" is -3, "i0e" is zero
//
// and "i03e", "i-0e" are invalid
func (b *Bencoder) extractInt() *BenStruct {
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
	if len(value) > 1 && string(value[0]) == "0" {
		logger.Error("integer cannot start with 0")
		return nil
	}
	// case of `-0` or `-0<anything>`
	if len(value) > 1 && (string(value[:2]) == "-0") {
		logger.Error("integer cannot have -0")
		return nil
	}

	// we have read till `e`, so move to next cursor
	b.increment()

	intValue, err := strconv.Atoi(valueString)

	if err != nil {
		logger.WithError(err).Error("failed to parse the int")
		return nil
	}
	return &BenStruct{DataType: IntType, IntValue: intValue, Raw: fmt.Sprintf("i%de", intValue)}
}

func (b *Bencoder) extractList() *BenStruct {
	return nil
}

func (b *Bencoder) extractDict() *BenStruct {
	return nil
}
