package cli

import (
	"context"
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/oci"
	re "github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type checkOptions struct {
	RecipeName          string
	UseDetailedExitCode bool
	option.Common
	option.WorkingDirectory
	option.OCIRepository
}

const DetailedExitCodeWhenUpdatesAvailable = 2

func NewCheckCmd() *cobra.Command {
	var opts checkOptions
	var cmd = &cobra.Command{
		Use:   "check",
		Short: "Check if there are new versions for recipes",
		Long:  "Check if there are newer versions available for recipes used in the project. By default it checks updates for all recipes, but it is possible to check updates for a specific recipe by using the `--recipe` flag.",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(cmd, opts)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			cmd.Root().SetContext(cmd.Context())
		},
		Example: `# Check updates for all recipes in the project
jalapeno check

# Check updates for a single recipe
jalapeno check --recipe my-recipe`,
	}

	cmd.Flags().StringVar(&opts.RecipeName, "recipe", "", "Name of the recipe to check for new versions")
	cmd.Flags().BoolVar(&opts.UseDetailedExitCode, "detailed-exitcode", false, fmt.Sprintf("Returns a detailed exit code when the command exits. When provided, this argument changes the exit codes and their meanings to provide more granular information about what the resulting plan contains: 0 = Succeeded with no updates available, 1 = Error, %d = Succeeded with updates available", DetailedExitCodeWhenUpdatesAvailable))

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCheck(cmd *cobra.Command, opts checkOptions) error {
	sauces, err := re.LoadSauces(opts.Dir)
	if err != nil {
		return fmt.Errorf("can not load sauces: %w", err)
	}

	if len(sauces) == 0 {
		return fmt.Errorf("working directory '%s' does not contain any sauces", opts.Dir)
	}

	// Check if we are looking updates for a specific recipe
	if opts.RecipeName != "" {
		filtered := make([]*re.Sauce, 0, len(sauces))
		for _, sauce := range sauces {
			if sauce.Recipe.Name == opts.RecipeName {
				filtered = append(filtered, sauce)
			}
		}

		if len(filtered) == 0 {
			return fmt.Errorf("project %s does not contain a sauce with recipe %s. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, opts.RecipeName)
		}

		sauces = filtered
	}

	cmd.Println("Checking for new versions...")

	updates := make([][]string, len(sauces))
	for i, sauce := range sauces {
		if sauce.CheckFrom == "" {
			return fmt.Errorf("source of the sauce with ID '%s' is undefined, can not check for new versions", sauce.ID)
		}

		repo, err := oci.NewRepository(opts.Repository(sauce.CheckFrom))

		if err != nil {
			return err
		}

		ctx := context.Background()
		err = repo.Tags(ctx, "", func(tags []string) error {
			updates[i] = make([]string, 0, len(tags))
			for _, tag := range tags {
				if semver.IsValid(tag) && semver.Compare(tag, sauce.Recipe.Version) > 0 {
					updates[i] = append(updates[i], tag)
				}
			}
			semver.Sort(updates[i])
			return nil
		})

		if err != nil {
			return err
		}
	}

	for i, u := range updates {
		if len(u) > 0 {
			cmd.Printf("New versions found for recipe '%s': %s\n", sauces[i].Recipe.Name, u)
			cmd.Println("Upgrade recipe with `jalapeno upgrade TODO`")

			if opts.UseDetailedExitCode {
				ctx := context.WithValue(cmd.Context(), ExitCodeContextKey{}, DetailedExitCodeWhenUpdatesAvailable)
				cmd.SetContext(ctx)
			}
		} else {
			cmd.Printf("No new versions found for recipe '%s'\n", sauces[i].Recipe.Name)
		}
	}

	return nil
}
