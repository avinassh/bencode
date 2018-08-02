package bencode

import "errors"

var ErrInvalidBenString = errors.New("Invalid Bencoded string")
var ErrSizeString = errors.New("Invalid size string")
var ErrBytesMissing = errors.New("No more bytes to read")
var ErrInvalidInteger = errors.New("Invalid integer value")
