package ui

import (
	"fmt"
	"strings"

	"golang-fileCmp/internal/differ"

	"github.com/charmbracelet/lipgloss"
)

// renderFileSelectView renders the file selection interface
func (m *Model) renderFileSelectView() string {
	var b strings.Builder

	// Title
	title := titleStyle.Width(m.windowWidth - 2).Render("File Comparison Tool")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Input fields - adapt width to terminal
	maxInputWidth := m.windowWidth - 20 // Leave space for labels and padding
	if maxInputWidth < 30 {
		maxInputWidth = 30
	}

	leftLabel := "Left Path: "
	rightLabel := "Right Path: "

	// Truncate input display if too long
	leftDisplay := m.inputLeft
	rightDisplay := m.inputRight
	if len(leftDisplay) > maxInputWidth {
		leftDisplay = "..." + leftDisplay[len(leftDisplay)-maxInputWidth+3:]
	}
	if len(rightDisplay) > maxInputWidth {
		rightDisplay = "..." + rightDisplay[len(rightDisplay)-maxInputWidth+3:]
	}

	var leftInput, rightInput string
	if m.focusLeft {
		leftInput = focusedInputStyle.Width(maxInputWidth).Render(leftDisplay + "│")
		rightInput = inputStyle.Width(maxInputWidth).Render(rightDisplay)
	} else {
		leftInput = inputStyle.Width(maxInputWidth).Render(leftDisplay)
		rightInput = focusedInputStyle.Width(maxInputWidth).Render(rightDisplay + "│")
	}

	b.WriteString(leftLabel + leftInput)
	b.WriteString("\n")

	// Show suggestions for left path if focused and available
	if m.focusLeft && m.showSuggestions && len(m.leftSuggestions) > 0 {
		b.WriteString(m.renderSuggestions(m.leftSuggestions, m.leftSuggIndex))
		b.WriteString("\n")
	}

	b.WriteString(rightLabel + rightInput)
	b.WriteString("\n")

	// Show suggestions for right path if focused and available
	if !m.focusLeft && m.showSuggestions && len(m.rightSuggestions) > 0 {
		b.WriteString(m.renderSuggestions(m.rightSuggestions, m.rightSuggIndex))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Status information - truncate paths if too long
	if m.leftFile != nil {
		leftPath := m.leftPath
		if len(leftPath) > maxInputWidth {
			leftPath = "..." + leftPath[len(leftPath)-maxInputWidth+3:]
		}
		b.WriteString(fmt.Sprintf("✓ Left: %s ", leftPath))
		if m.leftFile.IsDir {
			b.WriteString("(directory)")
		} else {
			b.WriteString("(file)")
		}
		b.WriteString("\n")
	}

	if m.rightFile != nil {
		rightPath := m.rightPath
		if len(rightPath) > maxInputWidth {
			rightPath = "..." + rightPath[len(rightPath)-maxInputWidth+3:]
		}
		b.WriteString(fmt.Sprintf("✓ Right: %s ", rightPath))
		if m.rightFile.IsDir {
			b.WriteString("(directory)")
		} else {
			b.WriteString("(file)")
		}
		b.WriteString("\n")
	}

	if len(m.commonFiles) > 0 {
		b.WriteString(fmt.Sprintf("\nFound %d common files:\n", len(m.commonFiles)))
		if m.selectedFile != "" {
			selectedPath := m.selectedFile
			if len(selectedPath) > maxInputWidth {
				selectedPath = "..." + selectedPath[len(selectedPath)-maxInputWidth+3:]
			}
			b.WriteString(selectedFileStyle.Render(fmt.Sprintf("► Selected: %s", selectedPath)))
			b.WriteString(" ")
			b.WriteString(helpStyle.Render("(Press Ctrl+D to compare)"))
			b.WriteString("\n")
		}
		b.WriteString(m.renderFileList())
	}

	// Error message
	if m.errorMsg != "" {
		b.WriteString("\n")
		errorMsg := m.errorMsg
		if len(errorMsg) > m.windowWidth-20 {
			errorMsg = errorMsg[:m.windowWidth-23] + "..."
		}
		b.WriteString(errorStyle.Width(m.windowWidth - 4).Render("Error: " + errorMsg))
		b.WriteString("\n")
	}

	// Help text - adapt to width
	b.WriteString("\n")
	var helpText string
	if m.showSuggestions {
		if m.windowWidth > 80 {
			helpText = "↑↓: Navigate suggestions • Tab: Next suggestion • Enter: Accept • Esc: Cancel"
		} else {
			helpText = "↑↓: Navigate • Tab: Next • Enter: Accept • Esc: Cancel"
		}
	} else {
		if m.windowWidth > 80 {
			helpText = "Tab: Switch input • Enter: Load path • ↑↓: Navigate files • Ctrl+D: Compare • ?: Help • Q: Quit"
		} else if m.windowWidth > 60 {
			helpText = "Tab: Switch • Enter: Load • ↑↓: Navigate • Ctrl+D: Compare • ?: Help • Q: Quit"
		} else {
			helpText = "Tab:Switch Enter:Load ↑↓:Navigate Ctrl+D:Compare ?:Help Q:Quit"
		}
	}
	b.WriteString(helpStyle.Width(m.windowWidth - 2).Render(helpText))

	return b.String()
}

// renderFileList renders the list of common files
func (m *Model) renderFileList() string {
	if len(m.commonFiles) == 0 {
		return ""
	}

	// Define styles for file status indicators
	identicalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true) // Green
	differentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true) // Red
	sizeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))                 // Gray

	// Calculate available space for file list
	usedHeight := 15 // Approximate height used by title, inputs, status, help
	if m.showSuggestions {
		usedHeight += 10 // More space used by suggestions
	}
	if m.errorMsg != "" {
		usedHeight += 3 // Space for error message
	}

	availableHeight := m.windowHeight - usedHeight
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Get all files and sort them for consistent ordering
	allFiles := make([]string, 0, len(m.commonFiles))
	for relPath := range m.commonFiles {
		allFiles = append(allFiles, relPath)
	}
	// Sort files alphabetically
	for i := 0; i < len(allFiles)-1; i++ {
		for j := i + 1; j < len(allFiles); j++ {
			if allFiles[i] > allFiles[j] {
				allFiles[i], allFiles[j] = allFiles[j], allFiles[i]
			}
		}
	}

	// Find current selection index
	selectedIndex := 0
	for i, relPath := range allFiles {
		if relPath == m.selectedFile {
			selectedIndex = i
			break
		}
	}

	// Calculate scroll offset to keep selected item visible
	startIndex := 0
	if len(allFiles) > availableHeight {
		startIndex = selectedIndex - availableHeight/2
		if startIndex < 0 {
			startIndex = 0
		}
		if startIndex+availableHeight > len(allFiles) {
			startIndex = len(allFiles) - availableHeight
		}
	}

	endIndex := startIndex + availableHeight
	if endIndex > len(allFiles) {
		endIndex = len(allFiles)
	}

	// Calculate max width for file display
	maxFileWidth := m.windowWidth - 8 // Account for borders and padding

	var items []string
	for i := startIndex; i < endIndex; i++ {
		relPath := allFiles[i]
		filePair := m.commonFiles[relPath]
		leftFile := filePair[0]
		rightFile := filePair[1]

		// Check if files are identical
		isIdentical := leftFile.Content == rightFile.Content
		var statusIndicator string
		if isIdentical {
			statusIndicator = identicalStyle.Render("✓")
		} else {
			statusIndicator = differentStyle.Render("✗")
		}

		// Format file size (show left file size, or both if different)
		var sizeInfo string
		if leftFile.Size == rightFile.Size {
			sizeInfo = sizeStyle.Render(fmt.Sprintf("(%s)", formatFileSize(leftFile.Size)))
		} else {
			sizeInfo = sizeStyle.Render(fmt.Sprintf("(L:%s R:%s)", formatFileSize(leftFile.Size), formatFileSize(rightFile.Size)))
		}

		// Truncate filename if too long
		displayPath := relPath
		baseWidth := 4 + len(sizeInfo) // Account for indicator and size
		if len(displayPath) > maxFileWidth-baseWidth {
			maxPathWidth := maxFileWidth - baseWidth - 3 // Account for "..."
			if maxPathWidth > 0 {
				displayPath = "..." + displayPath[len(displayPath)-maxPathWidth:]
			}
		}

		fileDisplay := fmt.Sprintf("%s %s %s", statusIndicator, displayPath, sizeInfo)

		if relPath == m.selectedFile {
			items = append(items, selectedFileStyle.Width(maxFileWidth).Render("► "+fileDisplay))
		} else {
			items = append(items, "  "+fileDisplay)
		}
	}

	content := strings.Join(items, "\n")

	// Add scroll indicator if needed
	if len(allFiles) > availableHeight {
		scrollInfo := fmt.Sprintf("\nShowing %d-%d of %d files", startIndex+1, endIndex, len(allFiles))
		content += helpStyle.Render(scrollInfo)
	}

	return fileListStyle.Width(m.windowWidth - 4).Height(availableHeight + 2).Render(content)
}

