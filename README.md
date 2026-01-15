# File Comparison TUI Tool

A powerful terminal-based user interface (TUI) application for visually comparing files and directories with beautiful color-coded diff highlighting. Built with Go using Bubble Tea framework.

## ✨ Features

- 🎨 **Visual Diff Highlighting**: Green backgrounds for additions, blue backgrounds for deletions
- 👥 **Side-by-Side View**: Compare files in unified or side-by-side layout modes
- 🔀 **Interactive Merge Mode**: Cherry-pick and apply specific changes between files
- 📋 **File Copy Mode**: Easily copy unique files between directories with selective control
- 📁 **Complete Directory Analysis**: Shows ALL files from both directories (common and unique)
- 🏷️ **Source Identification**: Clear indicators for files that exist in only one directory
- ⌨️ **Intuitive Controls**: Vim-like navigation (j/k) with full arrow key support
- 🔍 **Intelligent File Detection**: Automatically identifies 50+ text file types for comparison
- 📊 **Real-time Diff Statistics**: Live counts of equal, added, and deleted lines
- 🖥️ **Full Screen TUI**: Clean, distraction-free interface with proper scrolling
- 🚀 **Fast Performance**: Efficient diff algorithm with semantic cleanup
- 📱 **Responsive Design**: Adapts to terminal window size changes
- 🎯 **Multi-file Navigation**: Easy switching between multiple file comparisons
- 💾 **Selective Merging**: Save merged results with only the changes you want
- 🔄 **Directory Synchronization**: Copy unique files between directories for easy sync

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd golang-fileCmp

# Build the application
go build -o filecmp

# Or run directly
go run main.go
```

## Usage

### Command Line Options

```bash
# Start with interactive file selection
./filecmp

# Compare two files directly
./filecmp file1.txt file2.txt

# Compare two directories (finds common files automatically)  
./filecmp ./project-v1 ./project-v2

# Load left file, enter right path in TUI
./filecmp /path/to/file1.txt

# Show comprehensive help
./filecmp --help
```

### Quick Start with Make

```bash
# Install dependencies and build
make deps && make build

# Try the demo with example files
# Try with examples
make demo

# Demo the new merge functionality
./demo-merge.sh

# Quick test with sample files
make run-files

# Quick test with sample directories  
make run-dirs
```

### Interactive Controls

#### File Selection Mode
- **Tab**: Switch between left and right input fields
- **Enter**: Load the entered path (file or directory)
- **↑/↓**: Navigate through the list of all files (common and unique)
- **Ctrl+D**: Start comparing the selected files
- **c**: Enter copy mode (for directories with unique files)
- **?**: Show help screen
- **Q/Ctrl+C**: Quit application

#### Diff View Mode
- **↑/↓** or **j/k**: Navigate through diff lines
- **s**: Switch view mode (Unified ↔ Side-by-Side)
- **g**: Go to top of diff
- **G**: Go to bottom of diff
- **n**: Next common file
- **p**: Previous common file
- **m**: Enter merge mode
- **Esc**: Return to file selection
- **?**: Show help screen
- **Q/Ctrl+C**: Quit application

#### Side-by-Side View Mode
- **↑/↓** or **j/k**: Navigate through diff lines
- **h/l** or **←/→**: Visual focus left/right (for reference)
- **s**: Switch to Unified view mode
- **g**: Go to top of diff
- **G**: Go to bottom of diff
- **n**: Next common file
- **p**: Previous common file
- **m**: Enter merge mode
- **Esc**: Return to file selection
- **?**: Show help screen
- **Q/Ctrl+C**: Quit application

#### Merge Mode
- **↑/↓** or **j/k**: Navigate through diff lines
- **Space/Enter**: Toggle selection of current change
- **t**: Switch merge target (left/right file)
- **a**: Select all changes
- **n**: Select no changes
- **s**: Save merged result to file
- **Esc**: Return to diff view
- **?**: Show help screen
- **Q/Ctrl+C**: Quit application

#### Copy Mode (Directory Comparison Only)
- **↑/↓** or **j/k**: Navigate through unique files
- **Space/Enter**: Toggle selection of current file to copy
- **t**: Switch copy target (to-left ↔ to-right)
- **a**: Select all unique files
- **n**: Select no files
- **s**: Copy selected files to target directory
- **Esc**: Return to file selection
- **?**: Show help screen
- **Q/Ctrl+C**: Quit application

## Color Legend

### Diff Colors
- **Green background**: Added lines (+)
- **Blue background**: Deleted lines (-)
- **Gray text**: Unchanged lines
- **Yellow background**: Selected changes (merge mode)
- **Strikethrough text**: Unselected changes (merge mode)

### File Status Indicators
- **✓ Green checkmark**: Identical files (same content in both directories)
- **✗ Red X**: Different files (content differs between directories)
- **◄ Blue arrow**: File exists only in LEFT directory **[LEFT ONLY]**
- **► Orange arrow**: File exists only in RIGHT directory **[RIGHT ONLY]**

## 📄 Supported File Types

The tool intelligently detects and compares 50+ text file types:

### Programming Languages
`.go` `.js` `.ts` `.py` `.java` `.c` `.cpp` `.h` `.hpp` `.rs` `.swift` `.kt` `.cs` `.vb` `.fs` `.php` `.rb` `.scala` `.clj` `.hs` `.elm` `.pl` `.pm` `.r` `.R` `.m`

### Web & Data Files  
`.html` `.htm` `.css` `.xml` `.json` `.yaml` `.yml` `.toml` `.csv`

### Documentation & Text
`.txt` `.md` `.rst` `.tex` `.org`

### Configuration & Scripts
`.conf` `.ini` `.cfg` `.sh` `.bash` `.zsh` `.fish` `.ps1` `.bat` `.cmd`

### Build & Project Files
`.dockerfile` `.makefile` `.cmake` `.ninja` `.gradle` `.pom` `.gitignore` `.gitattributes` `.editorconfig`

### Special Files (no extension)
`README` `LICENSE` `CHANGELOG` `AUTHORS` `CONTRIBUTORS` `Dockerfile` `Gemfile` `Rakefile` `Procfile`

## How It Works

1. **Load Paths**: Enter file or directory paths in the input fields
2. **Analyze All Files**: For directories, the tool finds ALL text files from both locations
3. **Categorize Files**: Files are marked as common (both sides), left-only, or right-only
4. **Select File**: Choose any file to compare using the arrow keys
5. **View Diff**: See the comparison with color-coded changes
6. **Navigate**: Move through all files (common and unique) seamlessly
7. **Merge Changes**: Press 'm' to enter merge mode (only for common files)
8. **Save Results**: Choose which changes to keep and save the merged file

## Examples

### Comparing Two Files
```bash
./filecmp config.old.yaml config.new.yaml
```

### Comparing Project Directories
```bash
./filecmp ./project-v1 ./project-v2
```
This will find all files from both directories and allow you to compare them, merge changes, and copy unique files.

### Interactive Mode
```bash
./filecmp
```
Then enter paths interactively using the TUI.

### Copy Unique Files Between Directories
```bash
# Compare directories and copy unique files
./filecmp old-project/ new-project/
# Press 'c' to enter copy mode
# Select files to copy and press 's' to copy them
```

## 🏗️ Architecture

### Core Components
- **TUI Layer**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for responsive terminal interface
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) for beautiful colors and layouts  
- **Diff Engine**: [go-diff](https://github.com/sergi/go-diff) with semantic cleanup for accurate comparisons
- **File System**: Smart file detection and recursive directory traversal

### Project Structure
```
internal/
├── ui/         # TUI components and views
├── differ/     # Diff computation engine  
├── merge/      # Merge functionality and change selection
├── file/       # File operations and type detection
└── ...

