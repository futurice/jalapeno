package cli

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/changelog"
	"github.com/spf13/cobra"
)

type bumpVerOpts struct {
	RecipePath    string
	RecipeVersion string
	ChangelogMsg  string
	option.Common
	option.WorkingDirectory
}

func NewBumpVerCmd() *cobra.Command {
	var opts bumpVerOpts
	var cmd = &cobra.Command{
		Use:   "bumpver",
		Short: "Bump version number for recipe",
		Long:  "Bump version number for recipe. By default prompts user for update increment (patch/minor/major) and changelog messsage. These can also be specified directly with the -v and -m flags.",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				opts.RecipePath = args[0]
			} else {
				opts.RecipePath = "."
			}
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runBumpVer(cmd, opts)
			return errorHandler(cmd, err)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			cmd.Root().SetContext(cmd.Context())
		},
		Example: `# Prompt version increment and changelog message
jalapeno bumpver

# Specify recipe directory
jalapeno bumpver path/to/recipe

# Directly specify version and message
jalapeno bumpver -v v1.0.0 -m "Hello world"`,
	}

	cmd.Flags().StringVarP(&opts.RecipeVersion, "version", "v", "", "New semver number for recipe")
	cmd.Flags().StringVarP(&opts.ChangelogMsg, "message", "m", "", "Optional changelog message")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runBumpVer(cmd *cobra.Command, opts bumpVerOpts) error {
	var newVer semver.Version
	var changelogMsg string

	re, err := recipe.LoadRecipe(opts.RecipePath)
	if err != nil {
		return err
	}

	if opts.RecipeVersion == "" {
		currentVer, err := semver.NewVersion(re.Metadata.Version)
		if err != nil {
			return err
		}

		changelog, err := changelog.RunChangelog()
		if err != nil {
			return err
		}

		switch changelog.Increment {
		case "patch":
			newVer = currentVer.IncPatch()
		case "minor":
			newVer = currentVer.IncMinor()
		case "major":
			newVer = currentVer.IncMajor()
		}

		changelogMsg = changelog.Msg

	} else {
		optVer, err := semver.NewVersion(opts.RecipeVersion)
		if err != nil {
			if errors.Is(err, semver.ErrInvalidSemVer) {
				return fmt.Errorf("provided version is not valid semver: %s", opts.RecipeVersion)
			}
			return err
		}

		newVer = *optVer
		changelogMsg = opts.ChangelogMsg
	}

	newVerWithPrefix := "v" + newVer.String()
	prevVer := re.Metadata.Version

	re.Metadata.UpdateVersion(re, newVerWithPrefix, changelogMsg)

	err = re.Save(opts.RecipePath)
	if err != nil {
		return err
	}

	cmd.Printf("bumped version: %s => %s \n", prevVer, newVerWithPrefix)
	cmd.Printf("with changelog message: %s \n", changelogMsg)

	return nil
}
