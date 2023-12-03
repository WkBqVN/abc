package utils

import "errors"

var (
	ErrCannotCloseJsonfile = errors.New("can't open json file")
	ErrCannotFormatData    = errors.New("can't not format json data")
)
