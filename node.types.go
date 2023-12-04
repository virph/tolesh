package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type NodeLabel struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (nl *NodeLabel) BuildFromString(str string) *NodeLabel {
	fields := strings.Split(str, "=")
	if len(fields) != 2 {
		return nil
	}

	if nl == nil {
		nl = new(NodeLabel)
	}

	nl.Label = fields[0]
	nl.Value = fields[1]
	return nl
}

type Node struct {
	Name    string      `json:"name"`
	ID      string      `json:"id"`
	Address string      `json:"address"`
	Labels  []NodeLabel `json:"labels"`

	Hostgroup string `json:"hostgroup"`
}

func (n *Node) BuildFromString(str string) *Node {
	n = new(Node)

	fields := strings.Fields(str)
	if len(fields) < 3 {
		return n
	}

	n.Name = fields[0]
	n.ID = fields[1]
	n.Address = fields[2]
	if len(fields) > 3 {
		for _, strLabel := range strings.Split(fields[3], ",") {
			lbl := new(NodeLabel).BuildFromString(strLabel)
			if lbl != nil {
				n.Labels = append(n.Labels, *lbl)
			}
		}
	}
	if appLabel := n.GetLabel("app"); appLabel != nil {
		n.Hostgroup = appLabel.Value
	}

	return n
}

func (n *Node) GetLabel(str string) *NodeLabel {
	if n == nil || len(n.Labels) == 0 {
		return nil
	}

	for _, nl := range n.Labels {
		if str == nl.Label {
			return &nl
		}
	}

	return nil
}

type NodesMap map[string][]Node

type NodeData struct {
	Nodes    []Node   `json:"nodes"`
	NodesMap NodesMap `json:"nodes_map"`
}

func (nm NodesMap) toJSON() string {
	b, _ := json.Marshal(nm)
	return string(b)
}

func (nd *NodeData) toConfigString() string {
	buff := bytes.NewBufferString("")
	for hostgroup, nodes := range nd.NodesMap {
		buff.WriteString(fmt.Sprintf("[%s]\n", hostgroup))
		for _, n := range nodes {
			buff.WriteString(fmt.Sprintf("root@%s\n", n.ID))
		}
		buff.WriteString("\n")
	}

	str := buff.String()
	if param.Verbose {
		fmt.Println("Config String:", str)
	}

	return str
}

func (nd *NodeData) toJSON() string {
	b, _ := json.Marshal(nd)
	return string(b)
}

func (nd *NodeData) BuildEmpty() *NodeData {
	return &NodeData{
		Nodes: []Node{},
	}
}

func (nd *NodeData) BuildFromString(str string) *NodeData {
	nd = nd.BuildEmpty()
	nd = nd.buildFromString(str)
	nd = nd.mapNodesByHostgroup()

	return nd
}

func (nd *NodeData) mapNodesByHostgroup() *NodeData {
	nd.NodesMap = make(map[string][]Node, 0)
	for _, n := range nd.Nodes {
		nd.NodesMap[n.Hostgroup] = append(nd.NodesMap[n.Hostgroup], n)
	}

	return nd
}

func (nd *NodeData) buildFromString(str string) *NodeData {
	strs := strings.Split(str, "\n")
	if len(strs) == 0 {
		return nd
	}

	rexp, _ := regexp.Compile(`^\S+-\d{1,3}-\d{1,3}-\d{1,3}-\d{1,3}.+$`)

	for _, s := range strs {
		if !rexp.MatchString(s) {
			continue
		}
		nd.Nodes = append(nd.Nodes, *new(Node).BuildFromString(s))
	}

	return nd
}
