package kokoro

import "errors"

var (
	ErrEnvLoadFailed        error = errors.New("failed to load config from environment")
	ErrInitializationFailed error = errors.New("failed to initialize kokoro")
)
