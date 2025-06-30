# envsubst API Documentation

This document provides comprehensive API documentation for developers using the `envsubst` package.

## Overview

The `envsubst` package provides environment variable substitution functionality with support for various expansion formats, restrictions, and error handling modes. It supports both simple variable substitution and complex shell-style expansions.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/allex/envsubst"
)

func main() {
    result, err := envsubst.String("Hello $USER from $HOME")
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
}
```

## Main Package API

### Basic Functions

#### `String(s string) (string, error)`

Performs environment variable substitution on a string with no restrictions.

**Parameters:**
- `s` - Input string containing environment variable references

**Returns:**
- Processed string with variables substituted
- Error if parsing fails

**Example:**
```go
result, err := envsubst.String("Welcome $USER to $HOME")
```

#### `Bytes(b []byte) ([]byte, error)`

Performs environment variable substitution on a byte slice.

**Parameters:**
- `b` - Input byte slice containing environment variable references

**Returns:**
- Processed byte slice with variables substituted
- Error if parsing fails

**Example:**
```go
input := []byte("Config: ${CONFIG_FILE:-/etc/myapp.conf}")
result, err := envsubst.Bytes(input)
```

#### `ReadFile(filename string) ([]byte, error)`

Reads a file and performs environment variable substitution on its contents.

**Parameters:**
- `filename` - Path to the file to process

**Returns:**
- Processed file contents as byte slice
- Error if file reading or parsing fails

**Example:**
```go
result, err := envsubst.ReadFile("config.template")
```

### Restricted Functions

#### `StringRestricted(s string, noUnset, noEmpty bool) (string, error)`

Performs environment variable substitution with optional restrictions.

**Parameters:**
- `s` - Input string
- `noUnset` - If true, return error for unset variables
- `noEmpty` - If true, return error for empty variables

**Example:**
```go
// Fail if any variable is unset or empty
result, err := envsubst.StringRestricted("${USER} ${HOME}", true, true)
```

#### `BytesRestricted(b []byte, noUnset, noEmpty bool) ([]byte, error)`

Byte slice version of `StringRestricted`.

#### `ReadFileRestricted(filename string, noUnset, noEmpty bool) ([]byte, error)`

File reading version of `StringRestricted`.

### Advanced Functions

#### `StringRestrictedNoDigit(s string, noUnset, noEmpty, noDigit bool) (string, error)`

Most comprehensive substitution function with all available restrictions.

**Parameters:**
- `s` - Input string
- `noUnset` - Fail if variable is not set
- `noEmpty` - Fail if variable is set but empty  
- `noDigit` - Ignore variables starting with a digit (e.g., `$1`, `${2}`)

**Example:**
```go
// Process template but ignore numeric variables and fail on unset vars
result, err := envsubst.StringRestrictedNoDigit(
    "User: $USER, Arg1: $1, Config: ${CONFIG}", 
    true,  // noUnset
    false, // noEmpty
    true   // noDigit - will ignore $1
)
```

#### `BytesRestrictedNoDigit(b []byte, noUnset, noEmpty, noDigit bool) ([]byte, error)`

Byte slice version of `StringRestrictedNoDigit`.

#### `ReadFileRestrictedNoDigit(filename string, noUnset, noEmpty, noDigit bool) ([]byte, error)`

File reading version of `StringRestrictedNoDigit`.

## Advanced Usage with Parse Package

For more control over the parsing process, you can use the `parse` package directly.

### Types

#### `Parser`

The main parser type for processing templates.

```go
type Parser struct {
    Name     string        // Template name for error reporting
    Env      *Env          // Environment variable source
    Restrict *Restrictions // Processing restrictions
    Mode     Mode          // Error handling mode
}
```

#### `Restrictions`

Controls parsing behavior and validation.

```go
type Restrictions struct {
    NoUnset    bool       // Fail on unset variables
    NoEmpty    bool       // Fail on empty variables  
    NoDigit    bool       // Ignore numeric variables
    VarMatcher varMatcher // Custom variable matching (advanced)
}
```

#### `Mode`

Defines error handling strategy.

```go
const (
    Quick     Mode = iota // Stop on first error
    AllErrors             // Collect all errors
)
```

#### `Env`

Environment variable provider.

```go
type Env struct {
    // Internal fields
}

func NewEnv(env []string) *Env
func (e *Env) Get(key string) string
func (e *Env) Has(key string) bool
```

### Advanced Example

```go
package main

import (
    "fmt"
    "os"
    "github.com/allex/envsubst/parse"
)

