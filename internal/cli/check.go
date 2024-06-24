package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

type checkOptions struct {
	CheckFrom           string
	RecipeName          string
	UseDetailedExitCode bool
	Upgrade             bool
	ForceUpgrade        bool

	option.Common
	option.OCIRepository
	option.WorkingDirectory
	option.Timeout
}

const (
	ExitCodeOK               = 0
	ExitCodeError            = 1
	ExitCodeUpdatesAvailable = 2
)

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
			err := runCheck(cmd, opts)
			return errorHandler(cmd, err)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			cmd.Root().SetContext(cmd.Context())
		},
		Example: `# Check updates for all recipes in the project
jalapeno check

# Check updates for a single recipe
jalapeno check --recipe my-recipe

# Upgrade recipes to the latest version if new versions are found
jalapeno check --upgrade

# Add check URL for recipe which does not have it yet
jalapeno check --recipe my-recipe --from oci://my-registry.com/my-recipe`,
	}

	cmd.Flags().StringVarP(&opts.RecipeName, "recipe", "r", "", "Name of the recipe to check for new versions")
	cmd.Flags().StringVar(&opts.CheckFrom, "from", "", "Add or override the URL used for checking updates for the recipe. Works only with --recipe flag")
	cmd.Flags().BoolVar(&opts.Upgrade, "upgrade", false, "Upgrade recipes to the latest version if new versions are found")
	cmd.Flags().BoolVar(&opts.ForceUpgrade, "force-upgrade", false, "If upgrading, overwrite manual changes in the files with the new versions without prompting")
	cmd.Flags().BoolVar(
		&opts.UseDetailedExitCode,
		"detailed-exitcode",
		false,
		fmt.Sprintf("Returns a detailed exit code when the command exits. When provided, this argument changes the exit codes and their meanings to provide more granular information about what the resulting plan contains: %d = Succeeded with no updates available, %d = Error, %d = Succeeded with updates available", ExitCodeOK, ExitCodeUpdatesAvailable, ExitCodeUpdatesAvailable))

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCheck(cmd *cobra.Command, opts checkOptions) error {
	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return fmt.Errorf("can not load sauces: %w", err)
	}

	if len(sauces) == 0 {
		return fmt.Errorf("working directory '%s' does not contain any sauces", opts.Dir)
	}

	// Check if we are looking updates for a specific recipe
	if opts.RecipeName != "" {
		filtered := make([]*recipe.Sauce, 0, len(sauces))
		for _, sauce := range sauces {
			if sauce.Recipe.Name == opts.RecipeName {
				filtered = append(filtered, sauce)
			}
		}

		if len(filtered) == 0 {
			return fmt.Errorf("project %s does not contain a sauce with recipe %s. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, opts.RecipeName)
		}

		sauces = filtered

		if opts.CheckFrom != "" {
			for i := range sauces {
				// Save new check URL for sauces
				sauces[i].CheckFrom = opts.CheckFrom
				err = sauces[i].Save(opts.Dir)
				if err != nil {
					return err
				}
			}
		}
	} else if opts.CheckFrom != "" {
		return fmt.Errorf("can not use --from flag without --recipe flag")
	}

	if n := len(sauces); n > 1 {
		cmd.Printf("Checking new versions for %d recipes...", n)
	} else {
		cmd.Printf("Checking new versions for the recipe \"%s\"...", sauces[0].Recipe.Name)
	}

	errorsFound := false
	latestSauceVersions := make(map[*recipe.Sauce]string)
	for _, sauce := range sauces {
		versions, err := recipeutil.CheckForUpdates(sauce, opts.OCIRepository)
		if err != nil {
			errorsFound = true
			cmd.Printf("❌ %s: can not check for updates: %s\n", sauce.Recipe.Name, err)

		} else if len(versions) > 0 {
			cmd.Printf("🔄 %s: new versions found: %s\n", sauce.Recipe.Name, strings.Join(versions, ", "))
			latestSauceVersions[sauce] = versions[len(versions)-1]

		} else {
			cmd.Printf("👍 %s: no new versions found\n", sauce.Recipe.Name)
		}
	}

	cmd.Println()

	// Construct a list of upgradeable sauces so the order is deterministic when we list them
	upgradeableSauces := make([]*recipe.Sauce, 0, len(latestSauceVersions))
	for _, sauce := range sauces {
		if _, ok := latestSauceVersions[sauce]; ok {
			upgradeableSauces = append(upgradeableSauces, sauce)
		}
	}

	if !opts.Upgrade {
		if len(latestSauceVersions) > 0 {
			cmd.Println("To upgrade recipes to the latest version run:")
			for _, sauce := range upgradeableSauces {
				cmd.Printf("  %s upgrade %s:%s\n", os.Args[0], sauce.CheckFrom, latestSauceVersions[sauce])
			}
			cmd.Println("\nor rerun the command with '--upgrade' flag to upgrade all recipes to the latest version.")
		}

		var exitCode int
		switch {
		case errorsFound:
			exitCode = ExitCodeError
		case len(latestSauceVersions) > 0 && opts.UseDetailedExitCode:
			exitCode = ExitCodeUpdatesAvailable
		default:
			exitCode = ExitCodeOK
		}

		ctx := context.WithValue(cmd.Context(), ExitCodeContextKey{}, exitCode)
		cmd.SetContext(ctx)

		return nil
	}

	n := 0
	for _, sauce := range upgradeableSauces {
		err := runUpgrade(cmd, upgradeOptions{
			RecipeURL:        fmt.Sprintf("%s:%s", sauce.CheckFrom, latestSauceVersions[sauce]),
			TargetSauceID:    sauce.ID.String(),
			ReuseOldValues:   true,
			Force:            opts.ForceUpgrade,
			Common:           opts.Common,
			OCIRepository:    opts.OCIRepository,
			WorkingDirectory: opts.WorkingDirectory,
			Timeout:          opts.Timeout,
		})
		if err != nil {
			return err
		}

		n++
		if n <= len(latestSauceVersions) {
			cmd.Print("\n- - - - - - - - - -\n\n")
		}
	}

	cmd.Printf("All recipes with newer versions upgraded %s\n", ColorGreen.Render("successfully!"))
	return nil
}
