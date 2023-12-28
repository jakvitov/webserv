package static

import (
	"log"
	"os"
)

func printBanner(lg *log.Logger) *error {
	banner, err := os.ReadFile("../static/banner.txt")
	if err != nil {
		return &err
	}
	lg.Printf("%s\n", banner)
	return nil
}

func PrintBannerDecoration(logger *log.Logger) {
	err := printBanner(logger)
	if err != nil {
		logger.Printf("Error while opening banner: [%s]\n", (*err).Error())
	}
}
