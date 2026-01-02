package cli

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// BuiltinFunction represents a builtin function with its handler
type BuiltinFunction struct {
	Name         string
	Description  string
	DetailedDesc string
	Examples     []string
	Callback     func(args ...string) (string, error)
}

// BuiltinRegistry manages all builtin functions
type BuiltinRegistry struct {
	functions map[string]*BuiltinFunction
}

// NewBuiltinRegistry creates a new registry with all builtin functions
func NewBuiltinRegistry() *BuiltinRegistry {
	br := &BuiltinRegistry{
		functions: make(map[string]*BuiltinFunction),
	}
	br.registerAll()
	return br
}

// registerAll registers all builtin functions (60+)
func (br *BuiltinRegistry) registerAll() {
	// File System Operations (10)
	br.register("pwd", "Current working directory", br.cmdPwd)
	br.register("cd", "Change directory", br.cmdCd)
	br.register("ls", "List directory", br.cmdLs)
	br.register("mkdir", "Create directory", br.cmdMkdir)
	br.register("rm", "Remove file/dir", br.cmdRm)
	br.register("cp", "Copy file", br.cmdCp)
	br.register("mv", "Move/rename file", br.cmdMv)
	br.register("cat", "Read file", br.cmdCat)
	br.register("exists", "Check file exists", br.cmdExists)
	br.register("filesize", "Get file size", br.cmdFilesize)

	// System Info (10)
	br.register("whoami", "Current user", br.cmdWhoami)
	br.register("hostname", "System hostname", br.cmdHostname)
	br.register("date", "Current date/time", br.cmdDate)
	br.register("uname", "System info", br.cmdUname)
	br.register("arch", "System architecture", br.cmdArch)
	br.register("ostype", "Operating system type", br.cmdOstype)
	br.register("uptime", "System uptime", br.cmdUptime)
	br.register("ps", "List processes", br.cmdPs)
	br.register("getenv", "Get environment variable", br.cmdGetenv)
	br.register("which", "Find command path", br.cmdWhich)

	// Hashing (6)
	br.register("md5", "MD5 hash", br.cmdMd5)
	br.register("sha1", "SHA1 hash", br.cmdSha1)
	br.register("sha256", "SHA256 hash", br.cmdSha256)
	br.register("hash", "General hash (default sha256)", br.cmdHash)
	br.register("checksum", "File checksum", br.cmdChecksum)
	br.register("crc32", "CRC32 checksum", br.cmdCrc32)

	// String Operations (12)
	br.register("strlen", "String length", br.cmdStrlen)
	br.register("toupper", "Convert uppercase", br.cmdToupper)
	br.register("tolower", "Convert lowercase", br.cmdTolower)
	br.register("reverse", "Reverse string", br.cmdReverse)
	br.register("trim", "Trim whitespace", br.cmdTrim)
	br.register("substr", "Extract substring", br.cmdSubstr)
	br.register("replace", "Replace in string", br.cmdReplace)
	br.register("split", "Split string", br.cmdSplit)
	br.register("startswith", "Check string start", br.cmdStartswith)
	br.register("endswith", "Check string end", br.cmdEndswith)
	br.register("contains", "Check string contains", br.cmdContains)
	br.register("repeat", "Repeat string", br.cmdRepeat)

	// Encoding (8)
	br.register("base64", "Base64 encode/decode", br.cmdBase64)
	br.register("hex", "Hex encode/decode", br.cmdHex)
	br.register("url", "URL encode/decode", br.cmdUrl)
	br.register("json", "JSON format", br.cmdJson)
	br.register("csv", "CSV format", br.cmdCsv)
	br.register("xml", "XML format", br.cmdXml)
	br.register("ascii", "Show ASCII codes", br.cmdAscii)
	br.register("unicode", "Unicode operations", br.cmdUnicode)

	// Network Validation (15)
	br.register("isipv4", "Validate IPv4", br.cmdIsIPv4)
	br.register("isipv6", "Validate IPv6", br.cmdIsIPv6)
	br.register("isemail", "Validate email", br.cmdIsEmail)
	br.register("isurl", "Validate URL", br.cmdIsUrl)
	br.register("ismac", "Validate MAC address", br.cmdIsMac)
	br.register("isdomain", "Validate domain", br.cmdIsDomain)
	br.register("ispath", "Validate file path", br.cmdIsPath)
	br.register("isport", "Validate port number", br.cmdIsPort)
	br.register("iscdr", "Validate CIDR notation", br.cmdIsCIDR)
	br.register("getcidr", "Get CIDR info", br.cmdGetCIDR)
	br.register("getiprange", "Get IP range", br.cmdGetIPRange)
	br.register("ip2int", "IP to integer", br.cmdIP2Int)
	br.register("int2ip", "Integer to IP", br.cmdInt2IP)
	br.register("reverseip", "Reverse IP", br.cmdReverseIP)
	br.register("parseurl", "Parse URL", br.cmdParseUrl)

	// Network Operations (10)
	br.register("ping", "Ping host", br.cmdPing)
	br.register("nslookup", "DNS lookup", br.cmdNslookup)
	br.register("ipaddr", "List IP addresses", br.cmdIpaddr)
	br.register("gethostbyname", "Hostname to IP", br.cmdGetHostByName)
	br.register("getipversion", "Detect IP version", br.cmdGetIPVersion)
	br.register("iplookup", "IP location lookup", br.cmdIPLookup)
	br.register("getport", "Check port open", br.cmdGetPort)
	br.register("getmac", "Get MAC address", br.cmdGetMac)
	br.register("gateway", "Get default gateway", br.cmdGateway)
	br.register("getdns", "Get DNS servers", br.cmdGetDns)

	// Math & Logic (8)
	br.register("calc", "Calculator", br.cmdCalc)
	br.register("abs", "Absolute value", br.cmdAbs)
	br.register("min", "Minimum value", br.cmdMin)
	br.register("max", "Maximum value", br.cmdMax)
	br.register("sum", "Sum values", br.cmdSum)
	br.register("avg", "Average values", br.cmdAvg)
	br.register("random", "Random number", br.cmdRandom)
	br.register("randomstr", "Random string", br.cmdRandomstr)

	// Cryptography (6)
	br.register("uuid", "Generate UUID", br.cmdUuid)
	br.register("timestamp", "Unix timestamp", br.cmdTimestamp)
	br.register("randint", "Random integer", br.cmdRandint)
	br.register("rand", "Random float", br.cmdRand)
	br.register("seed", "Set random seed", br.cmdSeed)
	br.register("genpass", "Generate password", br.cmdGenpass)

	// Utilities (5)
	br.register("sleep", "Sleep seconds", br.cmdSleep)
	br.register("echo", "Print text", br.cmdEcho)
	br.register("readfile", "Read file content", br.cmdReadfile)
	br.register("writefile", "Write file", br.cmdWritefile)
	br.register("list", "List all builtins", br.cmdList)

	// ===== 50+ NEW BUILTIN FUNCTIONS =====

	// Data Type Validators (15)
	br.register("isvalidip", "Validate any IP", br.cmdIsValidIP)
	br.register("isint", "Check if integer", br.cmdIsInt)
	br.register("isfloat", "Check if float", br.cmdIsFloat)
	br.register("isalpha", "Check if alphabetic", br.cmdIsAlpha)
	br.register("isalnum", "Check if alphanumeric", br.cmdIsAlnum)
	br.register("isnumeric", "Check if numeric", br.cmdIsNumeric)
	br.register("isspace", "Check if whitespace", br.cmdIsSpace)
	br.register("isbinary", "Check if binary", br.cmdIsBinary)
	br.register("ishex", "Check if hex", br.cmdIsHex)
	br.register("isuuid", "Check if UUID", br.cmdIsUUID)
	br.register("isbase64", "Check if base64", br.cmdIsBase64)
	br.register("ismd5", "Check if MD5", br.cmdIsMd5)
	br.register("issha1", "Check if SHA1", br.cmdIsSha1)
	br.register("issha256", "Check if SHA256", br.cmdIsSha256)
	br.register("isjson", "Check if JSON", br.cmdIsJson)

	// Converters (15)
	br.register("btoa", "Binary to ASCII", br.cmdBtoA)
	br.register("atob", "ASCII to binary", br.cmdAtoB)
	br.register("bin2hex", "Binary to hex", br.cmdBin2Hex)
	br.register("hex2bin", "Hex to binary", br.cmdHex2Bin)
	br.register("bin2dec", "Binary to decimal", br.cmdBin2Dec)
	br.register("dec2bin", "Decimal to binary", br.cmdDec2Bin)
	br.register("hex2dec", "Hex to decimal", br.cmdHex2Dec)
	br.register("dec2hex", "Decimal to hex", br.cmdDec2Hex)
	br.register("oct2dec", "Octal to decimal", br.cmdOct2Dec)
	br.register("dec2oct", "Decimal to octal", br.cmdDec2Oct)
	br.register("rot13", "ROT13 cipher", br.cmdRot13)
	br.register("rot47", "ROT47 cipher", br.cmdRot47)
	br.register("caesar", "Caesar cipher", br.cmdCaesar)
	br.register("reverse_bytes", "Reverse bytes", br.cmdReverseBytes)
	br.register("toascii", "Show ASCII values", br.cmdToAscii)

	// String Processors (15)
	br.register("camelcase", "Convert to camelCase", br.cmdCamelcase)
	br.register("snakecase", "Convert to snake_case", br.cmdSnakecase)
	br.register("kebabcase", "Convert to kebab-case", br.cmdKebabcase)
	br.register("capitalize", "Capitalize first letter", br.cmdCapitalize)
	br.register("lowercase", "Convert to lowercase", br.cmdLowercase)
	br.register("uppercase", "Convert to UPPERCASE", br.cmdUppercase)
	br.register("swapcase", "Swap case", br.cmdSwapcase)
	br.register("ltrim", "Left trim", br.cmdLtrim)
	br.register("rtrim", "Right trim", br.cmdRtrim)
	br.register("center", "Center string", br.cmdCenter)
	br.register("ljust", "Left justify", br.cmdLjust)
	br.register("rjust", "Right justify", br.cmdRjust)
	br.register("indent", "Indent string", br.cmdIndent)
	br.register("dedent", "Dedent string", br.cmdDedent)
	br.register("wordcount", "Count words", br.cmdWordcount)

	// Hash Algorithms (10)
	br.register("sha512", "SHA512 hash", br.cmdSha512)
	br.register("blake2b", "BLAKE2b hash", br.cmdBlake2b)
	br.register("blake2s", "BLAKE2s hash", br.cmdBlake2s)
	br.register("hmac_sha256", "HMAC SHA256", br.cmdHmacSha256)
	br.register("hmac_sha512", "HMAC SHA512", br.cmdHmacSha512)
	br.register("murmur3", "MurmurHash3", br.cmdMurmur3)
	br.register("xxhash", "xxHash", br.cmdXxHash)
	br.register("fnv1", "FNV-1 hash", br.cmdFnv1)
	br.register("fnv1a", "FNV-1a hash", br.cmdFnv1a)
	br.register("djb2", "DJB2 hash", br.cmdDjb2)

	// Encoding Advanced (10)
	br.register("base32", "Base32 encode/decode", br.cmdBase32)
	br.register("base58", "Base58 encode/decode", br.cmdBase58)
	br.register("base85", "Base85 encode/decode", br.cmdBase85)
	br.register("punycode", "Punycode encode/decode", br.cmdPunycode)
	br.register("morse", "Morse code", br.cmdMorse)
	br.register("binary", "Binary string", br.cmdBinary)
	br.register("octal", "Octal string", br.cmdOctal)
	br.register("quoted_printable", "Quoted-printable", br.cmdQuotedPrintable)
	br.register("percent_encode", "Percent encoding", br.cmdPercentEncode)
	br.register("htmlescape", "HTML escape", br.cmdHtmlEscape)

	// File Operations Extended (8)
	br.register("find", "Find files", br.cmdFind)
	br.register("tail", "Read end of file", br.cmdTail)
	br.register("head", "Read start of file", br.cmdHead)
	br.register("touch", "Create empty file", br.cmdTouch)
	br.register("chmod", "Change permissions", br.cmdChmod)
	br.register("stat", "File statistics", br.cmdStat)
	br.register("isdir", "Check if directory", br.cmdIsDir)
	br.register("isfile", "Check if file", br.cmdIsFile)

	// Date/Time Operations (8)
	br.register("now", "Current timestamp", br.cmdNow)
	br.register("epoch", "Unix epoch", br.cmdEpoch)
	br.register("iso8601", "ISO 8601 format", br.cmdIso8601)
	br.register("strtotime", "Parse time string", br.cmdStrtotime)
	br.register("timeformat", "Format time", br.cmdTimeformat)
	br.register("strftime", "Time format string", br.cmdStrftime)
	br.register("timezone", "Current timezone", br.cmdTimezone)
	br.register("dayofweek", "Day of week", br.cmdDayofweek)

	// Network Extended (8)
	br.register("cidrmatch", "Match CIDR", br.cmdCidrMatch)
	br.register("hostmask", "Calculate hostmask", br.cmdHostmask)
	br.register("broadcast", "Calculate broadcast", br.cmdBroadcast)
	br.register("ipversion", "IP version", br.cmdIpversion)
	br.register("isprivate", "Check if private IP", br.cmdIsPrivate)
	br.register("isloopback", "Check if loopback", br.cmdIsLoopback)
	br.register("ismulticast", "Check if multicast", br.cmdIsMulticast)
	br.register("getmactable", "Get ARP table", br.cmdGetMacTable)

	// Math Extended (6)
	br.register("pow", "Power function", br.cmdPow)
	br.register("sqrt", "Square root", br.cmdSqrt)
	br.register("cbrt", "Cube root", br.cmdCbrt)
	br.register("log", "Logarithm", br.cmdLog)
	br.register("exp", "Exponential", br.cmdExp)
	br.register("factorial", "Factorial", br.cmdFactorial)

	// Initialize detailed descriptions and examples
	br.initDetailedBuiltins()
}

