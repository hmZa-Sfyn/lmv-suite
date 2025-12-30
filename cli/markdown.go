package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// MarkdownRenderer handles markdown formatting for display
type MarkdownRenderer struct {
	boldHeadings bool
	colorCode    bool
}

// NewMarkdownRenderer creates a new markdown renderer
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{
		boldHeadings: true,
		colorCode:    true,
	}
}

// Render processes markdown text and returns colored output
func (mr *MarkdownRenderer) Render(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		result = append(result, mr.renderLine(line))
	}

	return strings.Join(result, "\n")
}

// renderLine processes a single markdown line
func (mr *MarkdownRenderer) renderLine(line string) string {
	// Handle headings (#, ##, ###, etc)
	headingRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	if matches := headingRegex.FindStringSubmatch(line); matches != nil {
		level := len(matches[1])
		heading := matches[2]

		prefix := strings.Repeat("  ", level-1)

		if mr.boldHeadings {
			return fmt.Sprintf("%s%s", prefix, color.New(color.Bold, color.FgCyan).Sprint(heading))
		}
		return fmt.Sprintf("%s%s", prefix, color.CyanString(heading))
	}

	// Handle inline code (`code`)
	line = mr.renderInlineCode(line)

	// Handle bold (**text** or __text__)
	line = mr.renderBold(line)

	// Handle italic (*text* or _text_)
	line = mr.renderItalic(line)

	// Handle links [text](url)
	line = mr.renderLinks(line)

	return line
}

// renderInlineCode handles `code` syntax
func (mr *MarkdownRenderer) renderInlineCode(line string) string {
	codeRegex := regexp.MustCompile("`([^`]+)`")
	return codeRegex.ReplaceAllStringFunc(line, func(match string) string {
		code := strings.TrimPrefix(strings.TrimSuffix(match, "`"), "`")
		return color.New(color.BgBlack, color.FgYellow).Sprint(" " + code + " ")
	})
}

// renderBold handles **text** or __text__ syntax
func (mr *MarkdownRenderer) renderBold(line string) string {
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*|__([^_]+)__`)
	return boldRegex.ReplaceAllStringFunc(line, func(match string) string {
		text := strings.Trim(strings.Trim(match, "*"), "_")
		return color.New(color.Bold, color.FgWhite).Sprint(text)
	})
}

// renderItalic handles *text* or _text_ syntax
func (mr *MarkdownRenderer) renderItalic(line string) string {
	italicRegex := regexp.MustCompile(`\*([^*]+)\*|_([^_]+)_`)
	return italicRegex.ReplaceAllStringFunc(line, func(match string) string {
		text := strings.Trim(strings.Trim(match, "*"), "_")
		return color.New(color.Italic, color.FgGreen).Sprint(text)
	})
}

// renderLinks handles [text](url) syntax
func (mr *MarkdownRenderer) renderLinks(line string) string {
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	return linkRegex.ReplaceAllStringFunc(line, func(match string) string {
		parts := linkRegex.FindStringSubmatch(match)
		if len(parts) >= 3 {
			text := parts[1]
			url := parts[2]
			return color.BlueString(text) + " (" + color.MagentaString(url) + ")"
		}
		return match
	})
}

// RenderCodeBlock handles code blocks with ```language syntax
func (mr *MarkdownRenderer) RenderCodeBlock(code string, language string) string {
	lines := strings.Split(code, "\n")
	var result []string

	result = append(result, color.New(color.BgBlack, color.FgGreen).Sprint("╔ "+language))

	for _, line := range lines {
		if line == "" {
			result = append(result, "║")
		} else {
			result = append(result, color.New(color.FgYellow).Sprintf("║ %s", line))
		}
	}

	result = append(result, color.New(color.BgBlack, color.FgGreen).Sprint("╚"))

	return strings.Join(result, "\n")
}
