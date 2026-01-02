package cli

// initDetailedBuiltins initializes builtin functions with detailed descriptions and examples
func (br *BuiltinRegistry) initDetailedBuiltins() {
	// File System
	br.registerDetailed("pwd",
		"Print working directory",
		"Displays the absolute path of the current working directory. Useful for determining your location in the file system and in scripts that need to know the current path before performing operations.",
		[]string{"pwd()", "path=$(pwd)"},
		br.cmdPwd)

	br.registerDetailed("ls",
		"List directory contents",
		"Lists all files and directories in a specified path or current directory. Shows directory entries with / suffix, helping you see the complete structure of any folder quickly.",
		[]string{"ls()", "ls(/tmp)", "files=$(ls(/home))"},
		br.cmdLs)

	br.registerDetailed("cd",
		"Change directory",
		"Changes the current working directory to the specified path. Returns the new directory path for confirmation. Essential for navigation and file operations in different locations.",
		[]string{"cd(/tmp)", "newdir=$(cd(/home))"},
		br.cmdCd)

	br.registerDetailed("mkdir",
		"Create directory",
		"Creates a new directory at the specified path. Works like mkdir -p, creating parent directories as needed. Returns OK on success or error message on failure.",
		[]string{"mkdir(/tmp/test)", "mkdir(./my_folder)"},
		br.cmdMkdir)

	br.registerDetailed("cat",
		"Read file contents",
		"Reads and displays the complete contents of a text file. Useful for viewing configuration files, scripts, or any text document without opening an editor.",
		[]string{"cat(/etc/hostname)", "content=$(cat(./config.txt))"},
		br.cmdCat)

	br.registerDetailed("rm",
		"Remove file or directory",
		"Deletes a file or directory recursively. Be careful as this operation is permanent. Returns OK on success and cannot be undone.",
		[]string{"rm(/tmp/oldfile)", "rm(./backup_folder)"},
		br.cmdRm)

	br.registerDetailed("cp",
		"Copy file or directory",
		"Copies a file from source to destination path. Creates a duplicate of the file at the target location. Useful for backups or preparing files for different locations.",
		[]string{"cp(/source/file.txt,/dest/file.txt)", "cp(./original,./backup)"},
		br.cmdCp)

	br.registerDetailed("mv",
		"Move or rename file",
		"Moves a file from source to destination or renames it if in the same directory. Atomic operation that preserves file metadata and is faster than copy-then-delete.",
		[]string{"mv(/old/path,/new/path)", "mv(./file.old,./file.new)"},
		br.cmdMv)

	br.registerDetailed("exists",
		"Check if file exists",
		"Determines whether a file or directory exists at the specified path. Returns 'true' if exists, 'false' otherwise. Useful for conditional operations in scripts.",
		[]string{"exists(/etc/passwd)", "exists(./config.json)"},
		br.cmdExists)

	br.registerDetailed("filesize",
		"Get file size in bytes",
		"Returns the size of a file in bytes. Useful for checking large files, bandwidth calculations, or storage management. Returns the numeric value only.",
		[]string{"filesize(/var/log/syslog)", "size=$(filesize(./data.bin))"},
		br.cmdFilesize)

	// System Info
	br.registerDetailed("whoami",
		"Get current user",
		"Returns the username of the currently logged-in user. Essential for verifying permissions and ensuring scripts run under the correct user context.",
		[]string{"whoami()", "user=$(whoami())"},
		br.cmdWhoami)

	br.registerDetailed("hostname",
		"Get system hostname",
		"Returns the network hostname of the current machine. Used for system identification, network configuration, and logging purposes.",
		[]string{"hostname()", "server=$(hostname())"},
		br.cmdHostname)

	br.registerDetailed("date",
		"Get current date and time",
		"Returns the current date and time in specified format (default: 2006-01-02 15:04:05). Supports custom formats and timezone operations for scheduling.",
		[]string{"date()", "timestamp=$(date(2006-01-02))", "now=$(date(15:04:05))"},
		br.cmdDate)

	br.registerDetailed("uname",
		"Get system information",
		"Returns detailed system information including kernel name, version, hardware platform, and operating system details. Useful for compatibility checks.",
		[]string{"uname()", "info=$(uname())"},
		br.cmdUname)

	br.registerDetailed("arch",
		"Get system architecture",
		"Returns the CPU architecture (x86_64, arm64, i386, etc.). Essential for determining binary compatibility and choosing appropriate dependencies.",
		[]string{"arch()", "cpu=$(arch())"},
		br.cmdArch)

	br.registerDetailed("ostype",
		"Get operating system type",
		"Returns the operating system type (linux, darwin, windows, etc.). Enables cross-platform script compatibility by detecting the runtime environment.",
		[]string{"ostype()", "system=$(ostype())"},
		br.cmdOstype)

	// Hashing
	br.registerDetailed("md5",
		"MD5 hash function",
		"Computes MD5 hash of input string (32 hexadecimal characters). Legacy hash function - deprecated for security purposes but still used for checksums.",
		[]string{"md5(hello)", "hash=$(md5(password123))"},
		br.cmdMd5)

	br.registerDetailed("sha256",
		"SHA256 hash function",
		"Computes SHA256 cryptographic hash of input string (64 hexadecimal characters). Modern, secure hash suitable for passwords, data integrity, and signatures.",
		[]string{"sha256(secret)", "sig=$(sha256(data))"},
		br.cmdSha256)

	br.registerDetailed("sha1",
		"SHA1 hash function",
		"Computes SHA1 hash of input string (40 hexadecimal characters). Largely deprecated due to collision vulnerabilities but still used for Git commits.",
		[]string{"sha1(input)", "checksum=$(sha1(file))"},
		br.cmdSha1)

	// String Operations
	br.registerDetailed("strlen",
		"Get string length",
		"Returns the number of characters in a string. Useful for validation, buffer sizing, and conditional logic based on string length.",
		[]string{"strlen(hello)", "len=$(strlen(password))"},
		br.cmdStrlen)

	br.registerDetailed("toupper",
		"Convert to uppercase",
		"Converts all lowercase letters in a string to uppercase. Leaves numbers, special characters, and already-uppercase letters unchanged.",
		[]string{"toupper(hello)", "name=$(toupper(john))"},
		br.cmdToupper)

	br.registerDetailed("tolower",
		"Convert to lowercase",
		"Converts all uppercase letters in a string to lowercase. Useful for case-insensitive comparisons and normalizing user input.",
		[]string{"tolower(HELLO)", "text=$(tolower(MyData))"},
		br.cmdTolower)

	br.registerDetailed("reverse",
		"Reverse a string",
		"Reverses the order of characters in a string. Useful for palindrome checks, data obfuscation, and certain cryptographic operations.",
		[]string{"reverse(hello)", "backwards=$(reverse(12345))"},
		br.cmdReverse)

	br.registerDetailed("trim",
		"Remove leading/trailing whitespace",
		"Removes whitespace (spaces, tabs, newlines) from both ends of a string. Useful for cleaning user input and preparing data for processing.",
		[]string{"trim(  hello  )", "clean=$(trim(input))"},
		br.cmdTrim)

	br.registerDetailed("startswith",
		"Check string starts with substring",
		"Returns 'true' if string begins with the specified substring, 'false' otherwise. Case-sensitive comparison useful for prefix matching.",
		[]string{"startswith(hello,he)", "startswith(filename.txt,.txt)"},
		br.cmdStartswith)

	br.registerDetailed("endswith",
		"Check string ends with substring",
		"Returns 'true' if string ends with the specified substring, 'false' otherwise. Useful for file extension checking and suffix matching.",
		[]string{"endswith(hello.txt,.txt)", "endswith(domain.com,.com)"},
		br.cmdEndswith)

	br.registerDetailed("contains",
		"Check if string contains substring",
		"Returns 'true' if substring exists anywhere in the string, 'false' otherwise. Useful for searching and pattern detection in text.",
		[]string{"contains(hello world,world)", "contains(email@example.com,@)"},
		br.cmdContains)

	br.registerDetailed("substr",
		"Extract substring from string",
		"Extracts a portion of string starting at index (optional end). Returns substring between start and end positions or until end of string.",
		[]string{"substr(hello,1,4)", "substr(12345,2)"},
		br.cmdSubstr)

	br.registerDetailed("replace",
		"Replace text in string",
		"Replaces all occurrences of a substring with another substring. Useful for text transformation and data sanitization.",
		[]string{"replace(hello world,world,there)", "replace(path/to/file,/,\\\\)"},
		br.cmdReplace)

	// Network Validation
	br.registerDetailed("isipv4",
		"Validate IPv4 address",
		"Returns 'true' if input is a valid IPv4 address (0.0.0.0 to 255.255.255.255), 'false' otherwise. Essential for network configuration validation.",
		[]string{"isipv4(192.168.1.1)", "isipv4(999.999.999.999)"},
		br.cmdIsIPv4)

	br.registerDetailed("isipv6",
		"Validate IPv6 address",
		"Returns 'true' if input is a valid IPv6 address, 'false' otherwise. Supports compressed and full IPv6 notation for modern network validation.",
		[]string{"isipv6(2001:db8::1)", "isipv6(::1)"},
		br.cmdIsIPv6)

	br.registerDetailed("isemail",
		"Validate email address",
		"Returns 'true' if input matches email format pattern, 'false' otherwise. Basic validation useful for user input and form processing.",
		[]string{"isemail(user@example.com)", "isemail(invalid-email)"},
		br.cmdIsEmail)

	br.registerDetailed("isurl",
		"Validate URL format",
		"Returns 'true' if input is a valid URL with proper scheme and host, 'false' otherwise. Useful for link validation and web scraping.",
		[]string{"isurl(https://example.com)", "isurl(ftp://files.server.org)"},
		br.cmdIsUrl)

	br.registerDetailed("ismac",
		"Validate MAC address",
		"Returns 'true' if input is a valid MAC address (48-bit), 'false' otherwise. Supports colon and hyphen separated formats.",
		[]string{"ismac(00:1A:2B:3C:4D:5E)", "ismac(00-1A-2B-3C-4D-5E)"},
		br.cmdIsMac)

	br.registerDetailed("isdomain",
		"Validate domain name",
		"Returns 'true' if input is a valid domain name format, 'false' otherwise. Useful for DNS configuration and web filtering.",
		[]string{"isdomain(example.com)", "isdomain(sub.domain.co.uk)"},
		br.cmdIsDomain)

	br.registerDetailed("isport",
		"Validate port number",
		"Returns 'true' if input is valid port number (1-65535), 'false' otherwise. Ensures port is in valid range for network services.",
		[]string{"isport(8080)", "isport(99999)"},
		br.cmdIsPort)

	br.registerDetailed("iscdr",
		"Validate CIDR notation",
		"Returns 'true' if input is valid CIDR notation (e.g., 192.168.0.0/24), 'false' otherwise. Used for subnet validation.",
		[]string{"iscdr(192.168.0.0/24)", "iscdr(10.0.0.0/8)"},
		br.cmdIsCIDR)

	// Encoding
	br.registerDetailed("base64",
		"Base64 encode/decode",
		"Encodes string to Base64 or decodes Base64 input. Auto-detects if input is valid Base64 and returns decoded version, otherwise encodes.",
		[]string{"base64(hello)", "base64(aGVsbG8=)"},
		br.cmdBase64)

	br.registerDetailed("hex",
		"Hex encode/decode",
		"Encodes string to hexadecimal or decodes hex input. Auto-detects hex format and returns decoded version when applicable.",
		[]string{"hex(hello)", "hex(68656c6c6f)"},
		br.cmdHex)

	br.registerDetailed("url",
		"URL encode/decode",
		"URL-encodes a string (percent-encoding) for use in URLs. Replaces special characters with %XX notation for safe transmission.",
		[]string{"url(hello world)", "url(user@example.com)"},
		br.cmdUrl)

	br.registerDetailed("json",
		"JSON format/parse",
		"Parses JSON string and pretty-prints it with indentation. Validates JSON format and makes it human-readable.",
		[]string{"json({\"key\":\"value\"})", "json([1,2,3])"},
		br.cmdJson)

	// Math
	br.registerDetailed("calc",
		"Simple calculator",
		"Performs basic arithmetic operations: addition (+), subtraction (-), multiplication (*), division (/), modulo (%). Supports integer and float operations.",
		[]string{"calc(10,+,5)", "calc(20,-,8)", "calc(3,*,7)"},
		br.cmdCalc)

	br.registerDetailed("abs",
		"Absolute value",
		"Returns the absolute value (removes negative sign) of a number. Always returns non-negative number useful for distance calculations.",
		[]string{"abs(-5)", "abs(3.14)"},
		br.cmdAbs)

	br.registerDetailed("min",
		"Find minimum value",
		"Returns the smallest value from multiple number arguments. Useful for finding thresholds and limits.",
		[]string{"min(5,3,8,1)", "min(100,50,75)"},
		br.cmdMin)

	br.registerDetailed("max",
		"Find maximum value",
		"Returns the largest value from multiple number arguments. Useful for finding peaks and upper limits.",
		[]string{"max(5,3,8,1)", "max(100,50,75)"},
		br.cmdMax)

	br.registerDetailed("sum",
		"Sum multiple values",
		"Adds all number arguments together and returns the total. Useful for aggregation and financial calculations.",
		[]string{"sum(1,2,3,4,5)", "sum(10,20,30)"},
		br.cmdSum)

	br.registerDetailed("avg",
		"Calculate average",
		"Calculates the arithmetic mean of multiple number arguments. Returns average value with two decimal places.",
		[]string{"avg(1,2,3,4,5)", "avg(10,20,30)"},
		br.cmdAvg)

	br.registerDetailed("random",
		"Generate random number",
		"Generates a pseudo-random number between 0 and 999999. Seeded from current time, not cryptographically secure.",
		[]string{"random()", "id=$(random())"},
		br.cmdRandom)

	// Utilities
	br.registerDetailed("uuid",
		"Generate UUID",
		"Generates a pseudo-UUID (version 4 format). Returns unique identifier suitable for database keys and request IDs.",
		[]string{"uuid()", "id=$(uuid())"},
		br.cmdUuid)

	br.registerDetailed("timestamp",
		"Get Unix timestamp",
		"Returns current time as Unix timestamp in various formats: unix (seconds), milli (milliseconds), nano (nanoseconds), or custom format.",
		[]string{"timestamp()", "timestamp(unix)", "timestamp(2006-01-02)"},
		br.cmdTimestamp)

	br.registerDetailed("echo",
		"Print text output",
		"Outputs text string exactly as provided. Useful for concatenation and simple text manipulation in expressions.",
		[]string{"echo(hello world)", "echo(debugging info)"},
		br.cmdEcho)

	br.registerDetailed("sleep",
		"Sleep for seconds",
		"Pauses execution for specified number of seconds. Useful for rate limiting, retry delays, and timing-based operations.",
		[]string{"sleep(5)", "sleep(2)"},
		br.cmdSleep)

	br.registerDetailed("readfile",
		"Read file content",
		"Reads entire file content and returns as string. Useful for loading configuration, templates, or data files.",
		[]string{"readfile(/etc/hostname)", "content=$(readfile(./config.txt))"},
		br.cmdReadfile)

	br.registerDetailed("writefile",
		"Write content to file",
		"Creates or overwrites a file with specified content. Returns OK on success, useful for creating configurations.",
		[]string{"writefile(/tmp/test.txt,hello world)", "writefile(./config,settings)"},
		br.cmdWritefile)

	br.registerDetailed("randomstr",
		"Generate random string",
		"Generates random string of specified length (default 16). Useful for generating temporary passwords, tokens, and unique identifiers.",
		[]string{"randomstr()", "randomstr(32)", "token=$(randomstr(20))"},
		br.cmdRandomstr)

	br.registerDetailed("genpass",
		"Generate password",
		"Generates cryptographically random password with mixed characters, numbers, and symbols. Default length 16 characters.",
		[]string{"genpass()", "genpass(24)", "pwd=$(genpass(32))"},
		br.cmdGenpass)
}
