package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
)

func commandCopyConfig(ctx context.Context) error {
	log.Println("Running copy")

	// check is valid path
	fi, err := os.Stat(param.CopyParam.DestinationPath)
	if err != nil {
		pathError := err.(*os.PathError)
		return pathError
	}
	if fi.IsDir() {
		return errors.New("path can not be a directory")
	}

	// write to file
	fmt.Printf("Writing to [%s] ...\n", param.CopyParam.DestinationPath)
	return writeFile(param.CopyParam.DestinationPath, nodesMap.toJSON())
}
