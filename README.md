# Fabled Story Book

Write a story and run it anywhere

```jabl
{
    print("Would you like to write a simple story book?")

    choice("Yes", {
        print("Great! Let's check this out!")
    })

    choice("No", {
        print("No worries!")
    })
}
```

# Running a story

Stories that are written in JABL can be run from the command line, a web server or a web browser.

# JABL

JABL is a constrained language that allows you to write stories and run them equally well from a CLI, web server or a web browser.

Each section of a story is a single file in a directory, where the directory is the story itself.

```sh
$ tree
.
├── story
│   ├── entrypoint.jabl
│   ├── chapter1.jabl
│   └── chapter2.jabl
```

## Keywords

THere are only a handful of special keywords in JABL.

If something can be achieved without introducing a new keyword, it should be done.

- `print`: Print a message to the console
- `choice`: Present a choice to the user
- `goto`: Jump to another section
- `if`: Choose the branch to run based on a condition
- `set`: Set a variable
- `get`: Get a variable

### print

Adds text to the output buffer.

```jabl
print("Hello, World!")
```

### choice

Present a choice to the user. The first parameter is the identifier of this choice and the second parameter is a block of code to execute when the user selects it.

```jabl
choice("Yes", {
    print("Great! Let's check this out!")
})
```

### goto

Jump to another section. This causes the interpreter to load and execute the section with the given identifier.

It is up to the loader itself to determine the code for a particular identifier.

```jabl
goto("chapter1.jabl")
```

### if

Choose the branch to run based on a condition. If accepts any expression that evaluates to a boolean and must be wrapped in parentheses.

```jabl
if (true) {
    print("this is printed")
}
```
You can also use an else statement

```jabl
if (false) {
    print("this is not printed")
} else {
    print("this is printed")
}
```

### set

Set a variable. The first parameter is the name of the variable and the second parameter is the value.

This function returns the value that was set.

All variables are `float64` values.

```jabl
set("some-value", 123) # => 123
```

### get

Get a variable. The first parameter is the name of the variable.

This function returns the value of the variable.

All variables are `float64` values.

```jabl
get("some-value") # => 123
```

## Comments

You can add single-line coments by prefixing a line with `//`.

```jabl
// This is a comment
```

# Interpreter

The JABL interpreter is a simple program that reads a story and executes it. It is written in Go and can be run from the command line or compiled to WASM and invoked from Javascript.

The Interpreter has two key abstractions:

- `StateMapper`: Implements the `get` and `set` functions with a backing store for `float64` variables
- `SectionLoader`: Implements a strategy for loading sections from a story. This could be from a URL or a local file system.

It is (should be) safe to load JABL code from any source, as the interpreter does not execute any code until it has been parsed and validated.

## Execution

The interpreter reads a section by identifier from a `SectionLoader`, then parses the JABL code into an AST. The AST is then executed by the interpreter.

The JABL code is executed immediately and synchronously and results in a structure that has three components:

- `output`: A string that is printed to the console
- `choices`: A list of choices that the user can make
- `transition`: The identifier of the next section to load and execute

If there is a `transition` value, the interpreter will load and execute the next section. If there is no `transition` value, the interpreter will stop until the user selects a `choice`.

Each `choice` is a name and function that is executed when the user selects it. The function is a block of JABL code that is executed immediately and synchronously.

The `output` is an output buffer that is a collection of each `print` statement separated by a newline.

# Building

The JABL interpreter can be built for the command line or for the web.

## Command Line

```sh
go build -o bin cmd/cli/main.go
```

## Web

```sh
GOOS=js GOARCH=wasm go build -o ./web/jabl.wasm ./cmd/wasm/main.go
```

## Web UI

```sh
cd web
npm ci
npm run start
```