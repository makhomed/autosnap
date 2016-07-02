package config

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"path/filepath"
)

type filterLine struct {
	included bool   // true == include, false == exclude
	pattern  string // rules: https://golang.org/pkg/path/filepath/#Match
}

type Config struct {
	Interval map[string]int
	filter   []filterLine
}

func (config *Config) Included(dataset string) bool {
	for _, line := range config.filter {
		if line.pattern == "*" {
			return line.included
		}
		matched, err := filepath.Match(line.pattern, dataset);
		if err != nil {
			panic(fmt.Sprintf("pattern is malformed: '%s'", line.pattern))
		}
		if matched {
			return line.included
		}
	}
	panic(fmt.Sprintf("unexpected end of func config.Included for dataset '%s'", dataset))
}

func New(config string) (*Config, error) {
	conf := &Config{
		Interval: make(map[string]int),
		filter: make([]filterLine, 0),
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
			if interval == "clean" {
				return nil, fmt.Errorf("interval name '%s' not allowed", interval)
			}
			conf.Interval[interval] = count
		case "include":
			if spacePosition := strings.Index(value, " "); spacePosition >= 0 {
				return nil, fmt.Errorf("spaces not allowed: '%s'", value)
			}
			if _, err := filepath.Match(value, ""); err != nil {
				return nil, fmt.Errorf("pattern is malformed: '%s'", value)
			}
			conf.filter = append(conf.filter, filterLine{true, value})
		case "exclude":
			if spacePosition := strings.Index(value, " "); spacePosition >= 0 {
				return nil, fmt.Errorf("spaces not allowed: '%s'", value)
			}
			if _, err := filepath.Match(value, ""); err != nil {
				return nil, fmt.Errorf("pattern is malformed: '%s'", value)
			}
			conf.filter = append(conf.filter, filterLine{false, value})
		default:
			return nil, fmt.Errorf("unknown directive '%s'", name)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	conf.filter = append(conf.filter, filterLine{true, "*"}) // include all by default
	return conf, nil
}
