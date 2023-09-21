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
		Use:   "check RECIPE_NAME",
		Short: "Check if there are new versions for a recipe",
		Long:  "TODO",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeName = args[0]
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
	sauce, err := re.LoadSauce(opts.Dir, opts.RecipeName)
	if err != nil {
		if errors.Is(err, re.ErrSauceNotFound) {
			cmd.PrintErrf("Error: project %s does not contain sauce with recipe %s. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, opts.RecipeName)
		} else {
			cmd.PrintErrf("Error: %s", err)
		}
		return
	}

	if sauce.CheckFrom == "" {
		cmd.PrintErr("Error: source of the sauce is undefined, can not check for new versions")
		return
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

	newTags := []string{}
	err = repo.Tags(ctx, "", func(tags []string) error {
		for _, tag := range tags {
			if semver.IsValid(tag) && semver.Compare(tag, sauce.Recipe.Version) > 0 {
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
		cmd.Println("No new versions found")
	}
}
