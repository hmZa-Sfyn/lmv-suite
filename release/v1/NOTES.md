# LanManVan v1.0.0 Release Notes

**Release Date:** January 2, 2026  
**Build Author:** hmza

---

## 🚀 Major Features & Enhancements

### 1. **113+ Builtin Functions**
Massive expansion of the builtin function library, organized into 10+ specialized categories:

#### 📊 **Data Type Validators (15 Functions)**
- `isvalidip()` - Validates IPv4/IPv6 addresses
- `isint()` - Checks if value is integer
- `isfloat()` - Checks if value is floating point
- `isalpha()` - Validates alphabetic characters only
- `isalnum()` - Validates alphanumeric characters
- `isnumeric()` - Checks for numeric strings
- `isspace()` - Detects whitespace-only strings
- `isbinary()` - Validates binary format
- `ishex()` - Validates hexadecimal format
- `isuuid()` - Validates UUID format
- `isbase64()` - Validates Base64 encoding
- `ismd5()` - Validates MD5 hash format
- `issha1()` - Validates SHA1 hash format
- `issha256()` - Validates SHA256 hash format
- `isjson()` - Validates JSON format with proper parsing

#### 🔄 **Data Converters (15 Functions)**
- `btoa()` / `atob()` - Base64 encoding/decoding
- `bin2hex()` / `hex2bin()` - Binary to hexadecimal conversion
- `bin2dec()` / `dec2bin()` - Binary to decimal conversion
- `hex2dec()` / `dec2hex()` - Hexadecimal to decimal conversion
- `oct2dec()` / `dec2oct()` - Octal to decimal conversion
- `rot13()` - ROT13 cipher
- `rot47()` - ROT47 cipher
- `caesar()` - Caesar cipher with configurable shift
- `reverse_bytes()` - Reverse byte order
- `toascii()` - Convert bytes to ASCII

#### 📝 **String Processors (15 Functions)**
- `camelcase()` - Convert to camelCase
- `snakecase()` - Convert to snake_case
- `kebabcase()` - Convert to kebab-case
- `capitalize()` - Capitalize first letter
- `lowercase()` / `uppercase()` - Case conversion
- `swapcase()` - Toggle character case
- `ltrim()` / `rtrim()` - Remove leading/trailing whitespace
- `center()` - Center text with padding
- `ljust()` / `rjust()` - Justify text
- `indent()` / `dedent()` - Manage indentation
- `wordcount()` - Count words in text

#### 🔐 **Hash Algorithms (10 Functions)**
- `sha512()` - SHA-512 hashing
- `blake2b()` / `blake2s()` - BLAKE2 hashing
- `hmac_sha256()` / `hmac_sha512()` - HMAC hashing
- `murmur3()` - MurmurHash3
- `xxhash()` - xxHash algorithm
- `fnv1()` / `fnv1a()` - Fowler-Noll-Vo hash
- `djb2()` - Daniel J. Bernstein hash

#### 📦 **Advanced Encoding (10 Functions)**
- `base32()` - Base32 encoding
- `base58()` - Base58 encoding (Bitcoin-style)
- `base85()` - Base85 encoding (ASCII85)
- `punycode()` - Punycode encoding for internationalized domains
- `morse()` - Convert to Morse code
- `binary()` - Convert to binary representation
- `octal()` - Convert to octal representation
- `quoted_printable()` - Quoted-printable encoding
- `percent_encode()` - URL encoding
- `htmlescape()` - HTML entity escaping

#### 📂 **File Operations (9 Functions)**
- `find()` - Search for files matching patterns
- `tail()` - Get last N lines of file
- `head()` - Get first N lines of file
- `touch()` - Create/update file timestamps
- `chmod()` - Change file permissions
- `stat()` - Get file statistics
- `isdir()` - Check if path is directory
- `isfile()` - Check if path is file
- `file()` - **NEW** Write/append to files via pipes

#### ⏰ **Date & Time (8 Functions)**
- `now()` - Current timestamp
- `epoch()` - Unix epoch time
- `iso8601()` - ISO 8601 formatted date
- `strtotime()` - Parse time strings
- `timeformat()` - Format time values
- `strftime()` - Format using strftime syntax
- `timezone()` - Get timezone information
- `dayofweek()` - Get day of week

#### 🌐 **Network Functions (8 Functions)**
- `cidrmatch()` - Check IP against CIDR ranges
- `hostmask()` - Generate network hostmasks
- `broadcast()` - Calculate broadcast addresses
- `ipversion()` - Determine IPv4 or IPv6
- `isprivate()` - Check for private IP ranges
- `isloopback()` - Check for loopback addresses
- `ismulticast()` - Check for multicast addresses
- `getmactable()` - Retrieve MAC address table

#### 🧮 **Math Functions (6 Functions)**
- `pow()` - Exponentiation
- `sqrt()` - Square root
- `cbrt()` - Cube root
- `log()` - Natural logarithm
- `exp()` - Exponential function
- `factorial()` - Factorial calculation

