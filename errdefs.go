package kokoro

import (
	"fmt"
)

type externalErr string

var (
	ErrEnvLoadFailed        externalErr = "failed to load config from environment"
	ErrInitializationFailed externalErr = "failed to initialize kokoro"
)

func wrapErr(root externalErr, err error) error {
	root = root + ": %w"

	return fmt.Errorf(string(root), err)
}
