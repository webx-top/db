package factory

import (
	"net"
	"unicode"
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

// ToSnakeCase : WebxTop => webx_top
func ToSnakeCase(name string) string {
	bytes := []rune{}
	for i, char := range name {
		if 'A' <= char && 'Z' >= char {
			char = unicode.ToLower(char)
			if i > 0 {
				bytes = append(bytes, '_')
			}
		}
		bytes = append(bytes, char)
	}
	return string(bytes)
}
