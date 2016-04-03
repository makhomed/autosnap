package snapman

import (
	"config"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"bufio"
	"strings"
)

func Datasets(conf *config.Config) (map[string]bool, error) {
	datasets := make(map[string]bool)
	cmd := exec.Command("zfs", "list", "-H", "-o", "name")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(output)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		dataset := scanner.Text()
		dataset = strings.TrimSpace(dataset)
		if dataset == "" {
			continue
		}
		datasets[dataset] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return datasets, nil
}

func Snapshots(conf *config.Config) (map[string]bool, error) {
	datasets := make(map[string]bool)
	cmd := exec.Command("zfs", "list", "-H", "-p", "-o", "name,creation", "-t", "snap")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(output)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		dataset := scanner.Text()
		dataset = strings.TrimSpace(dataset)
		if dataset == "" {
			continue
		}
		datasets[dataset] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return datasets, nil
}



func CreateSnapshot(dataset string) {
	fmt.Fprintf(os.Stderr, "TODO: create shapshot for dataset: %v\n", dataset)
}

func Execute(conf *config.Config, command string) {
	fmt.Fprintf(os.Stderr, "conf: %v command: %s\n", conf, command)
	datasets, err := Datasets(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: can't read datasets: %v\n", err)
		os.Exit(1)
	}

	if _, ok := conf.Interval[command]; ok {
		for dataset := range datasets {
			if _, ok := conf.Exclude[dataset]; !ok {
				CreateSnapshot(dataset)
			}

		}
	}

	for dataset := range datasets {
		fmt.Fprintf(os.Stderr, "dataset: %v\n", dataset)
	}
}
