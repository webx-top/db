package factory

import (
	"net"
)

func NetError(err error) net.Error {
	if netErr, ok := err.(net.Error); ok {
		return netErr
	}
	return nil
}

func IsTimeoutError(err error) bool {
	if netErr := NetError(err); netErr != nil {
		return netErr.Timeout()
	}
	return false
}

func IsTemporaryError(err error) bool {
	if netErr := NetError(err); netErr != nil {
		return netErr.Temporary()
	}
	return false
}
