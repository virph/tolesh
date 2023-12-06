package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os/exec"
	"sync"
)

const strTemplate = ``

type ItermTemplateParam struct {
	NumberOfHost int
	Hosts        []struct {
		Session int
		Host    string
	}
}

func commandIterm(ctx context.Context) error {
	log.Println("Running iterm")

	if _, ok := nodesMap[param.ItermParam.Hostgroup]; !ok {
		return errors.New("hostgroup is not exists")
	}

	var (
		nodes         = nodesMap[param.ItermParam.Hostgroup]
		numOfHost     = len(nodes)
		maxSessions   = 9
		reqSessionNum = numOfHost
		reqHostNum    = 1

		requestedHosts = make([]Node, 0)
	)

	if param.ItermParam.IsPickHost {
		fmt.Printf("List host of [%s] hostgroup:\n", param.ItermParam.Hostgroup)
		fmt.Printf("| %-3s | %-36s | %-20s |\n", "Num", "Node ID", "Address")
		for i := 0; i < numOfHost; i++ {
			fmt.Printf("| %3d | %-36s | %20s |\n", (i + 1), nodes[i].ID, nodes[i].Address)
		}
		fmt.Printf("Input the host number [min:1, max:%d]: ", len(nodes))
		if _, err := fmt.Scanf("%d", &reqHostNum); err != nil {
			return err
		}
		log.Printf("Requested host number: %d\n", reqHostNum)
		if reqHostNum < 1 || reqHostNum > len(nodes) {
			return errors.New("invalid host number")
		}

		requestedHosts = append(requestedHosts, nodes[reqHostNum-1])
	} else {
		if numOfHost < maxSessions {
			maxSessions = numOfHost
		}
		fmt.Printf("Number of hosts in [%s]: %d\n", param.ItermParam.Hostgroup, numOfHost)
		fmt.Printf("Input number of session to run [min:1, max:%d]: ", maxSessions)

		if _, err := fmt.Scanf("%d", &reqSessionNum); err != nil {
			return err
		}
		log.Printf("Requested number of session: %d\n", reqSessionNum)
		if reqSessionNum < 1 || reqSessionNum > maxSessions {
			return errors.New("invalid session number")
		}

		requestedHosts = append(requestedHosts, nodes[:reqSessionNum]...)
	}

	strScriptTemplate, err := readFile("./iterm-script.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("script").Delims("[[", "]]").Parse(strScriptTemplate)
	if err != nil {
		return err
	}

	templateParam := &ItermTemplateParam{
		NumberOfHost: 0,
		Hosts: []struct {
			Session int
			Host    string
		}{},
	}
	for _, node := range requestedHosts {
		templateParam.Hosts = append(templateParam.Hosts, struct {
			Session int
			Host    string
		}{
			Host:    fmt.Sprintf("%s@%s", param.ItermParam.Username, node.ID),
			Session: templateParam.NumberOfHost + 1,
		})

		templateParam.NumberOfHost = templateParam.NumberOfHost + 1
	}

	buff := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buff, templateParam); err != nil {
		return err
	}

	// writing
	scriptPath := fmt.Sprintf("%s/iterm-script.temp.scpt", param.DataPath)
	log.Printf("Writing [%s]\n", scriptPath)
	if err := writeFile(scriptPath, buff.String()); err != nil {
		return err
	}

	var (
		wg     sync.WaitGroup
		errCmd error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		cmd := exec.CommandContext(ctx, "osascript", scriptPath)
		errCmd = cmd.Run()
		if errCmd != nil {
			return
		}
	}()
	log.Println("Running script for iTerm ...")
	wg.Wait()

	return errCmd
}
