package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/gofrs/uuid"
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
		Long: `Delete sauce(s) from the project. This will remove the rendered files and the sauce entry from sauces.yml.
If no sauce ID is provided and --all flag is not set, this command will fail.`,
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
jalapeno delete 21872763-f48e-4728-bc49-57f5898e098a

# Delete all sauces (same as 'jalapeno eject')
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

func deleteSauce(cmd *cobra.Command, opts deleteOptions) error {
	id, err := uuid.FromString(opts.SauceID)
	if err != nil {
		return fmt.Errorf("invalid sauce ID: %w", err)
	}

	sauce, err := recipe.LoadSauceByID(opts.Dir, id)
	if err != nil {
		if errors.Is(err, recipe.ErrSauceNotFound) {
			return fmt.Errorf("sauce with ID '%s' not found", opts.SauceID)
		}
		return err
	}

	// Delete rendered files
	for path := range sauce.Files {
		fullPath := filepath.Join(opts.Dir, path)
		err := os.Remove(fullPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to delete file '%s': %w", path, err)
		}
	}

	// Delete sauce entry
	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	filteredSauces := make([]*recipe.Sauce, 0, len(sauces))
	for _, s := range sauces {
		if s.ID != id {
			filteredSauces = append(filteredSauces, s)
		}
	}

	// If this was the last sauce, delete the entire .jalapeno directory
	if len(filteredSauces) == 0 {
		jalapenoPath := filepath.Join(opts.Dir, recipe.SauceDirName)
		err = os.RemoveAll(jalapenoPath)
		if err != nil {
			return err
		}
	} else {
		// Otherwise just save the filtered sauces
		err = recipe.SaveSauces(opts.Dir, filteredSauces)
		if err != nil {
			return err
		}
	}

	cmd.Printf("Sauce '%s' deleted %s\n", sauce.Recipe.Name, colors.Green.Render("successfully!"))
	return nil
}

func deleteAll(cmd *cobra.Command, opts deleteOptions) error {
	jalapenoPath := filepath.Join(opts.Dir, recipe.SauceDirName)

	if stat, err := os.Stat(jalapenoPath); os.IsNotExist(err) || !stat.IsDir() {
		return fmt.Errorf("'%s' is not a Jalapeno project", opts.Dir)
	}

	// Delete all rendered files first
	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	for _, sauce := range sauces {
		for path := range sauce.Files {
			fullPath := filepath.Join(opts.Dir, path)
			err := os.Remove(fullPath)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to delete file '%s': %w", path, err)
			}
		}
	}

	// Delete .jalapeno directory
	cmd.Printf("Deleting %s...\n", jalapenoPath)
	err = os.RemoveAll(jalapenoPath)
	if err != nil {
		return err
	}

	cmd.Printf("All sauces deleted %s\n", colors.Green.Render("successfully!"))
	return nil
}
