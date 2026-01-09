package cli

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"time"

	"github.com/fatih/color"
)

// GetPrompt returns the CLI prompt
func (cli *CLI) GetPrompt() string {
	user, _ := user.Current()
	hostname, _ := os.Hostname()

	// Simple, colorful prompt
	return fmt.Sprintf("%s%s%s%s ",
		color.CyanString(user.Username),
		color.WhiteString("@"),
		color.MagentaString(hostname),
		color.GreenString(" ❯"),
	)
}

// PrintBanner prints a random application banner
func (cli *CLI) PrintBanner() {
	// Seed random (only needed once)
	rand.Seed(time.Now().UnixNano())

	banners := []string{
		`
______                  
___  /______ ______   __
__  /__  __ '__ \_ | / /
_  / _  / / / / /_ |/ / 
/_/  /_/ /_/ /_/_____/  
                        
`,

		`
dP                     
88                     
88 88d8b.d8b. dP   .dP 
88 88''88''88 88   d8' 
88 88  88  88 88 .88'  
dP dP  dP  dP 8888P'
`,

		`
888                        
888                        
888                        
888 88888b.d88b.  888  888 
888 888 "888 "88b 888  888 
888 888  888  888 Y88  88P 
888 888  888  888  Y8bd8P  
888 888  888  888   Y88P
		`,

		`
$$\                         
$$ |                        
$$ |$$$$$$\$$$$\ $$\    $$\ 
$$ |$$  _$$  _$$\\$$\  $$  |
$$ |$$ / $$ / $$ |\$$\$$  / 
$$ |$$ | $$ | $$ | \$$$  /  
$$ |$$ | $$ | $$ |  \$  /   
\__|\__| \__| \__|   \_/
		`,
	}

	banner := banners[rand.Intn(len(banners))]

	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println(banner)
	fmt.Println()

	// Common footer for all banners
	color.New(color.FgGreen, color.Bold).Println("║   @ LANMANVAN v2.0 - Advanced Modular Tooling Framework @       ║")
	color.New(color.FgGreen, color.Bold).Println("║   Go Core | Python3/Bash Modules | Dynamic UI | Security Tools  ║")
	fmt.Println()

	fmt.Printf("Type %s for available commands, have fun!\n\n", color.CyanString("'help'"))
}

// ClearScreen clears the terminal
func (cli *CLI) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