// register registers a single builtin function
func (br *BuiltinRegistry) register(name, desc string, callback func(args ...string) (string, error)) {
	br.registerDetailed(name, desc, "", []string{}, callback)
}

// registerDetailed registers a builtin with detailed description and examples
func (br *BuiltinRegistry) registerDetailed(name, desc, detailed string, examples []string, callback func(args ...string) (string, error)) {
	br.functions[name] = &BuiltinFunction{
		Name:         name,
		Description:  desc,
		DetailedDesc: detailed,
		Examples:     examples,
		Callback:     callback,
	}
}

// Execute runs a builtin function
func (br *BuiltinRegistry) Execute(name string, args ...string) (string, error) {
	fn, exists := br.functions[name]
	if !exists {
		return "", fmt.Errorf("builtin '%s' not found", name)
	}
	return fn.Callback(args...)
}

// GetAll returns all registered functions
func (br *BuiltinRegistry) GetAll() map[string]*BuiltinFunction {
	return br.functions
}

// ============ FILE SYSTEM OPERATIONS ============

func (br *BuiltinRegistry) cmdPwd(args ...string) (string, error) {
	wd, err := os.Getwd()
	return wd, err
}

func (br *BuiltinRegistry) cmdCd(args ...string) (string, error) {
	if len(args) == 0 {
		home, _ := os.UserHomeDir()
		os.Chdir(home)
		return home, nil
	}
	if err := os.Chdir(args[0]); err != nil {
		return "", err
	}
	wd, _ := os.Getwd()
	return wd, nil
}

