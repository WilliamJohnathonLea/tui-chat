# AGENTS.md

## Overview

This repository is a Go (Golang) project. This file provides build and test commands, code style guidelines, and best practices for agents (and developers) making changes to this codebase.

---

## 1. Building, Linting, and Testing

### 1.1. Build Commands

- **Build the project:**
  ```sh
  go build ./...
  ```
- **Build specific package:**
  ```sh
  go build ./path/to/package
  ```

### 1.2. Lint Commands

- **Run linter (staticcheck or golangci-lint):**
  > NOTE: This repo does not have a linter configured out-of-box. If using lint tools, adopt one of the following:
  ```sh
  go vet ./...
  # (recommended: add 'golangci-lint' or 'staticcheck' for thorough linting)
  golangci-lint run ./...
  staticcheck ./...
  ```

### 1.3. Test Commands

- **Run all tests:**
  ```sh
  go test ./...
  ```
- **Run specific test file:**
  ```sh
  go test ./path/to/file_test.go
  ```
- **Run a specific test function:**
  ```sh
  go test -run '^TestFunctionName$' ./...
  ```
- **Test with coverage:**
  ```sh
  go test -cover ./...
  ```
- **Test with verbose output:**
  ```sh
  go test -v ./...
  ```
- **Common test flags**:
  - `-race` (detect race conditions)
  - `-bench` (run benchmarks)

### 1.4. Dependency Management

- **Add a dependency:**
  ```sh
  go get example.com/module
  ```
- **Tidy modules:**
  ```sh
  go mod tidy
  ```

### 1.5. Formatting

- **Format ALL code (canonical):**
  ```sh
  go fmt ./...
  ```
- **Format single file:**
  ```sh
  go fmt path/to/file.go
  ```

---

## 2. Code Style and Guidelines

### 2.1. Imports
- Standard library imports first, then third-party, then local modules.
- Each group separated by a blank line.
- Use explicit import paths; avoid import aliasing unless required.
  
### 2.2. Formatting
- All code MUST be formatted by `gofmt` or `go fmt`. Editors/plugins should apply this automatically.
- Use **tabs** for indentation. Avoid spaces.
- No enforced line length, but keep lines readable (~100 chars max is typical in practice).
  
### 2.3. Naming
- **Package:** Lowercase, short, meaningful.
- **Exported identifiers:** Capitalized (e.g., `Func`, `TypeName`).
- **Unexported:** CamelCase, starting with lowercase (e.g., `helperFunc`).
- **Types, structs, interfaces:** Use CamelCase. One-method interfaces use `-er` suffix (e.g., `Reader`).
- **Avoid underscores** in names. Prefer `MixedCaps`.
- **Getters/setters:** Do not use `Get`; prefer property name (`Owner()` not `GetOwner()`).

### 2.4. Doc Comments
- All exported functions/types must start with a doc comment that begins with the entity name and describes its purpose.
- Use Go-style doc comment format and full sentences.
- Examples:
  ```go
  // FooBar returns a FooBar configured for use.
  func FooBar() *Bar {...}
  ```

### 2.5. Types and Zero Values
- Design structs so zero-values are usable when possible.
- Prefer using composite literals and `make` for slices/maps/channels; use `new` only if pointer to zero value needed.

### 2.6. Control Structures
- No parentheses around conditionals.
- The opening `{` must be on the same line as `if`, `for`, etc.
- Mandatory braces for all control blocks, even single-line.
- Omit `else` blocks that follow a terminating statement (e.g., `return`).

### 2.7. Error Handling
- Return `error` as last return value in multi-valued returns.
- Handle errors immediately (do not ignore them).
- Use named error variables or wrapped errors (`fmt.Errorf("...")` or `errors.New` for static errors).
- Do not use exceptions; Go uses error values, not panic for ordinary control flow.
- Use `defer` for cleanup (e.g., closing files), placed immediately after successful open/acquire.
- Example:
  ```go
  f, err := os.Open(name)
  if err != nil {
      return err
  }
  defer f.Close()
  ```

### 2.8. Tests
- Test files named `*_test.go`.
- Test functions begin with `Test` (e.g., `func TestMyFeature(t *testing.T)`).
- Use `t.Helper()`, table-driven tests, and subtests (t.Run) for clarity and structure.
- Each package should aim for good coverage but avoid over-testing pure boilerplate.
- Benchmarks use `func BenchmarkXxx(b *testing.B)`.

### 2.9. Miscellaneous
- Use `context.Context` as the first parameter for functions that are cancelable or time-bound.
- When using third-party code, ensure modules and licenses are properly declared in `go.mod`/`go.sum`.
- Use Go modules (`go.mod`, `go.sum`) for project dependencies — do not vendor unless necessary.
- Avoid project-wide global variables where possible.

---

## 3. Useful Resources
- [Go Project Effective Go](https://go.dev/doc/effective_go)
- [Official Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [`go` command documentation](https://pkg.go.dev/cmd/go)

---

## 4. Task Guidance for Agents
- Use the above build/lint/test commands and ensure formatting before proposing code changes.
- Adhere strictly to the code style and error handling conventions described above.
- Always ensure `go build ./...` and `go test ./...` pass after any change.
- If in doubt, prefer idiomatic Go style as exemplified in the standard library!

---

_This AGENTS.md was generated for project guidance. Update it if project requirements or standards change._
