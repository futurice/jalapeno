package recipeutil

import (
	"slices"
	"strings"

	"github.com/xlab/treeprint"
)

func CreateFileTree(root string, filepaths []string) string {
	tree := treeprint.NewWithRoot(root)

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
