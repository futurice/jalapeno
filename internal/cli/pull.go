package cli

import (
	"context"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/oci"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	TargetRef string

	option.WorkingDirectory
	option.OCIRepository
	option.Common
}

func NewPullCmd() *cobra.Command {
	var opts pullOptions
	var cmd = &cobra.Command{
		Use:   "pull URL",
		Short: "Pull a recipe from OCI repository",
		Long:  "TODO",
		Example: `# asd
asd
`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.TargetRef = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runPull(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runPull(cmd *cobra.Command, opts pullOptions) {
	ctx := context.Background()

	err := oci.SaveRemoteRecipe(ctx, opts.Dir,
		oci.Repository{
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
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			cmd.PrintErrln("Error: recipe not found") // TODO: Give more descriptive error message
		} else {
			cmd.PrintErrf("Error: %s", err)
		}
		return
	}

	cmd.Println("Recipe pulled successfully")
}
