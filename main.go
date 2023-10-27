package main

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"sync"
)

var PARAM *Param

func init() {
	PARAM = &Param{}

	flag.BoolVar(&PARAM.UseSample, "usesample", false, "Use sample")
	flag.BoolVar(&PARAM.Verbose, "v", false, "Verbose")
	flag.StringVar(&PARAM.ConfigPath, "path", "config", "Config file path")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	if PARAM.Verbose {
		fmt.Println("Param:", PARAM.toJSON())
	}

	var wg sync.WaitGroup
	output := ""

	wg.Add(1)
	go func() {
		defer wg.Done()

		if PARAM.UseSample {
			output = sample
			return
		}

		cmd := exec.CommandContext(ctx, "tsh", "ls", "-v")
		o, err := cmd.Output()
		if err != nil {
			fmt.Println("Err: ", err)
			return
		}
		output = string(o)
	}()

	fmt.Println("Waiting \"tsh ls -v\" command result...")
	wg.Wait()

	// process
	fmt.Println("Processing data...")
	nodes := new(NodeData).BuildFromString(output)
	if PARAM.Verbose {
		fmt.Println("Result:", nodes.toJSON())
	}

	// write to file
	fmt.Printf("Writing to config file [%s] ...\n", PARAM.ConfigPath)
	writeConfig(PARAM.ConfigPath, nodes.toConfigString())

	fmt.Println("Exit.")
}
