package moddb

import (
	"encoding/json"
	"errors"
	"net/http"
)

const pnpURL = "https://raw.githubusercontent.com/CCDirectLink/CCModDB/master/npDatabase.json"

func getPackageDB() (PackageDB, error) {
	var result PackageDB

	resp, err := http.Get(pnpURL)
	if err != nil {
		return result, err
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

//Implements npDataBase format as per https://github.com/CCDirectLink/CCModDB/blob/master/npDatabaseFormat.md

//NodeOSPlatform identifys the operating system platform
type NodeOSPlatform string

//Possible values for NodeOSPlatform
const (
	NodeOSPlatformAix     NodeOSPlatform = "aix"
	NodeOSPlatformDarwin                 = "darwin"
	NodeOSPlatformFreeBSD                = "freebsd"
	NodeOSPlatformLinux                  = "linux"
	NodeOSPlatformOpenBSD                = "openbsd"
	NodeOSPlatformSunos                  = "sunos"
	NodeOSPlatformWindows                = "win32"
	NodeOSPlatformAndroid                = "android"
)

//Semver as used in nodejs
type Semver string

//SemverConstraint as used in nodejs
type SemverConstraint string

//PackageDBPackageType represents a kind of package.
//Please be aware that "base" is reserved for packages that absolutely require special-case local detection,
//and special-case UI to be user-friendly, such as CCLoader and NWJS upgrades.
//(In particular, CrossCode, CCLoader and NWJS upgrades require custom detection methods for their local copies.)
type PackageDBPackageType string

//Possible values for PackageDBPackageType
const (
	PackageDBPackageTypeMod  PackageDBPackageType = "mod"
	PackageDBPackageTypeTool                      = "tool"
	PackageDBPackageTypeBase                      = "base"
)

//PackageDBMetadataPage is a page relating to the mod.
type PackageDBMetadataPage struct {
	//The name of the page. For the canonical GitHub or GitLab page, this must be "GitHub" / "GitLab".
	Name string `json:"name"`
	URL  string `json:"url"`
}

//PackageDBPackageMetadata  is related to the supported package metadata for mods, on purpose.
//Note, however, that the 'dependencies' key is NOT supported.
//Also note the care to try not to reinvent NPM fields, but also to avoid them when inappropriate.
//Some mods will use NPM packages, have NPM build commands (See: TypeScript mods), etc.
//So it's very important to keep the package metadata format safe for NPM to read,
// and that means ensuring all package metadata is either avoided by NPM or understood by it.
type PackageDBPackageMetadata struct {
	Name              string                      `json:"name"`
	CCModType         PackageDBPackageType        `json:"ccmodType"`
	Version           Semver                      `json:"version"`
	CCModDependencies map[string]SemverConstraint `json:"ccmodDependencies"`
	CCModHumanName    string                      `json:"ccmodHumanName"`
	Description       string                      `json:"description"`
	Licence           string                      `json:"licence"`
	Homepage          string                      `json:"homepage"`
	Hidden            bool                        `json:"hidden"`
}

func (meta PackageDBPackageMetadata) niceName() string {
	if meta.CCModHumanName != "" {
		return meta.CCModHumanName
	}

	return meta.Name
}

//PackageDBHash represents some set of hashes for something.
type PackageDBHash struct {
	//Lowercase hexadecimal-encoded SHA-256 hash of the data.
	Sha256 string `json:"sha256"`
}

//ErrWrongInstallationMethod is returned when wrong PackageDBInstallationMethod type is parsed.
var ErrWrongInstallationMethod = errors.New("moddb: Wrong installation method")

//PackageDBInstallationMethod represents a method of installing the package.
type PackageDBInstallationMethod struct {
	PackageDBInstallationMethodCommon
	PackageDBInstallationMethodModZip
}

//common struct to all installation methods.
func (m PackageDBInstallationMethod) common() PackageDBInstallationMethodCommon {
	return m.PackageDBInstallationMethodCommon
}

//modZip struct. Returns error if type is wrong.
func (m PackageDBInstallationMethod) modZip() (PackageDBInstallationMethodModZip, error) {
	common := m.common()
	m.PackageDBInstallationMethodModZip.PackageDBInstallationMethodCommon = common
	if common.Type != PackageDBInstallationMethodTypeModZip {
		return m.PackageDBInstallationMethodModZip, ErrWrongInstallationMethod
	}
	return m.PackageDBInstallationMethodModZip, nil
}

//PackageDBInstallationMethodType represents the possible types of a PackageDBInstallationMethod
type PackageDBInstallationMethodType string

//Possible values for PackageDBInstallationMethodType
const (
	PackageDBInstallationMethodTypeModZip = "modZip"
)

//PackageDBInstallationMethodCommon contains common fields between all PackageDBInstallationMethods.
type PackageDBInstallationMethodCommon struct {
	Type     PackageDBInstallationMethodType `json:"type"`
	Platform NodeOSPlatform                  `json:"platform"`
}

//PackageDBInstallationMethodModZip contains the data required for modzip installation
type PackageDBInstallationMethodModZip struct {
	PackageDBInstallationMethodCommon

	URL    string        `json:"url"`
	Hash   PackageDBHash `json:"hash"`
	Source string        `json:"source"`
}

//PackageDBPackage represents a package in the database.
type PackageDBPackage struct {
	//Metadata for the package.
	Metadata PackageDBPackageMetadata `json:"metadata"`
	//Installation methods (try in order)
	Installation []PackageDBInstallationMethod `json:"installation"`
}

//PackageDB represents the database. Keys in this Record MUST match their respective `value.metadata.name`
type PackageDB map[string]PackageDBPackage
