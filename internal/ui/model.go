package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang-fileCmp/internal/differ"
	"golang-fileCmp/internal/file"
	"golang-fileCmp/internal/merge"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewMode represents the current view state
type ViewMode int

const (
	ViewModeFileSelect ViewMode = iota
	ViewModeDiff
	ViewModeMerge
	ViewModeHelp
)

// Model represents the main application state
type Model struct {
	// Application state
	viewMode     ViewMode
	windowWidth  int
	windowHeight int

	// File selection
	leftPath       string
	rightPath      string
	leftFile       *file.FileInfo
	rightFile      *file.FileInfo
	commonFiles    map[string][2]*file.FileInfo
	allFiles       map[string]*file.FileComparison
	selectedFile   string
	fileListScroll int

	// Diff view
	currentDiff  *differ.FileDiff
	scrollOffset int
	cursor       int

	// Merge view
	changeSelection *merge.ChangeSelection
	mergeTarget     string // "left" or "right"
	mergePreview    string

	// Services
	fileManager *file.Manager
	differ      *differ.Differ
	merger      *merge.Merger

	// UI state
	inputLeft   string
	inputRight  string
	focusLeft   bool
	showingHelp bool
	errorMsg    string

	// Path suggestions
	leftSuggestions  []string
	rightSuggestions []string
	leftSuggIndex    int
	rightSuggIndex   int
	showSuggestions  bool
}

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#555555")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	focusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFFF00")).
				Padding(0, 1)

	equalLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	insertLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0000FF")).
			Bold(true)

	deleteLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF0000")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF0000")).
			Padding(0, 1).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	fileListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	selectedFileStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	suggestionStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#888888")).
			Padding(0, 1).
			MaxWidth(80)

	selectedSuggestionStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#DDDDDD")).
				Foreground(lipgloss.Color("#000000"))

	// Merge mode styles
	mergeHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#00AA00")).
				Padding(0, 1).
				Bold(true)

	selectedChangeStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#FFFF00")).
				Foreground(lipgloss.Color("#000000")).
				Bold(true)

	unselectedChangeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666")).
				Strikethrough(true)

	previewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00AA00")).
			Padding(1)
)

// New creates a new model
func New() *Model {
	return &Model{
		viewMode:        ViewModeFileSelect,
		fileManager:     file.New(),
		differ:          differ.New(),
		merger:          merge.New(),
		focusLeft:       true,
		commonFiles:     make(map[string][2]*file.FileInfo),
		allFiles:        make(map[string]*file.FileComparison),
		leftSuggIndex:   -1,
		rightSuggIndex:  -1,
		showSuggestions: false,
		fileListScroll:  0,
		mergeTarget:     "left",
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// View renders the current view
func (m *Model) View() string {
	switch m.viewMode {
	case ViewModeFileSelect:
		return m.renderFileSelectView()
	case ViewModeDiff:
		return m.renderDiffView()
	case ViewModeMerge:
		return m.renderMergeView()
	case ViewModeHelp:
		return m.renderHelpView()
	default:
		return "Unknown view mode"
	}
}

// handleKeyPress processes keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.viewMode {
	case ViewModeFileSelect:
		return m.handleFileSelectKeys(msg)
	case ViewModeDiff:
		return m.handleDiffKeys(msg)
	case ViewModeMerge:
		return m.handleMergeKeys(msg)
	case ViewModeHelp:
		return m.handleHelpKeys(msg)
	}
	return m, nil
}

