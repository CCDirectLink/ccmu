// +build !windows

package registry

import "errors"

const Supported = false

func RegisterProtocol(ccmu, game string) error {
	return errors.New("registry: Not supported")
}

func UnregisterProtocol() error {
	return errors.New("registry: Not supported")
}

func ProtocolInstalled() string {
	return ""
}
