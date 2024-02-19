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
	FileUnchanged FileStatus = iota
	FileAdded
	FileModified
	FileDeleted
)

func CreateFileTree(root string, files map[string]FileStatus) string {
	tree := treeprint.NewWithRoot(root)

	filepaths := maps.Keys(files)
	slices.Sort(filepaths)

	for _, filepath := range filepaths {
		filepathSegments := strings.Split(filepath, "/")
		cursor := tree
		for i, filepathSegment := range filepathSegments {
			if node := cursor.FindByValue(filepathSegment); node == nil {
				branchName := filepathSegment
				status := files[filepath]
				if i == len(filepathSegments)-1 {
					if status == FileDeleted {
						filepathSegment = lipgloss.NewStyle().Strikethrough(true).Render(filepathSegment)
					}
					branchName = fmt.Sprintf("%s (%s)", filepathSegment, status)
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