// handleFileSelectKeys handles keys in file selection mode
func (m *Model) handleFileSelectKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		// If suggestions are showing and we have suggestions, cycle through them
		if m.showSuggestions {
			if m.focusLeft && len(m.leftSuggestions) > 0 {
				m.leftSuggIndex = (m.leftSuggIndex + 1) % len(m.leftSuggestions)
				return m, nil
			} else if !m.focusLeft && len(m.rightSuggestions) > 0 {
				m.rightSuggIndex = (m.rightSuggIndex + 1) % len(m.rightSuggestions)
				return m, nil
			}
		}
		m.focusLeft = !m.focusLeft
		m.clearSuggestions()
		return m, nil

	case "enter":
		// If suggestions are showing, accept the selected suggestion
		if m.showSuggestions {
			if m.focusLeft && len(m.leftSuggestions) > 0 && m.leftSuggIndex >= 0 {
				m.inputLeft = m.leftSuggestions[m.leftSuggIndex]
				m.clearSuggestions()
				return m, nil
			} else if !m.focusLeft && len(m.rightSuggestions) > 0 && m.rightSuggIndex >= 0 {
				m.inputRight = m.rightSuggestions[m.rightSuggIndex]
				m.clearSuggestions()
				return m, nil
			}
		}

		// Load the path
		if m.focusLeft {
			if m.inputLeft != "" {
				m.leftPath = m.inputLeft
				m.loadLeftPath()
			}
		} else {
			if m.inputRight != "" {
				m.rightPath = m.inputRight
				m.loadRightPath()
			}
		}
		m.clearSuggestions()
		m.updateCommonFiles()
		return m, nil

	case "ctrl+d":
		if len(m.allFiles) > 0 {
			// If no file is selected, select the first one
			if m.selectedFile == "" {
				files := m.getSortedFiles()
				if len(files) > 0 {
					m.selectedFile = files[0]
				}
			}
			// Load diff for the currently selected file
			if m.selectedFile != "" {
				m.loadDiff()
				m.viewMode = ViewModeDiff
			}
		}
		return m, nil

	case "?":
		m.viewMode = ViewModeHelp
		return m, nil

	case "backspace":
		if m.focusLeft {
			if len(m.inputLeft) > 0 {
				m.inputLeft = m.inputLeft[:len(m.inputLeft)-1]
				m.updateSuggestions()
			}
		} else {
			if len(m.inputRight) > 0 {
				m.inputRight = m.inputRight[:len(m.inputRight)-1]
				m.updateSuggestions()
			}
		}
		return m, nil

	case "esc":
		m.clearSuggestions()
		return m, nil

	case "up":
		if m.showSuggestions {
			if m.focusLeft && len(m.leftSuggestions) > 0 {
				if m.leftSuggIndex <= 0 {
					m.leftSuggIndex = len(m.leftSuggestions) - 1
				} else {
					m.leftSuggIndex--
				}
				return m, nil
			} else if !m.focusLeft && len(m.rightSuggestions) > 0 {
				if m.rightSuggIndex <= 0 {
					m.rightSuggIndex = len(m.rightSuggestions) - 1
				} else {
					m.rightSuggIndex--
				}
				return m, nil
			}
		}
		if len(m.allFiles) > 0 {
			m.selectPreviousFile()
			// Only auto-load diff if we're in file selection mode
			// In diff mode, user needs to press Ctrl+D or navigate with n/p
		}
		return m, nil

	case "down":
		if m.showSuggestions {
			if m.focusLeft && len(m.leftSuggestions) > 0 {
				m.leftSuggIndex = (m.leftSuggIndex + 1) % len(m.leftSuggestions)
				return m, nil
			} else if !m.focusLeft && len(m.rightSuggestions) > 0 {
				m.rightSuggIndex = (m.rightSuggIndex + 1) % len(m.rightSuggestions)
				return m, nil
			}
		}
		if len(m.allFiles) > 0 {
			m.selectNextFile()
			// Only auto-load diff if we're in file selection mode
			// In diff mode, user needs to press Ctrl+D or navigate with n/p
		}
		return m, nil

	default:
		// Add character to appropriate input
		if len(msg.String()) == 1 {
			if m.focusLeft {
				m.inputLeft += msg.String()
			} else {
				m.inputRight += msg.String()
			}
			m.updateSuggestions()
		}
		return m, nil
	}
}

