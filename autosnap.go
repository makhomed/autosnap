package main

import (
	"fmt"
	"os"
	"config"
	"snapman"
	"flag"
)

//const Config = "/opt/autosnap/conf/autosnap.conf"

var configName = flag.String("c", "/opt/autosnap/conf/autosnap.conf", "config")

func main() {
	flag.Parse()
	conf, err := config.New(*configName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config '%s' : %v\n", *configName, err)
		os.Exit(2)
	}
	if len(flag.Args()) == 0 || len(flag.Args()) > 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-c path/to/config] <command>|clean\n", os.Args[0])
		os.Exit(2)
	}
	command := flag.Arg(0)
	if _, ok := conf.Interval[command]; !ok && command != "clean" {
		fmt.Fprintf(os.Stderr, "unknown command '%s', interval not defined\n", command)
		os.Exit(2)
	}
	snapman.Execute(conf, command)
}
