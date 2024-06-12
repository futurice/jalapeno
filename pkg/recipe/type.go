package recipe

import (
	"os"
	"strings"
)

type URLType int

const (
	UnknownType URLType = iota
	LocalType
	ManifestType
	OCIType
)

func DetermineRecipeURLType(url string) URLType {
	if strings.HasPrefix(url, "oci://") {
		return OCIType
	}

	fileinfo, _ := os.Stat(url)

	// If the file exists and it is not a directory, assume manifest
	if fileinfo != nil && !fileinfo.IsDir() {
		return ManifestType
	}

	// Else assume that the URL points to a local recipe
	return LocalType
}
