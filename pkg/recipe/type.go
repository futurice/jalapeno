package recipe

import (
	"os"
	"path/filepath"
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

	if strings.HasPrefix(url, "file://") || filepath.IsAbs(url) || filepath.IsLocal(url) {
		fileinfo, _ := os.Stat(url)

		// If the file exists and it is not a directory, assume manifest
		if fileinfo != nil && !fileinfo.IsDir() {
			return ManifestType
		}

		// Else assume that the URL points to a recipe
		return LocalType
	}

	return UnknownType
}
