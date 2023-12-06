package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	param    *Param
	nodesMap NodesMap

	validCommands   = []string{"fetch", "list", "copy", "iterm"}
	defaultDataPath = fmt.Sprintf("%s/.tolesh", os.Getenv("HOME"))
)

const (
	cachePath = "node-cache.json"
)

func initParam() error {
	if len(os.Args) < 2 {
		return errors.New("invalid input")
	}

	param = &Param{
		Command: os.Args[1],
	}
	isValidCommand := false
	for _, c := range validCommands {
		if param.Command == c {
			isValidCommand = true
			break
		}
	}
	if !isValidCommand {
		return errors.New("unknown command")
	}

	// default attributes
	flag.BoolVar(&param.Verbose, "v", false, "Print verbose log")
	flag.StringVar(&param.DataPath, "p", defaultDataPath, "Data path")

	// command-specific attributes
	switch param.Command {
	case "fetch":
		param.FetchParam = new(FetchParam)
		flag.BoolVar(&param.FetchParam.UseSample, "s", false, "Use sample")
	case "list":
		param.ListParam = new(ListParam)
		flag.StringVar(&param.ListParam.Hostgroup, "g", "all", "Hostgroup")
		flag.StringVar(&param.ListParam.Type, "t", "all", "Type")
	case "copy":
		param.CopyParam = new(CopyParam)
		flag.StringVar(&param.CopyParam.DestinationPath, "d", "", "Destination path")
	case "iterm":
		param.ItermParam = new(ItermParam)
		flag.StringVar(&param.ItermParam.Hostgroup, "g", "", "Hostgroup")
		flag.StringVar(&param.ItermParam.Username, "u", "root", "Username")
		flag.BoolVar(&param.ItermParam.IsPickHost, "h", false, "Is pick host")
	}

	return flag.CommandLine.Parse(os.Args[2:])
}

func initNodesMap() error {
	strNodesMap, err := readFile(fmt.Sprintf("%s/%s", param.DataPath, cachePath))
	if err != nil {
		return err
	}

	nodesMap = make(NodesMap, 0)
	if err := json.Unmarshal([]byte(strNodesMap), &nodesMap); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()

	if err := initParam(); err != nil {
		log.Fatalln(err)
	}
	if param.Verbose {
		log.Println("Requested param:", param.toJSON())
	}

	// init default path
	if param.DataPath == defaultDataPath {
		if err := os.MkdirAll(param.DataPath, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	initNodesMap()
	// run the application command
	if param.Command == "fetch" {
		// fetch
		if err := commandFetchServers(ctx); err != nil {
			log.Fatalln(err)
		}

	} else {
		// check nodes map data first
		if nodesMap == nil || len(nodesMap) == 0 {
			log.Fatalln("Failed to get node data from cache. Try to do \"fetch\" first")
		}

		switch param.Command {
		case "list":
			if err := commandListConfig(ctx); err != nil {
				log.Fatalln(err)
			}
		case "copy":
			if err := commandCopyConfig(ctx); err != nil {
				log.Fatalln(err)
			}
		case "iterm":
			if err := commandIterm(ctx); err != nil {
				log.Fatalln(err)
			}
		}
	}

	log.Println("Exit.")
}
