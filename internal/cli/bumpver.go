package cli

import (
	"errors"
	"strconv"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/changelog"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type bumpVerOpts struct {
	RecipeName    string
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
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
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
		
# Directly specify version and message
jalapeno bumpver -v v1.0.0 -m "Hello world"`,
	}

	cmd.Flags().StringVarP(&opts.RecipeName, "recipe", "r", "", "Name of the recipe to upgrade")
	cmd.Flags().StringVarP(&opts.RecipeVersion, "version", "v", "", "New semver number for recipe")
	cmd.Flags().StringVarP(&opts.ChangelogMsg, "message", "m", "", "Optional changelog message")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runBumpVer(cmd *cobra.Command, opts bumpVerOpts) error {
	var increment string
	newVer := opts.RecipeVersion
	changelogMsg := opts.ChangelogMsg

	if newVer != "" && !semver.IsValid(newVer) {
		return errors.New("invalid version format, please enter valid Semantic Version")
	}

	re, err := recipe.LoadRecipe(opts.RecipeName)
	if err != nil {
		return err
	}

	if opts.RecipeVersion == "" {
		cmd.Println()
		prompt, err := changelog.RunChangelog()

		if err != nil {
			return err
		}

		increment, changelogMsg = prompt[0], prompt[1]
		bumpedVer, err := BumpSemVer(re.Metadata.Version, increment)

		if err != nil {
			return err
		}

		newVer = bumpedVer
	}

	err = re.Metadata.Update(re, newVer, changelogMsg)
	if err != nil {
		return err
	}

	err = re.Save(opts.WorkingDirectory.Dir)
	if err != nil {
		return err
	}

	cmd.Printf("bumped version %s => %s \n", re.Metadata.Version, newVer)
	cmd.Printf("with changelog message %s \n", changelogMsg)

	return nil
}

// BumpSemVer takes SemVer as string and an increment as string "patch", "minor", or "major"
// and bumps the version by the increment and returns the new version number
func BumpSemVer(ver string, by string) (string, error) {
	trim := strings.TrimPrefix(ver, "v")
	parsed := strings.Split(trim, ".")

	intArr := make([]int, len(parsed))
	resultArr := []string{"", "", ""}

	for i, v := range parsed {
		conv, err := strconv.Atoi(v)
		if err != nil {
			return "", err
		}
		intArr[i] = conv
	}

	switch by {
	case "patch":
		intArr[2]++
	case "minor":
		intArr[2] = 0
		intArr[1]++
	case "major":
		intArr[2] = 0
		intArr[1] = 0
		intArr[0]++
	}

	for i, v := range intArr {
		conv := strconv.Itoa(v)
		resultArr[i] = conv
	}

	resultArr[0] = "v" + resultArr[0]
	newVer := strings.Join(resultArr, ".")

	if semver.IsValid(newVer) {
		return newVer, nil
	} else {
		return "", errors.New("failed to create valid semantic version")
	}
}
