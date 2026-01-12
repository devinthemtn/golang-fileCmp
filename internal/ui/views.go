package ui

import (
	"fmt"
	"strings"

	"golang-fileCmp/internal/differ"
)

// renderFileSelectView renders the file selection interface
func (m *Model) renderFileSelectView() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("File Comparison Tool")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Input fields
	leftLabel := "Left Path: "
	rightLabel := "Right Path: "

	var leftInput, rightInput string
	if m.focusLeft {
		leftInput = focusedInputStyle.Render(m.inputLeft + "│")
		rightInput = inputStyle.Render(m.inputRight)
	} else {
		leftInput = inputStyle.Render(m.inputLeft)
		rightInput = focusedInputStyle.Render(m.inputRight + "│")
	}

	b.WriteString(leftLabel + leftInput)
	b.WriteString("\n")
	b.WriteString(rightLabel + rightInput)
	b.WriteString("\n\n")

	// Status information
	if m.leftFile != nil {
		b.WriteString(fmt.Sprintf("✓ Left: %s ", m.leftPath))
		if m.leftFile.IsDir {
			b.WriteString("(directory)")
		} else {
			b.WriteString("(file)")
		}
		b.WriteString("\n")
	}

	if m.rightFile != nil {
		b.WriteString(fmt.Sprintf("✓ Right: %s ", m.rightPath))
		if m.rightFile.IsDir {
			b.WriteString("(directory)")
		} else {
			b.WriteString("(file)")
		}
		b.WriteString("\n")
	}

	if len(m.commonFiles) > 0 {
		b.WriteString(fmt.Sprintf("\nFound %d common files:\n", len(m.commonFiles)))
		b.WriteString(m.renderFileList())
	}

	// Error message
	if m.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("Error: " + m.errorMsg))
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Tab: Switch input • Enter: Load path • ↑↓: Select file • Ctrl+D: Compare • ?: Help • Q: Quit"))

	return b.String()
}

// renderFileList renders the list of common files
func (m *Model) renderFileList() string {
	if len(m.commonFiles) == 0 {
		return ""
	}

	var items []string
	for relPath := range m.commonFiles {
		if relPath == m.selectedFile {
			items = append(items, selectedFileStyle.Render("→ "+relPath))
		} else {
			items = append(items, "  "+relPath)
		}
	}

	content := strings.Join(items, "\n")
	return fileListStyle.Render(content)
}

// renderDiffView renders the diff comparison view
func (m *Model) renderDiffView() string {
	if m.currentDiff == nil {
		return "No diff loaded"
	}

	var b strings.Builder

	// Header with file names
	header := fmt.Sprintf("%s vs %s", m.currentDiff.LeftFile, m.currentDiff.RightFile)
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n\n")

	// Stats
	equal, inserted, deleted := m.currentDiff.GetStats()
	stats := fmt.Sprintf("Lines: %d equal, %d inserted (+), %d deleted (-)",
		equal, inserted, deleted)
	b.WriteString(helpStyle.Render(stats))
	b.WriteString("\n\n")

	// Diff content
	b.WriteString(m.renderDiffContent())

	// Navigation help
	b.WriteString("\n")
	help := "↑↓/j/k: Navigate • g/G: Top/Bottom • n/p: Next/Prev file • Esc: Back • ?: Help • Q: Quit"
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

// renderDiffContent renders the actual diff lines with syntax highlighting
func (m *Model) renderDiffContent() string {
	if m.currentDiff == nil || len(m.currentDiff.Lines) == 0 {
		return "No differences found"
	}

	var b strings.Builder
	maxVisible := m.windowHeight - 10 // Account for header, stats, and help text

	start := m.scrollOffset
	end := start + maxVisible
	if end > len(m.currentDiff.Lines) {
		end = len(m.currentDiff.Lines)
	}

	for i := start; i < end; i++ {
		line := m.currentDiff.Lines[i]
		prefix := "  "

		if i == m.cursor {
			prefix = "▶ "
		}

		lineNum := fmt.Sprintf("%4d", line.LineNum)
		content := line.Content

		// Truncate long lines to fit screen width
		maxContentWidth := m.windowWidth - 20 // Account for prefix and line numbers
		if maxContentWidth > 0 && len(content) > maxContentWidth {
			content = content[:maxContentWidth-3] + "..."
		}

		var renderedLine string
		switch line.Type {
		case differ.DiffEqual:
			renderedLine = equalLineStyle.Render(fmt.Sprintf("%s%s %s", prefix, lineNum, content))
		case differ.DiffInsert:
			renderedLine = insertLineStyle.Render(fmt.Sprintf("%s%s +%s", prefix, lineNum, content))
		case differ.DiffDelete:
			renderedLine = deleteLineStyle.Render(fmt.Sprintf("%s%s -%s", prefix, lineNum, content))
		}

		b.WriteString(renderedLine)
		b.WriteString("\n")
	}

	// Show scroll indicator if needed
	if len(m.currentDiff.Lines) > maxVisible {
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d lines", start+1, end, len(m.currentDiff.Lines))
		b.WriteString(helpStyle.Render(scrollInfo))
		b.WriteString("\n")
	}

	return b.String()
}

// renderHelpView renders the help screen
func (m *Model) renderHelpView() string {
	var b strings.Builder

	title := titleStyle.Render("Help - File Comparison Tool")
	b.WriteString(title)
	b.WriteString("\n\n")

	help := `File Selection Mode:
  Tab              Switch between left/right input fields
  Enter            Load the entered path (file or directory)
  ↑/↓              Navigate through common files list
  Ctrl+D           Start comparing selected files
  ?                Show this help screen
  Q/Ctrl+C         Quit application

Diff View Mode:
  ↑/↓ or j/k       Navigate through diff lines
  g                Go to top of diff
  G                Go to bottom of diff
  n                Next common file
  p                Previous common file
  Esc              Return to file selection
  ?                Show this help screen
  Q/Ctrl+C         Quit application

Color Legend:
  `

	b.WriteString(help)

	// Color examples
	b.WriteString(insertLineStyle.Render("Green background: Added lines (+)"))
	b.WriteString("\n")
	b.WriteString(deleteLineStyle.Render("Blue background: Deleted lines (-)"))
	b.WriteString("\n")
	b.WriteString(equalLineStyle.Render("Gray text: Unchanged lines"))
	b.WriteString("\n\n")

	instructions := `Instructions:
1. Enter paths to files or directories in the input fields
2. Press Enter to load each path
3. Use ↑/↓ to select which common file to compare
4. Press Ctrl+D to start comparing
5. Navigate through the diff using arrow keys or j/k
6. Use n/p to switch between different files

Note: Only text files will be compared. The tool automatically
detects common text file extensions and filenames.`

	b.WriteString(instructions)
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press Esc or ? to return to previous view"))

	return b.String()
}
