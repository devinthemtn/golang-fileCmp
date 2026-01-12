package main

import (
	"fmt"
	"log"
	"os"

	"golang-fileCmp/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the model
	model := ui.New()

	// Handle command line arguments
	args := os.Args[1:]
	if len(args) >= 1 {
		model.SetLeftPath(args[0])
	}
	if len(args) >= 2 {
		model.SetRightPath(args[1])
	}

	// Show usage if help is requested
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		showUsage()
		return
	}

	// Create the program
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func showUsage() {
	fmt.Printf(`File Comparison TUI Tool

Usage:
  %s [left_path] [right_path]

Arguments:
  left_path   Path to left file or directory (optional)
  right_path  Path to right file or directory (optional)

Examples:
  %s                           # Start with empty inputs
  %s file1.txt file2.txt       # Compare two files
  %s ./dir1 ./dir2             # Compare two directories
  %s /path/to/file             # Load left file, input right path in TUI

Interactive Controls:
  Tab              Switch between input fields
  Enter            Load entered path
  ↑/↓              Navigate file list
  Ctrl+D           Start diff comparison
  j/k              Navigate diff (vim-style)
  n/p              Next/previous file
  g/G              Go to top/bottom
  Esc              Go back
  ?                Show help
  Q/Ctrl+C         Quit

The tool will automatically find common files between directories
and highlight differences with colors:
  - Green background: Added lines
  - Blue background:  Deleted lines
  - Gray text:        Unchanged lines

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
