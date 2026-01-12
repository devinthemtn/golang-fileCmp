package differ

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffType represents the type of difference
type DiffType int

const (
	DiffEqual DiffType = iota
	DiffInsert
	DiffDelete
)

// DiffLine represents a single line in the diff output
type DiffLine struct {
	Type    DiffType
	Content string
	LineNum int
}

// FileDiff represents the complete diff between two files
type FileDiff struct {
	LeftFile  string
	RightFile string
	Lines     []DiffLine
}

// Differ handles file comparison operations
type Differ struct {
	dmp *diffmatchpatch.DiffMatchPatch
}

// New creates a new Differ instance
func New() *Differ {
	dmp := diffmatchpatch.New()
	return &Differ{
		dmp: dmp,
	}
}

// CompareFiles compares two files and returns a structured diff
func (d *Differ) CompareFiles(leftPath, rightPath string, leftContent, rightContent io.Reader) (*FileDiff, error) {
	leftLines, err := readLines(leftContent)
	if err != nil {
		return nil, fmt.Errorf("error reading left file: %w", err)
	}

	rightLines, err := readLines(rightContent)
	if err != nil {
		return nil, fmt.Errorf("error reading right file: %w", err)
	}

	return d.compareLines(leftPath, rightPath, leftLines, rightLines), nil
}

// CompareStrings compares two strings and returns a structured diff
func (d *Differ) CompareStrings(leftPath, rightPath, leftContent, rightContent string) *FileDiff {
	leftLines := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	return d.compareLines(leftPath, rightPath, leftLines, rightLines)
}

// compareLines performs line-by-line comparison
func (d *Differ) compareLines(leftPath, rightPath string, leftLines, rightLines []string) *FileDiff {
	diff := &FileDiff{
		LeftFile:  leftPath,
		RightFile: rightPath,
		Lines:     make([]DiffLine, 0),
	}

	// Join lines with newlines for diff algorithm
	leftText := strings.Join(leftLines, "\n")
	rightText := strings.Join(rightLines, "\n")

	// Compute diffs
	diffs := d.dmp.DiffMain(leftText, rightText, false)
	diffs = d.dmp.DiffCleanupSemantic(diffs)

	leftLineNum := 1
	rightLineNum := 1

	for _, diffOp := range diffs {
		lines := strings.Split(diffOp.Text, "\n")

		// Handle empty splits
		if len(lines) == 1 && lines[0] == "" {
			continue
		}

		switch diffOp.Type {
		case diffmatchpatch.DiffEqual:
			for _, line := range lines {
				if line == "" && len(lines) == 1 {
					continue
				}
				diff.Lines = append(diff.Lines, DiffLine{
					Type:    DiffEqual,
					Content: line,
					LineNum: leftLineNum,
				})
				leftLineNum++
				rightLineNum++
			}

		case diffmatchpatch.DiffDelete:
			for _, line := range lines {
				if line == "" && len(lines) == 1 {
					continue
				}
				diff.Lines = append(diff.Lines, DiffLine{
					Type:    DiffDelete,
					Content: line,
					LineNum: leftLineNum,
				})
				leftLineNum++
			}

		case diffmatchpatch.DiffInsert:
			for _, line := range lines {
				if line == "" && len(lines) == 1 {
					continue
				}
				diff.Lines = append(diff.Lines, DiffLine{
					Type:    DiffInsert,
					Content: line,
					LineNum: rightLineNum,
				})
				rightLineNum++
			}
		}
	}

	return diff
}

// readLines reads all lines from a reader
func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// GetStats returns statistics about the diff
func (fd *FileDiff) GetStats() (int, int, int) {
	equal, inserted, deleted := 0, 0, 0

	for _, line := range fd.Lines {
		switch line.Type {
		case DiffEqual:
			equal++
		case DiffInsert:
			inserted++
		case DiffDelete:
			deleted++
		}
	}

	return equal, inserted, deleted
}
