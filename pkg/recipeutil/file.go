package recipeutil

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/xlab/treeprint"
	"golang.org/x/exp/maps"
)

type FileStatus int

const (
	FileUnknown FileStatus = iota
	FileUnchanged
	FileAdded
	FileModified
	FileDeleted
)

func CreateFileTree(root string, files map[string]FileStatus) string {
	tree := treeprint.NewWithRoot(lipgloss.NewStyle().Bold(true).Render(root))

	filepaths := maps.Keys(files)
	slices.Sort(filepaths)

	for _, filepath := range filepaths {
		filepathSegments := strings.Split(filepath, "/")

		// If the last segment is empty, it means the filepath should represent a directory
		if filepathSegments[len(filepathSegments)-1] == "" {
			filepathSegments = filepathSegments[:len(filepathSegments)-1]
			filepathSegments[len(filepathSegments)-1] += "/"
		}

		cursor := tree
		for i, filepathSegment := range filepathSegments {
			if node := cursor.FindByValue(filepathSegment); node == nil {
				branchName := filepathSegment
				status := files[filepath]
				isLeaf := i == len(filepathSegments)-1

				if isLeaf {
					switch status {
					case FileUnknown:
						break
					case FileDeleted:
						filepathSegment = lipgloss.NewStyle().Strikethrough(true).Render(filepathSegment)
						fallthrough
					default:
						branchName = fmt.Sprintf("%s (%s)", filepathSegment, status)
					}
				}
				cursor = cursor.AddBranch(branchName)
			} else {
				cursor = node
			}
		}
	}

	return tree.String()
}

func (f FileStatus) String() string {
	switch f {
	case FileUnchanged:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render("unchanged")
	case FileAdded:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#26A568")).Render("added")
	case FileModified:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FADA5E")).Render("modified")
	case FileDeleted:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4136")).Render("deleted")
	default:
		return ""
	}
}
