package main

import (
	"io/ioutil"
	"os"
)

func writeFile(path, data string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	return f.Sync()
}

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
