package recipe

import (
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
		// TODO: Possible validation errors hides the type
		if _, err := LoadManifest(url); err == nil {
			return ManifestType
		}

		if _, err := LoadRecipe(url); err == nil {
			return LocalType
		}
	}

	return UnknownType
}
