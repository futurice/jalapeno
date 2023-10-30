package cli

import (
	"context"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/oci"
	"github.com/spf13/cobra"
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
		Long:  "Push a recipe to OCI repository (e.g. Docker registry). You can authenticate by using the --username and --password flags or logging in first with `docker login`.",
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

	err := oci.PushRecipe(ctx, opts.RecipePath, oci.Repository{
		Reference: opts.TargetRef,
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
		cmd.PrintErrf("Error: %s\n", err)
		return
	}

	cmd.Println("Recipe pushed successfully")
}
