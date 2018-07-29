package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// pass either of one. encoded takes the precedence
func NewBenString(encoded string) (*BenStruct, error) {
	logger := log.WithField("method", "NewBenString")
	if encoded == "" {
		logger.Error("received empty encoded string")
		return nil, ErrInvalidBenString
	}

	logger = logger.WithField("encoded", encoded)

	values := strings.Split(encoded, ":")

	if len(values) == 1 {
		if values[0] == "0" {
			return &BenStruct{
				Raw:         encoded, // which is "0:"
				StringValue: "",
			}, nil
		} else {
			logger.Error("invalid bencoded string")
			return nil, ErrInvalidBenString
		}
	}

	// values is array like ["4", "spam"]
	// string length encoded in base ten ASCII>
	stringLength, err := strconv.ParseUint(values[0], 10, 0)
	if err != nil {
		logger.Error("failed to parse the size from bencoded string")
		return nil, err
	}

	if stringLength != uint64(len(values[1])) {
		logger.Error("invalid length in the bencoded string")
		return nil, ErrInvalidBenString
	}

	return &BenStruct{Raw: encoded, StringValue: values[1]}, nil
}

func NewBenStringFromValue(decoded string) (*BenStruct, error) {
	return &BenStruct{
		Raw: fmt.Sprintf("%d:%s", len(decoded), decoded), StringValue: decoded}, nil
}
