package shield

import "errors"

var (
	ErrInvalidSite = errors.New("invalid site name")
	ErrEmptyConfig = errors.New("shield config is empty")
)
