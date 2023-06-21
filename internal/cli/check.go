package cli

import (
	"context"
	"errors"

	"github.com/futurice/jalapeno/internal/cli/internal/option"
	re "github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type checkOptions struct {
	ProjectPath string
	RecipeName  string
	option.Common
	option.OCIRepository
}

func newCheckCmd() *cobra.Command {
	var opts checkOptions
	var cmd = &cobra.Command{
		Use:   "check PROJECT_PATH RECIPE_NAME",
		Short: "Check if there are new versions for a recipe",
		Long:  "", // TODO
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectPath = args[0]
			opts.RecipeName = args[1]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCheck(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCheck(cmd *cobra.Command, opts checkOptions) {
	sauce, err := re.LoadSauce(opts.ProjectPath, opts.RecipeName)
	if err != nil {
		if errors.Is(err, re.ErrSauceNotFound) {
			cmd.PrintErrf("Error: project %s does not contain sauce with recipe %s. Recipe name used in the project should match the recipe which is used for upgrading", opts.ProjectPath, opts.RecipeName)
		} else {
			cmd.PrintErrf("Error: %s", err)
		}
		return
	}

	recipe := sauce.Recipe

	if len(recipe.Sources) == 0 {
		cmd.PrintErr("Error: source of the recipe is undefined, can not check for new versions")
		return
	}

	ctx := context.Background()

	// TODO: How to handle multiple sources?

	repo, err := opts.NewRepository(recipe.Sources[0], opts.Common)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	newTags := []string{}
	err = repo.Tags(ctx, "", func(tags []string) error {
		for _, tag := range tags {
			if semver.IsValid(tag) && semver.Compare(tag, recipe.Version) > 0 {
				newTags = append(newTags, tag)
			}
		}
		semver.Sort(newTags)
		return nil
	})

	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	if len(newTags) > 0 {
		cmd.Println("New versions found:")
		cmd.Println(newTags)
	} else {
		// TODO: Use different exit code
		cmd.PrintErrln("No new versions found")
	}
}
