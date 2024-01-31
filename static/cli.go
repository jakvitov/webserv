package static

import (
	"cz/jakvitov/webserv/sharedlogger"
	"fmt"
	"os"
)

const SECTION_PREFIX string = "static/resources/section_"

var (
	VERSION         string
	BUILD_TIMESTAMP string
	COMMIT_HASH     string
)

func printBanner(lg *sharedlogger.SharedLogger) *error {
	banner, err := os.ReadFile("static/resources/banner.txt")
	if err != nil {
		return &err
	}

	bnr := "\033[96m" + string(banner) + "\033[0m"
	lg.Info(string(bnr))
	return nil
}

func PrintBannerDecoration(logger *sharedlogger.SharedLogger) {
	err := printBanner(logger)
	if err != nil {
		logger.Finfo("Error while opening banner: [%s]\n", (*err).Error())
	}
}

func PrintVersionInfo() {
	fmt.Printf("Webserv build info: \n")
	fmt.Printf("\t - version: [%s]\n", VERSION)
	fmt.Printf("\t - build timestamp: [%s]\n", BUILD_TIMESTAMP)
	fmt.Printf("\t - build comit hash:  [%s]\n", COMMIT_HASH)
}

// Verify user choice
// Return false if choice is invalid
func verifyMenuChoice(input int) bool {
	return input < 8
}

// Print menu section, that the user chose to see
func printChosenSection(section int) {
	file, err := os.ReadFile(fmt.Sprintf("%s%d.txt", SECTION_PREFIX, section))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while scanning from the stdin [%s]\n.", err.Error())
		return
	}
	fmt.Println(string(file))
	HelpMenu()
}

// Print the help menu and capture user input
func HelpMenu() {
	data, err := os.ReadFile("static/resources/help.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error,  cannot open the help menu resources!")
		return
	}
	fmt.Println(string(data))
	var choice int
	_, err = fmt.Scanf("%d", &choice)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while scanning from the stdin [%s]\n.", err.Error())
		return
	}
	if !verifyMenuChoice(choice) {
		fmt.Fprintf(os.Stderr, "Invalid choice [%d]\n", choice)
		return
	}
	//User wants to quit
	if choice == 0 {
		return
	}
	printChosenSection(choice)
}
