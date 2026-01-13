package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang-fileCmp/internal/differ"
	"golang-fileCmp/internal/file"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewMode represents the current view state
type ViewMode int

const (
	ViewModeFileSelect ViewMode = iota
	ViewModeDiff
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
	selectedFile   string
	fileListScroll int

	// Diff view
	currentDiff  *differ.FileDiff
	scrollOffset int
	cursor       int

	// Services
	fileManager *file.Manager
	differ      *differ.Differ

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
)

// New creates a new model
func New() *Model {
	return &Model{
		viewMode:        ViewModeFileSelect,
		fileManager:     file.New(),
		differ:          differ.New(),
		focusLeft:       true,
		commonFiles:     make(map[string][2]*file.FileInfo),
		leftSuggIndex:   -1,
		rightSuggIndex:  -1,
		showSuggestions: false,
		fileListScroll:  0,
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
		if len(m.commonFiles) > 0 {
			// Select first common file by default
			for relPath := range m.commonFiles {
				m.selectedFile = relPath
				break
			}
			m.loadDiff()
			m.viewMode = ViewModeDiff
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
		if len(m.commonFiles) > 0 {
			m.selectPreviousFile()
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
		if len(m.commonFiles) > 0 {
			m.selectNextFile()
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
	}

	return m, nil
}

// handleHelpKeys handles keys in help view mode
func (m *Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "?":
		if len(m.commonFiles) > 0 && m.currentDiff != nil {
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
		m.fileListScroll = 0 // Reset scroll when files change

		// Select first file by default if none selected
		if len(m.commonFiles) > 0 && m.selectedFile == "" {
			files := m.getSortedFiles()
			if len(files) > 0 {
				m.selectedFile = files[0]
			}
		}

		// Ensure selected file still exists in the new common files
		if m.selectedFile != "" {
			if _, exists := m.commonFiles[m.selectedFile]; !exists {
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
	files := make([]string, 0, len(m.commonFiles))
	for relPath := range m.commonFiles {
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
	if len(m.commonFiles) == 0 {
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
	if len(m.commonFiles) == 0 {
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
		return
	}

	filePair, exists := m.commonFiles[m.selectedFile]
	if !exists {
		return
	}

	leftFile := filePair[0]
	rightFile := filePair[1]

	diff := m.differ.CompareStrings(
		leftFile.Path,
		rightFile.Path,
		leftFile.Content,
		rightFile.Content,
	)

	m.currentDiff = diff
	m.cursor = 0
	m.scrollOffset = 0
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
}

// SetRightPath sets the right path and loads it
func (m *Model) SetRightPath(path string) {
	m.inputRight = path
	m.rightPath = path
	m.loadRightPath()
	m.updateCommonFiles()
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
