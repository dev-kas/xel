# Xel (/z·eh·l/)

Xel is a runtime for VirtLang, a lightweight programming language designed for simplicity and extensibility. Xel provides an environment for executing VirtLang code.

## Features

- **VirtLang Execution**: Run VirtLang scripts with a simple command
- **Cross-Platform**: Available for macOS, Linux, and Windows
- **Lightweight**: Minimal dependencies and fast startup time
- **Extensible**: Designed to be extended with additional functionality

## Installation

### Quick Install (Linux, macOS, Windows with WSL)

You can install Xel with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/install.sh | sh
```

This will automatically detect your operating system and architecture, download the appropriate binary, and install it to your system.

#### Update Xel

To update to the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/update.sh | sh
```

#### Uninstall Xel

To remove Xel from your system:

```bash
curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/uninstall.sh | sh
```

To uninstall without confirmation (useful for automated scripts):

```bash
curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/uninstall.sh | sh -s -- -y
```

### Manual Installation

1. Download the appropriate binary for your platform from the [releases page](https://github.com/dev-kas/xel/releases/latest)
2. Rename it to `xel` (or `xel.exe` on Windows)
3. Make it executable (on Unix-like systems): `chmod +x xel`
4. Move it to a directory in your PATH:
   - Linux/macOS: `/usr/local/bin` or `~/.local/bin`
   - Windows: Create a directory and add it to your PATH

### From Source

```bash
# Clone the repository
git clone https://github.com/dev-kas/xel.git
cd xel

# Build for your platform
make build

# Or build for a specific platform
make build-mac    # macOS (arm64)
make build-linux  # Linux (amd64)
make build-windows # Windows (amd64)
```

The compiled binaries will be available in the `bin` directory.

## Usage

### Basic Commands

```bash
# Check version
xel --version

# Show help
xel --help
```

### Running VirtLang Scripts

```bash
# Basic usage
xel run script.xel

# With arguments (for future implementation)
xel run script.xel arg1 arg2 arg3
```

### Example VirtLang Script

Create a file named `example.xel` with the following content:

```
// Define a function
fn add(a, b) {
  return a + b
}

// Call the function
add(10, 20)
```

Run it with:

```bash
xel run example.xel
```

Output:
```
Result: 30
```

### Future Commands

```bash
# Initialize a new project (coming soon)
xel init

# Install a package (coming soon)
xel install package-name
```

## Development

### Running Tests

```bash
make test
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## Release Process

Xel uses GitHub Actions for continuous integration and automated releases:

1. Every push to the master branch and pull requests are automatically built and tested
2. To create a new release:
   - Create and push a new tag with the version number: `git tag v1.0.0 && git push origin v1.0.0`
   - GitHub Actions will automatically build the binaries for all platforms
   - A new GitHub Release will be created with the binaries attached
   - The installation scripts will automatically use the latest release

## Acknowledgements

Special thanks to [@dev-kas](https://github.com/dev-kas) for creating the VirtLang project, without which Xel would not be possible. VirtLang serves as the foundation for this runtime environment.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.