package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

func readWowConfig() (map[string]string, error) {
	ret := map[string]string{}
	fieldre := regexp.MustCompile(`^SET ([a-zA-Z0-9-]+)\s+"(.*)"`)
	filename := filepath.Join(wowDir, "WTF", "Config.wtf")
	file, err := os.Open(filename)
	if err != nil {
		return ret, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := fieldre.FindStringSubmatch(scanner.Text())
		if match != nil {
			ret[match[1]] = match[2]
		}
	}
	err = scanner.Err()
	if err != nil {
		return ret, err
	}
	return ret, nil
}