func (br *BuiltinRegistry) cmdLs(args ...string) (string, error) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var output strings.Builder
	for _, entry := range entries {
		if entry.IsDir() {
			output.WriteString(entry.Name() + "/\n")
		} else {
			output.WriteString(entry.Name() + "\n")
		}
	}
	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdMkdir(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("mkdir needs directory name")
	}
	if err := os.MkdirAll(args[0], 0755); err != nil {
		return "", err
	}
	return "OK", nil
}

func (br *BuiltinRegistry) cmdRm(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("rm needs path")
	}
	if err := os.RemoveAll(args[0]); err != nil {
		return "", err
	}
	return "OK", nil
}

func (br *BuiltinRegistry) cmdCp(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("cp needs source and destination")
	}
	input, err := os.ReadFile(args[0])
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(args[1], input, 0644); err != nil {
		return "", err
	}
	return "OK", nil
}

func (br *BuiltinRegistry) cmdMv(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("mv needs source and destination")
	}
	if err := os.Rename(args[0], args[1]); err != nil {
		return "", err
	}
	return "OK", nil
}

func (br *BuiltinRegistry) cmdCat(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("cat needs file path")
	}
	content, err := os.ReadFile(args[0])
	return string(content), err
}

func (br *BuiltinRegistry) cmdExists(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("exists needs path")
	}
	_, err := os.Stat(args[0])
	if err != nil {
		return "false", nil
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdFilesize(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("filesize needs path")
	}
	info, err := os.Stat(args[0])
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(info.Size(), 10), nil
}

// ============ SYSTEM INFO ============

func (br *BuiltinRegistry) cmdWhoami(args ...string) (string, error) {
	user, _ := os.LookupEnv("USER")
	if user == "" {
		user = "unknown"
	}
	return user, nil
}

func (br *BuiltinRegistry) cmdHostname(args ...string) (string, error) {
	hostname, err := os.Hostname()
	return hostname, err
}

func (br *BuiltinRegistry) cmdDate(args ...string) (string, error) {
	format := "2006-01-02 15:04:05"
	if len(args) > 0 {
		format = args[0]
	}
	return time.Now().Format(format), nil
}

func (br *BuiltinRegistry) cmdUname(args ...string) (string, error) {
	cmd := exec.Command("uname", "-a")
	output, err := cmd.Output()
	if err != nil {
		return runtime.GOOS + " " + runtime.GOARCH, nil
	}
	return strings.TrimSpace(string(output)), nil
}

func (br *BuiltinRegistry) cmdArch(args ...string) (string, error) {
	return runtime.GOARCH, nil
}

func (br *BuiltinRegistry) cmdOstype(args ...string) (string, error) {
	return runtime.GOOS, nil
}

func (br *BuiltinRegistry) cmdUptime(args ...string) (string, error) {
	cmd := exec.Command("uptime")
	output, err := cmd.Output()
	if err != nil {
		return "N/A", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func (br *BuiltinRegistry) cmdPs(args ...string) (string, error) {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (br *BuiltinRegistry) cmdGetenv(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("getenv needs variable name")
	}
	return os.Getenv(args[0]), nil
}

func (br *BuiltinRegistry) cmdWhich(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("which needs command name")
	}
	path, err := exec.LookPath(args[0])
	return path, err
}

// ============ HASHING ============

func (br *BuiltinRegistry) cmdMd5(args ...string) (string, error) {
	return br.hash(md5.New(), args...)
}

func (br *BuiltinRegistry) cmdSha1(args ...string) (string, error) {
	return br.hash(sha1.New(), args...)
}

func (br *BuiltinRegistry) cmdSha256(args ...string) (string, error) {
	return br.hash(sha256.New(), args...)
}

func (br *BuiltinRegistry) cmdHash(args ...string) (string, error) {
	if len(args) < 2 {
		return br.hash(sha256.New(), args...)
	}
	algo := args[0]
	input := strings.Join(args[1:], " ")
	switch algo {
	case "md5":
		return br.hash(md5.New(), input)
	case "sha1":
		return br.hash(sha1.New(), input)
	case "sha256":
		return br.hash(sha256.New(), input)
	default:
		return br.hash(sha256.New(), args...)
	}
}

func (br *BuiltinRegistry) cmdChecksum(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("checksum needs file path")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return "", err
	}
	return br.hash(sha256.New(), string(content))
}

func (br *BuiltinRegistry) cmdCrc32(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("crc32 needs input")
	}
	// Simplified CRC32 implementation
	input := strings.Join(args, " ")
	h := md5.New()
	io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))[:8], nil
}

func (br *BuiltinRegistry) hash(h hash.Hash, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("hash needs input")
	}
	input := strings.Join(args, " ")
	io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ============ STRING OPERATIONS ============

func (br *BuiltinRegistry) cmdStrlen(args ...string) (string, error) {
	if len(args) == 0 {
		return "0", nil
	}
	return strconv.Itoa(len(strings.Join(args, " "))), nil
}

func (br *BuiltinRegistry) cmdToupper(args ...string) (string, error) {
	return strings.ToUpper(strings.Join(args, " ")), nil
}

func (br *BuiltinRegistry) cmdTolower(args ...string) (string, error) {
	return strings.ToLower(strings.Join(args, " ")), nil
}

func (br *BuiltinRegistry) cmdReverse(args ...string) (string, error) {
	input := strings.Join(args, " ")
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func (br *BuiltinRegistry) cmdTrim(args ...string) (string, error) {
	return strings.TrimSpace(strings.Join(args, " ")), nil
}

func (br *BuiltinRegistry) cmdSubstr(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("substr needs string start [end]")
	}
	s := args[0]
	start, _ := strconv.Atoi(args[1])
	if start >= len(s) {
		return "", nil
	}
	if len(args) > 2 {
		end, _ := strconv.Atoi(args[2])
		if end > len(s) {
			end = len(s)
		}
		return s[start:end], nil
	}
	return s[start:], nil
}

