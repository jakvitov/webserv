package err

import (
	"fmt"
	"os"
)

/*
This file includes functions to print user error prompts or manuals when application
runs into error during launch before loggers and other more complex utils are initialized.
Instead of loggers, we aim for short and descriptive prompts and expect application exit afterwards
*/

// Prompt the user with an error during initialization
func ErrorPrompt(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	PrintHelp()
}

// Displays the server usage prompt to the screen
func PrintHelp() {
	fmt.Printf("For info about the application usage use the --help falg!\n")
}