#### 📋 **JSON Processing (3 Functions)**
- `jsonpretty()` - Pretty-print JSON with indentation
- `jsonminify()` - Minify JSON by removing whitespace
- `jsonformat()` - Format JSON to specified structure

---

### 2. **Pipe Operator (`|>`) Syntax**

Revolutionary pipe syntax for command chaining with full data flow:

```bash
# Basic piping
whoami() |> sha256()

# Multi-stage pipelines
echo("hello") |> sha256() |> rot13() |> base64()

# Module output to builtin
echo("data") |> module_name |> sha256()

# File writing with separators
whoami() |> "\n" |> file("~/logs.txt")
```

**Features:**
- Unlimited pipeline depth
- Data flows through each stage
- String literals with escape sequences support
- Module integration
- Result capture and display

---

### 3. **For Loop Implementation**

Native for loop syntax for batch operations and iterative commands:

```bash
# Basic loop
for x in 0..255 -> echo($x)

# With piping
for x in 1..10 -> echo($x) |> sha256()

# Module execution with variable substitution
for y in 0..256 -> ip-geolocation ip=192.168.1.$y

# Complex pipelines in loops
for i in 0..256 -> echo($i) |> sha256() |> "\n" |> file("~/hashes.txt")
```

**Features:**
- Integer range support (`START..END`)
- Variable substitution with `$variable`
- Pipe integration within loops
- Result collection and display
- Clean progress indication

---

### 4. **String Literals in Pipes**

Support for escape sequences within pipe chains:

```bash
# Newline separators
echo("output") |> "\n" |> file("~/file.txt")

# Tab indentation
result() |> "\t\t" |> file("~/output.txt")

# Supported escape sequences
\n  - Newline
\t  - Tab character
\r  - Carriage return
\\  - Literal backslash
```

---

### 5. **File I/O Through Pipes**

New `file()` builtin for persistent output:

```bash
# Write to file with automatic append
whoami() |> file("~/output.log")

# Home directory expansion
echo("data") |> file("~/logs/app.log")

# Automatic directory creation
result() |> file("~/data/nested/path/file.txt")

# Multiple appends
cmd1() |> file("~/log.txt")
cmd2() |> "\n" |> file("~/log.txt")
```

**Features:**
- Append mode (never overwrites)
- Home directory (`~/`) expansion
- Automatic parent directory creation
- Clean error handling

---

### 6. **Module Argument Piping**

Pass piped output as module arguments dynamically:

```bash
# Pattern: variable=$piped_value |> module
echo("192.168.1.1") |> ip-geolocation ip=$ip

# Module receives piped data
somedata() |> my-module arg=$data

# Integrates with loops
for x in 0..256 -> echo(192.168.1.$x) |> ip-geolocation ip=$data
```

---

## 📈 Technical Improvements

### Code Quality
- Modular builtin function registry
- Clean separation of concerns
- Comprehensive error handling
- Proper Go idioms throughout

### Performance
- Efficient pipe execution
- Optimized loop processing
- Smart result caching
- Memory-conscious operations

### Developer Experience
- Clear error messages
- Consistent function naming
- Intuitive syntax
- Full command chainability

---

## 🔧 Internal Architecture

### New Files & Modifications

**cli/builtins.go**
- Expanded from 60 to 113+ functions
- Organized into 10+ categories
- Full error handling for all operations
- Proper argument parsing

**cli/cli.go**
- Enhanced command parser
- Pipe operator implementation
- For loop engine
- Module integration layer
- String literal support

**cli/display.go**
- Updated help system
- Better formatting for 113+ builtins

---

## 🎯 Use Cases

### Network Administration
```bash
for ip in 192.168.1.1..192.168.1.254 -> check_port ip=$ip port=22 |> "\n" |> file("~/open_ports.txt")
```

### Data Processing
```bash
cat("data.json") |> jsonpretty() |> sha256() |> file("~/checksums.log")
```

### Security & Hashing
```bash
for hash in md5|sha256|sha512 -> echo("password") |> $hash() |> file("~/hash_tests.txt")
```

### Module Orchestration
```bash
enumerate() |> process() |> analyze() |> file("~/results.json")
```

---

## ✅ Quality Assurance

All features have been tested and validated:
- Multi-stage pipes working correctly
- For loops with variable substitution functioning
- String literals with escape sequences operational
- File I/O append mode verified
- Module argument piping integrated
- All 113+ builtins functional

---

## 🎉 Summary

LanManVan v1.0.0 represents a major milestone with:
- **113+ builtin functions** for data manipulation
- **Pipe operator** for elegant command chaining
- **For loops** for batch operations
- **File I/O** for persistent output
- **Full module integration** with dynamic arguments

This release provides a complete, production-ready platform for network administration, penetration testing, and data processing with powerful scripting capabilities.

---

**Build Information:**
- Version: 1.0.0
- Build Date: January 2, 2026
- Author: hmza
