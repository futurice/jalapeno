package recipeutil

import (
	"context"
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/oci"
	re "github.com/futurice/jalapeno/pkg/recipe"
	"golang.org/x/mod/semver"
)

func CheckForUpdates(sauce *re.Sauce, opts option.OCIRepository) ([]string, error) {
	if sauce.CheckFrom == "" {
		return nil, fmt.Errorf("source of the sauce with ID '%s' is undefined, can not check for new versions", sauce.ID)
	}

	repo, err := oci.NewRepository(opts.Repository(sauce.CheckFrom))
	if err != nil {
		return nil, err
	}

	var versions []string
	ctx := context.Background()
	err = repo.Tags(ctx, "", func(tags []string) error {
		versions = make([]string, 0, len(tags))
		for _, tag := range tags {
			if semver.IsValid(tag) && semver.Compare(tag, sauce.Recipe.Version) > 0 {
				versions = append(versions, tag)
			}
		}
		semver.Sort(versions)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return versions, nil
}
