package config

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

type Config struct {
	Interval map[string]int
	Exclude  map[string]bool
}

func New(config string) (*Config, error) {
	conf := &Config{
		Interval: make(map[string]int),
		Exclude:make(map[string]bool),
	}
	configFile, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()
		commentPosition := strings.Index(line, "#")
		if commentPosition >= 0 {
			line = line[0:commentPosition]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.Replace(line, "\t", " ", -1)
		split := strings.SplitN(line, " ", 2)
		name, value := split[0], split[1]
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)

		switch name {
		case "interval":
			secondSplit := strings.SplitN(value, " ", 2)
			interval, countStr := secondSplit[0], secondSplit[1]
			interval = strings.TrimSpace(interval)
			countStr = strings.TrimSpace(countStr)
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return nil, fmt.Errorf("bad %s count value '%s' : %s", interval, countStr, err)
			}
			if count <= 0 {
				return nil, fmt.Errorf("bad %s count value '%s' : must be positive integer", interval, countStr)
			}
			if _, ok := conf.Interval[interval]; ok {
				return nil, fmt.Errorf("duplicate interval '%s'", interval)
			}
			if interval=="clean" {
				return nil, fmt.Errorf("interval name '%s' not allowed", interval)
			}
			conf.Interval[interval] = count
		case "exclude":
			if spacePosition := strings.Index(value, " "); spacePosition >= 0 {
				return nil, fmt.Errorf("spaces not allowed: '%s'", value)
			}
			if _, ok := conf.Exclude[value]; ok {
				return nil, fmt.Errorf("duplicate exclude '%s'", value)
			}
			conf.Exclude[value] = true
		default:
			return nil, fmt.Errorf("unknown directive '%s'", name)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return conf, nil
}