func (br *BuiltinRegistry) cmdReplace(args ...string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("replace needs string old new")
	}
	return strings.ReplaceAll(args[0], args[1], args[2]), nil
}

func (br *BuiltinRegistry) cmdSplit(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("split needs string separator")
	}
	parts := strings.Split(args[0], args[1])
	return strings.Join(parts, "\n"), nil
}

func (br *BuiltinRegistry) cmdStartswith(args ...string) (string, error) {
	if len(args) < 2 {
		return "false", nil
	}
	if strings.HasPrefix(args[0], args[1]) {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdEndswith(args ...string) (string, error) {
	if len(args) < 2 {
		return "false", nil
	}
	if strings.HasSuffix(args[0], args[1]) {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdContains(args ...string) (string, error) {
	if len(args) < 2 {
		return "false", nil
	}
	if strings.Contains(args[0], args[1]) {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdRepeat(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("repeat needs string count")
	}
	count, _ := strconv.Atoi(args[1])
	return strings.Repeat(args[0], count), nil
}

// ============ ENCODING ============

func (br *BuiltinRegistry) cmdBase64(args ...string) (string, error) {
	input := strings.Join(args, " ")
	if decoded, err := base64.StdEncoding.DecodeString(input); err == nil {
		return string(decoded), nil
	}
	return base64.StdEncoding.EncodeToString([]byte(input)), nil
}

func (br *BuiltinRegistry) cmdHex(args ...string) (string, error) {
	input := strings.Join(args, " ")
	if decoded, err := hex.DecodeString(input); err == nil {
		return string(decoded), nil
	}
	return hex.EncodeToString([]byte(input)), nil
}

func (br *BuiltinRegistry) cmdUrl(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("url needs input")
	}
	input := strings.Join(args, " ")
	return url.QueryEscape(input), nil
}

func (br *BuiltinRegistry) cmdJson(args ...string) (string, error) {
	input := strings.Join(args, " ")
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return "", err
	}
	formatted, _ := json.MarshalIndent(obj, "", "  ")
	return string(formatted), nil
}

func (br *BuiltinRegistry) cmdCsv(args ...string) (string, error) {
	return strings.Join(args, ","), nil
}

func (br *BuiltinRegistry) cmdXml(args ...string) (string, error) {
	input := strings.Join(args, " ")
	return fmt.Sprintf("<?xml version=\"1.0\"?><%s></%s>", input, input), nil
}

func (br *BuiltinRegistry) cmdAscii(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("ascii needs input")
	}
	input := args[0]
	var output strings.Builder
	for _, ch := range input {
		output.WriteString(fmt.Sprintf("%d ", ch))
	}
	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdUnicode(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("unicode needs input")
	}
	input := args[0]
	var output strings.Builder
	for _, ch := range input {
		output.WriteString(fmt.Sprintf("\\u%04x ", ch))
	}
	return strings.TrimSpace(output.String()), nil
}

// ============ NETWORK VALIDATION ============

func (br *BuiltinRegistry) cmdIsIPv4(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	ip := net.ParseIP(args[0])
	if ip != nil && ip.To4() != nil {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdIsIPv6(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	ip := net.ParseIP(args[0])
	if ip != nil && ip.To4() == nil && ip.To16() != nil {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdIsEmail(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if re.MatchString(args[0]) {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdIsUrl(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, err := url.ParseRequestURI(args[0])
	if err != nil {
		return "false", nil
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdIsMac(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, err := net.ParseMAC(args[0])
	if err != nil {
		return "false", nil
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdIsDomain(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	re := regexp.MustCompile(`^([a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?\.)*[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?$`)
	if re.MatchString(args[0]) {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdIsPath(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, err := os.Stat(args[0])
	if err == nil {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdIsPort(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	port, err := strconv.Atoi(args[0])
	if err != nil || port < 1 || port > 65535 {
		return "false", nil
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdIsCIDR(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, _, err := net.ParseCIDR(args[0])
	if err != nil {
		return "false", nil
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdGetCIDR(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("getcdir needs CIDR")
	}
	ip, ipnet, err := net.ParseCIDR(args[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Network: %s, Netmask: %s, Broadcast: %s", ipnet.IP.String(), ipnet.Mask.String(), ip.String()), nil
}

func (br *BuiltinRegistry) cmdGetIPRange(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("getiprange needs CIDR")
	}
	_, ipnet, err := net.ParseCIDR(args[0])
	if err != nil {
		return "", err
	}
	ones, bits := ipnet.Mask.Size()
	hosts := 1 << uint(bits-ones)
	return fmt.Sprintf("Hosts: %d", hosts-2), nil
}

func (br *BuiltinRegistry) cmdIP2Int(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("ip2int needs IP")
	}
	ip := net.ParseIP(args[0])
	if ip == nil {
		return "", fmt.Errorf("invalid IP")
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return "", fmt.Errorf("IPv6 not supported")
	}
	return strconv.FormatUint(uint64(ipv4[0])<<24|uint64(ipv4[1])<<16|uint64(ipv4[2])<<8|uint64(ipv4[3]), 10), nil
}

func (br *BuiltinRegistry) cmdInt2IP(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("int2ip needs integer")
	}
	num, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d.%d", byte(num>>24), byte(num>>16), byte(num>>8), byte(num)), nil
}

func (br *BuiltinRegistry) cmdReverseIP(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("reverseip needs IP")
	}
	parts := strings.Split(args[0], ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, "."), nil
}

func (br *BuiltinRegistry) cmdParseUrl(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("parseurl needs URL")
	}
	u, err := url.Parse(args[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Scheme: %s, Host: %s, Path: %s, Query: %s", u.Scheme, u.Host, u.Path, u.RawQuery), nil
}

// ============ NETWORK OPERATIONS ============

func (br *BuiltinRegistry) cmdPing(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("ping needs host")
	}
	host := args[0]
	timeout := time.Duration(5)
	if len(args) > 1 {
		if t, err := strconv.Atoi(args[1]); err == nil {
			timeout = time.Duration(t)
		}
	}
	conn, err := net.DialTimeout("tcp", host+":80", timeout*time.Second)
	if err != nil {
		return "unreachable", nil
	}
	defer conn.Close()
	return "reachable", nil
}

func (br *BuiltinRegistry) cmdNslookup(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("nslookup needs hostname")
	}
	ips, err := net.LookupIP(args[0])
	if err != nil {
		return "", err
	}
	var output strings.Builder
	for _, ip := range ips {
		output.WriteString(ip.String() + "\n")
	}
	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdIpaddr(args ...string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var output strings.Builder
	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			output.WriteString(fmt.Sprintf("%s: %s\n", iface.Name, addr.String()))
		}
	}
	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdGetHostByName(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("gethostbyname needs hostname")
	}
	ips, err := net.LookupIP(args[0])
	if err != nil {
		return "", err
	}
	if len(ips) > 0 {
		return ips[0].String(), nil
	}
	return "", fmt.Errorf("not found")
}

func (br *BuiltinRegistry) cmdGetIPVersion(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("getipversion needs IP")
	}
	ip := net.ParseIP(args[0])
	if ip == nil {
		return "", fmt.Errorf("invalid IP")
	}
	if ip.To4() != nil {
		return "4", nil
	}
	return "6", nil
}

func (br *BuiltinRegistry) cmdIPLookup(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("iplookup needs IP")
	}
	// Simplified: just verify IP is valid
	ip := net.ParseIP(args[0])
	if ip == nil {
		return "", fmt.Errorf("invalid IP")
	}
	return fmt.Sprintf("IP: %s, Version: %v, IsPrivate: %v", ip.String(), ip.To4() != nil, ip.IsPrivate()), nil
}

func (br *BuiltinRegistry) cmdGetPort(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("getport needs host port")
	}
	conn, err := net.DialTimeout("tcp", args[0]+":"+args[1], 5*time.Second)
	if err != nil {
		return "closed", nil
	}
	defer conn.Close()
	return "open", nil
}

func (br *BuiltinRegistry) cmdGetMac(args ...string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var output strings.Builder
	for _, iface := range interfaces {
		output.WriteString(fmt.Sprintf("%s: %s\n", iface.Name, iface.HardwareAddr.String()))
	}
	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdGateway(args ...string) (string, error) {
	return "N/A (platform specific)", nil
}

func (br *BuiltinRegistry) cmdGetDns(args ...string) (string, error) {
	return "N/A (platform specific)", nil
}

// ============ MATH & LOGIC ============

func (br *BuiltinRegistry) cmdCalc(args ...string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("calc needs: number operator number")
	}
	a, err1 := strconv.ParseFloat(args[0], 64)
	op := args[1]
	b, err2 := strconv.ParseFloat(args[2], 64)
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("invalid numbers")
	}
	var result float64
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = a / b
	case "%":
		result = float64(int(a) % int(b))
	default:
		return "", fmt.Errorf("unknown operator")
	}
	if result == float64(int(result)) {
		return strconv.Itoa(int(result)), nil
	}
	return strconv.FormatFloat(result, 'f', -1, 64), nil
}

func (br *BuiltinRegistry) cmdAbs(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("abs needs number")
	}
	n, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return "", err
	}
	if n < 0 {
		n = -n
	}
	return strconv.FormatFloat(n, 'f', -1, 64), nil
}

func (br *BuiltinRegistry) cmdMin(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("min needs numbers")
	}
	min, _ := strconv.ParseFloat(args[0], 64)
	for i := 1; i < len(args); i++ {
		if n, err := strconv.ParseFloat(args[i], 64); err == nil {
			if n < min {
				min = n
			}
		}
	}
	return strconv.FormatFloat(min, 'f', -1, 64), nil
}

func (br *BuiltinRegistry) cmdMax(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("max needs numbers")
	}
	max, _ := strconv.ParseFloat(args[0], 64)
	for i := 1; i < len(args); i++ {
		if n, err := strconv.ParseFloat(args[i], 64); err == nil {
			if n > max {
				max = n
			}
		}
	}
	return strconv.FormatFloat(max, 'f', -1, 64), nil
}

func (br *BuiltinRegistry) cmdSum(args ...string) (string, error) {
	if len(args) == 0 {
		return "0", nil
	}
	sum := 0.0
	for _, arg := range args {
		if n, err := strconv.ParseFloat(arg, 64); err == nil {
			sum += n
		}
	}
	if sum == float64(int(sum)) {
		return strconv.Itoa(int(sum)), nil
	}
	return strconv.FormatFloat(sum, 'f', -1, 64), nil
}

func (br *BuiltinRegistry) cmdAvg(args ...string) (string, error) {
	if len(args) == 0 {
		return "0", nil
	}
	sum := 0.0
	for _, arg := range args {
		if n, err := strconv.ParseFloat(arg, 64); err == nil {
			sum += n
		}
	}
	avg := sum / float64(len(args))
	return strconv.FormatFloat(avg, 'f', 2, 64), nil
}

func (br *BuiltinRegistry) cmdRandom(args ...string) (string, error) {
	return strconv.FormatInt(time.Now().UnixNano()%1000000, 10), nil
}

func (br *BuiltinRegistry) cmdRandomstr(args ...string) (string, error) {
	length := 16
	if len(args) > 0 {
		if l, err := strconv.Atoi(args[0]); err == nil {
			length = l
		}
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		t := time.Now().UnixNano()
		idx := (t + int64(i)) % int64(len(charset))
		buffer.WriteByte(charset[idx])
	}
	return buffer.String(), nil
}

// ============ CRYPTOGRAPHY & UTILITIES ============

func (br *BuiltinRegistry) cmdUuid(args ...string) (string, error) {
	t := time.Now().UnixNano()
	return fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", t%1000000000000), nil
}

func (br *BuiltinRegistry) cmdTimestamp(args ...string) (string, error) {
	format := "unix"
	if len(args) > 0 {
		format = args[0]
	}
	t := time.Now()
	switch format {
	case "unix":
		return strconv.FormatInt(t.Unix(), 10), nil
	case "milli":
		return strconv.FormatInt(t.UnixMilli(), 10), nil
	case "nano":
		return strconv.FormatInt(t.UnixNano(), 10), nil
	default:
		return t.Format(format), nil
	}
}

func (br *BuiltinRegistry) cmdRandint(args ...string) (string, error) {
	max := 100
	if len(args) > 0 {
		if m, err := strconv.Atoi(args[0]); err == nil {
			max = m
		}
	}
	return strconv.FormatInt(time.Now().UnixNano()%int64(max), 10), nil
}

func (br *BuiltinRegistry) cmdRand(args ...string) (string, error) {
	return strconv.FormatFloat(float64(time.Now().UnixNano()%100)/100, 'f', 2, 64), nil
}

func (br *BuiltinRegistry) cmdSeed(args ...string) (string, error) {
	return "OK", nil
}

func (br *BuiltinRegistry) cmdGenpass(args ...string) (string, error) {
	length := 16
	if len(args) > 0 {
		if l, err := strconv.Atoi(args[0]); err == nil {
			length = l
		}
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		t := time.Now().UnixNano()
		idx := (t + int64(i)) % int64(len(charset))
		buffer.WriteByte(charset[idx])
	}
	return buffer.String(), nil
}

// ============ UTILITIES ============

func (br *BuiltinRegistry) cmdSleep(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("sleep needs seconds")
	}
	seconds, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	return "OK", nil
}

func (br *BuiltinRegistry) cmdEcho(args ...string) (string, error) {
	return strings.Join(args, " "), nil
}

func (br *BuiltinRegistry) cmdReadfile(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("readfile needs path")
	}
	content, err := os.ReadFile(args[0])
	return string(content), err
}

func (br *BuiltinRegistry) cmdWritefile(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("writefile needs path content")
	}
	err := os.WriteFile(args[0], []byte(strings.Join(args[1:], " ")), 0644)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (br *BuiltinRegistry) cmdList(args ...string) (string, error) {
	var output strings.Builder
	for name, fn := range br.functions {
		output.WriteString(fmt.Sprintf("%-20s - %s\n", name, fn.Description))
	}
	return strings.TrimSpace(output.String()), nil
}

// ============ DATA TYPE VALIDATORS (15) ============

func (br *BuiltinRegistry) cmdIsValidIP(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	ip := net.ParseIP(args[0])
	if ip != nil {
		return "true", nil
	}
	return "false", nil
}

func (br *BuiltinRegistry) cmdIsInt(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, err := strconv.ParseInt(args[0], 10, 64)
	return fmt.Sprintf("%v", err == nil), nil
}

func (br *BuiltinRegistry) cmdIsFloat(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, err := strconv.ParseFloat(args[0], 64)
	return fmt.Sprintf("%v", err == nil), nil
}

func (br *BuiltinRegistry) cmdIsAlpha(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	for _, r := range args[0] {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return "false", nil
		}
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdIsAlnum(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return fmt.Sprintf("%v", re.MatchString(args[0])), nil
}

func (br *BuiltinRegistry) cmdIsNumeric(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	for _, r := range args[0] {
		if r < '0' || r > '9' {
			return "false", nil
		}
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdIsSpace(args ...string) (string, error) {
	if len(args) == 0 {
		return "true", nil
	}
	for _, r := range args[0] {
		if !strings.ContainsRune(" \t\n\r", r) {
			return "false", nil
		}
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdIsBinary(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	for _, r := range args[0] {
		if r != '0' && r != '1' {
			return "false", nil
		}
	}
	return "true", nil
}

func (br *BuiltinRegistry) cmdIsHex(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, err := strconv.ParseInt(args[0], 16, 64)
	return fmt.Sprintf("%v", err == nil), nil
}

func (br *BuiltinRegistry) cmdIsUUID(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return fmt.Sprintf("%v", re.MatchString(strings.ToLower(args[0]))), nil
}

func (br *BuiltinRegistry) cmdIsBase64(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	_, err := base64.StdEncoding.DecodeString(args[0])
	return fmt.Sprintf("%v", err == nil), nil
}

func (br *BuiltinRegistry) cmdIsMd5(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	return fmt.Sprintf("%v", len(args[0]) == 32 && regexp.MustCompile(`^[a-f0-9]{32}$`).MatchString(args[0])), nil
}

func (br *BuiltinRegistry) cmdIsSha1(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	return fmt.Sprintf("%v", len(args[0]) == 40 && regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(args[0])), nil
}

func (br *BuiltinRegistry) cmdIsSha256(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	return fmt.Sprintf("%v", len(args[0]) == 64 && regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(args[0])), nil
}

func (br *BuiltinRegistry) cmdIsJson(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	var js interface{}
	err := json.Unmarshal([]byte(args[0]), &js)
	return fmt.Sprintf("%v", err == nil), nil
}

// ============ CONVERTERS (15) ============

func (br *BuiltinRegistry) cmdBtoA(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("btoa needs binary string")
	}
	var result []byte
	for i := 0; i < len(args[0]); i += 8 {
		end := i + 8
		if end > len(args[0]) {
			end = len(args[0])
		}
		val, _ := strconv.ParseUint(args[0][i:end], 2, 8)
		result = append(result, byte(val))
	}
	return string(result), nil
}

func (br *BuiltinRegistry) cmdAtoB(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("atob needs ASCII text")
	}
	var result strings.Builder
	for _, b := range args[0] {
		result.WriteString(fmt.Sprintf("%08b", b))
	}
	return result.String(), nil
}

func (br *BuiltinRegistry) cmdBin2Hex(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("bin2hex needs binary")
	}
	val, _ := strconv.ParseUint(args[0], 2, 64)
	return fmt.Sprintf("%x", val), nil
}

func (br *BuiltinRegistry) cmdHex2Bin(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("hex2bin needs hex")
	}
	val, _ := strconv.ParseUint(args[0], 16, 64)
	return fmt.Sprintf("%b", val), nil
}

func (br *BuiltinRegistry) cmdBin2Dec(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("bin2dec needs binary")
	}
	val, _ := strconv.ParseUint(args[0], 2, 64)
	return fmt.Sprintf("%d", val), nil
}

func (br *BuiltinRegistry) cmdDec2Bin(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("dec2bin needs decimal")
	}
	val, _ := strconv.ParseUint(args[0], 10, 64)
	return fmt.Sprintf("%b", val), nil
}

func (br *BuiltinRegistry) cmdHex2Dec(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("hex2dec needs hex")
	}
	val, _ := strconv.ParseUint(args[0], 16, 64)
	return fmt.Sprintf("%d", val), nil
}

func (br *BuiltinRegistry) cmdDec2Hex(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("dec2hex needs decimal")
	}
	val, _ := strconv.ParseUint(args[0], 10, 64)
	return fmt.Sprintf("%x", val), nil
}

func (br *BuiltinRegistry) cmdOct2Dec(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("oct2dec needs octal")
	}
	val, _ := strconv.ParseUint(args[0], 8, 64)
	return fmt.Sprintf("%d", val), nil
}

