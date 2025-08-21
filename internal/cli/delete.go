package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	SauceID string
	All     bool

	option.Common
	option.WorkingDirectory
}

func NewDeleteCmd() *cobra.Command {
	var opts deleteOptions
	var cmd = &cobra.Command{
		Use:   "delete [SAUCE_ID]",
		Short: "Delete sauce(s) from the project",
		Long:  fmt.Sprintf(`Delete sauce(s) from the project. This will remove the rendered files and the sauce entry from %s%s.`, recipe.SaucesFileName, recipe.YAMLExtension),
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.SauceID = args[0]
			}
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runDelete(cmd, opts)
			return errorHandler(cmd, err)
		},
		Example: `# Delete a specific sauce
jalapeno list
jalapeno delete 21872763-f48e-4728-bc49-57f5898e098a

# Delete all sauces
jalapeno delete --all`,
	}

	cmd.Flags().BoolVar(&opts.All, "all", false, "Delete all sauces from the project")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runDelete(cmd *cobra.Command, opts deleteOptions) error {
	if !opts.All && opts.SauceID == "" {
		return errors.New("either provide a sauce ID or use --all flag")
	}

	if opts.All {
		return deleteAll(cmd, opts)
	}

	return deleteSauce(cmd, opts)
}

// TODO: The operation should be atomic
func deleteSauce(cmd *cobra.Command, opts deleteOptions) error {
	sauce, err := recipe.LoadSauceByID(opts.Dir, opts.SauceID)
	if err != nil {
		if errors.Is(err, recipe.ErrSauceNotFound) {
			return fmt.Errorf("sauce with ID '%s' not found", opts.SauceID)
		}
		return err
	}

	if err := deleteSauceFiles(sauce, opts.Dir); err != nil {
		return err
	}

	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	// Delete sauce entry
	filteredSauces := make([]*recipe.Sauce, 0, len(sauces))
	for _, s := range sauces {
		if s.ID.String() != opts.SauceID {
			filteredSauces = append(filteredSauces, s)
		}
	}

	if err := deleteSauceDir(opts.Dir); err != nil {
		return err
	}

	// Save remaining sauces
	for _, sauce := range filteredSauces {
		err := sauce.Save(opts.Dir)
		if err != nil {
			return err
		}
	}

	cmd.Printf("Sauce '%s' (from recipe '%s') deleted %s\n", sauce.ID, sauce.Recipe.Name, colors.Green.Render("successfully!"))
	return nil
}

// TODO: The operation should be atomic
func deleteAll(cmd *cobra.Command, opts deleteOptions) error {
	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	if len(sauces) == 0 {
		return fmt.Errorf("the directory '%s' did not contain any sauces to delete", opts.Dir)
	}

	for _, sauce := range sauces {
		if err := deleteSauceFiles(sauce, opts.Dir); err != nil {
			return err
		}
	}

	if err := deleteSauceDir(opts.Dir); err != nil {
		return err
	}

	cmd.Printf("All sauces deleted %s\n", colors.Green.Render("successfully!"))
	return nil
}

func deleteSauceFiles(sauce *recipe.Sauce, dir string) error {
	for path := range sauce.Files {
		fullPath := filepath.Join(dir, path)
		err := os.Remove(fullPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to delete file '%s': %w", path, err)
		}
	}
	return nil
}

func deleteSauceDir(dir string) error {
	return os.RemoveAll(filepath.Join(dir, recipe.SauceDirName))
}
