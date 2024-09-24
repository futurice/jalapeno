package conflict

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/diff"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	uiutil "github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/muesli/termenv"
)

type ConflictResolution int

const (
	UseOld ConflictResolution = iota
	UseNew
	UseDiffFile
)

type Model struct {
	resolution ConflictResolution
	filePath   string
	fileA      []byte
	fileB      []byte
	diff       diff.Diff
	ready      bool
	err        error
	submitted  bool
	viewport   viewport.Model
	lineCount  int
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

var _ tea.Model = Model{}

func Solve(in io.Reader, out io.Writer, filePath string, fileA, fileB []byte) ([]byte, error) {
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	p := tea.NewProgram(NewModel(filePath, fileA, fileB), tea.WithInput(in), tea.WithOutput(out))
	if m, err := p.Run(); err != nil {
		return []byte{}, err
	} else {
		m, ok := m.(Model)
		if !ok {
			return []byte{}, errors.New("internal error: unexpected model type")
		}

		if m.err != nil {
			return []byte{}, m.err
		}

		return m.Result(), nil
	}
}

func NewModel(filePath string, fileA, fileB []byte) Model {

	fileDiff := diff.New(string(fileA), string(fileB))

	m := Model{
		filePath:   filePath,
		fileA:      fileA,
		fileB:      fileB,
		diff:       fileDiff,
		ready:      false,
		resolution: UseOld,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func calculateViewportHeight(windowHeight, verticalMarginal int) int {
	viewportHeight := windowHeight - verticalMarginal
	return viewportHeight
}

func (m Model) headerView() string {
	title := titleStyle.Render(m.filePath)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	viewPortHeader := lipgloss.JoinHorizontal(lipgloss.Center, title, line)
	conflictPrompt := fmt.Sprintf("There are conflicts in the following file: %s, what do you want to do?", m.filePath)
	return fmt.Sprintf("%s\n%s", conflictPrompt, viewPortHeader)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	progressLine := lipgloss.JoinHorizontal(lipgloss.Center, line, info)

	overwriteOpt := "Overwrite old"
	keepOpt := "Keep old"
	diffOpt := "Write the file with diffs"
	var resolutionSelector string

	selectedStyle := lipgloss.NewStyle().Bold(true)

	if m.resolution == UseOld {
		resolutionSelector = fmt.Sprintf("> %s / %s / %s", selectedStyle.Render(keepOpt), overwriteOpt, diffOpt)
	} else if m.resolution == UseNew {
		resolutionSelector = fmt.Sprintf("%s / > %s / %s", keepOpt, selectedStyle.Render(overwriteOpt), diffOpt)
	} else {
		resolutionSelector = fmt.Sprintf("%s / %s / > %s", keepOpt, overwriteOpt, selectedStyle.Render(diffOpt))
	}

	instructions := "Use up and down arrows to move up and down the diff file.\nPage down to go the the next conflict and page up to go to the previous conflict."

	return lipgloss.JoinVertical(0, progressLine, instructions, resolutionSelector)
}

func combineDiffLinesToColorizedDiffString(lines []string) string {
	colorizedLines := make([]string, 0)
	for _, line := range lines {
		switch line[0] {
		case ' ':
			colorizedLines = append(colorizedLines, line)
		case '+':
			colorizedLines = append(colorizedLines, colors.Green.Render(line))
		case '-':
			colorizedLines = append(colorizedLines, colors.Red.Render(line))
		}
	}
	colorizedLines = append(colorizedLines, colors.Yellow.Render("-- End of file --"))
	return strings.Join(colorizedLines, "\n")
}

func findNextDiffLineNumber(diffLineNums []int, curLine int, lastLine int) int {
	for _, lineNum := range diffLineNums {
		if lineNum > curLine {
			return lineNum
		}
	}
	// Did not find any diff that was after current line, go to last line.
	// We can set the yOffset to the last line, because it still has the
	// file end marker that will be shown.
	return lastLine
}

func findPrevDiffLineNumber(diffLineNums []int, curLine int) int {
	result := 0

	for _, lineNum := range diffLineNums {
		if lineNum < curLine {
			result = lineNum
		}
	}

	// if there is no previous diff, go the the 0th line
	return result
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// We need to init the viewport here so we know it's size.
			m.viewport = viewport.New(msg.Width, calculateViewportHeight(msg.Height, verticalMarginHeight))
			m.viewport.YPosition = headerHeight
			diffLines := m.diff.GetUnifiedDiffLines()
			m.viewport.SetContent(combineDiffLinesToColorizedDiffString(diffLines))
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
			m.lineCount = len(diffLines)
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = calculateViewportHeight(msg.Height, verticalMarginHeight)
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = uiutil.ErrUserAborted
			return m, tea.Quit
		case tea.KeyEnter:
			m.submitted = true
			return m, tea.Quit
		case tea.KeyRight:
			m.resolution = min(m.resolution+1, UseDiffFile)
		case tea.KeyLeft:
			m.resolution = max(m.resolution-1, UseOld)
		case tea.KeyPgDown:
			m.viewport.YOffset = findNextDiffLineNumber(m.diff.GetUnifiedDiffConflictIndices(), m.viewport.YOffset, m.lineCount)
		case tea.KeyPgUp:
			m.viewport.YOffset = findPrevDiffLineNumber(m.diff.GetUnifiedDiffConflictIndices(), m.viewport.YOffset)
		case tea.KeyDown:
			m.viewport.YOffset = min(m.viewport.TotalLineCount()-1, m.viewport.YOffset+1)
		case tea.KeyUp:
			m.viewport.YOffset = max(0, m.viewport.YOffset-1)
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	if m.submitted || m.err != nil {
		s.WriteString(fmt.Sprintf("%s: ", m.filePath))
		if m.resolution == UseOld {
			s.WriteString("keep old")
		} else if m.resolution == UseNew {
			s.WriteString("use new")
		} else {
			s.WriteString("write diff file")
		}

		s.WriteRune('\n')

		return s.String()
	}

	s.WriteString(m.headerView())
	s.WriteString("\n")
	s.WriteString(m.viewport.View())
	s.WriteString("\n")
	s.WriteString(m.footerView())

	return s.String()
}

func (m Model) Result() []byte {
	if m.resolution == UseOld {
		return m.fileA
	} else if m.resolution == UseNew {
		return m.fileB
	} else {
		return []byte(m.diff.GetConflictResolutionTemplate())
	}
}