// handleDiffKeys handles keys in diff view mode
func (m *Model) handleDiffKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		m.viewMode = ViewModeFileSelect
		return m, nil

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.scrollOffset {
				m.scrollOffset = m.cursor
			}
		}
		return m, nil

	case "down", "j":
		if m.currentDiff != nil && m.cursor < len(m.currentDiff.Lines)-1 {
			m.cursor++
			maxVisible := m.windowHeight - 10 // Account for header and footer
			if m.cursor >= m.scrollOffset+maxVisible {
				m.scrollOffset = m.cursor - maxVisible + 1
			}
		}
		return m, nil

	case "g":
		m.cursor = 0
		m.scrollOffset = 0
		return m, nil

	case "G":
		if m.currentDiff != nil {
			m.cursor = len(m.currentDiff.Lines) - 1
			maxVisible := m.windowHeight - 10
			m.scrollOffset = max(0, m.cursor-maxVisible+1)
		}
		return m, nil

	case "n":
		m.selectNextFile()
		m.loadDiff()
		return m, nil

	case "p":
		m.selectPreviousFile()
		m.loadDiff()
		return m, nil

	case "?":
		m.viewMode = ViewModeHelp
		return m, nil

	case "m":
		// Enter merge mode if we have a diff loaded
		if m.currentDiff != nil {
			// Check if this is a valid file for merging
			if m.selectedFile != "" {
				if fileComparison, exists := m.allFiles[m.selectedFile]; exists {
					if fileComparison.Source == file.SourceBoth {
						m.initializeMergeMode()
						m.viewMode = ViewModeMerge
					} else {
						// Cannot merge files that only exist on one side
						if fileComparison.Source == file.SourceLeft {
							m.errorMsg = "Cannot merge: File exists only in LEFT directory"
						} else {
							m.errorMsg = "Cannot merge: File exists only in RIGHT directory"
						}
					}
				}
			}
		}
		return m, nil
	}

	return m, nil
}

// handleMergeKeys handles keys in merge view mode
func (m *Model) handleMergeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		m.viewMode = ViewModeDiff
		return m, nil

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.scrollOffset {
				m.scrollOffset = m.cursor
			}
		}
		return m, nil

	case "down", "j":
		if m.currentDiff != nil && m.cursor < len(m.currentDiff.Lines)-1 {
			m.cursor++
			maxVisible := m.windowHeight - 15 // Account for header and footer
			if m.cursor >= m.scrollOffset+maxVisible {
				m.scrollOffset = m.cursor - maxVisible + 1
			}
		}
		return m, nil

	case " ", "enter":
		// Toggle selection of current change
		if m.currentDiff != nil && m.cursor < len(m.currentDiff.Lines) {
			line := m.currentDiff.Lines[m.cursor]
			switch line.Type {
			case differ.DiffInsert:
				m.changeSelection.ToggleInsertion(m.cursor)
			case differ.DiffDelete:
				m.changeSelection.ToggleDeletion(m.cursor)
			}
			m.updateMergePreview()
		}
		return m, nil

	case "a":
		// Select all changes
		m.changeSelection.SelectAll(m.currentDiff)
		m.updateMergePreview()
		return m, nil

	case "n":
		// Select no changes
		m.changeSelection.SelectNone(m.currentDiff)
		m.updateMergePreview()
		return m, nil

	case "t":
		// Toggle merge target (left/right)
		if m.mergeTarget == "left" {
			m.mergeTarget = "right"
		} else {
			m.mergeTarget = "left"
		}
		m.updateMergePreview()
		return m, nil

	case "s":
		// Save merged result
		return m, m.saveMergedFile()

	case "?":
		m.viewMode = ViewModeHelp
		return m, nil
	}

	return m, nil
}

// handleHelpKeys handles keys in help view mode
func (m *Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "?":
		if m.viewMode == ViewModeMerge {
			m.viewMode = ViewModeMerge
		} else if len(m.commonFiles) > 0 && m.currentDiff != nil {
			m.viewMode = ViewModeDiff
		} else {
			m.viewMode = ViewModeFileSelect
		}
		return m, nil
	}
	return m, nil
}

// Helper functions

func (m *Model) loadLeftPath() {
	leftFile, err := m.fileManager.LoadPath(m.leftPath)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Error loading left path: %s", err.Error())
		return
	}
	m.leftFile = leftFile
	m.errorMsg = ""
}

func (m *Model) loadRightPath() {
	rightFile, err := m.fileManager.LoadPath(m.rightPath)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Error loading right path: %s", err.Error())
		return
	}
	m.rightFile = rightFile
	m.errorMsg = ""
}