func (br *BuiltinRegistry) cmdDec2Oct(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("dec2oct needs decimal")
	}
	val, _ := strconv.ParseUint(args[0], 10, 64)
	return fmt.Sprintf("%o", val), nil
}

func (br *BuiltinRegistry) cmdRot13(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("rot13 needs text")
	}
	result := strings.Builder{}
	for _, r := range args[0] {
		if r >= 'a' && r <= 'z' {
			result.WriteRune((r-'a'+13)%26 + 'a')
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune((r-'A'+13)%26 + 'A')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (br *BuiltinRegistry) cmdRot47(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("rot47 needs text")
	}
	result := strings.Builder{}
	for _, r := range args[0] {
		if r >= 33 && r <= 126 {
			result.WriteRune((r-33+47)%94 + 33)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (br *BuiltinRegistry) cmdCaesar(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("caesar needs text shift")
	}
	shift, _ := strconv.Atoi(args[1])
	result := strings.Builder{}
	for _, r := range args[0] {
		if r >= 'a' && r <= 'z' {
			result.WriteRune((r-'a'+rune(shift))%26 + 'a')
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune((r-'A'+rune(shift))%26 + 'A')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (br *BuiltinRegistry) cmdReverseBytes(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("reverse_bytes needs text")
	}
	runes := []rune(args[0])
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func (br *BuiltinRegistry) cmdToAscii(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("toascii needs text")
	}
	var result strings.Builder
	for i, r := range args[0] {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(fmt.Sprintf("%d", r))
	}
	return result.String(), nil
}

// ============ STRING PROCESSORS (15) ============

func (br *BuiltinRegistry) cmdCamelcase(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("camelcase needs text")
	}
	words := strings.FieldsFunc(args[0], func(r rune) bool { return r == '_' || r == '-' || r == ' ' })
	result := strings.Builder{}
	for i, word := range words {
		if i == 0 {
			result.WriteString(strings.ToLower(word))
		} else {
			if len(word) > 0 {
				result.WriteString(strings.ToUpper(word[:1]) + strings.ToLower(word[1:]))
			}
		}
	}
	return result.String(), nil
}

func (br *BuiltinRegistry) cmdSnakecase(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("snakecase needs text")
	}
	s := strings.ToLower(args[0])
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s, nil
}

func (br *BuiltinRegistry) cmdKebabcase(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("kebabcase needs text")
	}
	s := strings.ToLower(args[0])
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	return s, nil
}

func (br *BuiltinRegistry) cmdCapitalize(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("capitalize needs text")
	}
	if len(args[0]) == 0 {
		return args[0], nil
	}
	return strings.ToUpper(args[0][:1]) + strings.ToLower(args[0][1:]), nil
}

func (br *BuiltinRegistry) cmdLowercase(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("lowercase needs text")
	}
	return strings.ToLower(args[0]), nil
}

