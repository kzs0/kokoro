package errdefs

import (
	"fmt"
)

type externalErr string

var (
	ErrEnvLoadFailed        externalErr = "failed to load config from environment"
	ErrInitializationFailed externalErr = "failed to initialize kokoro"
)

func WrapErr(root externalErr, err error) error {
	return fmt.Errorf(string(root+": %w"), err)
}
