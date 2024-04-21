package utils

import (
	"errors"
	"io"
	"net"
)

func IsNetClosedErr(err error) bool {
	var netErr *net.OpError
	return errors.As(err, &netErr) && errors.Is(netErr.Err, net.ErrClosed)
}

func IsEof(err error) bool {
	return errors.Is(err, io.EOF)
}
