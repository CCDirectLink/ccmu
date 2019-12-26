package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//ErrorReason that is included in Error.
type ErrorReason int

//Reasons for an installation, uninstallation or update to fail
const (
	ReasonUnknown ErrorReason = iota
	ReasonNoInternet
	ReasonAlreadyInstalled
	ReasonInvalidFormat
	ReasonNotFound
	ReasonAccess
)

//Mode of the installation.
type Mode int

//Modes of operations
const (
	ModeUnknown Mode = iota
	ModeInstall
	ModeUninstall
	ModeUpdate
)

//Error contains details about the errors that occured while installing or uninstalling a package.
type Error struct {
	Reason ErrorReason
	Mode   Mode
	Pkg    Package
	Err    error
}

//NewError with given mode and package.
func NewError(mode Mode, pkg Package, err error) Error {
	var reason = ReasonUnknown

	if _, ok := err.(*http.ProtocolError); ok {
		reason = ReasonNoInternet
	} else if _, ok := err.(*json.SyntaxError); ok {
		reason = ReasonInvalidFormat
	} else if _, ok := err.(*os.PathError); ok {
		reason = ReasonAccess
	}

	return Error{reason, mode, pkg, err}
}

//NewErrorReason with given reason, mode and package.
func NewErrorReason(reason ErrorReason, mode Mode, pkg Package, err error) Error {
	return Error{reason, mode, pkg, err}
}

func (p Error) Error() string {
	return "pkg: " + p.String()
}

//Unwrap the error underneath.
func (p Error) Unwrap() error {
	return p.Err
}

func (p Error) String() string {
	info, _ := p.Pkg.Info()
	prefix := fmt.Sprintf("Could not %s %s because ", p.Mode, info.NiceName)

	switch p.Reason {
	case ReasonNotFound:
		return prefix + "it was not found"
	case ReasonAlreadyInstalled:
		return prefix + "it was already installed"
	case ReasonAccess:
		return prefix + "the access was denied"
	case ReasonUnknown:
		fallthrough
	default:
		return prefix + fmt.Sprintf("an unknown error occured in %s", p.Err)
	}
}

func (m Mode) String() string {
	switch m {
	case ModeInstall:
		return "install"
	case ModeUninstall:
		return "uninstall"
	case ModeUpdate:
		return "update"
	default:
		return "manipulate"
	}
}
