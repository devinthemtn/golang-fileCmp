# File Comparison TUI Tool

A powerful terminal-based user interface (TUI) application for visually comparing files and directories with beautiful color-coded diff highlighting. Built with Go using Bubble Tea framework.

## ✨ Features

- 🎨 **Visual Diff Highlighting**: Green backgrounds for additions, blue backgrounds for deletions
- 📁 **Smart Directory Comparison**: Automatically finds and compares common files between directories
- ⌨️ **Intuitive Controls**: Vim-like navigation (j/k) with full arrow key support
- 🔍 **Intelligent File Detection**: Automatically identifies 50+ text file types for comparison
- 📊 **Real-time Diff Statistics**: Live counts of equal, added, and deleted lines
- 🖥️ **Full Screen TUI**: Clean, distraction-free interface with proper scrolling
- 🚀 **Fast Performance**: Efficient diff algorithm with semantic cleanup
- 📱 **Responsive Design**: Adapts to terminal window size changes
- 🎯 **Multi-file Navigation**: Easy switching between multiple file comparisons

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
make demo

# Quick test with sample files
make run-files

# Quick test with sample directories  
make run-dirs
```

### Interactive Controls

#### File Selection Mode
- **Tab**: Switch between left and right input fields
- **Enter**: Load the entered path (file or directory)
- **↑/↓**: Navigate through the list of common files
- **Ctrl+D**: Start comparing the selected files
- **?**: Show help screen
- **Q/Ctrl+C**: Quit application

#### Diff View Mode
- **↑/↓** or **j/k**: Navigate through diff lines
- **g**: Go to top of diff
- **G**: Go to bottom of diff
- **n**: Next common file
- **p**: Previous common file
- **Esc**: Return to file selection
- **?**: Show help screen
- **Q/Ctrl+C**: Quit application

## Color Legend

- **Green background**: Added lines (+)
- **Blue background**: Deleted lines (-)
- **Gray text**: Unchanged lines

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
2. **Find Common Files**: For directories, the tool finds files that exist in both locations with the same relative path
3. **Select File**: Choose which common file to compare using the arrow keys
4. **View Diff**: See the side-by-side comparison with color-coded changes
5. **Navigate**: Move through the diff and switch between files seamlessly

## Examples

### Comparing Two Files
```bash
./filecmp config.old.yaml config.new.yaml
```

### Comparing Project Directories
```bash
./filecmp ./project-v1 ./project-v2
```
This will find all common files between the two project directories and allow you to compare them one by one.

### Interactive Mode
```bash
./filecmp
```
Then enter paths interactively using the TUI.

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
├── file/       # File operations and type detection
└── ...

examples/       # Sample files for testing
```

## 🎛️ Advanced Usage

### Keyboard Shortcuts Summary
| Key | File Selection | Diff View | Description |
|-----|----------------|-----------|-------------|
| `Tab` | ✅ | ❌ | Switch input fields |
| `Enter` | ✅ | ❌ | Load entered path |
| `↑/↓` | ✅ | ✅ | Navigate lists/lines |
| `j/k` | ❌ | ✅ | Vim-style navigation |
| `g/G` | ❌ | ✅ | Jump to top/bottom |
| `n/p` | ✅ | ✅ | Next/previous file |
| `Ctrl+D` | ✅ | ❌ | Start comparison |
| `Esc` | ❌ | ✅ | Return to file selection |
| `?` | ✅ | ✅ | Show help screen |
| `Q/Ctrl+C` | ✅ | ✅ | Quit application |

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
- Export diff results
- Configuration file support

Please feel free to submit issues, feature requests, or pull requests.

## 📋 Roadmap

- [ ] Syntax highlighting for code diffs
- [ ] Side-by-side comparison view
- [ ] Export diffs to HTML/PDF
- [ ] Configuration file support
- [ ] Plugin system for custom file types
- [ ] Integration with Git for commit diffs