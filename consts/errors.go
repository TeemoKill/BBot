package consts

import "errors"

var (
	ConfigPathNotFileErr      = errors.New("config path is not a file")
	ExampleConfigGeneratedErr = errors.New("example config generated")
)
