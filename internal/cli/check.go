package cli

import (
	"context"
	"errors"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/oci"
	re "github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type checkOptions struct {
	RecipeName string
	option.Common
	option.WorkingDirectory
	option.OCIRepository
}

func NewCheckCmd() *cobra.Command {
	var opts checkOptions
	var cmd = &cobra.Command{
		Use:   "check",
		Short: "Check if there are new versions for a recipe",
		Long:  "TODO",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCheck(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.RecipeName, "recipe", "", "Name of the recipe to check for new versions")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCheck(cmd *cobra.Command, opts checkOptions) {
	sauces, err := re.LoadSauces(opts.Dir)
	if err != nil {
		if errors.Is(err, re.ErrSauceNotFound) {
			cmd.PrintErrf("Error: project %s does not contain sauce with recipe %s. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, opts.RecipeName)
		} else {
			cmd.PrintErrf("Error: %s", err)
		}
		return
	}

	if opts.RecipeName != "" {
		found := false
		for _, sauce := range sauces {
			if sauce.Recipe.Name == opts.RecipeName {
				found = true
				sauces = []*re.Sauce{sauce}
				break
			}
		}

		if !found {
			cmd.PrintErrf("Error: project %s does not contain sauce with recipe %s. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, opts.RecipeName)
			return
		}
	}

	cmd.Println("Checking for new versions...")

	updates := make([][]string, len(sauces))
	for i, sauce := range sauces {
		if sauce.CheckFrom == "" {
			cmd.PrintErrf("Error: source of the sauce with ID '%s' is undefined, can not check for new versions\n", sauce.ID)
			continue
		}

		ctx := context.Background()

		repo, err := oci.NewRepository(oci.Repository{
			Reference: strings.TrimPrefix(sauce.CheckFrom, "oci://"),
			PlainHTTP: opts.PlainHTTP,
			Credentials: oci.Credentials{
				Username:      opts.Username,
				Password:      opts.Password,
				DockerConfigs: opts.Configs,
			},
			TLS: oci.TLSConfig{
				CACertFilePath: opts.CACertFilePath,
				Insecure:       opts.Insecure,
			},
		})

		if err != nil {
			cmd.PrintErrf("Error: %s", err)
			return
		}

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
			cmd.PrintErrf("Error: %s", err)
			return
		}
	}

	for i, u := range updates {
		if len(u) > 0 {
			cmd.Println()
			cmd.Printf("New versions found for recipe '%s': %s\n", sauces[i].Recipe.Name, u)
			cmd.Println("Upgrade recipe with `jalapeno upgrade ...`")
		} else {
			// TODO: Use different exit code
			cmd.Printf("No new versions found for recipe '%s'\n", sauces[i].Recipe.Name)
		}
	}
}
