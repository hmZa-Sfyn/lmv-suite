package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// TableRow represents a row in a table
type TableRow struct {
	Cols []string
}

// Table represents a formatted table
type Table struct {
	Headers []string
	Rows    []TableRow
	Widths  []int
}

// ANSI color regex for width calculation
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// NewTable creates a new table
func NewTable(headers []string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h) + 2
	}
	return &Table{
		Headers: headers,
		Rows:    []TableRow{},
		Widths:  widths,
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(cols ...string) {
	for i, col := range cols {
		if i < len(t.Widths) {
			colLen := len(stripANSI(col)) + 2
			if colLen > t.Widths[i] {
				t.Widths[i] = colLen
			}
		}
	}
	t.Rows = append(t.Rows, TableRow{Cols: cols})
}

// stripANSI removes ANSI color codes from string
func stripANSI(str string) string {
	var result strings.Builder
	inEscape := false

	for _, char := range str {
		if char == '\x1b' {
			inEscape = true
			continue
		}

		if inEscape {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == 'm' {
				inEscape = false
			}
			continue
		}

		result.WriteRune(char)
	}

	return result.String()
}

// Render renders the table as a string
func (t *Table) Render() string {
	var sb strings.Builder

	// Top border
	sb.WriteString(t.drawBorder("╔", "╤", "╗"))

	// Header
	sb.WriteString(t.drawRow(t.Headers, true))

	// Header border
	sb.WriteString(t.drawBorder("╟", "┼", "╢"))

	// Rows
	for _, row := range t.Rows {
		sb.WriteString(t.drawRow(row.Cols, false))
	}

	// Bottom border
	sb.WriteString(t.drawBorder("╚", "╧", "╝"))

	return sb.String()
}

// drawBorder draws a border line
func (t *Table) drawBorder(left, mid, right string) string {
	var sb strings.Builder
	sb.WriteString(left)

	for i, width := range t.Widths {
		sb.WriteString(strings.Repeat("═", width))
		if i < len(t.Widths)-1 {
			sb.WriteString(mid)
		}
	}

	sb.WriteString(right)
	sb.WriteString("\n")
	return sb.String()
}

// drawRow draws a content row
func (t *Table) drawRow(cols []string, isHeader bool) string {
	var sb strings.Builder
	sb.WriteString("║")

	for i, col := range cols {
		width := t.Widths[i]
		cleanCol := stripANSI(col)
		padding := width - len(cleanCol)

		sb.WriteString(" ")
		if isHeader {
			sb.WriteString(color.CyanString(cleanCol))
		} else {
			sb.WriteString(col)
		}
		sb.WriteString(strings.Repeat(" ", padding-1))
		sb.WriteString("║")
	}

	sb.WriteString("\n")
	return sb.String()
}

// NmapBox creates nmap-style output box
func NmapBox(title string) string {
	return color.GreenString("|_ ") + color.CyanString(title)
}

// NmapSubBox creates nmap-style sub output
func NmapSubBox(title string) string {
	return color.GreenString("   \\_ ") + color.WhiteString(title)
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	fmt.Printf("%s %s\n", color.GreenString("[+]"), msg)
}

// PrintError prints an error message
func PrintError(msg string) {
	fmt.Printf("%s %s\n", color.RedString("[!]"), msg)
}

// PrintInfo prints an info message
func PrintInfo(msg string) {
	fmt.Printf("%s %s\n", color.YellowString("[*]"), msg)
}

// PrintDebug prints a debug message
func PrintDebug(msg string) {
	fmt.Printf("%s %s\n", color.MagentaString("[~]"), msg)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	fmt.Printf("%s %s\n", color.YellowString("[w]"), msg)
}

// CenterText centers text within a width
func CenterText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", padding)
}

// DrawBox draws a beautiful box
func DrawBox(title string, content string, boxType string) {
	titleLen := len(title)
	contentLen := len(strings.Split(content, "\n")[0])
	maxLen := titleLen
	if contentLen > maxLen {
		maxLen = contentLen
	}
	boxWidth := maxLen + 4

	// Top
	fmt.Print("┌")
	fmt.Print(strings.Repeat("─", boxWidth-2))
	fmt.Println("┐")

	// Title
	fmt.Print("│ ")
	fmt.Print(color.CyanString(title))
	fmt.Print(strings.Repeat(" ", boxWidth-titleLen-3))
	fmt.Println("│")

	// Divider
	fmt.Print("├")
	fmt.Print(strings.Repeat("─", boxWidth-2))
	fmt.Println("┤")

	// Content
	for _, line := range strings.Split(content, "\n") {
		if line != "" {
			fmt.Print("│ ")
			fmt.Print(line)
			fmt.Print(strings.Repeat(" ", boxWidth-len(line)-3))
			fmt.Println("│")
		}
	}

	// Bottom
	fmt.Print("└")
	fmt.Print(strings.Repeat("─", boxWidth-2))
	fmt.Println("┘")
}

// ProgressBar creates a simple progress bar
func ProgressBar(current, total int, width int) string {
	if total <= 0 {
		total = 1
	}
	percentage := float64(current) / float64(total)
	filledWidth := int(float64(width) * percentage)

	bar := "["
	bar += strings.Repeat("=", filledWidth)
	bar += strings.Repeat(" ", width-filledWidth)
	bar += "]"

	return color.GreenString(bar) + fmt.Sprintf(" %.1f%%", percentage*100)
}

// Color returns a colored string
func Color(colorName string, text string) string {
	switch colorName {
	case "red":
		return color.RedString(text)
	case "green":
		return color.GreenString(text)
	case "yellow":
		return color.YellowString(text)
	case "blue":
		return color.BlueString(text)
	case "cyan":
		return color.CyanString(text)
	case "magenta":
		return color.MagentaString(text)
	case "white":
		return color.WhiteString(text)
	default:
		return text
	}
}
