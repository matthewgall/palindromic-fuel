# AGENTS.md - Development Guidelines for Palindromic Fuel Calculator

## Build Commands
- **Build binary**: `go build -o palindromic-fuel main.go`
- **Format code**: `gofmt -w main.go`
- **Check formatting**: `gofmt -d main.go` (should return empty)
- **Vet code**: `go vet main.go`
- **Run program**: `./palindromic-fuel [flags]`
- **Start web server**: `./palindromic-fuel -web [-port=8080]`

## Test Commands
- **Run tests**: `go test ./...`
- **Run tests with verbose output**: `go test -v ./...`
- **Run single test**: `go test -run TestName`
- **Run tests with coverage**: `go test -cover ./...`

## Code Style Guidelines

### File Structure
- Single `main.go` file (legacy Go project, no modules)
- No external dependencies - uses only standard library
- Keep all code in main.go unless significant refactoring needed

### Naming Conventions
- **Types/Structs**: PascalCase (e.g., `Result`, `PalindromeData`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase (e.g., `pricePerLitre`, `maxLitres`)
- **Constants**: PascalCase or ALL_CAPS if exported

### Code Formatting
- Use `gofmt` for consistent formatting
- 4-space indentation (Go standard)
- Max line length: reasonable, break long lines
- Group imports by standard library, then third-party

### Error Handling
- Return errors from functions: `func doSomething() error`
- Use `fmt.Errorf` with `%w` for error wrapping
- Handle errors immediately or return them up the chain
- No panics in production code

### Documentation
- Add comments for all exported functions/types
- Use `// FunctionName does X` format for function docs
- Keep comments concise but descriptive

### Performance
- Prefer efficiency in algorithms (current code is optimized for palindrome generation)
- Use appropriate data types (float64 for prices, int for counts)
- Consider memory usage for large datasets

### Testing
- Write table-driven tests for core functions
- Test edge cases (negative numbers, zero, large values)
- Verify mathematical correctness of palindrome calculations

## Common Patterns
- Use `math` package for floating-point operations
- String manipulation with `strconv` and `strings` packages
- CSV export using `encoding/csv` package
- Command-line flags using `flag` package