func main() {
    // Create custom environment
    env := parse.NewEnv([]string{
        "USER=developer",
        "HOME=/home/developer",
        "DEBUG=true",
    })
    
    // Create restrictions
    restrictions := &parse.Restrictions{
        NoUnset: true,
        NoEmpty: false,
        NoDigit: true,
    }
    
    // Create parser
    parser := parse.New("mytemplate", env, restrictions)
    parser.Mode = parse.AllErrors // Collect all errors
    
    // Process template
    template := "User: ${USER}, Home: ${HOME}, Missing: ${MISSING:-default}"
    result, err := parser.Parse(template)
    
    if err != nil {
        fmt.Printf("Errors: %v\n", err)
        return
    }
    
    fmt.Printf("Result: %s\n", result)
}
```

## Variable Expansion Formats

| Expression | Description |
|------------|-------------|
| `$VAR` or `${VAR}` | Simple variable substitution |
| `${VAR-default}` | Use default if VAR is unset |
| `${VAR:-default}` | Use default if VAR is unset or empty |
| `${VAR=default}` | Set and use default if VAR is unset |
| `${VAR:=default}` | Set and use default if VAR is unset or empty |
| `${VAR+alternate}` | Use alternate if VAR is set |
| `${VAR:+alternate}` | Use alternate if VAR is set and non-empty |
| `$$VAR` | Literal `$VAR` (escaped) |

## Error Handling

### Error Types

The package returns descriptive errors for various failure conditions:

- **Parse errors**: Invalid syntax in template
- **NoUnset errors**: Required variable not set
- **NoEmpty errors**: Variable set but empty when not allowed

### Error Modes

#### Quick Mode (Default)
```go
parser.Mode = parse.Quick
```
Stops processing on first error and returns immediately.

#### All Errors Mode
```go
parser.Mode = parse.AllErrors
```
Continues processing and collects all errors, returning them as a combined error message.

## Best Practices

### 1. Use Appropriate Restriction Level

```go
// For configuration files - ensure all variables are set
result, err := envsubst.StringRestricted(template, true, true)

// For templates with optional variables
result, err := envsubst.String(template)
```

### 2. Handle File Processing Errors

```go
content, err := envsubst.ReadFile("config.template")
if err != nil {
    if os.IsNotExist(err) {
        // Handle missing file
    } else {
        // Handle substitution error
    }
}
```

### 3. Validate Templates

```go
// Test with empty environment to catch syntax errors
parser := parse.New("test", parse.NewEnv([]string{}), &parse.Restrictions{})
_, err := parser.Parse(template)
if err != nil {
    // Template has syntax errors
}
```

### 4. Custom Environment

```go
// Use custom environment instead of os.Environ()
customEnv := []string{
    "API_KEY=secret123",
    "DB_HOST=localhost",
    "DB_PORT=5432",
}
env := parse.NewEnv(customEnv)
parser := parse.New("config", env, restrictions)
```

## Testing

### Unit Testing Templates

```go
func TestTemplate(t *testing.T) {
    // Set test environment
    os.Setenv("TEST_VAR", "test_value")
    defer os.Unsetenv("TEST_VAR")
    
    template := "Value: ${TEST_VAR}"
    result, err := envsubst.String(template)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    expected := "Value: test_value"
    if result != expected {
        t.Errorf("expected %q, got %q", expected, result)
    }
}
```

### Testing with Custom Environment

```go
func TestWithCustomEnv(t *testing.T) {
    env := parse.NewEnv([]string{"CUSTOM=value"})
    parser := parse.New("test", env, &parse.Restrictions{})
    
    result, err := parser.Parse("${CUSTOM}")
    
    assert.NoError(t, err)
    assert.Equal(t, "value", result)
}
```

## Performance Considerations

1. **Reuse Parsers**: Create parser instances once and reuse them for multiple templates with the same configuration.

2. **Environment Caching**: The `Env` type caches environment variable lookups for better performance.

3. **Error Mode**: Use `Quick` mode for better performance when you only need to know if processing succeeds.

4. **Template Validation**: Pre-validate templates during application startup rather than at runtime.

## Migration from os.ExpandEnv

If you're migrating from `os.ExpandEnv`, note these differences:

```go
// os.ExpandEnv - basic functionality
result := os.ExpandEnv("$HOME/config")

// envsubst equivalent
result, err := envsubst.String("$HOME/config")

// envsubst with shell-style expansions (not supported by os.ExpandEnv)
result, err := envsubst.String("${HOME:-/tmp}/config")
```

## Common Use Cases

### Configuration File Processing

```go
// Process YAML/JSON config templates
configData, err := envsubst.ReadFileRestricted("app.config.yaml", true, true)
if err != nil {
    log.Fatal("Config processing failed:", err)
}
```

### Docker Compose Templates

```go
// Process docker-compose templates with defaults
template := `
version: '3'
services:
  app:
    image: myapp:${VERSION:-latest}
    environment:
      - DB_HOST=${DB_HOST:-localhost}
      - DB_PORT=${DB_PORT:-5432}
`
result, err := envsubst.String(template)
```

### Kubernetes Manifest Templates

```go
// Process k8s manifests with environment-specific values
manifest, err := envsubst.ReadFileRestricted("deployment.yaml", true, false)
```

## Changelog and Versioning

See [CHANGELOG.md](CHANGELOG.md) for version history and breaking changes.

## License

MIT License - see [LICENSE](LICENSE) file for details. 
