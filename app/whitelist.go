package app

import (
	"bufio"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	once      sync.Once
	whitelist map[string]bool // do not use directly, use getWhitelist
)

func IsWhitelisted(appId string) bool {
	whitelist := getWhitelist()
	_, ok := whitelist[appId]
	return ok
}

func getWhitelist() map[string]bool {
	once.Do(func() {
		whitelist = map[string]bool{}

		file, err := os.Open("whitelist.txt")
		if err != nil {
			log.Fatal("Failed to read the whitelist")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			whitelist[getAppIdWithoutComments(scanner.Text())] = true
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	})
	return whitelist
}

func getAppIdWithoutComments(line string) string {
	idx := strings.Index(line, "#")
	if idx > 0 {
		return line[0:idx]
	}
	return line
}
