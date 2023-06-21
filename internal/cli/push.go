package cli

import (
	"context"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
)

type pushOptions struct {
	RecipePath string
	TargetRef  string
	option.OCIRepository
	option.Common
}

func NewPushCmd() *cobra.Command {
	var opts pushOptions
	var cmd = &cobra.Command{
		Use:   "push RECIPE_PATH TARGET_URL",
		Short: "Push a recipe to OCI repository",
		Long:  "TODO",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			opts.TargetRef = args[1]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runPush(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runPush(cmd *cobra.Command, opts pushOptions) {
	ctx := context.Background()

	re, err := recipe.LoadRecipe(opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("Error: can't load the recipe: %s\n", err)
		return
	}

	store, err := file.New("")
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	defer store.Close()

	desc, err := store.Add(ctx, re.Name, "application/x.futurice.jalapeno.recipe.v1", opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	root, err := oras.Pack(ctx, store, "", []v1.Descriptor{desc}, oras.PackOptions{PackImageManifest: true})
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	err = store.Tag(ctx, root, re.Version)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	repo, err := opts.NewRepository(opts.TargetRef, opts.Common)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	_, err = oras.Copy(ctx, store, re.Version, repo, re.Version, oras.DefaultCopyOptions)
	if err != nil {
		if strings.Contains(err.Error(), "credential required") {
			cmd.PrintErrln("Error: failed to authorize: 401 Unauthorized")
		} else {
			cmd.PrintErrf("Error: unexpected error happened: %s\n", err)
		}
		return
	}

	cmd.Println("Recipe pushed successfully!")
}