examples/       # Sample files for testing
```

## 🎛️ Advanced Usage

### Keyboard Shortcuts Summary
| Key | File Selection | Diff View | Side-by-Side | Merge Mode | Description |
|-----|----------------|-----------|--------------|------------|-------------|
| `Tab` | ✅ | ❌ | ❌ | ❌ | Switch input fields |
| `Enter` | ✅ | ❌ | ❌ | ✅ | Load entered path / Toggle change |
| `↑/↓` | ✅ | ✅ | ✅ | ✅ | Navigate lists/lines |
| `j/k` | ❌ | ✅ | ✅ | ✅ | Vim-style navigation |
| `h/l` | ❌ | ❌ | ✅ | ❌ | Visual left/right focus |
| `g/G` | ❌ | ✅ | ✅ | ❌ | Jump to top/bottom |
| `s` | ❌ | ✅ | ✅ | ✅ | Switch view mode / Save result |
| `n/p` | ✅ | ✅ | ✅ | ❌ | Next/previous file |
| `m` | ❌ | ✅ | ✅ | ❌ | Enter merge mode |
| `t` | ❌ | ❌ | ❌ | ✅ | Switch merge target |
| `a` | ❌ | ❌ | ❌ | ✅ | Select all changes |
| `Space` | ❌ | ❌ | ❌ | ✅ | Toggle current change |
| `Ctrl+D` | ✅ | ❌ | ❌ | ❌ | Start comparison |
| `Esc` | ❌ | ✅ | ✅ | ✅ | Return to previous view |
| `?` | ✅ | ✅ | ✅ | ✅ | Show help screen |
| `Q/Ctrl+C` | ✅ | ✅ | ✅ | ✅ | Quit application |

### Performance Tips
- Large files (>10MB) may take a moment to process
- Directory comparisons are optimized to only load text files
- The diff algorithm uses semantic cleanup for better readability

## License

MIT License - see LICENSE file for details.

## 🧪 Testing

Run the included test suite:
```bash
# Test the diff engine
make test-diff

# Run Go unit tests  
make test

# Try with examples
make demo
```

## 🚀 Building from Source

```bash
# Clone and build
git clone <repository-url>
cd golang-fileCmp
make deps
make build

# Or build for multiple platforms
make build-all

# Install system-wide (optional)
make install
```

## 🤝 Contributing

Contributions are welcome! Areas for improvement:
- Additional file type support
- Syntax highlighting within diffs  
- Side-by-side view mode
- Advanced merge conflict resolution
- Copy operation undo/rollback
- Export diff results
- Configuration file support
- Undo/redo for merge operations
- Directory structure visualization
- File filtering and search capabilities

Please feel free to submit issues, feature requests, or pull requests.

## 📋 Roadmap

- [x] Interactive merge mode with selective change application
- [x] Complete directory analysis (all files, not just common ones)
- [x] File source identification with clear indicators
- [x] Copy mode for easily copying unique files between directories
- [x] Side-by-side comparison view with unified/split toggle
- [ ] Syntax highlighting for code diffs
- [ ] Three-way merge support
- [ ] Merge conflict resolution
- [ ] Export diffs to HTML/PDF
- [ ] Configuration file support
- [ ] Plugin system for custom file types
- [ ] Integration with Git for commit diffs
- [ ] Undo/redo functionality in merge mode
- [ ] File filtering and search capabilities
- [ ] Directory tree visualization