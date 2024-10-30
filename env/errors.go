package env

import (
	"errors"
)

var (
	ErrParseValue           = errors.New("unable to parse")
	ErrStructPtr            = errors.New("not struct pointer error")
	ErrNoSupportedTagOption = errors.New("tag option not supported")
	ErrVarIsNotSet          = errors.New("required environment variable is not set")
	ErrEmptyVar             = errors.New("environment variable should not be empty")
	ErrLoadFileContent      = errors.New("could not load content of file from variable")
	ErrNoParser             = errors.New("no parser found")
)
