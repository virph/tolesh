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

	numOfHost := len(nodesMap[param.ItermParam.Hostgroup])
	maxSessions := 9
	if numOfHost < maxSessions {
		maxSessions = numOfHost
	}
	fmt.Printf("Number of hosts in [%s]: %d\n", param.ItermParam.Hostgroup, numOfHost)
	fmt.Printf("Input number of session to run [min:1, max:%d]: ", maxSessions)

	reqSessionNum := numOfHost
	if _, err := fmt.Scanf("%d", &reqSessionNum); err != nil {
		return err
	}
	log.Printf("Requested number of session: %d\n", reqSessionNum)
	if reqSessionNum < 1 || reqSessionNum > maxSessions {
		return errors.New("invalid session number")
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
	for _, node := range nodesMap[param.ItermParam.Hostgroup] {
		if templateParam.NumberOfHost >= reqSessionNum {
			break
		}

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