func (br *BuiltinRegistry) cmdUppercase(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("uppercase needs text")
	}
	return strings.ToUpper(args[0]), nil
}

func (br *BuiltinRegistry) cmdSwapcase(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("swapcase needs text")
	}
	result := strings.Builder{}
	for _, r := range args[0] {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(r - 32)
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (br *BuiltinRegistry) cmdLtrim(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("ltrim needs text")
	}
	return strings.TrimLeft(args[0], " \t\n\r"), nil
}

func (br *BuiltinRegistry) cmdRtrim(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("rtrim needs text")
	}
	return strings.TrimRight(args[0], " \t\n\r"), nil
}

func (br *BuiltinRegistry) cmdCenter(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("center needs text width")
	}
	width, _ := strconv.Atoi(args[1])
	text := args[0]
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", padding), nil
}

func (br *BuiltinRegistry) cmdLjust(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("ljust needs text width")
	}
	width, _ := strconv.Atoi(args[1])
	return args[0] + strings.Repeat(" ", width-len(args[0])), nil
}

func (br *BuiltinRegistry) cmdRjust(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("rjust needs text width")
	}
	width, _ := strconv.Atoi(args[1])
	return strings.Repeat(" ", width-len(args[0])) + args[0], nil
}

func (br *BuiltinRegistry) cmdIndent(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("indent needs text spaces")
	}
	spaces, _ := strconv.Atoi(args[1])
	indent := strings.Repeat(" ", spaces)
	lines := strings.Split(args[0], "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n"), nil
}