func (m *Model) updateCommonFiles() {
	if m.leftFile != nil && m.rightFile != nil {
		m.commonFiles = file.FindCommonFiles(m.leftFile, m.rightFile)
		m.allFiles = file.FindAllFiles(m.leftFile, m.rightFile)
		m.fileListScroll = 0 // Reset scroll when files change

		// Select first file by default if none selected
		if len(m.allFiles) > 0 && m.selectedFile == "" {
			files := m.getSortedFiles()
			if len(files) > 0 {
				m.selectedFile = files[0]
			}
		}

		// Ensure selected file still exists in the new file list
		if m.selectedFile != "" {
			if _, exists := m.allFiles[m.selectedFile]; !exists {
				files := m.getSortedFiles()
				if len(files) > 0 {
					m.selectedFile = files[0]
				} else {
					m.selectedFile = ""
				}
			}
		}
	}
}

func (m *Model) getSortedFiles() []string {
	files := make([]string, 0, len(m.allFiles))
	for relPath := range m.allFiles {
		files = append(files, relPath)
	}

	// Sort files alphabetically for consistent ordering
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i] > files[j] {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	return files
}

func (m *Model) selectNextFile() {
	if len(m.allFiles) == 0 {
		return
	}

	files := m.getSortedFiles()

	currentIndex := -1
	for i, f := range files {
		if f == m.selectedFile {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 || currentIndex == len(files)-1 {
		m.selectedFile = files[0]
	} else {
		m.selectedFile = files[currentIndex+1]
	}
}

func (m *Model) selectPreviousFile() {
	if len(m.allFiles) == 0 {
		return
	}

	files := m.getSortedFiles()

	currentIndex := -1
	for i, f := range files {
		if f == m.selectedFile {
			currentIndex = i
			break
		}
	}

	if currentIndex <= 0 {
		m.selectedFile = files[len(files)-1]
	} else {
		m.selectedFile = files[currentIndex-1]
	}
}

func (m *Model) loadDiff() {
	if m.selectedFile == "" {
		m.errorMsg = "No file selected for comparison"
		return
	}

	fileComparison, exists := m.allFiles[m.selectedFile]
	if !exists {
		m.errorMsg = fmt.Sprintf("Selected file '%s' no longer exists", m.selectedFile)
		return
	}

	var leftContent, rightContent, leftPath, rightPath string

	switch fileComparison.Source {
	case file.SourceBoth:
		leftFile := fileComparison.LeftFile
		rightFile := fileComparison.RightFile
		leftContent = leftFile.Content
		rightContent = rightFile.Content
		leftPath = leftFile.Path
		rightPath = rightFile.Path
	case file.SourceLeft:
		leftFile := fileComparison.LeftFile
		leftContent = leftFile.Content
		rightContent = ""
		leftPath = leftFile.Path
		rightPath = "<file not found>"
	case file.SourceRight:
		rightFile := fileComparison.RightFile
		leftContent = ""
		rightContent = rightFile.Content
		leftPath = "<file not found>"
		rightPath = rightFile.Path
	}

	diff := m.differ.CompareStrings(
		leftPath,
		rightPath,
		leftContent,
		rightContent,
	)

	m.currentDiff = diff
	m.cursor = 0
	m.scrollOffset = 0
	m.errorMsg = "" // Clear any previous errors
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SetLeftPath sets the left path and loads it
func (m *Model) SetLeftPath(path string) {
	m.inputLeft = path
	m.leftPath = path
	m.loadLeftPath()

	// If we already have a right file, update files and select first
	if m.rightFile != nil {
		m.updateCommonFiles()
		if len(m.allFiles) > 0 && m.selectedFile == "" {
			files := m.getSortedFiles()
			if len(files) > 0 {
				m.selectedFile = files[0]
			}
		}
	}
}

// SetRightPath sets the right path and loads it
func (m *Model) SetRightPath(path string) {
	m.inputRight = path
	m.rightPath = path
	m.loadRightPath()
	m.updateCommonFiles()

	// Ensure first file is selected when both paths are loaded
	if m.leftFile != nil && m.rightFile != nil && len(m.allFiles) > 0 && m.selectedFile == "" {
		files := m.getSortedFiles()
		if len(files) > 0 {
			m.selectedFile = files[0]
		}
	}
}

// Path suggestion methods

func (m *Model) updateSuggestions() {
	if m.focusLeft {
		m.leftSuggestions = m.generateSuggestions(m.inputLeft)
		m.leftSuggIndex = 0
		if len(m.leftSuggestions) == 0 {
			m.leftSuggIndex = -1
		}
	} else {
		m.rightSuggestions = m.generateSuggestions(m.inputRight)
		m.rightSuggIndex = 0
		if len(m.rightSuggestions) == 0 {
			m.rightSuggIndex = -1
		}
	}

	m.showSuggestions = len(m.leftSuggestions) > 0 || len(m.rightSuggestions) > 0
}

func (m *Model) clearSuggestions() {
	m.leftSuggestions = nil
	m.rightSuggestions = nil
	m.leftSuggIndex = -1
	m.rightSuggIndex = -1
	m.showSuggestions = false
}

func (m *Model) generateSuggestions(input string) []string {
	if len(input) == 0 {
		return nil
	}

	var suggestions []string

	// Determine the directory to search and the prefix to match
	var searchDir, prefix string

	if strings.HasSuffix(input, "/") || strings.HasSuffix(input, "\\") {
		// Input ends with separator, search in that directory
		searchDir = input
		prefix = ""
	} else {
		// Input is a partial path, split into directory and filename parts
		searchDir = filepath.Dir(input)
		prefix = filepath.Base(input)

		if searchDir == "." && !strings.Contains(input, "/") && !strings.Contains(input, "\\") {
			searchDir = ""
		}
	}

	// Handle empty or current directory
	if searchDir == "" || searchDir == "." {
		searchDir = "."
	}

	// Read directory contents
	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return nil
	}

	// Filter and collect matching entries
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless explicitly requested
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(prefix, ".") {
			continue
		}

		// Check if the name matches the prefix
		if prefix == "" || strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			var suggestion string
			if searchDir == "." {
				suggestion = name
			} else {
				suggestion = filepath.Join(searchDir, name)
			}

			// Add trailing slash for directories
			if entry.IsDir() {
				suggestion += string(filepath.Separator)
			}

			suggestions = append(suggestions, suggestion)
		}
	}

	// Sort suggestions
	sort.Strings(suggestions)

	// Limit number of suggestions
	maxSuggestions := 8
	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}

	return suggestions
}

