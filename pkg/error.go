package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Masterminds/semver"
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
	ReasonNotAvailable
	ReasonAccess
	ReasonDependant
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

//ErrNotFound is returned when the given package is not found.
var ErrNotFound = errors.New("pkg: Mod not found")

//NewError with given mode and package.
func NewError(mode Mode, pkg Package, err error) Error {
	var reason = ReasonUnknown

	var pkgError Error
	if errors.As(err, &pkgError) && pkgError.Pkg != pkg {
		pkgError.Mode = mode
		return pkgError
	} else if _, ok := err.(*http.ProtocolError); ok {
		reason = ReasonNoInternet
	} else if _, ok := err.(*json.SyntaxError); ok {
		reason = ReasonInvalidFormat
	} else if _, ok := err.(*os.PathError); ok {
		reason = ReasonAccess
	} else if errors.Is(err, semver.ErrInvalidMetadata) || errors.Is(err, semver.ErrInvalidPrerelease) || errors.Is(err, semver.ErrInvalidSemVer) {
		reason = ReasonInvalidFormat
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
	var prefix string
	if p.Pkg != nil {
		info, _ := p.Pkg.Info()
		prefix = fmt.Sprintf("Could not %s %s because ", p.Mode, info.NiceName)
	} else {
		prefix = fmt.Sprintf("Could not %s mod because ", p.Mode)
	}

	switch p.Reason {
	case ReasonNotFound:
		return prefix + "it was not found"
	case ReasonNotAvailable:
		return prefix + "it was not available"
	case ReasonAlreadyInstalled:
		return prefix + "it was already installed"
	case ReasonAccess:
		return prefix + "the access was denied"
	case ReasonInvalidFormat:
		return prefix + "the mod was malformated"
	case ReasonDependant:
		return prefix + "there is a mod that depends on it"
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