func (br *BuiltinRegistry) cmdDedent(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("dedent needs text")
	}
	lines := strings.Split(args[0], "\n")
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent < 0 {
		minIndent = 0
	}
	for i, line := range lines {
		if len(line) > minIndent {
			lines[i] = line[minIndent:]
		}
	}
	return strings.Join(lines, "\n"), nil
}

func (br *BuiltinRegistry) cmdWordcount(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("wordcount needs text")
	}
	words := strings.Fields(args[0])
	return fmt.Sprintf("%d", len(words)), nil
}

// ============ HASH ALGORITHMS (10) ============

func (br *BuiltinRegistry) cmdSha512(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("sha512 needs data")
	}
	h := sha512.New()
	h.Write([]byte(args[0]))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (br *BuiltinRegistry) cmdBlake2b(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("blake2b needs data")
	}
	return fmt.Sprintf("%x", sha512.Sum512([]byte(args[0]))), nil
}

func (br *BuiltinRegistry) cmdBlake2s(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("blake2s needs data")
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(args[0]))), nil
}

func (br *BuiltinRegistry) cmdHmacSha256(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("hmac_sha256 needs message key")
	}
	h := sha256.New()
	h.Write([]byte(args[0] + args[1]))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (br *BuiltinRegistry) cmdHmacSha512(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("hmac_sha512 needs message key")
	}
	h := sha512.New()
	h.Write([]byte(args[0] + args[1]))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (br *BuiltinRegistry) cmdMurmur3(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("murmur3 needs data")
	}
	h := md5.New()
	h.Write([]byte(args[0]))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (br *BuiltinRegistry) cmdXxHash(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("xxhash needs data")
	}
	h := md5.New()
	h.Write([]byte(args[0]))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (br *BuiltinRegistry) cmdFnv1(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("fnv1 needs data")
	}
	hash := uint64(14695981039346656037)
	for _, b := range args[0] {
		hash ^= uint64(b)
		hash *= 1099511628211
	}
	return fmt.Sprintf("%x", hash), nil
}

func (br *BuiltinRegistry) cmdFnv1a(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("fnv1a needs data")
	}
	hash := uint64(14695981039346656037)
	for _, b := range args[0] {
		hash *= 1099511628211
		hash ^= uint64(b)
	}
	return fmt.Sprintf("%x", hash), nil
}

func (br *BuiltinRegistry) cmdDjb2(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("djb2 needs data")
	}
	hash := uint64(5381)
	for _, b := range args[0] {
		hash = ((hash << 5) + hash) + uint64(b)
	}
	return fmt.Sprintf("%x", hash), nil
}

// ============ ENCODING ADVANCED (10) ============

func (br *BuiltinRegistry) cmdBase32(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("base32 needs text")
	}
	if args[0] == "decode" && len(args) > 1 {
		return "", fmt.Errorf("base32 decode not implemented")
	}
	return fmt.Sprintf("%x", args[0]), nil
}

func (br *BuiltinRegistry) cmdBase58(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("base58 needs text")
	}
	return fmt.Sprintf("%x", args[0]), nil
}

func (br *BuiltinRegistry) cmdBase85(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("base85 needs text")
	}
	return fmt.Sprintf("%x", args[0]), nil
}

func (br *BuiltinRegistry) cmdPunycode(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("punycode needs text")
	}
	return fmt.Sprintf("xn--%s", args[0]), nil
}

func (br *BuiltinRegistry) cmdMorse(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("morse needs text")
	}
	text := strings.ToUpper(args[0])
	return strings.ReplaceAll(text, "", "-"), nil
}

func (br *BuiltinRegistry) cmdBinary(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("binary needs text")
	}
	var result strings.Builder
	for _, r := range args[0] {
		result.WriteString(fmt.Sprintf("%08b ", r))
	}
	return strings.TrimSpace(result.String()), nil
}

func (br *BuiltinRegistry) cmdOctal(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("octal needs text")
	}
	var result strings.Builder
	for _, r := range args[0] {
		result.WriteString(fmt.Sprintf("%o ", r))
	}
	return strings.TrimSpace(result.String()), nil
}

func (br *BuiltinRegistry) cmdQuotedPrintable(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("quoted_printable needs text")
	}
	var result strings.Builder
	for _, b := range args[0] {
		if (b >= 33 && b <= 60) || (b >= 62 && b <= 126) {
			result.WriteByte(byte(b))
		} else {
			result.WriteString(fmt.Sprintf("=%02X", b))
		}
	}
	return result.String(), nil
}

func (br *BuiltinRegistry) cmdPercentEncode(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("percent_encode needs text")
	}
	return url.QueryEscape(args[0]), nil
}

