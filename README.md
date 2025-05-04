# Xel

Xel is a runtime for VirtLang, a lightweight programming language designed for simplicity and extensibility. Xel provides an environment for executing VirtLang code.

## Features

- **VirtLang Execution**: Run VirtLang scripts with a simple command
- **Cross-Platform**: Available for macOS, Linux, and Windows
- **Lightweight**: Minimal dependencies and fast startup time
- **Extensible**: Designed to be extended with additional functionality

## Installation

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

## Acknowledgements

Special thanks to [@dev-kas](https://github.com/dev-kas) for creating the VirtLang project, without which Xel would not be possible. VirtLang serves as the foundation for this runtime environment.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.