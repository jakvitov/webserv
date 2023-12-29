package static

import (
	"cz/jakvitov/webserv/sharedlogger"
	"fmt"
	"os"
)

func printBanner(lg *sharedlogger.SharedLogger) *error {
	banner, err := os.ReadFile("../static/banner.txt")
	if err != nil {
		return &err
	}
	lg.Info(string(banner))
	return nil
}

func PrintBannerDecoration(logger *sharedlogger.SharedLogger) {
	err := printBanner(logger)
	if err != nil {
		logger.Info(fmt.Sprintf("Error while opening banner: [%s]\n", (*err).Error()))
	}
}