// formatFileSize formats file size in human readable format
func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1fGB", float64(size)/(1024*1024*1024))
	}
}

// renderDiffView renders the diff comparison view
func (m *Model) renderDiffView() string {
	if m.currentDiff == nil {
		return "No diff loaded"
	}

	var b strings.Builder

	// Header with file names - truncate if too long
	leftFile := m.currentDiff.LeftFile
	rightFile := m.currentDiff.RightFile
	maxFileNameWidth := (m.windowWidth - 4) / 2 // Split width for both names
	if len(leftFile) > maxFileNameWidth {
		leftFile = "..." + leftFile[len(leftFile)-maxFileNameWidth+3:]
	}
	if len(rightFile) > maxFileNameWidth {
		rightFile = "..." + rightFile[len(rightFile)-maxFileNameWidth+3:]
	}

	header := fmt.Sprintf("%s vs %s", leftFile, rightFile)
	b.WriteString(headerStyle.Width(m.windowWidth).Render(header))
	b.WriteString("\n\n")

	// Stats
	equal, inserted, deleted := m.currentDiff.GetStats()
	var stats string
	if m.windowWidth > 60 {
		stats = fmt.Sprintf("Lines: %d equal, %d inserted (+), %d deleted (-)", equal, inserted, deleted)
	} else {
		stats = fmt.Sprintf("%d equal, %d added, %d deleted", equal, inserted, deleted)
	}
	b.WriteString(helpStyle.Width(m.windowWidth).Render(stats))
	b.WriteString("\n\n")

	// Diff content
	b.WriteString(m.renderDiffContent())

	// Navigation help - adapt to width
	b.WriteString("\n")
	var helpText string
	if m.windowWidth > 80 {
		helpText = "↑↓/j/k: Navigate • g/G: Top/Bottom • n/p: Next/Prev file • m: Merge • Esc: Back • ?: Help • Q: Quit"
	} else if m.windowWidth > 60 {
		helpText = "↑↓/j/k: Navigate • g/G: Top/Bottom • n/p: Files • m: Merge • Esc: Back • ?: Help • Q: Quit"
	} else {
		helpText = "↑↓:Nav g/G:Top/Bot n/p:Files m:Merge Esc:Back ?:Help Q:Quit"
	}
	b.WriteString(helpStyle.Width(m.windowWidth).Render(helpText))

	return b.String()
}

