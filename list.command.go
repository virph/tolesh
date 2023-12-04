package main

import (
	"context"
	"fmt"
	"log"
)

func commandListConfig(ctx context.Context) error {
	log.Println("Running list")

	fmt.Printf("%s, %s, %s\n", "Node Name", "Node ID", "Address")
	if param.ListParam.Hostgroup == "all" {
		for _, nodes := range nodesMap {
			for _, node := range nodes {
				t := param.ListParam.Type
				if l := node.GetLabel("type"); t != "all" && (l == nil || l.Value != t) {
					continue
				}

				fmt.Printf("%s, %s, %s\n", node.Name, node.ID, node.Address)
			}
		}
	} else {
		for _, node := range nodesMap[param.ListParam.Hostgroup] {
			t := param.ListParam.Type
			if l := node.GetLabel("type"); t != "all" && (l == nil || l.Value != t) {
				continue
			}

			fmt.Printf("%s,%s,%s\n", node.Name, node.ID, node.Address)
		}
	}

	return nil
}
