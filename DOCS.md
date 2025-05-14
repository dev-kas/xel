# Xel Language Documentation

Welcome to Xel, a dynamic, modern scripting language built for simplicity and power. Xel is implemented in Go (powered by VirtLang-Go v2 engine) and offers a familiar syntax with robust features.

## Table of Contents

1.  [Overview](#overview)
2.  [Getting Started](#getting-started)
    *   [Installation](#installation)
    *   [Running Xel Scripts](#running-xel-scripts)
    *   [The REPL](#the-repl)
3.  [Language Basics](#language-basics)
    *   [Comments](#comments)
    *   [Case Sensitivity](#case-sensitivity)
    *   [Semicolons](#semicolons)
    *   [Data Types](#data-types)
    *   [Variables and Constants](#variables-and-constants)
    *   [Operators](#operators)
4.  [Control Flow](#control-flow)
    *   [Conditional Statements (`if`, `else if`, `else`)](#conditional-statements)
    *   [Looping (`while`)](#looping)
    *   [Loop Control (`break`, `continue`)](#loop-control)
5.  [Functions](#functions)
    *   [Defining Functions](#defining-functions)
    *   [Anonymous Functions](#anonymous-functions)
    *   [Closures](#closures)
    *   [Return Values (`return`)](#return-values)
6.  [Classes and Objects](#classes-and-objects)
    *   [Defining Classes](#defining-classes)
    *   [Constructors](#constructors)
    *   [Properties and Methods](#properties-and-methods)
    *   [Public and Private Members](#public-and-private-members)
    *   [Instantiation](#instantiation)
    *   [Object Literals](#object-literals)
    *   [Array Literals](#array-literals)
7.  [Error Handling (`try`, `catch`)](#error-handling)
8.  [Modules and Imports](#modules-and-imports)
    *   [Importing Modules](#importing-modules)
    *   [Creating Modules (Exports)](#creating-modules-exports)
    *   [Native Modules](#native-modules)
9.  [Built-in Globals and Functions](#built-in-globals-and-functions)
    *   [Constants](#constants)
    *   [Utility Functions](#utility-functions)
    *   [Environment Variables](#environment-variables)
10. [Standard Library](#standard-library)
    *   [`xel:math`](#xelmath)
    *   [`xel:strings`](#xelstrings)
    *   [`xel:array`](#xelarray)
11. [Command-Line Interface (CLI)](#command-line-interface-cli)

## Overview

Xel is designed to be an easy-to-learn scripting language that is both expressive and efficient. It supports procedural, object-oriented, and functional programming paradigms.

Key features include:
*   Dynamic Typing
*   Lexical Scoping and Closures
*   Classes with Public/Private Members
*   A Simple Module System
*   Built-in Data Structures (Objects, Arrays)
*   Comprehensive Error Handling
*   A Growing Standard Library
*   An Interactive REPL

## Getting Started

### Installation

Xel provides an update script for easy installation and updates:
```sh
curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/install.sh | sh
```
The runtime will also notify you if a new version is available when you run it.

### Running Xel Scripts

To execute a Xel script file (which should have a `.xel` extension):
```sh
xel run your_script.xel [arg1 arg2 ...]
```
Any arguments passed after the script name will be available within the script via the `proc.args` array.

**Example:**
`your_script.xel`:
```xel
print("Script name:", __filename__)
print("Arguments:", proc.args)
```
Running `xel run your_script.xel hello world` would output:
```
Script name: /path/to/your_script.xel
Arguments: ["hello", "world"]
```

### The REPL

Running `xel` without any arguments starts the Read-Eval-Print Loop (REPL):
```sh
xel
```
This provides an interactive environment to experiment with Xel code.
*   Type Xel expressions or statements and press Enter.
*   Use `!exit` to quit the REPL.
*   Command history is available (use Up/Down arrows).

## Language Basics

### Comments

Xel supports two types of comments:
*   Single-line comments: `// This is a comment`
*   Multi-line comments: `/* This is a
   multi-line comment */`

### Case Sensitivity

Xel is case-sensitive. `myVariable` and `myvariable` are treated as two different identifiers.

### Separation

Multiple statements or expressions can be simply separated by a space. They are automatically distinguished appropriately by the parser. They can be used to separate multiple statements on a single line.
```xel
let a = 1 let b = 2 // Space separating two statements
print(a + b)
```

### Data Types

Xel is dynamically typed. The main data types are:

*   **Number:** Represents both integers and floating-point numbers (internally `float64`).
    ```xel
    let integer = 10
    let floatNum = 3.14
    ```
    <!--let scientific = 1.2e3-->
*   **String:** A sequence of characters. Can be defined using single (`'`) or double (`"`) quotes.
    ```xel
    let greeting = "Hello, Xel!"
    let name = 'World'
    ```
*   **Boolean:** Represents `true` or `false`.
    ```xel
    let isActive = true
    let isDone = false
    ```
*   **Nil:** Represents the absence of a value.
    ```xel
    let noValue = nil
    ```
*   **Object:** An unordered collection of key-value pairs (like dictionaries or maps). Keys are strings.
    ```xel
    let person = { name: "Jedi", age: 30 }
    ```
*   **Array:** An ordered list of values.
    ```xel
    let numbers = [1, 2, 3, 4]
    let mixed = [1, "two", true, nil]
    ```
*   **Function:** A block of code that can be called.
*   **Class:** A blueprint for creating objects.
*   **ClassInstance:** An object created from a class.

### Variables and Constants

*   **Variables (`let`):** Used to declare mutable variables.
    ```xel
    let message = "Hello"
    message = "Hi" // Allowed
    ```
*   **Constants (`const`):** Used to declare immutable variables. Once assigned, their value cannot be changed.
    ```xel
    const PI = 3.14159
    // PI = 3.14 // This would cause an error
    ```
Scope is lexical (block-scoped).

### Operators

Xel supports a standard set of operators:

*   **Arithmetic Operators:**
    *   `+` (Addition, String Concatenation)
    *   `-` (Subtraction)
    *   `*` (Multiplication)
    *   `/` (Division)
    *   `%` (Modulo)
    Precedence follows standard mathematical rules (e.g., `*` and `/` before `+` and `-`). Use parentheses `()` to control evaluation order.
    ```xel
    let sum = 5 + 3          // 8
    let product = 4 * 2      // 8
    let result = (2 + 3) * 4 // 20
    let greeting = "Hello" + " " + "World" // "Hello World"
    ```

*   **Comparison Operators:**
    *   `==` (Equal to)
    *   `!=` (Not equal to)
    *   `<`  (Less than)
    *   `>`  (Greater than)
    *   `<=` (Less than or equal to)
    *   `>=` (Greater than or equal to)
    These operators evaluate to a boolean (`true` or `false`). Currently, they primarily operate on numbers.
    ```xel
    let isEqual = (5 == 5)   // true
    let isGreater = (10 > 3) // true
    ```

*   **Assignment Operator:**
    *   `=` (Assign value)
    ```xel
    let x = 10
    x = 20 // x is now 20
    ```

*   **Member Access Operators:**
    *   `.` (Dot notation for object properties): `myObject.property`
    *   `[]` (Bracket notation for object properties and array elements): `myObject["property"]`, `myArray[0]`

## Control Flow

### Conditional Statements

Xel uses `if`, `else if`, and `else` for conditional execution.
```xel
let num = 10

if (num > 10) {
  print("Number is greater than 10")
} else if (num < 10) {
  print("Number is less than 10")
} else {
  print("Number is exactly 10")
}
```
The condition in parentheses is evaluated. If "truthy", the corresponding block is executed.
*   **Truthy values:** `true`, non-zero numbers, non-empty strings, objects, arrays, functions.
*   **Falsy values:** `false`, `0`, `""` (empty string), `nil`.

### Looping

Xel provides a `while` loop for repeated execution of a block of code as long as a condition is true.
```xel
let count = 0
while (count < 5) {
  print(count)
  count = count + 1
}
// Output: 0, 1, 2, 3, 4
```

### Loop Control

*   **`break`**: Exits the innermost `while` loop immediately.
    ```xel
    let i = 0
    while (true) {
      if (i == 3) {
        break // Exit loop when i is 3
      }
      print(i)
      i = i + 1
    }
    // Output: 0, 1, 2
    ```
*   **`continue`**: Skips the rest of the current iteration and proceeds to the next iteration of the innermost `while` loop.
    ```xel
    let i = 0
    while (i < 5) {
      i = i + 1
      if (i == 3) {
        continue // Skip printing 3
      }
      print(i)
    }
    // Output: 1, 2, 4, 5
    ```

## Functions

Functions are first-class citizens in Xel. They can be assigned to variables, passed as arguments, and returned from other functions.

### Defining Functions

Use the `fn` keyword to define a function.
```xel
fn greet(name) {
  return "Hello, " + name + "!"
}

let message = greet("Xel") // "Hello, Xel!"
print(message)
```

### Anonymous Functions

Functions can also be defined without a name (anonymously) and assigned to variables.
```xel
let add = fn(a, b) {
  return a + b
}

let sum = add(5, 3) // 8
print(sum)
```
Anonymous functions are often used for callbacks or Immediately Invoked Function Expressions (IIFEs), though Xel doesn't require explicit IIFE syntax for simple expressions.

### Closures

Xel supports lexical scoping and closures. A function can access variables from its containing (enclosing) scopes, even after the outer function has finished executing.
```xel
fn outer(x) {
  fn inner(y) {
    return x + y // 'x' is captured from outer's scope
  }
  return inner
}

let add5 = outer(5)
let result = add5(3) // 8
print(result)
```

### Return Values

Functions use the `return` keyword to send a value back to the caller. If `return` is omitted or used without a value, the function implicitly returns `nil`.
```xel
fn getNumber() {
  return 42
}

fn doNothing() {
  // implicitly returns nil
}

print(getNumber()) // 42
print(doNothing()) // nil
```

## Classes and Objects

Xel supports object-oriented programming with classes.

### Defining Classes

Use the `class` keyword to define a class.
```xel
class Rectangle {
  public width
  public height

  public constructor(w, h) {
    width = w
    height = h
  }

  public area() {
    return width * height
  }

  private describe() { // Private method
    return "Rectangle: " + width + "x" + height
  }

  public getDetails() {
    return describe() // Can call private method from within the class
  }
}
```

### Constructors

A special method named `constructor` is called when a new instance of the class is created. It's used to initialize the object's properties. If no explicit constructor is provided, a default one is assumed.
Properties and methods defined directly in the class body (not inside `constructor`) are also part of the class definition.

### Properties and Methods

*   **Properties:** Variables associated with an object (e.g., `width`, `height`).
*   **Methods:** Functions associated with an object that can operate on its data (e.g., `area()`).

### Public and Private Members

Class members (properties and methods) can be declared as `public` or `private`.
*   `public`: Accessible from anywhere. This is the default if no access modifier is specified (though explicitly using `public` is good practice).
*   `private`: Accessible only from within the class itself.
    ```xel
    class MyClass {
      public publicVar
      private privateVar

      public constructor(val) {
        publicVar = val
        privateVar = val * 2
      }

      public getPrivate() {
        return privateVar // Accessible
      }
    }

    let obj = MyClass(10)
    print(obj.publicVar) // 10
    print(obj.privateVar) // nil (cannot access privateVar from outside)
    print(obj.getPrivate()) // 20
    ```

### Instantiation

Create an instance of a class by calling the class name as if it were a function, passing arguments to its constructor.
```xel
let rect = Rectangle(10, 5)
print(rect.width)      // 10
print(rect.area())     // 50
print(rect.getDetails()) // Rectangle: 10x5
// print(rect.describe()) // This would result in an error or nil, as describe is private
```

### Object Literals

You can create objects directly using object literal syntax:
```xel
let car = {
  make: "XelMotors",
  model: "SedanX",
  year: 2024,
  start: fn() {
    print(make + " " + model + " started.")
  }
}

print(car.make) // XelMotors
car.start()     // XelMotors SedanX started.

// Accessing properties
print(car["model"]) // SedanX

// Adding new properties
car.color = "blue"
print(car.color) // blue

// Modifying properties
car.year = 2025
```
Property shorthand is also supported if a variable with the same name exists in the current scope:
```xel
let name = "Jedi"
let age = 30
let user = { name, age } // equivalent to { name: name, age: age }
print(user.name) // Jedi
```

### Array Literals

Arrays are ordered collections of items:
```xel
let numbers = [10, 20, 30]
let colors = ["red", "green", "blue"]

print(numbers[0])      // 10
print(colors[1])       // "green"

colors[1] = "yellow" // Modify element
print(colors)          // ["red", "yellow", "blue"]

// Adding elements (though `xel:array.push` is preferred for clarity)
colors[3] = "purple"
print(colors)          // ["red", "yellow", "blue", "purple"]
print(len(colors))     // 4
```

## Error Handling

Xel provides `try` and `catch` blocks for handling runtime errors.
```xel
try {
  // Code that might throw an error
  let result = riskyOperation()
  print("Operation successful:", result)
} catch e {
  // Code to handle the error
  print("An error occurred:", e) // 'e' will contain the error message string
}
```
If an error occurs within the `try` block, execution jumps to the `catch` block. The variable specified after `catch` (e.g., `e`) will hold a string describing the error.

## Modules and Imports

Xel supports a module system for organizing code into reusable units.

### Importing Modules

Use the `import()` function to load modules.
*   **Relative file paths:**
    ```xel
    // Assuming myModule.xel is in the same directory or a subdirectory
    let myMod = import("./myModule")
    let otherMod = import("./lib/otherModule")
    ```
*   **Native Xel modules:**
    ```xel
    let math = import("xel:math")
    let strings = import("xel:strings")
    ```
The `import()` function returns the `exports` object from the imported module. Xel caches imported modules, so subsequent imports of the same module will return the cached instance. Cycle detection is in place to prevent infinite loops from circular dependencies.

### Creating Modules (Exports)

To expose functionality from a module file, Simply return the value as the last expression. This value becomes the value returned by `import()`.

**Example `myMath.xel`:**
```xel
fn add(a, b) {
  return a + b
}

fn subtract(a, b) {
  return a - b
}

// This object will be what's returned when importing myMath.xel
return {
    adder: add,
    subtractor: subtract,
    PI: 3.14159
}
```

**Example `main.xel`:**
```xel
let myMath = import("./myMath")

print(myMath.PI)           // 3.14159
print(myMath.adder(5, 3))  // 8
```

### Native Modules

Xel comes with pre-built native modules (implemented in Go) for common tasks. These are imported using the `xel:` prefix. See the [Standard Library](#standard-library) section.

## Built-in Globals and Functions

Xel provides several global constants and functions:

### Constants

*   `true`: The boolean true value.
*   `false`: The boolean false value.
*   `nil`: The nil value.
*   `NaN`: Not-a-Number.
*   `inf`: Infinity.

### Utility Functions

*   `print(...args)`: Prints its arguments to the console, space-separated.
    ```xel
    print("Hello", "Xel", 2024) // Hello Xel 2024
    ```
*   `len(value)`: Returns the length of a string (character count) or an array (element count).
    ```xel
    print(len("hello")) // 5
    print(len([1, 2, 3])) // 3
    ```
*   `typeof(value)`: Returns a string representing the type of the value.
    ```xel
    print(typeof(10))        // "number"
    print(typeof("xel"))     // "string"
    print(typeof(true))      // "boolean"
    print(typeof(nil))       // "nil"
    print(typeof({}))        // "object"
    print(typeof([]))        // "array"
    print(typeof(fn(){}))    // "function"
    ```
*   `import(pathOrName)`: Loads a module. (See [Modules and Imports](#modules-and-imports))

### Environment Variables

These variables are automatically available in scripts run via `xel run`:

*   `__filename__`: A string containing the absolute path to the currently executing script file.
*   `__dirname__`: A string containing the absolute path to the directory of the currently executing script file.
*   `proc`: An object containing process-related information:
    *   `proc.args`: An array of strings representing command-line arguments passed to the script.
    *   `proc.runtime_version`: String, the version of the Xel runtime.
    *   `proc.engine_version`: String, the version of the VirtLang-Go engine.

## Standard Library

Xel includes a standard library accessible via native modules.

### `xel:math`

Provides mathematical functions and constants.
```xel
let math = import("xel:math")
print(math.PI)
print(math.sqrt(16))    // 4
print(math.random())    // A random float between 0.0 (inclusive) and 1.0 (exclusive)
print(math.random(10))  // A random float between 0.0 and 10.0
print(math.random(5, 10))// A random float between 5.0 and 10.0
print(math.max(1, 5, 2)) // 5
print(math.abs(-10))     // 10
// ... and many more (sin, cos, floor, ceil, pow, log, etc.)
```
Some notable functions: `abs`, `round`, `floor`, `ceil`, `sign`, `sin`, `cos`, `tan`, `trunc`, `sum` (for arrays), `mean` (for arrays), `median` (for arrays), `random`, `degToRad`, `radToDeg`, `atan`, `atan2`, `acos`, `asin`, `exp`, `log`, `log2`, `log10`, `pow`, `sqrt`, `cbrt`, `clamp`, `max`, `min`. Constants: `PI`, `E`.

### `xel:strings`

Provides functions for string manipulation.
```xel
let strings = import("xel:strings")
let s = "  Hello Xel!  "
print(strings.trim(s))                 // "Hello Xel!"
print(strings.upper("hello"))          // "HELLO"
print(strings.includes("Xel World", "Xel")) // true
print(strings.split("a,b,c", ","))     // ["a", "b", "c"]
print(strings.charAt("Xel", 0))        // "X"
print(strings.slice("XelLang", 0, 3))  // "Xel"
print(strings.replace("foo bar foo", "foo", "baz")) // "baz bar foo"
print(strings.format("Name: %v, Age: %v", "Jedi", 30)) // "Name: Jedi, Age: 30"
// ... and many more
```
Includes: `charAt`, `charCodeAt`, `includes`, `startsWith`, `endsWith`, `indexOf`, `lastIndexOf`, `concat`, `slice`, `substring`, `substr`, `lower`, `upper`, `trim`, `trimStart`, `trimEnd`, `padStart`, `padEnd`, `repeat`, `replace`, `replaceAll`, `split`, `toArray`, `format`.

### `xel:array`

Provides functions for array manipulation.
```xel
let array = import("xel:array")
let nums = [1, 2, 3, 4]

let doubled = array.map(nums, fn(x) { return x * 2 })
print(doubled) // [2, 4, 6, 8]

let evens = array.filter(nums, fn(x) { return x % 2 == 0 })
print(evens)   // [2, 4]

let sum = array.reduce(nums, fn(acc, x) { return acc + x }, 0)
print(sum)     // 10

print(array.includes(nums, 3)) // true
nums = array.push(nums, 5, 6)  // Returns the new array: [1, 2, 3, 4, 5, 6]
print(nums)
// ... and many more
```
Includes: `push`, `pop`, `shift`, `unshift`, `slice`, `splice`, `fill`, `reverse`, `sort`, `map`, `filter`, `forEach`, `reduce`, `reduceRight`, `includes`, `indexOf`, `lastIndexOf`, `find`, `findIndex`, `every`, `some`, `join`, `concat`, `from`, `of`, `create`.

## Command-Line Interface (CLI)

The `xel` executable provides the following commands:

*   `xel run <filepath.xel> [args...]`: Executes the specified Xel script.
*   `xel init`: (Currently Unimplemented) Intended to initialize a new Xel project.
*   `xel install`: (Currently Unimplemented) Intended to install Xel packages.
*   `xel`: Starts the REPL if no command is given.
*   `xel --version`: Displays the Xel runtime and VirtLang engine versions.
*   `xel --help`: Displays help information.