// initializeMergeMode sets up merge mode with default selections
func (m *Model) initializeMergeMode() {
	if m.currentDiff == nil {
		return
	}

	// Only initialize merge mode for files that exist on both sides
	if m.selectedFile != "" {
		if fileComparison, exists := m.allFiles[m.selectedFile]; exists && fileComparison.Source == file.SourceBoth {
			m.changeSelection = merge.NewChangeSelection(m.currentDiff)
			m.updateMergePreview()
			m.cursor = 0
			m.scrollOffset = 0
			m.errorMsg = "" // Clear any previous error
		}
	}
}

// updateMergePreview updates the merge preview text
func (m *Model) updateMergePreview() {
	if m.currentDiff == nil || m.changeSelection == nil {
		return
	}

	m.mergePreview = m.merger.CreateMergePreview(m.currentDiff, m.changeSelection, m.mergeTarget)
}

// saveMergedFile saves the merged result to a file
func (m *Model) saveMergedFile() tea.Cmd {
	return func() tea.Msg {
		if m.currentDiff == nil || m.changeSelection == nil {
			return nil
		}

		var result *merge.MergeResult
		var targetPath string

		if m.mergeTarget == "left" {
			result = m.merger.ApplyToLeft(m.currentDiff, m.changeSelection)
			targetPath = m.currentDiff.LeftFile + ".merged"
		} else {
			result = m.merger.ApplyToRight(m.currentDiff, m.changeSelection)
			targetPath = m.currentDiff.RightFile + ".merged"
		}

		err := os.WriteFile(targetPath, []byte(result.Content), 0644)
		if err != nil {
			return fmt.Sprintf("Error saving merged file: %s", err.Error())
		}

		return fmt.Sprintf("Saved merged result to %s (%d changes applied, %d skipped)",
			targetPath, result.Applied, result.Skipped)
	}
}
