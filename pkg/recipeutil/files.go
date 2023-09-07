package recipeutil

import (
	"slices"
	"strings"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/xlab/treeprint"
)

func CreateFileTree(root string, files map[string]recipe.File) string {
	tree := treeprint.NewWithRoot(root)

	filepaths := make([]string, len(files))

	i := 0
	for f := range files {
		filepaths[i] = f
		i++
	}

	slices.Sort(filepaths)

	for _, filepath := range filepaths {
		filepathSegments := strings.Split(filepath, "/")
		cursor := tree
		for _, filepathSegment := range filepathSegments {
			if node := cursor.FindByValue(filepathSegment); node == nil {
				cursor = cursor.AddBranch(filepathSegment)
			} else {
				cursor = node
			}
		}
	}

	return tree.String()
}
