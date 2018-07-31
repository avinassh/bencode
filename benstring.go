package bencode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var RegExpr = regexp.MustCompile("^(\\d+):.*")

func isValidEndDelimeter(lastChar string) bool {
	if lastChar == "e" {
		return true
	}
	return false
}

func isBenString(encoded string) bool {
	return RegExpr.MatchString(encoded)
}

func getLength(encoded string) (uint64, error) {
	match := RegExpr.FindStringSubmatch(encoded)
	if len(match) == 0 || len(match) != 2 {
		return 0, ErrInvalidBenString
	}
	stringLength, err := strconv.ParseUint(match[1], 10, 0)
	if err != nil {
		return 0, err
	}
	return stringLength, nil

}

func extract(encoded string) (*BenStruct, error) {
	response := &BenStruct{}

	firstChar := string(encoded[0])
	lastChar := string(encoded[len(encoded)-1])
	buffer := string(encoded[1 : len(encoded)-1])

	if isBenString(encoded) {
		length, err := getLength(encoded)
		if err != nil {
			return nil, err
		}

		if length == 0 {
			return &BenStruct{
				StringValue: "",
				Raw:         "0:",
			}, nil
		}
		return &BenStruct{
			StringValue: encoded[:length],
			Raw:         fmt.Sprintf("%d:%s", length, encoded[:length]),
		}, nil
	}

	if !isValidEndDelimeter(lastChar) {
		return nil, ErrInvalidBenString
	}

	switch firstChar {
	case "d":
		return extract(buffer)
	case "i":
		return extract(buffer)
	case "l":
		return extract(buffer)
	}

	return response, nil
}

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
