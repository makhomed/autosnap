package main

import (
	"fmt"
	"os"
	"config"
	"snapman"
)

const Config = "/opt/autosnap/conf/autosnap.conf"

func main() {
	conf, err := config.New(Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config '%s' : %v\n", Config, err)
		os.Exit(2)
	}
	if len(os.Args) == 1 || len(os.Args) == 2 && os.Args[1] == "-h" || len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: %s command|clean\n", os.Args[0])
		os.Exit(2)
	}
	command := os.Args[1]
	if _, ok := conf.Interval[command]; !ok && command != "clean" {
		fmt.Fprintf(os.Stderr, "unknown command '%s', interval not defined\n", command)
		os.Exit(2)
	}
	snapman.Execute(conf, command)
}