// renderDiffContent renders the actual diff lines with syntax highlighting
func (m *Model) renderDiffContent() string {
	if m.currentDiff == nil || len(m.currentDiff.Lines) == 0 {
		return "No differences found"
	}

	var b strings.Builder
	maxVisible := m.windowHeight - 10 // Account for header, stats, and help text
	if maxVisible < 5 {
		maxVisible = 5
	}

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
		maxContentWidth := m.windowWidth - 12 // Account for prefix, line numbers, and padding
		if maxContentWidth < 20 {
			maxContentWidth = 20
		}
		if len(content) > maxContentWidth {
			content = content[:maxContentWidth-3] + "..."
		}

		var renderedLine string
		lineText := fmt.Sprintf("%s%s %s", prefix, lineNum, content)

		switch line.Type {
		case differ.DiffEqual:
			renderedLine = equalLineStyle.Width(m.windowWidth - 2).Render(lineText)
		case differ.DiffInsert:
			renderedLine = insertLineStyle.Width(m.windowWidth - 2).Render(fmt.Sprintf("%s%s +%s", prefix, lineNum, content))
		case differ.DiffDelete:
			renderedLine = deleteLineStyle.Width(m.windowWidth - 2).Render(fmt.Sprintf("%s%s -%s", prefix, lineNum, content))
		}

		b.WriteString(renderedLine)
		b.WriteString("\n")
	}

	// Show scroll indicator if needed
	if len(m.currentDiff.Lines) > maxVisible {
		var scrollInfo string
		if m.windowWidth > 50 {
			scrollInfo = fmt.Sprintf("Showing %d-%d of %d lines", start+1, end, len(m.currentDiff.Lines))
		} else {
			scrollInfo = fmt.Sprintf("%d-%d of %d", start+1, end, len(m.currentDiff.Lines))
		}
		b.WriteString(helpStyle.Width(m.windowWidth).Render(scrollInfo))
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
  Tab              Switch between left/right input fields (or cycle suggestions)
  Enter            Load the entered path / Accept selected suggestion
  ↑/↓              Navigate through common files list / Navigate suggestions
  Esc              Clear path suggestions
  Ctrl+D           Start comparing selected files
  ?                Show this help screen
  Q/Ctrl+C         Quit application

Path Suggestions:
  Type any path    Smart suggestions appear automatically
  ↑/↓              Navigate through suggestions
  Tab              Cycle to next suggestion
  Enter            Accept selected suggestion
  Esc              Cancel suggestions

Diff View Mode:
  ↑/↓ or j/k       Navigate through diff lines
  g                Go to top of diff
  G                Go to bottom of diff
  n                Next common file
  p                Previous common file
  m                Enter merge mode
  Esc              Return to file selection
  ?                Show this help screen
  Q/Ctrl+C         Quit application

Merge Mode:
  ↑/↓ or j/k       Navigate through diff lines
  Space/Enter      Toggle selection of current change
  t                Switch merge target (left/right)
  a                Select all changes
  n                Select no changes
  s                Save merged result to file
  Esc              Return to diff view
  ?                Show this help screen
  Q/Ctrl+C         Quit application

Color Legend:
  `

	b.WriteString(help)

	// Color examples
	b.WriteString(insertLineStyle.Render("Blue background: Added lines (+)"))
	b.WriteString("\n")
	b.WriteString(deleteLineStyle.Render("Red background: Deleted lines (-)"))
	b.WriteString("\n")
	b.WriteString(equalLineStyle.Render("Gray text: Unchanged lines"))
	b.WriteString("\n")
	b.WriteString(selectedChangeStyle.Render("Yellow background: Selected changes (merge mode)"))
	b.WriteString("\n")
	b.WriteString(unselectedChangeStyle.Render("Strikethrough: Unselected changes (merge mode)"))
	b.WriteString("\n\n")

	instructions := `Instructions:
1. Enter paths to files or directories in the input fields
   - Path suggestions appear automatically as you type
   - Use ↑/↓ or Tab to navigate through suggestions
   - Press Enter to accept a suggestion or Esc to dismiss
2. Press Enter to load each path
3. Use ↑/↓ to select which common file to compare
4. Press Ctrl+D to start comparing
5. Navigate through the diff using arrow keys or j/k
6. Use n/p to switch between different files
7. Press 'm' in diff view to enter merge mode
8. In merge mode, select which changes to apply and press 's' to save

Merge Workflow:
- Enter merge mode from diff view with 'm'
- Navigate with ↑/↓ or j/k through changes
- Toggle individual changes with Space/Enter
- Use 'a' to select all or 'n' to select none
- Switch target file with 't' (left or right)
- Save merged result with 's' - creates .merged file

Note: Only text files will be compared. The tool automatically
detects common text file extensions and filenames. Path suggestions
include files and directories from your current working directory.`

	b.WriteString(instructions)
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press Esc or ? to return to previous view"))

	return b.String()
}

// renderSuggestions renders path suggestions
func (m *Model) renderSuggestions(suggestions []string, selectedIndex int) string {
	if len(suggestions) == 0 {
		return ""
	}

	// Limit suggestions based on available space
	maxSuggestions := 6
	availableHeight := m.windowHeight - 20
	if availableHeight > 0 && availableHeight < maxSuggestions {
		maxSuggestions = availableHeight
	}
	if maxSuggestions < 3 {
		maxSuggestions = 3
	}

	displaySuggestions := suggestions
	if len(displaySuggestions) > maxSuggestions {
		displaySuggestions = displaySuggestions[:maxSuggestions]
	}

	maxSuggestionWidth := m.windowWidth - 10
	if maxSuggestionWidth < 30 {
		maxSuggestionWidth = 30
	}

	var suggestionLines []string
	for i, suggestion := range displaySuggestions {
		// Truncate suggestion if too long
		displaySuggestion := suggestion
		if len(displaySuggestion) > maxSuggestionWidth {
			displaySuggestion = "..." + displaySuggestion[len(displaySuggestion)-maxSuggestionWidth+3:]
		}

		if i == selectedIndex && selectedIndex < len(displaySuggestions) {
			suggestionLines = append(suggestionLines, selectedSuggestionStyle.Width(maxSuggestionWidth).Render("  → "+displaySuggestion))
		} else {
			suggestionLines = append(suggestionLines, "    "+displaySuggestion)
		}
	}

	content := strings.Join(suggestionLines, "\n")
	if len(suggestions) > maxSuggestions {
		content += helpStyle.Render(fmt.Sprintf("\n... and %d more", len(suggestions)-maxSuggestions))
	}
	return suggestionStyle.Width(m.windowWidth - 4).Render(content)
}

// renderMergeView renders the merge interface
func (m *Model) renderMergeView() string {
	if m.currentDiff == nil {
		return "No diff loaded for merging"
	}

	var b strings.Builder

	// Header
	header := fmt.Sprintf("Merge Mode - Target: %s", strings.ToUpper(m.mergeTarget))
	b.WriteString(mergeHeaderStyle.Width(m.windowWidth).Render(header))
	b.WriteString("\n\n")

	// Statistics
	if m.changeSelection != nil {
		selIns, totIns, selDel, totDel := m.changeSelection.GetSelectedStats(m.currentDiff)
		stats := fmt.Sprintf("Selected: %d/%d insertions, %d/%d deletions", selIns, totIns, selDel, totDel)
		b.WriteString(helpStyle.Width(m.windowWidth).Render(stats))
		b.WriteString("\n\n")
	}

	// Diff content with selection indicators
	b.WriteString(m.renderMergeContent())

	// Help text
	b.WriteString("\n")
	var helpText string
	if m.windowWidth > 80 {
		helpText = "Space/Enter: Toggle • t: Switch target • a: Select all • n: Select none • s: Save • Esc: Back • ?: Help"
	} else if m.windowWidth > 60 {
		helpText = "Space: Toggle • t: Target • a: All • n: None • s: Save • Esc: Back"
	} else {
		helpText = "Space:Toggle t:Target a:All n:None s:Save Esc:Back"
	}
	b.WriteString(helpStyle.Width(m.windowWidth).Render(helpText))

	return b.String()
}

// renderMergeContent renders the diff lines with selection indicators
func (m *Model) renderMergeContent() string {
	if m.currentDiff == nil || len(m.currentDiff.Lines) == 0 {
		return "No differences found"
	}

	var b strings.Builder
	maxVisible := m.windowHeight - 12 // Account for header, stats, and help text
	if maxVisible < 5 {
		maxVisible = 5
	}

	start := m.scrollOffset
	end := start + maxVisible
	if end > len(m.currentDiff.Lines) {
		end = len(m.currentDiff.Lines)
	}

	for i := start; i < end; i++ {
		line := m.currentDiff.Lines[i]
		prefix := "  "
		cursor := "  "

		if i == m.cursor {
			cursor = "▶ "
		}

		lineNum := fmt.Sprintf("%4d", line.LineNum)
		content := line.Content

		// Truncate long lines to fit screen width
		maxContentWidth := m.windowWidth - 16 // Account for prefix, cursor, line numbers, and padding
		if maxContentWidth < 20 {
			maxContentWidth = 20
		}
		if len(content) > maxContentWidth {
			content = content[:maxContentWidth-3] + "..."
		}

		var renderedLine string
		var selected bool

		switch line.Type {
		case differ.DiffEqual:
			lineText := fmt.Sprintf("%s%s%s %s", cursor, prefix, lineNum, content)
			renderedLine = equalLineStyle.Width(m.windowWidth - 2).Render(lineText)

		case differ.DiffInsert:
			if m.changeSelection != nil {
				selected = m.changeSelection.IsInsertionSelected(i)
			}

			if selected {
				prefix = "[✓]"
				lineText := fmt.Sprintf("%s%s%s +%s", cursor, prefix, lineNum, content)
				renderedLine = selectedChangeStyle.Width(m.windowWidth - 2).Render(lineText)
			} else {
				prefix = "[ ]"
				lineText := fmt.Sprintf("%s%s%s +%s", cursor, prefix, lineNum, content)
				renderedLine = unselectedChangeStyle.Width(m.windowWidth - 2).Render(lineText)
			}

		case differ.DiffDelete:
			if m.changeSelection != nil {
				selected = m.changeSelection.IsDeletionSelected(i)
			}

			if selected {
				prefix = "[✓]"
				lineText := fmt.Sprintf("%s%s%s -%s", cursor, prefix, lineNum, content)
				renderedLine = selectedChangeStyle.Width(m.windowWidth - 2).Render(lineText)
			} else {
				prefix = "[ ]"
				lineText := fmt.Sprintf("%s%s%s -%s", cursor, prefix, lineNum, content)
				renderedLine = unselectedChangeStyle.Width(m.windowWidth - 2).Render(lineText)
			}
		}

		b.WriteString(renderedLine)
		b.WriteString("\n")
	}

	// Show scroll indicator if needed
	if len(m.currentDiff.Lines) > maxVisible {
		var scrollInfo string
		if m.windowWidth > 50 {
			scrollInfo = fmt.Sprintf("Showing %d-%d of %d lines", start+1, end, len(m.currentDiff.Lines))
		} else {
			scrollInfo = fmt.Sprintf("%d-%d of %d", start+1, end, len(m.currentDiff.Lines))
		}
		b.WriteString(helpStyle.Width(m.windowWidth).Render(scrollInfo))
		b.WriteString("\n")
	}

	return b.String()
}
