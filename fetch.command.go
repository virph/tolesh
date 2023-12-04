package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

func commandFetchServers(ctx context.Context) error {
	var (
		wg       sync.WaitGroup
		errFetch error
	)
	output := ""

	wg.Add(1)
	go func() {
		defer wg.Done()

		if param.FetchParam.UseSample {
			output = sample
			return
		}

		var o []byte
		cmd := exec.CommandContext(ctx, "tsh", "ls", "-v")
		o, errFetch = cmd.Output()
		if errFetch != nil {
			return
		}
		output = string(o)
	}()
	log.Println("Waiting \"tsh ls -v\" command result...")
	wg.Wait()

	if errFetch != nil {
		return errFetch
	}

	// process
	log.Println("Processing data...")
	nodes := new(NodeData).BuildFromString(output)
	nodesMap = nodes.NodesMap

	// save to cache
	if err := writeFile(fmt.Sprintf("%s/%s", param.DataPath, cachePath), nodes.NodesMap.toJSON()); err != nil {
		return err
	}

	if param.Verbose {
		log.Println("Nodes result JSON:", nodes.toJSON())
	}

	writeFile("./tsh-ls_result.txt", output)

	return nil
}
