package main

import "os"

func writeConfig(path, data string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	return f.Sync()
}
