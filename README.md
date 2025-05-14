# Xel (/z·eh·l/)

[![Go Report Card](https://goreportcard.com/badge/github.com/dev-kas/xel)](https://goreportcard.com/report/github.com/dev-kas/xel)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<!-- Add other badges, e.g., build status, release version -->

Xel is a runtime environment for **VirtLang**, a dynamic, modern scripting language designed for simplicity, power, and extensibility. Xel provides the tools to write, execute, and manage VirtLang code, complete with a REPL, a standard library, and a module system.

## Features

**Language (VirtLang via Xel):**
*   **Dynamic Typing:** Flexible type system.
*   **Modern Syntax:** Familiar C-like syntax with support for:
    *   Variables (`let`) and Constants (`const`).
    *   Functions (named, anonymous, closures).
    *   Classes with `public`/`private` members and `constructor`s.
*   **Rich Data Types:** Includes `Number`, `String`, `Boolean`, `Nil`, `Object`, and `Array`.
*   **Control Flow:** `if/else if/else`, `while` loops, `break`, `continue`.
*   **Error Handling:** `try/catch` blocks.
*   **Module System:** Simple `import` and `exports` for code organization.

**Runtime (Xel):**
*   **Script Execution:** Run `.xel` script files directly from the command line.
*   **Interactive REPL:** Experiment with VirtLang code in real-time.
*   **Standard Library:** Built-in modules for math (`xel:math`), strings (`xel:strings`), and arrays (`xel:array`).
*   **CLI Tools:** Intuitive command-line interface for running scripts and managing projects.
*   **Informative Error Reporting:** Errors include line and column numbers for easier debugging.
*   **Environment Variables:** Access script path (`__filename__`, `__dirname__`) and command-line arguments (`proc.args`).
*   **Automatic Version Checking:** Get notified of new Xel releases.
*   **Cross-Platform:** Available for macOS, Linux, and Windows.
*   **Lightweight:** Minimal dependencies and fast startup time.

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

1.  Download the appropriate binary for your platform from the [releases page](https://github.com/dev-kas/xel/releases/latest).
2.  Rename it to `xel` (or `xel.exe` on Windows).
3.  Make it executable (on Unix-like systems): `chmod +x xel`.
4.  Move it to a directory in your PATH (e.g., `/usr/local/bin` on Linux/macOS).

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

### Running Xel Scripts

Create a file with a `.xel` extension (e.g., `myscript.xel`).
```bash
# Basic usage
xel run myscript.xel

# With command-line arguments
xel run myscript.xel arg1 "another argument"
```
Arguments are accessible within the script via the `proc.args` array.

### Interactive REPL (Read-Eval-Print Loop)

Run `xel` without any arguments to start the REPL:
```bash
xel
```
```
Welcome to Xel vX.Y.Z REPL (VirtLang vA.B.C)!
Type '!exit' to exit the REPL.
> let message = "Hello from REPL!"
< "Hello from REPL!"
> print(message)
Hello from REPL!
< nil
> 10 + 20
< 30
```

### Example Xel Script

Create a file named `example.xel`:
```xel
// example.xel

// Import the strings module from the standard library
const strings = import("xel:strings")

// Define a function
fn greet(name) {
  return "Hello, " + name + "!"
}

let personName = "Xel User"
let greeting = greet(personName)

print(strings.upper(greeting))

if (len(proc.args) > 0) {
  print("You passed these arguments:")
  array.forEach(proc.args, fn(arg, index) {
    print(strings.format("Arg %v: %v", index + 1, arg))
  })
} else {
  print("Try running with arguments: xel run example.xel test1 test2")
}
```

Run it with:
```bash
xel run example.xel "first arg" 42
```

Expected Output:
```
HELLO, XEL USER!
You passed these arguments:
Arg 1: first arg
Arg 2: 42
```

## Documentation

For detailed documentation on the Xel runtime and the VirtLang language features, syntax, and standard library, please refer to:
 **[DOCS.md](DOCS.md)**

## Standard Library Highlights

Xel comes with a useful set of built-in native modules:

*   **`xel:math`**: For mathematical operations like `sqrt`, `random`, `sin`, `cos`, `PI`, `E`, aggregation functions (`sum`, `mean`, `median`), and more.
*   **`xel:strings`**: For string manipulation like `trim`, `split`, `upper`, `lower`, `includes`, `format`, `slice`, and more.
*   **`xel:array`**: For array operations like `map`, `filter`, `reduce`, `push`, `pop`, `sort`, `slice`, and more.

See [DOCS.md](DOCS.md) for full details on available functions.

## Development

### Running Tests

Ensure you have Go installed. Then, from the project root:
```bash
make test
```

## Contributing

Contributions are welcome and highly appreciated! Whether it's bug fixes, feature enhancements, documentation improvements, or new standard library modules, please feel free to submit a Pull Request.

Please see our **[CONTRIBUTING.md](CONTRIBUTING.md)** for detailed guidelines on how to contribute.

## Release Process

Xel uses GitHub Actions for continuous integration and automated releases:

1.  Every push to the `master` branch and pull requests are automatically built and tested.
2.  To create a new release:
    *   Create and push a new tag with the version number (e.g., `git tag v1.0.0 && git push origin v1.0.0`).
    *   GitHub Actions will automatically build the binaries for all supported platforms.
    *   A new GitHub Release will be drafted with the binaries attached.
    *   Once the release is published, the installation scripts will automatically use the latest release.

## Acknowledgements

Xel builds upon the **VirtLang-Go v2 engine**. Special thanks to [@dev-kas](https://github.com/dev-kas) for creating and maintaining the VirtLang project, which serves as the core foundation for this runtime environment.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
