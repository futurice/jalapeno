package diff

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kylelemons/godebug/diff"
)

type Diff struct {
	a string
	b string

	chunks []diff.Chunk

	// lazily created cache, can only be created if there are some differences
	unifiedDiffDigest *unifiedDiffDigest
}

// New creates a new Diff that might be queried further.
func New(a, b string) Diff {
	linesA := strings.Split(a, "\n")
	linesB := strings.Split(b, "\n")

	chunks := diff.DiffChunks(linesA, linesB)

	return Diff{
		a: a,
		b: b,

		chunks: chunks,
	}
}

// IsDifferent returns true if there were any differences between a and b.
func (d *Diff) IsDifferent() bool {
	return d.chunks != nil
}

// GetUnifiedDiffLines returns the lines of the diff. The lines are not terminated with \n.
// Returns nil if there were zero differences.
func (d *Diff) GetUnifiedDiffLines() []string {
	if !d.IsDifferent() {
		return nil
	}

	return d.getUnifiedDiffDigest().lines
}

// GetUnifiedDiffConflictIndices returns the indices of the lines of the diff returned by either GetUnifiedDiffLines
// or GetUnifiedDiff. Indices are zero-based.
// Returns nil if there were zero differences. Returns at least one index otherwise.
func (d *Diff) GetUnifiedDiffConflictIndices() []int {
	if !d.IsDifferent() {
		return nil
	}

	return d.getUnifiedDiffDigest().conflictLineIndices
}

// GetUnifiedDiff returns the entire diff as a single string.
// Returns an empty string if there are no differences.
func (d *Diff) GetUnifiedDiff() string {
	if !d.IsDifferent() {
		return ""
	}

	return strings.Join(d.GetUnifiedDiffLines(), "\n") + "\n"
}

// GetConflictResolutionTemplate returns a single content with differences formatted in a manner similar to how
// merge conflicts are formatted by Git. Differences will be preceded by `<<<<<<< Deleted` followed by any lines
// present in `a` but not `b`, followed by `=======`, followed by any lines present in `b` but not `a`, followed by
// `>>>>>>> Added`. The lines in those blocks are not prefixed or otherwise altered. Any equal lines are included as is.
// Returns an empty string if there are no differences.
func (d *Diff) GetConflictResolutionTemplate() string {
	buf := new(bytes.Buffer)

	// diff.Chunk is documented not to contain both Deleted and Added lines in the same Chunk
	// this is used to combine adjacent chunks with Deleted and Added without intervening Equal lines
	currentConflict := newConflict()

	for chunkIndex, chunk := range d.chunks {
		if len(chunk.Deleted) > 0 {
			if !currentConflict.canAcceptDeletedLine() {
				// must start a new conflict block in order not to re-order lines
				currentConflict.flushTo(buf)
			}

			for _, deleted := range chunk.Deleted {
				currentConflict.appendDeletedLine(deleted)
			}
		}

		for _, added := range chunk.Added {
			currentConflict.appendAddedLine(added)
		}

		if len(chunk.Equal) > 0 {
			currentConflict.flushTo(buf)

			for equalIndex, equal := range chunk.Equal {
				buf.WriteString(equal)

				if chunkIndex < len(d.chunks)-1 || equalIndex < len(chunk.Equal)-1 {
					// we terminate each equal line with a \n, except for the last one
					// this is because we originally fed lines separated by \n to diff.DiffChunks, so we
					// only need to put the \n in between lines, not after the last one
					buf.WriteString("\n")
				}
			}
		}
	}

	currentConflict.flushTo(buf)

	return buf.String()
}

func (d *Diff) computeUnifiedDiffLineCount() int {
	lineCount := 0

	for _, chunk := range d.chunks {
		lineCount += len(chunk.Added) + len(chunk.Deleted) + len(chunk.Equal)
	}

	return lineCount
}

func (d *Diff) getUnifiedDiffDigest() *unifiedDiffDigest {
	if !d.IsDifferent() {
		return nil
	}

	if d.unifiedDiffDigest != nil {
		return d.unifiedDiffDigest
	}

	digest := d.createUnifiedDiffDigest()
	d.unifiedDiffDigest = digest

	return digest
}

func (d *Diff) createUnifiedDiffDigest() *unifiedDiffDigest {
	lineCount := d.computeUnifiedDiffLineCount()

	lines := make([]string, lineCount)
	conflictLineIndices := make([]int, 0, 10)

	lineIndex := 0
	isInsideConflict := false
	for _, chunk := range d.chunks {
		for _, deleted := range chunk.Deleted {
			lines[lineIndex] = fmt.Sprintf("-%s", deleted)
			if !isInsideConflict {
				conflictLineIndices = append(conflictLineIndices, lineIndex)
				isInsideConflict = true
			}
			lineIndex++
		}
		for _, added := range chunk.Added {
			lines[lineIndex] = fmt.Sprintf("+%s", added)
			if !isInsideConflict {
				conflictLineIndices = append(conflictLineIndices, lineIndex)
				isInsideConflict = true
			}
			lineIndex++
		}
		for _, equal := range chunk.Equal {
			lines[lineIndex] = fmt.Sprintf(" %s", equal)
			isInsideConflict = false
			lineIndex++
		}
	}

	return &unifiedDiffDigest{
		lines:               lines,
		conflictLineIndices: conflictLineIndices,
	}
}

type unifiedDiffDigest struct {
	lines               []string
	conflictLineIndices []int
}

const conflictHeader = "<<<<<<< Deleted\n"
const conflictDivider = "=======\n"
const conflictFooter = ">>>>>>> Added\n"

type conflict struct {
	deletedLines []string
	addedLines   []string
}

func newConflict() conflict {
	return conflict{}
}

func (c *conflict) hasAnyLines() bool {
	return len(c.deletedLines) > 0 || len(c.addedLines) > 0
}

func (c *conflict) canAcceptDeletedLine() bool {
	return len(c.addedLines) == 0
}

func (c *conflict) appendDeletedLine(line string) {
	c.deletedLines = append(c.deletedLines, line)
}

func (c *conflict) appendAddedLine(line string) {
	c.addedLines = append(c.addedLines, line)
}

func (c *conflict) flushTo(buf *bytes.Buffer) bool {
	if !c.hasAnyLines() {
		return false
	}

	buf.WriteString(conflictHeader)

	for _, deleted := range c.deletedLines {
		// we can ignore errors from Fprintf since bytes.Buffer panics on out of memory
		fmt.Fprintf(buf, "%s\n", deleted)
	}

	buf.WriteString(conflictDivider)

	for _, added := range c.addedLines {
		fmt.Fprintf(buf, "%s\n", added)
	}

	buf.WriteString(conflictFooter)

	c.deletedLines = nil
	c.addedLines = nil

	return true
}