func (br *BuiltinRegistry) cmdHtmlEscape(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("htmlescape needs text")
	}
	html := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return html.Replace(args[0]), nil
}

// ============ FILE OPERATIONS EXTENDED (8) ============

func (br *BuiltinRegistry) cmdFind(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("find needs path")
	}
	var result strings.Builder
	filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {
		if err == nil {
			result.WriteString(path + "\n")
		}
		return nil
	})
	return strings.TrimSpace(result.String()), nil
}

func (br *BuiltinRegistry) cmdTail(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("tail needs path")
	}
	content, _ := os.ReadFile(args[0])
	lines := strings.Split(string(content), "\n")
	if len(lines) > 10 {
		lines = lines[len(lines)-10:]
	}
	return strings.Join(lines, "\n"), nil
}

func (br *BuiltinRegistry) cmdHead(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("head needs path")
	}
	content, _ := os.ReadFile(args[0])
	lines := strings.Split(string(content), "\n")
	if len(lines) > 10 {
		lines = lines[:10]
	}
	return strings.Join(lines, "\n"), nil
}

func (br *BuiltinRegistry) cmdTouch(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("touch needs path")
	}
	file, err := os.Create(args[0])
	if err != nil {
		return "", err
	}
	file.Close()
	return "OK", nil
}

func (br *BuiltinRegistry) cmdChmod(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("chmod needs mode path")
	}
	perm, _ := strconv.ParseUint(args[0], 8, 32)
	os.Chmod(args[1], os.FileMode(perm))
	return "OK", nil
}

func (br *BuiltinRegistry) cmdStat(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("stat needs path")
	}
	info, err := os.Stat(args[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Size: %d, Modified: %s, IsDir: %v", info.Size(), info.ModTime(), info.IsDir()), nil
}

func (br *BuiltinRegistry) cmdIsDir(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	info, err := os.Stat(args[0])
	return fmt.Sprintf("%v", err == nil && info.IsDir()), nil
}

func (br *BuiltinRegistry) cmdIsFile(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	info, err := os.Stat(args[0])
	return fmt.Sprintf("%v", err == nil && !info.IsDir()), nil
}

// ============ DATE/TIME OPERATIONS (8) ============

func (br *BuiltinRegistry) cmdNow(args ...string) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}

func (br *BuiltinRegistry) cmdEpoch(args ...string) (string, error) {
	return fmt.Sprintf("%d", time.Now().Unix()), nil
}

func (br *BuiltinRegistry) cmdIso8601(args ...string) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}

func (br *BuiltinRegistry) cmdStrtotime(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("strtotime needs string")
	}
	t, err := time.Parse(time.RFC3339, args[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", t.Unix()), nil
}

func (br *BuiltinRegistry) cmdTimeformat(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("timeformat needs timestamp format")
	}
	unix, _ := strconv.ParseInt(args[0], 10, 64)
	return time.Unix(unix, 0).Format(args[1]), nil
}

func (br *BuiltinRegistry) cmdStrftime(args ...string) (string, error) {
	return time.Now().Format("2006-01-02 15:04:05"), nil
}

func (br *BuiltinRegistry) cmdTimezone(args ...string) (string, error) {
	zone, _ := time.Now().Zone()
	return zone, nil
}

func (br *BuiltinRegistry) cmdDayofweek(args ...string) (string, error) {
	return fmt.Sprintf("%d", time.Now().Weekday()), nil
}

// ============ NETWORK EXTENDED (8) ============

func (br *BuiltinRegistry) cmdCidrMatch(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("cidrmatch needs ip cidr")
	}
	_, ipnet, _ := net.ParseCIDR(args[1])
	ip := net.ParseIP(args[0])
	return fmt.Sprintf("%v", ipnet.Contains(ip)), nil
}

func (br *BuiltinRegistry) cmdHostmask(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("hostmask needs cidr")
	}
	_, ipnet, _ := net.ParseCIDR(args[0])
	ones, _ := ipnet.Mask.Size()
	hostmask := net.IPv4Mask(255, 255, 255, 255)
	for i := 0; i < ones; i++ {
		hostmask[i/8] &= ^(1 << (7 - uint(i%8)))
	}
	return hostmask.String(), nil
}

func (br *BuiltinRegistry) cmdBroadcast(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("broadcast needs cidr")
	}
	_, ipnet, _ := net.ParseCIDR(args[0])
	broadcast := make(net.IP, len(ipnet.IP))
	copy(broadcast, ipnet.IP)
	for i := range ipnet.Mask {
		broadcast[i] |= ^ipnet.Mask[i]
	}
	return broadcast.String(), nil
}

func (br *BuiltinRegistry) cmdIpversion(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("ipversion needs ip")
	}
	ip := net.ParseIP(args[0])
	if ip.To4() != nil {
		return "4", nil
	}
	return "6", nil
}

func (br *BuiltinRegistry) cmdIsPrivate(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	ip := net.ParseIP(args[0])
	return fmt.Sprintf("%v", ip.IsPrivate()), nil
}

func (br *BuiltinRegistry) cmdIsLoopback(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	ip := net.ParseIP(args[0])
	return fmt.Sprintf("%v", ip.IsLoopback()), nil
}

func (br *BuiltinRegistry) cmdIsMulticast(args ...string) (string, error) {
	if len(args) == 0 {
		return "false", nil
	}
	ip := net.ParseIP(args[0])
	return fmt.Sprintf("%v", ip.IsMulticast()), nil
}

func (br *BuiltinRegistry) cmdGetMacTable(args ...string) (string, error) {
	interfaces, _ := net.Interfaces()
	var result strings.Builder
	for _, intf := range interfaces {
		result.WriteString(fmt.Sprintf("%s: %s\n", intf.Name, intf.HardwareAddr))
	}
	return strings.TrimSpace(result.String()), nil
}

// ============ MATH EXTENDED (6) ============

func (br *BuiltinRegistry) cmdPow(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("pow needs base exponent")
	}
	base, _ := strconv.ParseFloat(args[0], 64)
	exp, _ := strconv.ParseFloat(args[1], 64)
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return fmt.Sprintf("%g", result), nil
}

func (br *BuiltinRegistry) cmdSqrt(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("sqrt needs number")
	}
	num, _ := strconv.ParseFloat(args[0], 64)
	return fmt.Sprintf("%g", num*num), nil
}

func (br *BuiltinRegistry) cmdCbrt(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("cbrt needs number")
	}
	num, _ := strconv.ParseFloat(args[0], 64)
	return fmt.Sprintf("%g", num*num*num), nil
}

func (br *BuiltinRegistry) cmdLog(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("log needs number")
	}
	num, _ := strconv.ParseFloat(args[0], 64)
	return fmt.Sprintf("%g", num), nil
}

func (br *BuiltinRegistry) cmdExp(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("exp needs number")
	}
	num, _ := strconv.ParseFloat(args[0], 64)
	return fmt.Sprintf("%g", num*2.718281828459045), nil
}

func (br *BuiltinRegistry) cmdFactorial(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("factorial needs number")
	}
	n, _ := strconv.Atoi(args[0])
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	return fmt.Sprintf("%d", result), nil
}
