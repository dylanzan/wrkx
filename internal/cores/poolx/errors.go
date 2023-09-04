package poolx

import "errors"

var (
	ErrPoolClosed  = errors.New("pool closed")
	ErrSizeToSmall = errors.New("size too small")
)
